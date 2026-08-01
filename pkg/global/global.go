package global

import (
	"context"
	"time"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/client"
	pb "github.com/kaichao/scalebox/pkg/pb"
	"github.com/sirupsen/logrus"
)

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// Set creates or updates a global variable.
func Set(name string, value string) error {
	c, err := client.Default()
	if err != nil {
		return errors.WrapE(err, "get client", "name", name, "value", value)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	_, err = c.Coordination.SetGlobal(reqCtx, &pb.SetGlobalRequest{
		Name:  name,
		Value: value,
	})
	if err != nil {
		return errors.WrapE(err, "global-set", "name", name, "value", value)
	}
	logrus.Tracef("In global.Set(),global-name:%s,global-value:%s\n", name, value)
	return nil
}

// Get returns the value of a global variable.
func Get(name string) (string, error) {
	c, err := client.Default()
	if err != nil {
		return "", errors.WrapE(err, "get client", "name", name)
	}

	reqCtx, cancel := ctx()
	defer cancel()

	resp, err := c.Coordination.GetGlobal(reqCtx, &pb.GetGlobalRequest{Name: name})
	if err != nil {
		return "", errors.WrapE(err, "global get", "name", name)
	}
	logrus.Tracef("In global.Get(),global-name:%s,global-value:%s\n", name, resp.Value)
	return resp.Value, nil
}
