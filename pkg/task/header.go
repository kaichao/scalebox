package task

import (
	"database/sql"
	"fmt"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetTaskHeader ...
func GetTaskHeader(taskID int64, name string) (string, error) {
	sqlText := `SELECT headers->>$1 FROM t_task WHERE id=$2`
	var value string
	err := postgres.GetDB().QueryRow(sqlText, name, taskID).Scan(&value)
	if err == nil {
		return value, nil
	}
	if err == sql.ErrNoRows {
		err = errors.WrapE(err, "header not found", "task-id", taskID, "header", name)
	} else {
		err = errors.WrapE(err, "get-task-header", "task-id", taskID, "header", name)
	}
	return "", err

}

// SetTaskHeader ...
func SetTaskHeader(taskID int64, name string, value string) error {
	sqlText := `
		UPDATE t_task
		SET headers = jsonb_set(headers, $1, $2)
		WHERE id = $3
	`
	jsonPath := fmt.Sprintf(`{%s}`, name)
	newValue := fmt.Sprintf(`"%s"`, value)
	_, err := postgres.GetDB().Exec(sqlText, jsonPath, newValue, taskID)
	if err != nil {
		return errors.WrapE(err, "set-task-header",
			"task-id", taskID, "header", name, "value", value)
	}
	return nil
}

// RemoveTaskHeader ...
func RemoveTaskHeader(taskID int64, name string) error {
	sqlText := `
		UPDATE t_task
		SET headers = headers - $1
		WHERE id = $2
	`
	_, err := postgres.GetDB().Exec(sqlText, name, taskID)
	if err != nil {
		return errors.WrapE(err, "remove-task-header",
			"task-id", taskID, "header", name)
	}
	return nil
}
