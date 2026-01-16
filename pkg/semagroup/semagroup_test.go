package semagroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kaichao/scalebox/pkg/semagroup"
	"github.com/kaichao/scalebox/pkg/semaphore"
	"github.com/sirupsen/logrus"
)

func TestDiffMax(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 1
	vtaskID := 0
	semaphore.Create("node_progress:my-mod:n-00", 2, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-01", 2, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-10", 3, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-11", 3, vtaskID, appID)

	v, err := semagroup.DiffMax("(node_progress:node-progress-gap:n-0)0", appID)
	if err != nil {
		logrus.Errorf("err:%v\n", err)
	}
	fmt.Printf("diff=%d\n", v)
}

func TestIncrement(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 1
	vtaskID := 0
	semaphore.Create("node_progress:my-mod:n-00", 1, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-01", 2, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-10", 3, vtaskID, appID)
	semaphore.Create("node_progress:my-mod:n-11", 4, vtaskID, appID)

	for i := 0; i < 10; i++ {
		v, err := semagroup.Increment("node_progress:my-mod:n-0", appID)
		if err != nil {
			logrus.Errorf("err:%v\n", err)
		}
		fmt.Printf("sema=%s\n", v)
	}
}
