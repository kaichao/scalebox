package semaphore_test

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestCreateJSONSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 3
	jsonText := `{"sema-3":0,"sema-4":3}`

	semaphore.CreateJSONSemaphores(jsonText, appID, 10)
}
