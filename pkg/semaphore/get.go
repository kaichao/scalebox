package semaphore

import (
	"encoding/json"
	"regexp"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	"github.com/kaichao/scalebox/pkg/common"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

// GetValue returns the current value of a semaphore.
func GetValue(name string, appID int) (int, error) {
	c, err := client.Default()
	if err != nil {
		return -1, errors.WrapE(err, "get client", "app-id", appID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.GetSemaphore(reqCtx, &pb.GetSemaphoreRequest{
		AppId: int32(appID),
		Name:  name,
	})
	if err != nil {
		return -1, errors.WrapE(err, "get semaphore", "app-id", appID, "sema-name", name)
	}
	logrus.Tracef("In semaphore.GetValue(),name=%s,app-id:%d,value:%d\n",
		name, appID, resp.Value)
	return int(resp.Value), nil
}

// GetJSON returns semaphore values matching a regex as a JSON object string.
func GetJSON(name string, appID int) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "{}", errors.WrapE(err, "get client", "app-id", appID, "sema-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.ListSemaphores(reqCtx, &pb.ListSemaphoresRequest{
		AppId:    int32(appID),
		LeafOnly: true,
	})
	if err != nil {
		return "{}", errors.WrapE(err, "list semaphores for get-json",
			"app-id", appID, "sema-name", name)
	}

	if !common.IsRegexString(name[0:1]) {
		name = "^" + name
	}
	re, reErr := regexp.Compile(name)
	if reErr != nil {
		return "{}", errors.WrapE(reErr, "compile regex", "pattern", name)
	}

	result := make(map[string]int32)
	for _, n := range resp.Nodes {
		if re.MatchString(n.Name) {
			result[n.Name] = n.Value
		}
	}

	packed, _ := json.Marshal(result)
	v := string(packed)
	logrus.Tracef("In semaphore.GetJSON(),name=%s,app-id:%d,json-value:%s\n",
		name, appID, v)
	return v, nil
}
