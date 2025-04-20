package task

import (
	"database/sql"
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetTaskHeader ...
func GetTaskHeader(taskID int, name string) (string, error) {
	sqlText := `SELECT headers->>$1 FROM t_task WHERE id=$2`
	var value string
	err := postgres.GetDB().QueryRow(sqlText, name, taskID).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("header %s not-found in task-id:%d\n", name, taskID)
		} else {
			logrus.Errorf("Unable to get-task-header, task-id:%d, header=%s, err=%T, err=%v\n",
				taskID, name, err, err)
		}
		return "", err
	}
	return value, nil
}

// SetTaskHeader ...
func SetTaskHeader(taskID int, name string, value string) error {
	sqlText := `
		UPDATE t_task
		SET headers = jsonb_set(headers, $1, $2)
		WHERE id = $3
	`
	jsonPath := fmt.Sprintf(`{%s}`, name)
	newValue := fmt.Sprintf(`"%s"`, value)
	_, err := postgres.GetDB().Exec(sqlText, jsonPath, newValue, taskID)
	if err != nil {
		logrus.Errorf("Unable to set-task-header, task-id=%d, name=%s, value=%s, err=%T, err=%v\n",
			taskID, name, value, err, err)
		return err
	}
	return nil
}

// RemoveTaskHeader ...
func RemoveTaskHeader(taskID int, name string) error {
	sqlText := `
		UPDATE t_task
		SET headers = headers - $1
		WHERE id = $2
	`
	_, err := postgres.GetDB().Exec(sqlText, name, taskID)
	if err != nil {
		logrus.Errorf("Unable to remove-task-header, task-id=%d, name=%s,  err=%T, err=%v\n",
			taskID, name, err, err)
		return err
	}
	return nil
}
