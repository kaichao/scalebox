package vtask

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// CreateSemaphore creates a vtask-scoped semaphore.
func CreateSemaphore(name string, value int, vtaskID int64, appID int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.App.CreateVtaskSemaphore(reqCtx, &pb.CreateVtaskSemaphoreRequest{
		AppId:   int32(appID),
		VtaskId: vtaskID,
		Name:    name,
		Value:   int32(value),
	})
	if err != nil {
		return errors.WrapE(err, "semaphore-create",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "value", value)
	}
	logrus.Tracef("vtask semaphore-create: name=%s,value=%d,vtask-id=%d,app-id=%d\n",
		name, value, vtaskID, appID)
	return nil
}

// Sema holds a semaphore name-value pair for batch creation.
type Sema struct {
	Name  string
	Value int
}

// CreateSemaphores creates multiple vtask-scoped semaphores from parsed lines.
func CreateSemaphores(lines []string, vtaskID int64, appID int, batchSize int) error {
	var semas []*Sema
	re := regexp.MustCompile(`"([^"]+)":(\d+)`)
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			key := matches[1]
			var value int
			fmt.Sscanf(matches[2], "%d", &value)
			semas = append(semas, &Sema{Name: key, Value: value})
		} else {
			logrus.Warnf("Not matched semaphore :%s,\n", line)
		}
	}
	return createSemaphores(semas, vtaskID, appID, batchSize)
}

func createSemaphores(ordered []*Sema, vtaskID int64, appID int, batchSize int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID, "vtask-id", vtaskID)
	}

	total := len(ordered)
	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		for _, v := range ordered[i:end] {
			reqCtx, cancel := ctx()
			_, err := c.App.CreateVtaskSemaphore(reqCtx, &pb.CreateVtaskSemaphoreRequest{
				AppId:   int32(appID),
				VtaskId: vtaskID,
				Name:    v.Name,
				Value:   int32(v.Value),
			})
			cancel()
			if err != nil {
				if os.Getenv("CONFLICT_ACTION") != "IGNORE" {
					return errors.WrapE(err, "create prepared-semaphore",
						"app-id", appID, "vtask-id", vtaskID, "sema-name", v.Name, "sema-value", v.Value)
				}
			}
		}
		logrus.Infof("[%d..%d], %d row(s) inserted.\n", i, end-1, end-i)
	}
	return nil
}
