package misc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// ExecShellCommand ...
// Deprecated
//
//	return stdout
func ExecShellCommand(myCmd string) string {
	cmd := exec.Command("bash", "-c", myCmd)
	output, err := cmd.Output()
	logrus.Infof("IN execCmd(), cmd=%s,stdout=%s\n", myCmd, string(output))
	if err != nil {
		logrus.Errorf("ERROR in execCmd(): cmd=%s,err=%v\n", myCmd, err)
		return ""
	}
	// remove tail \n
	return strings.Replace(string(output), "\n", "", -1)
}

// ExecShellCommandWithExitCode ...
// Deprecated
// if timeout <= 0  then no timeout
func ExecShellCommandWithExitCode(command string, timeout int) (int, string, string) {
	var cmd *exec.Cmd
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", command)
	} else {
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		errMsg := fmt.Sprintf("start command %s failed with error:%v\n", command, err)
		logrus.Errorln(errMsg)
		return 103, "", errMsg
	}
	exitCode := 0
	var errMsg string
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// timeout : exit_code = -1
			errMsg = fmt.Sprintf("Exit Status: %d,exit err_message:%s\ncmd:%s.\n",
				exitErr.ExitCode(), exitErr.Error(), command)
			logrus.Warnln(errMsg)
			exitCode = exitErr.ExitCode()
			if exitCode == -1 {
				// timeout !
				exitCode = 100
			}
		} else {
			errMsg = fmt.Sprintf("wait command '%s' failed with error:%v\n", command, err)
			logrus.Errorln(errMsg)
			exitCode = 105
		}
	}

	return exitCode, string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()) + errMsg
}

// ExecCommandReturnAll ...
//
//	params:
//		command : command string
//		timeout : timeout seconds for golang-timout, 0 for none
//	return (exit-code, stdout, stderr)
func ExecCommandReturnAll(command string, timeout int) (int, string, string) {
	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} // 支持进程组终止

	// 获取输出管道并检查错误
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("capture stdout pipe failed: %v", err)
		return 125, "", fmt.Sprintf("[EXEC INTERNAL ERROR]: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("capture stderr pipe failed: %v", err)
		return 125, "", fmt.Sprintf("[EXEC INTERNAL ERROR]: %v", err)
	}

	// 异步捕获输出并同步等待
	var stdoutBuf, stderrBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(io.MultiWriter(os.Stdout, &stdoutBuf), stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		io.Copy(io.MultiWriter(os.Stderr, &stderrBuf), stderrPipe)
	}()

	// 超时后终止整个进程组
	if timeout > 0 {
		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded && cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) // 终止进程组
			}
		}()
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		errMsg := fmt.Sprintf("[START FAILED]: %v", err)
		logrus.Errorln(errMsg)
		return 125, "", errMsg
	}

	// 等待命令结束
	var exitCode int
	waitErr := cmd.Wait()
	// 确保所有输出复制完成
	wg.Wait()
	// 处理退出码
	if waitErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124 // 明确标记超时
		} else if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			// 处理信号终止场景（如非超时的 SIGKILL）
			if exitCode == -1 {
				exitCode = 128 + int(exitErr.Sys().(syscall.WaitStatus).Signal())
			}
		} else {
			exitCode = 125 // 通用错误码
		}
	}

	return exitCode, stdoutBuf.String(), stderrBuf.String()
}

// ExecCommandReturnExitCode ...
func ExecCommandReturnExitCode(command string, timeout int) int {
	code, stdout, stderr := ExecCommandReturnAll(command, timeout)
	fmt.Printf("exec command:%s\n stdout:\n%s\n", command, stdout)
	fmt.Fprintf(os.Stderr, "exec command: %s\n stderr:\n%s\n", command, stderr)
	return code
}

// ExecCommandReturnStdout ...
func ExecCommandReturnStdout(command string, timeout int) string {
	code, stdout, stderr := ExecCommandReturnAll(command, timeout)
	if code != 0 {
		fmt.Fprintf(os.Stderr, "exec command:%s\nexit-code=%d\n", command, code)
		fmt.Fprintf(os.Stderr, "stdout:\n%s\n", stdout)
		stdout = ""
	}
	fmt.Fprintf(os.Stderr, "exec command:\n%s\n%s\n", command, stderr)

	// remove leading/tail space
	return strings.TrimSpace(stdout)
}

// ExecWithRetries ...
func ExecWithRetries(cmd string, numRetries int, timeout int) int {
	delay := 10 * time.Second
	var code int
	for i := 0; i < numRetries; i++ {
		code = ExecCommandReturnExitCode(cmd, timeout)
		if code == 0 {
			return code
		}
		fmt.Printf("num-of-retries:%d,cmd=%s\n", i+1, cmd)
		time.Sleep(delay)
		delay *= 2
		timeout *= 2
	}
	return code
}
