package semaphore

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kaichao/scalebox/pkg/misc"
	"github.com/sirupsen/logrus"
)

// Decrement ...
func Decrement(sema string) (int, error) {
	cmd := "scalebox semaphore decrement " + sema
	code, stdout, stderr, err := misc.ExecCommandReturnAll(cmd, 20)
	logrus.Errorf("stcerr:\n%s\n", stderr)
	fmt.Printf("stdout:\n%s\n", stdout)
	if err != nil {
		return math.MinInt, err
	}
	if code > 0 {
		return code, fmt.Errorf("[ERROR]semaphore decrement")
	}
	v, err := strconv.Atoi(strings.TrimSpace(stdout))
	if err != nil {
		logrus.Errorf("semaphore-value not a integer, value=%s\n", stdout)
		return -1, err
	}
	return v, nil
}

// DecrementExpr ...
func DecrementExpr(semExpr string) (map[string]int, error) {
	return map[string]int{}, nil
}
