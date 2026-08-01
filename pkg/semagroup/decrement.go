package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// Decrement subtracts 1 from the semaphore with maximum value in the group.
func Decrement(semaExpr string, appID int) (string, int, error) {
	c, err := client.Default()
	if err != nil {
		return "", 0, errors.WrapE(err, "get client", "app-id", appID, "sema-expr", semaExpr)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.DecrementSemagroup(reqCtx, &pb.SemagroupRequest{
		AppId:    int32(appID),
		SemaExpr: semaExpr,
	})
	if err != nil {
		return "", 0, errors.WrapE(err, "semagroup decrement",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return resp.Name, int(resp.Value), nil
}
