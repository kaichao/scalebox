package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// AddSubtask creates a subtask in an existing vtask, establishing
// the vtask column parent-child relationship.
//
// Parameters:
//   - appID: application ID
//   - module: target module name (e.g. "vtask-core")
//   - body: task body
//   - headers: task headers (must include _vtask_id for vtask linking)
//
// Returns the newly created task ID.
func AddSubtask(appID int, module string, body string, headers map[string]string) (int64, error) {
	c, err := client.Default()
	if err != nil {
		return 0, errors.WrapE(err, "get client", "app-id", appID, "module", module)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.AddVtaskSubtask(reqCtx, &pb.AddVtaskSubtaskRequest{
		AppId:   int32(appID),
		Module:  module,
		Body:    body,
		Headers: headers,
	})
	if err != nil {
		return 0, errors.WrapE(err, "add vtask subtask",
			"app-id", appID, "module", module, "body", body)
	}

	logrus.Debugf("Added vtask subtask: app=%d module=%s task=%d", appID, module, resp.TaskId)
	return resp.TaskId, nil
}
