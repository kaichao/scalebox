package task

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kaichao/gopkg/exec"
	"github.com/kaichao/scalebox/pkg/common"
	"github.com/sirupsen/logrus"
)

// Add ...
func Add(body string, headers string, envVars map[string]string) int {
	parts := make([]string, 0, len(envVars))
	for k, v := range envVars {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
	}

	if headers == "" {
		headers = "{}"
	}
	cmd := fmt.Sprintf(`%s scalebox task add --headers='%s' %s`,
		strings.Join(parts, " "), headers, body)
	code, err := exec.RunReturnExitCode(cmd, 15)
	if err != nil {
		logrus.Errorf("tasks-add, err-info:%v", err)
		return -1
	}
	return code
}

// AddWithMapHeaders ...
func AddWithMapHeaders(body string, headers map[string]string, envVars map[string]string) int {
	return Add(body, mapToCleanJSON(headers), envVars)
}

// AddTasks 增加一组task
// 环境变量：
// - SINK_MODULE:
// - MODULE_ID:
// - APP_ID:
// - REMOTE_SERVER:
// - TIMEOUT_SECONDS
func AddTasks(bodies []string, headers string, envVars map[string]string) int {
	parts := make([]string, 0, len(envVars))
	for k, v := range envVars {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
	}

	taskFile := "my-tasks.txt"
	for _, m := range bodies {
		common.AppendToFile(taskFile, m)
	}
	if headers == "" {
		headers = "{}"
	}
	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT_SECONDS"))
	if timeout <= 0 {
		timeout = 60
	}
	cmd := fmt.Sprintf(`%s scalebox task add --headers='%s' --task-file=my-tasks.txt`,
		strings.Join(parts, " "), headers)
	code, err := exec.RunReturnExitCode(cmd, timeout)
	if err != nil {
		logrus.Errorf("tasks-add, err-info:%v", err)
		return -2
	}
	if code != 0 {
		return code
	}

	if err := os.Remove(taskFile); err != nil {
		logrus.Errorf("remove file %s\n", taskFile)
		return 1
	}
	return 0
}

// AddTasksWithMapHeaders ...
func AddTasksWithMapHeaders(bodies []string, headers map[string]string, envVars map[string]string) int {
	return AddTasks(bodies, mapToCleanJSON(headers), envVars)
}

func mapToCleanJSON(m map[string]string) string {
	// 直接 Marshal，不去除 key/value 中的空格
	jsonBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		logrus.Errorf("map to json, err-info:%v\n", err)
		return "{}"
	}

	jsonStr := string(jsonBytes)

	// 去除所有双引号之外的空白字符
	re := regexp.MustCompile(`"(?:\\.|[^"\\])*"|[\s]+`)
	cleaned := re.ReplaceAllStringFunc(jsonStr, func(s string) string {
		if len(s) > 0 && s[0] == '"' {
			return s // 保留字符串字面量（引号包裹部分）
		}
		return "" // 去掉引号外的空白
	})

	return cleaned
}
