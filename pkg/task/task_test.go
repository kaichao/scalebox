package task_test

import (
	"os"
	"testing"

	"github.com/kaichao/scalebox/pkg/task"
)

func TestAddTask(t *testing.T) {
	envVars := map[string]string{
		"SINK_JOB":    "scatter",
		"JOB_ID":      "48",
		"APP_ID":      "",
		"GRPC_SERVER": "10.0.6.100",
	}
	task.Add("001", "", envVars)

	envVars = map[string]string{
		"SINK_JOB":    "scatter",
		"JOB_ID":      "",
		"APP_ID":      "28",
		"GRPC_SERVER": "10.0.6.100",
	}
	task.Add("002", "", envVars)

	envVars = map[string]string{
		"SINK_JOB":    "",
		"JOB_ID":      "46",
		"APP_ID":      "",
		"GRPC_SERVER": "10.0.6.100",
	}
	task.Add("003", "", envVars)

	os.Setenv("JOB_ID", "47")
	envVars = map[string]string{
		"SINK_JOB": "",
		// "JOB_ID":      "46",
		"APP_ID":      "",
		"GRPC_SERVER": "10.0.6.100",
	}
	task.Add("004", "", envVars)
}
func TestAddTasks(t *testing.T) {
	envVars := map[string]string{
		"SINK_JOB":    "",
		"JOB_ID":      "46",
		"APP_ID":      "",
		"GRPC_SERVER": "10.0.6.100",
	}
	bodies := []string{"100", "101"}
	task.AddTasks(bodies, "", envVars)
}
