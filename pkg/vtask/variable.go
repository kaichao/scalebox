package vtask

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

// GetVariable returns the value of a vtask-scoped variable.
func GetVariable(name string, vtaskID int64, appID int) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "", errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID, "var-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.GetVtaskVariable(reqCtx, &pb.GetVtaskVariableRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
	})
	if err != nil {
		return "", errors.WrapE(err, "get variable",
			"app-id", appID, "vtask-id", vtaskID, "var-name", name)
	}
	return resp.Value, nil
}

// SetVariable creates or updates a vtask-scoped variable.
func SetVariable(name string, value string, vtaskID int64, appID int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID, "var-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.SetVtaskVariable(reqCtx, &pb.SetVtaskVariableRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
		Value:   value,
	})
	if err != nil {
		return errors.WrapE(err, "set-variable",
			"app-id", appID, "vtask-id", vtaskID, "var-name", name, "var-value", value)
	}
	return nil
}
