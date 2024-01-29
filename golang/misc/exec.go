package misc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// ExecShellCommand ...
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
