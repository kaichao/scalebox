package task

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/gopkg/exec"
	"github.com/kaichao/scalebox/pkg/common"
	"github.com/sirupsen/logrus"
)

// Add 增加单个task
// 环境变量：
// - SINK_MODULE:
// - MODULE_ID:
// - APP_ID:
// - REMOTE_SERVER:
// - TIMEOUT_SECONDS:
// - DIRECT_WRITE_IN_AGENT:
func Add(body string, headers string, envVars map[string]string) (taskID int64, err error) {
	parts := make([]string, 0, len(envVars))
	for k, v := range envVars {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
	}

	if headers == "" {
		headers = "{}"
	}
	// Escape single quotes in headers for shell command
	// Replace ' with '\'' (end single quote, escaped single quote, start single quote)
	escapedHeaders := strings.ReplaceAll(headers, "'", "'\\''")
	cmd := fmt.Sprintf(`%s scalebox task add --headers='%s' %s`,
		strings.Join(parts, " "), escapedHeaders, body)
	code, stdout, stderr, err := exec.RunReturnAll(cmd, 30)
	logrus.Tracef("In task.Add(),headers:%s\ncmd:'%s'\nexit-code:%d\nstdout:%s\nstderr:%s\nerr:%v\n",
		headers, cmd, code, stdout, stderr, err)

	if err != nil || code != 0 {
		return -1, errors.WrapE(err, "exec-cmd",
			"cmd", cmd, "code", code, "stdout", stdout, "stderr", stderr)
	}

	if inAgent() && envVars["DIRECT_WRITE_IN_AGENT"] != "yes" {
		return 0, nil
	}

	num, err := fmt.Sscanf(strings.TrimSpace(stdout), `{"task_id":%d}`, &taskID)
	if err != nil || num != 1 {
		return -2, errors.WrapE(err, "parse-task-id by fmt.Sscanf()",
			"stdout", stdout, "num-parsed", num)
	}
	return taskID, nil
}

// AddWithMapHeaders ...
func AddWithMapHeaders(body string, headers map[string]string, envVars map[string]string) (int64, error) {
	return Add(body, mapToCleanJSON(headers), envVars)
}

// AddTasks 增加一组task
// 环境变量：
// - SINK_MODULE:
// - MODULE_ID:
// - APP_ID:
// - REMOTE_SERVER:
// - TIMEOUT_SECONDS:
// - DIRECT_WRITE_IN_AGENT:
func AddTasks(bodies []string, headers string, envVars map[string]string) (int, error) {
	parts := make([]string, 0, len(envVars))
	for k, v := range envVars {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, v))
	}

	// Create a unique temporary file for tasks
	taskFile, err := os.CreateTemp("", "scalebox-tasks-*.txt")
	if err != nil {
		logrus.Errorf("create temp file error: %v\n", err)
		return -1, err
	}
	taskFilePath := taskFile.Name()
	defer os.Remove(taskFilePath) // Clean up file after function returns
	taskFile.Close()

	// Write tasks to file
	for _, m := range bodies {
		common.AppendToFile(taskFilePath, m)
	}

	if headers == "" {
		headers = "{}"
	}
	// Escape single quotes in headers for shell command
	// Replace ' with '\'' (end single quote, escaped single quote, start single quote)
	escapedHeaders := strings.ReplaceAll(headers, "'", "'\\''")
	timeout, _ := strconv.Atoi(os.Getenv("TIMEOUT_SECONDS"))
	if timeout <= 0 {
		timeout = 60
	}
	cmd := fmt.Sprintf(`%s scalebox task add --headers='%s' --task-file=%s`,
		strings.Join(parts, " "), escapedHeaders, taskFilePath)
	code, stdout, stderr, err := exec.RunReturnAll(cmd, 15)
	logrus.Tracef("In task.AddTask(),cmd:'%s'\ntask-body:%v\nexit-code:%d\nstdout:%s\nstderr:%s\nerr:%v\n",
		cmd, bodies, code, stdout, stderr, err)

	if err != nil || code != 0 {
		return -1, errors.WrapE(err, "exec-cmd",
			"cmd", cmd, "code", code, "stdout", stdout, "stderr", stderr)
	}

	if inAgent() && envVars["DIRECT_WRITE_IN_AGENT"] != "yes" {
		return -9, nil
	}

	var numTasks int
	num, err := fmt.Sscanf(strings.TrimSpace(stdout), `{"num_tasks":%d}`, &numTasks)
	if err != nil || num != 1 {
		return -2, errors.WrapE(err, "parse-num-tasks by fmt.Sscanf()",
			"stdout", stdout, "num-parsed", num)
	}

	return numTasks, nil
}

// AddTasksWithMapHeaders ...
func AddTasksWithMapHeaders(bodies []string, headers map[string]string, envVars map[string]string) (int, error) {
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

func inAgent() bool {
	return (os.Getenv("PLAT_MODULE_NAME") != "")
}
