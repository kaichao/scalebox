package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// BindResource atomically binds a vtask to a compute resource group.
//
// Parameters:
//   - appID: application ID
//   - moduleName: target head module (empty = auto-detect vtask_role=head module)
//
// Returns:
//   - resource: hostname (HOST-BOUND) / slot_seq (GROUP-BOUND) / "" (DEFAULT)
//   - mode: task_dist_mode actually used
func BindResource(appID int, moduleName string) (resource string, mode string, err error) {
	c, err := client.Default()
	if err != nil {
		return "", "", errors.WrapE(err, "get client", "app-id", appID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.BindVtaskResource(reqCtx, &pb.BindVtaskResourceRequest{
		AppId:      int32(appID),
		ModuleName: moduleName,
	})
	if err != nil {
		return "", "", errors.WrapE(err, "bind vtask resource",
			"app-id", appID, "module", moduleName)
	}

	logrus.Debugf("Bound vtask: app=%d module=%s mode=%s resource=%s",
		appID, moduleName, resp.Mode, resp.Resource)

	return resp.Resource, resp.Mode, nil
}
