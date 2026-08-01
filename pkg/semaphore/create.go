package semaphore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/common"

	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Create creates a semaphore with the given name and initial value.
func Create(name string, value int, appID int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID, "sema-name", name)
	}

	conflictAction := os.Getenv("CONFLICT_ACTION")

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.Coordination.CreateSemaphore(reqCtx, &pb.CreateSemaphoreRequest{
		AppId:          int32(appID),
		Name:           name,
		Value:          int32(value),
		ConflictAction: conflictAction,
	})
	if err != nil {
		return errors.WrapE(err, "semaphore-create",
			"app-id", appID, "sema-name", name, "value", value, "conflict-action", conflictAction)
	}
	logrus.Tracef("semaphore-create: name=%s,value=%d,app-id=%d,conflict-action=%s\n",
		name, value, appID, conflictAction)
	return nil
}

// Sema holds a semaphore name-value pair for batch creation.
type Sema struct {
	Name  string
	Value int
}

// CreateSemaphores creates multiple semaphores from parsed lines in batch.
func CreateSemaphores(lines []string, appID int, batchSize int) error {
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
	return createSemaphores(semas, appID, batchSize)
}

// CreateFileSemaphores creates semaphores from a file, one JSON key-value pair per line.
func CreateFileSemaphores(fileName string, appID int, batchSize int) error {
	lines, err := common.GetTextFileLines(fileName)
	if err != nil {
		return errors.WrapE(err, "get-file-lines", "filename", fileName)
	}
	if len(lines) == 0 {
		return errors.E("null sema-file", "file-name", fileName)
	}
	return CreateSemaphores(lines, appID, batchSize)
}

// CreateJSONSemaphores creates semaphores from a JSON string.
func CreateJSONSemaphores(jsonText string, appID int, batchSize int) error {
	type semaItem struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	var items []semaItem

	if err := json.Unmarshal([]byte(jsonText), &items); err == nil {
		// Successfully parsed as array
	} else {
		var kvPairs map[string]int
		if err := json.Unmarshal([]byte(jsonText), &kvPairs); err == nil {
			items = make([]semaItem, 0, len(kvPairs))
			for name, value := range kvPairs {
				items = append(items, semaItem{Name: name, Value: value})
			}
		} else {
			var jsonData struct {
				Semaphores map[string]int `json:"semaphores"`
			}
			if err := json.Unmarshal([]byte(jsonText), &jsonData); err != nil {
				logrus.Errorf("Invalid JSON format, err-info:%v, json-text:%s\n", err, jsonText)
				return err
			}
			items = make([]semaItem, 0, len(jsonData.Semaphores))
			for name, value := range jsonData.Semaphores {
				items = append(items, semaItem{Name: name, Value: value})
			}
		}
	}

	ordered := make([]*Sema, 0, len(items))
	for _, item := range items {
		ordered = append(ordered, &Sema{Name: item.Name, Value: item.Value})
	}

	logrus.Tracef("Unmarshalled %d semaphores from JSON text", len(ordered))
	if err := createSemaphores(ordered, appID, batchSize); err != nil {
		return errors.WrapE(err, "createSemaphores", "app-id", appID, "semas", ordered)
	}
	return nil
}

func createSemaphores(ordered []*Sema, appID int, batchSize int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID)
	}

	conflictAction := os.Getenv("CONFLICT_ACTION")
	total := len(ordered)

	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		for _, v := range ordered[i:end] {
			reqCtx, cancel := ctx()
			_, err := c.Coordination.CreateSemaphore(reqCtx, &pb.CreateSemaphoreRequest{
				AppId:          int32(appID),
				Name:           v.Name,
				Value:          int32(v.Value),
				ConflictAction: conflictAction,
			})
			cancel()
			if err != nil {
				if conflictAction != "IGNORE" {
					return errors.WrapE(err, "create prepared-semaphore",
						"app-id", appID, "sema-name", v.Name, "sema-value", v.Value)
				}
			}
		}
		logrus.Infof("[%d..%d], %d row(s) inserted.\n", i, end-1, end-i)
	}

	return nil
}
