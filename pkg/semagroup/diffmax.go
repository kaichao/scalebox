package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// DiffMax returns the difference between the group maximum and current semaphore.
func DiffMax(semaExpr string, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return 0, errors.WrapE(err, "get client", "app-id", appID, "sema-expr", semaExpr)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.DiffSemagroupMax(reqCtx, &pb.SemagroupRequest{
		AppId:    int32(appID),
		SemaExpr: semaExpr,
	})
	if err != nil {
		return 0, errors.WrapE(err, "semagroup diffmax",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return int(resp.Value), nil
}
