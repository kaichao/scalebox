package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// DeleteSemaphore deletes a vtask-scoped semaphore.
func DeleteSemaphore(name string, vtaskID int64, appID int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.DeleteVtaskSemaphore(reqCtx, &pb.DeleteVtaskSemaphoreRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
	})
	if err != nil {
		return errors.WrapE(err, "delete semaphore",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}

	logrus.Debugf("Deleted vtask semaphore: name=%s vtask=%d app=%d",
		name, vtaskID, appID)
	return nil
}
