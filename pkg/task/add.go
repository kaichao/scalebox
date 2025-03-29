package task

import (
	"fmt"
	"os"

	"github.com/kaichao/scalebox/pkg/misc"
	"github.com/sirupsen/logrus"
)

// Add ...
func Add(sinkJob string, message string, headers string) int {
	if headers == "" {
		headers = "{}"
	}
	cmd := fmt.Sprintf(`scalebox task add --sink-job=%s --headers='%s' %s`,
		sinkJob, headers, message)
	code, err := misc.ExecCommandReturnExitCode(cmd, 15)
	if err != nil {
		logrus.Errorf("tasks-add, err-info:%v", err)
		return -1
	}
	return code
}

// AddWithMapHeaders ...
func AddWithMapHeaders(sinkJob string, message string, headers map[string]string) int {
	return 0
}

// AddTasks ...
func AddTasks(sinkJob string, messages []string, headers string, timeout int) int {
	taskFile := "my-tasks.txt"
	for _, m := range messages {
		misc.AppendToFile(taskFile, m)
	}
	if headers == "" {
		headers = "{}"
	}
	cmd := fmt.Sprintf(`scalebox task add --headers='%s' --sink-job=%s --task-file=my-tasks.txt`, headers, sinkJob)
	code, err := misc.ExecCommandReturnExitCode(cmd, timeout)
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
