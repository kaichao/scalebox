package semaphore_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestAddListValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 24
	semaphore.Create("a", 5, appID)
	semaphore.Create("b", 10, appID)
	semaphore.Create("c", 15, appID)
	// 测试用例：正常输入
	names := []string{"a", "b", "c"}
	delta := 5

	// 调用函数
	result, err := semaphore.AddListValue(names, appID, delta)

	// 期望结果
	want := `{"a":10,"b":15,"c":20}`

	// 验证
	if err != nil {
		t.Errorf("AddListValue() error = %v, want nil", err)
	}
	if result != want {
		t.Errorf("AddListValue() = %q, want %q", result, want)
	}
}

func TestAddValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	appID := 25
	v, err := semaphore.AddValue("node_progress:node-progress-gap:n-01.inline", appID, 3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("v=%s.", v)
}
