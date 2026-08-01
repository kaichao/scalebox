package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// DiffMin returns the difference between the current semaphore and group minimum.
func DiffMin(semaExpr string, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return 0, errors.WrapE(err, "get client", "app-id", appID, "sema-expr", semaExpr)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.DiffSemagroupMin(reqCtx, &pb.SemagroupRequest{
		AppId:    int32(appID),
		SemaExpr: semaExpr,
	})
	if err != nil {
		return 0, errors.WrapE(err, "semagroup diffmin",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return int(resp.Value), nil
}
