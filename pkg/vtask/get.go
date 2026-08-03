package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// GetInfo returns detailed information about a vtask,
// including derived status, bound resource, and subtask statistics.
func GetInfo(vtaskID int64) (*Info, error) {
	c, err := client.Default()
	if err != nil {
		return nil, errors.WrapE(err, "get client", "vtask-id", vtaskID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.GetVtask(reqCtx, &pb.GetVtaskRequest{VtaskId: vtaskID})
	if err != nil {
		return nil, errors.WrapE(err, "get vtask", "vtask-id", vtaskID)
	}

	return &Info{
		ID:            resp.Id,
		AppID:         resp.AppId,
		Body:          resp.Body,
		Status:        resp.Status,
		ModuleName:    resp.ModuleName,
		BoundResource: resp.BoundResource,
		SemaName:      resp.SemaName,
		SubtaskCount:  resp.SubtaskCount,
		FinishedCount: resp.FinishedCount,
		CreatedAt:     resp.CreatedAt.AsTime(),
	}, nil
}
