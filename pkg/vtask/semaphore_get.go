package vtask

import (
	"os"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// GetSemaphore returns the current value of a vtask-scoped semaphore.
// If not found and SEMAPHORE_AUTO_CREATE=yes, creates it with value 0.
func GetSemaphore(name string, vtaskID int64, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return -1, errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.GetVtaskSemaphore(reqCtx, &pb.GetVtaskSemaphoreRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
	})
	if err != nil {
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			if createErr := CreateSemaphore(name, 0, vtaskID, appID); createErr != nil {
				return -1, errors.WrapE(createErr, "create semaphore",
					"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
			}
			return 0, nil
		}
		return -1, errors.WrapE(err, "semaphore not found",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}
	logrus.Tracef("In vtask.GetSemaphore(),name=%s,vtask-id:%d,app-id:%d,value:%d\n",
		name, vtaskID, appID, resp.Value)
	return int(resp.Value), nil
}
