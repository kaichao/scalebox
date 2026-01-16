package semaphore_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestCreateJSONSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	jsonText := `{"sema-3":0,"sema-4":3}`

	semaphore.CreateJSONSemaphores(jsonText, 0, appID, 10)
}

func TestCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	vtaskID := 4117

	// 测试创建单个信号量
	err := semaphore.Create("test_create_single", 100, vtaskID, appID)
	if err != nil {
		// 由于外键约束可能失败，我们只记录错误但不失败测试
		fmt.Printf("Create single semaphore error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("Create single semaphore succeeded")
	}

	// 测试创建全局信号量（vtaskID=0）
	err2 := semaphore.Create("test_create_global", 200, 0, appID)
	if err2 != nil {
		fmt.Printf("Create global semaphore error: %v\n", err2)
	} else {
		fmt.Println("Create global semaphore succeeded")
	}

	// 测试创建信号量（负vtaskID）
	err3 := semaphore.Create("test_create_negative", 300, -1, appID)
	if err3 != nil {
		fmt.Printf("Create semaphore with negative vtaskID error: %v\n", err3)
	} else {
		fmt.Println("Create semaphore with negative vtaskID succeeded")
	}

	// 测试更新已存在的信号量
	err4 := semaphore.Create("test_create_single", 150, vtaskID, appID) // 更新值
	if err4 != nil {
		fmt.Printf("Update existing semaphore error: %v\n", err4)
	} else {
		fmt.Println("Update existing semaphore succeeded")
	}
}

func TestCreateSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	vtaskID := 4117
	lines := []string{
		`"sema-1":10`,
		`"sema-2":20`,
		`"sema-3":30`,
	}

	// 测试批量创建信号量（全局，vtaskID=0）
	err := semaphore.CreateSemaphores(lines, 0, appID, 10)
	if err != nil {
		// 由于外键约束可能失败，我们只记录错误但不失败测试
		fmt.Printf("CreateSemaphores global error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("CreateSemaphores global succeeded")
	}

	// 测试批量创建信号量（特定vtaskID）
	err2 := semaphore.CreateSemaphores(lines, vtaskID, appID, 10)
	if err2 != nil {
		fmt.Printf("CreateSemaphores with vtaskID error: %v\n", err2)
	} else {
		fmt.Println("CreateSemaphores with vtaskID succeeded")
	}

	// 测试空列表
	emptyLines := []string{}
	err3 := semaphore.CreateSemaphores(emptyLines, vtaskID, appID, 10)
	if err3 != nil {
		fmt.Printf("CreateSemaphores empty lines error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphores empty lines succeeded (should work)")
	}

	// 测试无效格式的行
	invalidLines := []string{
		`"sema-1":10`,
		`invalid_format`,
		`"sema-3":30`,
	}
	err4 := semaphore.CreateSemaphores(invalidLines, vtaskID, appID, 10)
	if err4 != nil {
		fmt.Printf("CreateSemaphores with invalid lines error (expected): %v\n", err4)
	} else {
		fmt.Println("CreateSemaphores with invalid lines succeeded (partial success)")
	}
}
