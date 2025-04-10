package semaphore

import (
	"fmt"
	"os"

	"github.com/kaichao/gopkg/common"
	"github.com/kaichao/gopkg/exec"
)

// Create ...
func Create(semaLines string) error {
	common.AppendToFile("my-sema.txt", semaLines)
	defer os.Remove("my-sema.txt")

	cmd := "scalebox semaphore create --sema-file my-sema.txt"
	code, err := exec.RunReturnExitCode(cmd, 600)
	if code > 0 {
		return fmt.Errorf("[ERROR]semaphore-create,code=%d", code)
	}
	return err
}
