package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// UnbindResource releases a vtask's bound compute resource.
//
// At least one of vtaskID or semaName must be specified:
//   - vtaskID: server looks up _vtask_size_sema from root task headers
//   - semaName: direct semaphore name (e.g. ":host_vtask_size:vtask-head:n0-0")
func UnbindResource(appID int, vtaskID int64, semaName string) error {
	if vtaskID == 0 && semaName == "" {
		return errors.E("either vtask_id or sema_name is required",
			"app-id", appID)
	}

	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.UnbindVtaskResource(reqCtx, &pb.UnbindVtaskResourceRequest{
		AppId:    int32(appID),
		VtaskId:  vtaskID,
		SemaName: semaName,
	})
	if err != nil {
		return errors.WrapE(err, "unbind vtask resource",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", semaName)
	}

	logrus.Debugf("Unbound vtask: app=%d vtask=%d sema=%s", appID, vtaskID, semaName)
	return nil
}
