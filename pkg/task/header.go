package task

import (
	"context"
	"time"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// GetTaskHeader returns the value of a task header field.
func GetTaskHeader(taskID int64, name string) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "", errors.WrapE(err, "get client", "task-id", taskID, "header", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.GetTaskHeader(reqCtx, &pb.GetTaskHeaderRequest{
		TaskId: taskID,
		Name:   name,
	})
	if err != nil {
		return "", errors.WrapE(err, "get-task-header",
			"task-id", taskID, "header", name)
	}
	return resp.Value, nil
}

// SetTaskHeader sets a task header field to the given value.
func SetTaskHeader(taskID int64, name string, value string) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "task-id", taskID, "header", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.SetTaskHeader(reqCtx, &pb.SetTaskHeaderRequest{
		TaskId: taskID,
		Name:   name,
		Value:  value,
	})
	if err != nil {
		return errors.WrapE(err, "set-task-header",
			"task-id", taskID, "header", name, "value", value)
	}
	return nil
}

// RemoveTaskHeader removes a task header field.
func RemoveTaskHeader(taskID int64, name string) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "task-id", taskID, "header", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.RemoveTaskHeader(reqCtx, &pb.RemoveTaskHeaderRequest{
		TaskId: taskID,
		Name:   name,
	})
	if err != nil {
		return errors.WrapE(err, "remove-task-header",
			"task-id", taskID, "header", name)
	}
	return nil
}
