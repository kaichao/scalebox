package variable

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

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// GetValue returns the value of an app-level variable.
func GetValue(name string, appID int) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "", errors.WrapE(err, "get client", "app-id", appID, "var-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.GetVariable(reqCtx, &pb.GetVariableRequest{
		AppId: int32(appID),
		Name:  name,
	})
	if err != nil {
		return "", errors.WrapE(err, "get variable",
			"app-id", appID, "var-name", name)
	}
	return resp.Value, nil
}

// GetJSON returns variable values matching a regex as a JSON object string.
func GetJSON(name string, appID int) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "{}", errors.WrapE(err, "get client", "app-id", appID, "var-name", name)
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

	resp, err := c.Coordination.ListVariables(reqCtx, &pb.ListVariablesRequest{
		AppId:    int32(appID),
		LeafOnly: true,
	})
	if err != nil {
		return "{}", errors.WrapE(err, "list variables for get-json",
			"app-id", appID, "var-name", name)
	}

	result := make(map[string]string)
	for _, n := range resp.Nodes {
		if re.MatchString(n.Name) {
			result[n.Name] = n.Value
		}
	}

	packed, _ := json.Marshal(result)
	return regexp.MustCompile(`\s+`).ReplaceAllString(string(packed), ""), nil
}

// Set creates or updates an app-level variable.
func Set(name string, value string, appID int) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "app-id", appID, "var-name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.Coordination.SetVariable(reqCtx, &pb.SetVariableRequest{
		AppId: int32(appID),
		Name:  name,
		Value: value,
	})
	if err != nil {
		return errors.WrapE(err, "set-variable",
			"app-id", appID, "var-name", name, "var-value", value)
	}
	logrus.Tracef("In variable.Set(),name=%s,value=%s,app-id:%d\n",
		name, value, appID)
	return nil
}
