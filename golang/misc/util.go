package misc

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// AppendToFile ...
func AppendToFile(fileName string, line string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file %s error,err-info:%v\n", fileName, err)
		fmt.Fprintln(os.Stderr, os.Args)
		os.Exit(3)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(line + "\n")
	writer.Flush()
}

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
	// 删除尾部的\n
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

// GetTextFileLines ...
func GetTextFileLines(textFile string) ([]string, error) {
	if _, err := os.Stat(textFile); err != nil {
		_, ok := err.(*os.PathError)
		if ok && strings.Contains(err.Error(), "no such file or directory") {
			// file not exists
			return []string{}, nil
		}
		return []string{}, err
	}

	fileContent, err := ioutil.ReadFile(textFile)
	if err != nil {
		return []string{}, fmt.Errorf("Read file error, filename:%s, err:%v", textFile, err)
	}
	var lines []string
	for _, line := range strings.Split(string(fileContent), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// SplitCommaWithEscapeSupport ..
func SplitCommaWithEscapeSupport(s string) []string {
	var ret []string
	i0 := 0
	i := 0
	for ; i < len(s); i++ {
		if s[i] == ',' && i > 0 && s[i-1] != '\\' {
			ret = append(ret, strings.ReplaceAll(s[i0:i], "\\,", ","))
			i0 = i + 1
		}
	}
	if i0 < i {
		ret = append(ret, strings.ReplaceAll(s[i0:i], "\\,", ","))
	}
	return ret
}

// GetFunctionName ...
func GetFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

// IsRunnable ...
func IsRunnable(runFile string) bool {
	stat, err := os.Stat(runFile)
	if err != nil {
		return false
	}
	return stat.Mode()&0111 != 0
}
