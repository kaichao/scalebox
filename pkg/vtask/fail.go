package vtask

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// Fail force-terminates a vtask:
//   - releases bound resource (unbind)
//   - marks root task status_code=1
//   - cascade-marks unfinished subtasks as failed
func Fail(vtaskID int64) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "vtask-id", vtaskID)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.FailVtask(reqCtx, &pb.FailVtaskRequest{VtaskId: vtaskID})
	if err != nil {
		return errors.WrapE(err, "fail vtask", "vtask-id", vtaskID)
	}

	logrus.Infof("Failed vtask: id=%d", vtaskID)
	return nil
}
