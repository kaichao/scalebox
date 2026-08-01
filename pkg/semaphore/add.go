package semaphore

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	"github.com/kaichao/scalebox/pkg/common"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// AddValue atomically adds delta to a semaphore and returns the new value.
func AddValue(name string, delta int, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return -1, errors.WrapE(err, "get client", "app-id", appID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.UpdateSemaphore(reqCtx, &pb.UpdateSemaphoreRequest{
		AppId: int32(appID),
		Name:  name,
		Delta: int32(delta),
	})
	if err != nil {
		return -1, errors.WrapE(err, "update semaphore",
			"app-id", appID, "sema-name", name, "delta", delta)
	}
	logrus.Tracef("In semaphore.AddValue(),name=%s,app-id:%d,delta:%d,ret-value:%d\n",
		name, appID, delta, resp.Value)
	return int(resp.Value), nil
}

// AddRegexValue adds delta to all semaphores matching the regex name pattern.
// Returns a JSON object string mapping matched names to their new values.
func AddRegexValue(name string, delta int, appID int) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "{}", errors.WrapE(err, "get client", "app-id", appID, "sema-name", name)
	}

	if !common.IsRegexString(name[0:1]) {
		name = "^" + name
	}
	re, reErr := regexp.Compile(name)
	if reErr != nil {
		return "{}", errors.WrapE(reErr, "compile regex", "pattern", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.ListSemaphores(reqCtx, &pb.ListSemaphoresRequest{
		AppId:    int32(appID),
		LeafOnly: true,
	})
	if err != nil {
		return "{}", errors.WrapE(err, "list semaphores for add-regex",
			"app-id", appID, "sema-name", name, "delta", delta)
	}

	result := make(map[string]int32)
	for _, n := range resp.Nodes {
		if re.MatchString(n.Name) {
			updCtx, updCancel := context.WithTimeout(context.Background(), 10*time.Second)
			updResp, updErr := c.Coordination.UpdateSemaphore(updCtx, &pb.UpdateSemaphoreRequest{
				AppId: int32(appID),
				Name:  n.Name,
				Delta: int32(delta),
			})
			updCancel()
			if updErr != nil {
				logrus.Warnf("AddRegexValue: update %s failed: %v", n.Name, updErr)
				continue
			}
			result[n.Name] = updResp.Value
		}
	}

	packed, _ := json.Marshal(result)
	v := string(packed)
	logrus.Tracef("In semaphore.AddRegexValue(),name=%s,app-id:%d,delta:%d,ret-value:%s\n",
		name, appID, delta, v)
	return v, nil
}

// AddMapValues adds deltas to multiple semaphores and returns the new values.
func AddMapValues(pairs map[string]int, appID int) (map[string]int, error) {
	if len(pairs) == 0 {
		return map[string]int{}, nil
	}

	result := make(map[string]int, len(pairs))

	for name, delta := range pairs {
		v, err := AddValue(name, delta, appID)
		if err != nil {
			return result, errors.WrapE(err, "semaphore-op failed",
				"app-id", appID, "sema-name", name, "delta", delta)
		}
		result[name] = v
	}

	return result, nil
}
