package misc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"

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
//
//	func ExecShellCommand(cmdText string) {
//		cmd := exec.Command("bash", "-c", cmdText)
//		if out, err := cmd.Output(); err != nil {
//			fmt.Fprintln(os.Stderr, err)
//		} else {
//			fmt.Println(string(out))
//		}
//	}
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

// GetTextFileLines ...
func GetTextFileLines(textFile string) ([]string, error) {
	if _, err := os.Stat(textFile); err != nil {
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
