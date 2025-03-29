package semaphore

import (
	"fmt"
	"os"

	"github.com/kaichao/scalebox/pkg/misc"
)

// Create ...
func Create(semaLines string) error {
	misc.AppendToFile("my-sema.txt", semaLines)
	defer os.Remove("my-sema.txt")

	cmd := "scalebox semaphore create --sema-file my-sema.txt"
	code, err := misc.ExecCommandReturnExitCode(cmd, 600)
	if code > 0 {
		return fmt.Errorf("[ERROR]semaphore-create,code=%d", code)
	}
	return err
}
