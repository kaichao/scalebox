package semagroup

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

// GetMin returns the semaphore with minimum value in the group matching semaExpr.
func GetMin(semaExpr string, appID int) (string, int, error) {
	c, err := client.Default()
	if err != nil {
		return "", 0, errors.WrapE(err, "get client", "app-id", appID, "sema-expr", semaExpr)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.GetSemagroupMin(reqCtx, &pb.SemagroupRequest{
		AppId:    int32(appID),
		SemaExpr: semaExpr,
	})
	if err != nil {
		return "", 0, errors.WrapE(err, "semagroup min",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return resp.Name, int(resp.Value), nil
}
