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
// func ExecShellCommand(myCmd string) string {
// 	cmd := exec.Command("bash", "-c", myCmd)
// 	output, err := cmd.Output()
// 	logrus.Infof("IN execCmd(), cmd=%s,stdout=%s\n", myCmd, string(output))
// 	if err != nil {
// 		logrus.Errorf("ERROR in execCmd(): cmd=%s,err=%v\n", myCmd, err)
// 		return ""
// 	}
// 	// remove tail \n
// 	return strings.Replace(string(output), "\n", "", -1)
// }

// ExecShellCommandWithExitCode ...
// Deprecated
// if timeout <= 0  then no timeout
// func ExecShellCommandWithExitCode(command string, timeout int) (int, string, string) {
// 	var cmd *exec.Cmd
// 	if timeout > 0 {
// 		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
// 		defer cancel()
// 		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", command)
// 	} else {
// 		cmd = exec.Command("/bin/bash", "-c", command)
// 	}
// 	var stdoutBuf, stderrBuf bytes.Buffer
// 	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
// 	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
// 	if err := cmd.Start(); err != nil {
// 		errMsg := fmt.Sprintf("start command %s failed with error:%v\n", command, err)
// 		logrus.Errorln(errMsg)
// 		return 103, "", errMsg
// 	}
// 	exitCode := 0
// 	var errMsg string
// 	if err := cmd.Wait(); err != nil {
// 		if exitErr, ok := err.(*exec.ExitError); ok {
// 			// timeout : exit_code = -1
// 			errMsg = fmt.Sprintf("Exit Status: %d,exit err_message:%s\ncmd:%s.\n",
// 				exitErr.ExitCode(), exitErr.Error(), command)
// 			logrus.Warnln(errMsg)
// 			exitCode = exitErr.ExitCode()
// 			if exitCode == -1 {
// 				// timeout !
// 				exitCode = 100
// 			}
// 		} else {
// 			errMsg = fmt.Sprintf("wait command '%s' failed with error:%v\n", command, err)
// 			logrus.Errorln(errMsg)
// 			exitCode = 105
// 		}
// 	}
// 	return exitCode, string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()) + errMsg
// }

// ExecCommandReturnAll executes a command and returns its exit code, stdout, stderr, and any error.
//
// Params:
//   - command: the command string to execute
//   - timeout: timeout in seconds (0 for no timeout)
//
// Returns: (exitCode, stdout, stderr, err)
//   - exitCode：命令的退出码（0 表示成功，非零表示命令失败或超时等）
//   - stdout：标准输出
//   - stderr：标准错误
//   - err：执行过程中遇到的错误（如管道创建失败、命令启动失败、超时等）。若命令以非零退出码结束，err 为 nil
//   - 管道创建或命令启动失败时，返回退出码 125 和具体的 error
//   - 超时情况下，返回退出码 124 和 err = "command timed out"
//   - 命令以非零退出码结束时，返回该退出码，err 为 nil
//   - 其他未预期的错误通过 err 返回，退出码为 125
func ExecCommandReturnAll(command string, timeout int) (int, string, string, error) {
	baseCtx := context.Background()
	ctx := baseCtx
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(baseCtx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	// 创建命令并支持进程组终止
	cmd := exec.CommandContext(ctx, "/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 获取输出管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("capture stdout pipe failed: %v", err)
		return 125, "", "", fmt.Errorf("capture stdout pipe failed: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("capture stderr pipe failed: %v", err)
		return 125, "", "", fmt.Errorf("capture stderr pipe failed: %v", err)
	}

	// 异步捕获输出
	var stdoutBuf, stderrBuf bytes.Buffer
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(&stdoutBuf, stdoutPipe) // 移除 os.Stdout
		// io.Copy(io.MultiWriter(os.Stdout, &stdoutBuf), stdoutPipe)
	}()
	go func() {
		defer wg.Done()
		io.Copy(&stderrBuf, stderrPipe) // 移除 os.Stderr
		// io.Copy(io.MultiWriter(os.Stderr, &stderrBuf), stderrPipe)
	}()

	// 超时后终止进程组
	if timeout > 0 {
		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded && cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
		}()
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		logrus.Errorf("start command failed: %v", err)
		return 125, "", "", fmt.Errorf("start command failed: %v", err)
	}

	// 等待命令结束
	waitErr := cmd.Wait()
	// 确保输出复制完成
	wg.Wait()

	// 处理退出码和错误
	var exitCode int
	var retErr error
	if waitErr != nil {
		if ctx.Err() == context.DeadlineExceeded {
			exitCode = 124
			retErr = fmt.Errorf("command timed out")
		} else if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			// 处理信号终止
			if exitCode == -1 {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					if status.Signaled() {
						exitCode = 128 + int(status.Signal())
					}
				}
			}
			// 命令以非零退出码结束，不是错误
			retErr = nil
		} else {
			exitCode = 125
			retErr = waitErr
		}
	} else {
		exitCode = 0
		retErr = nil
	}

	return exitCode, stdoutBuf.String(), stderrBuf.String(), retErr
}

// ExecCommandReturnExitCode ...
func ExecCommandReturnExitCode(command string, timeout int) (int, error) {
	code, stdout, stderr, err := ExecCommandReturnAll(command, timeout)
	fmt.Printf("exec command:%s\n stdout:\n%s\n", command, stdout)
	fmt.Fprintf(os.Stderr, "exec command: %s\n stderr:\n%s\n", command, stderr)
	return code, err
}

// ExecCommandReturnStdout ...
func ExecCommandReturnStdout(command string, timeout int) (string, error) {
	code, stdout, stderr, err := ExecCommandReturnAll(command, timeout)
	if code != 0 {
		fmt.Fprintf(os.Stderr, "exec command:%s\nexit-code=%d\n", command, code)
		fmt.Fprintf(os.Stderr, "stdout:\n%s\n", stdout)
		stdout = ""
	}
	fmt.Fprintf(os.Stderr, "exec command:\n%s\n%s\n", command, stderr)

	// remove leading/tail space
	return strings.TrimSpace(stdout), err
}

// ExecWithRetries ...
func ExecWithRetries(cmd string, numRetries int, timeout int) int {
	delay := 10 * time.Second
	var code int
	for i := 0; i < numRetries; i++ {
		code, _ = ExecCommandReturnExitCode(cmd, timeout)
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
