package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
)

// ListItem is a simplified view of a vtask list entry.
type ListItem struct {
	ID         int64
	ModuleName string
	FromModule string
	Status     string
	Body       string
}

// List returns all vtasks for an app (WHERE vtask = id).
func List(appID int) ([]ListItem, error) {
	c, err := client.Default()
	if err != nil {
		return nil, errors.WrapE(err, "get client", "app-id", appID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.ListVTasks(reqCtx, &pb.ListVTasksRequest{
		AppId: int32(appID),
	})
	if err != nil {
		return nil, errors.WrapE(err, "list vtasks", "app-id", appID)
	}

	items := make([]ListItem, len(resp.Items))
	for i, t := range resp.Items {
		items[i] = ListItem{
			ID:         t.Id,
			ModuleName: t.ModuleName,
			FromModule: t.FromModule,
			Status:     t.Status,
			Body:       t.Body,
		}
	}
	return items, nil
}

// ListSubtasks returns all subtasks for a vtask (WHERE vtask = $1 AND vtask <> id).
func ListSubtasks(vtaskID int64) ([]ListItem, error) {
	c, err := client.Default()
	if err != nil {
		return nil, errors.WrapE(err, "get client", "vtask-id", vtaskID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.ListVTaskSubtasks(reqCtx, &pb.ListVTaskSubtasksRequest{
		VtaskId: vtaskID,
	})
	if err != nil {
		return nil, errors.WrapE(err, "list vtask subtasks", "vtask-id", vtaskID)
	}

	items := make([]ListItem, len(resp.Items))
	for i, t := range resp.Items {
		items[i] = ListItem{
			ID:         t.Id,
			ModuleName: t.ModuleName,
			FromModule: t.FromModule,
			Status:     t.Status,
			Body:       t.Body,
		}
	}
	return items, nil
}
