package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// Increment adds 1 to the semaphore with minimum value in the group.
func Increment(semaExpr string, appID int) (string, int, error) {
	c, err := client.Default()
	if err != nil {
		return "", 0, errors.WrapE(err, "get client", "app-id", appID, "sema-expr", semaExpr)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.IncrementSemagroup(reqCtx, &pb.SemagroupRequest{
		AppId:    int32(appID),
		SemaExpr: semaExpr,
	})
	if err != nil {
		return "", 0, errors.WrapE(err, "semagroup increment",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return resp.Name, int(resp.Value), nil
}
