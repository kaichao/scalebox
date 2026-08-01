package vtask

import (
	"os"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// AddSemaphoreValue adds delta to a vtask-scoped semaphore and returns the new value.
// If not found and SEMAPHORE_AUTO_CREATE=yes, creates it with delta as initial value.
func AddSemaphoreValue(name string, delta int, vtaskID int64, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return -1, errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.App.AddVtaskSemaphore(reqCtx, &pb.AddVtaskSemaphoreRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
		Delta:   int32(delta),
	})
	if err != nil {
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			if createErr := CreateSemaphore(name, delta, vtaskID, appID); createErr != nil {
				return -1, errors.WrapE(createErr, "create semaphore",
					"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
			}
			return 0, nil
		}
		return -1, errors.WrapE(err, "semaphore not found",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
	}
	logrus.Tracef("In vtask.AddSemaphoreValue(),name=%s,vtask-id:%d,app-id:%d,delta:%d,ret-value:%d\n",
		name, vtaskID, appID, delta, resp.Value)
	return int(resp.Value), nil
}

// AddSemaphoreMapValues atomically adds deltas to multiple vtask-scoped semaphores.
func AddSemaphoreMapValues(pairs map[string]int, vtaskID int64, appID int) (map[string]int, error) {
	if len(pairs) == 0 {
		return map[string]int{}, nil
	}

	result := make(map[string]int, len(pairs))
	for name, delta := range pairs {
		v, err := AddSemaphoreValue(name, delta, vtaskID, appID)
		if err != nil {
			return result, errors.WrapE(err, "semaphore-op failed",
				"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
		}
		result[name] = v
	}
	return result, nil
}
