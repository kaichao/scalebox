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

	appID := 17
	semaphore.Create("node_progress:my-mod:n-00", 2, appID)
	semaphore.Create("node_progress:my-mod:n-01", 2, appID)
	semaphore.Create("node_progress:my-mod:n-10", 3, appID)
	semaphore.Create("node_progress:my-mod:n-11", 3, appID)

	v, err := semagroup.DiffMax("(node_progress:my-mod:n-0)1", appID)
	if err != nil {
		logrus.Errorf("err:%v\n", err)
	}
	fmt.Printf("diff=%d\n", v)
}

func TestIncrement(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 17
	semaphore.Create("node_progress:my-mod:n-00", 1, appID)
	semaphore.Create("node_progress:my-mod:n-01", 2, appID)
	semaphore.Create("node_progress:my-mod:n-10", 3, appID)
	semaphore.Create("node_progress:my-mod:n-11", 4, appID)

	for i := 0; i < 10; i++ {
		v, err := semagroup.Increment("node_progress:my-mod:n-0", appID)
		if err != nil {
			logrus.Errorf("err:%v\n", err)
		}
		fmt.Printf("sema=%s\n", v)
	}
}
