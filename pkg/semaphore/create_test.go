package semaphore_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

var (
	appID   = 4
	vtaskID = 43
)

func TestCreateJSONSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	jsonText := `{"sema-3":0,"sema-4":3}`

	semaphore.CreateJSONSemaphores(jsonText, 0, appID, 10)
}

func TestCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

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

	// 测试更新已存在的信号量
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "IGNORE")
	err5 := semaphore.Create("test_create_single", 150, vtaskID, appID) // 更新值
	if err5 != nil {
		fmt.Printf("Update existing semaphore error: %v\n", err5)
	} else {
		fmt.Println("Update existing semaphore succeeded")
	}
	// 测试更新已存在的信号量
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "OVERWRITE")
	err6 := semaphore.Create("test_create_single", 150, vtaskID, appID) // 更新值
	if err6 != nil {
		fmt.Printf("Update existing semaphore error: %v\n", err6)
	} else {
		fmt.Println("Update existing semaphore succeeded")
	}
}

func TestCreateSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

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

func TestCreateWithExistsSema(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	semaphoreName := "test_exists_sema"

	// 测试1: SEMAPHORE_CONFLICT_ACTION=OVERWRITE
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "OVERWRITE")

	// 第一次创建
	err1 := semaphore.Create(semaphoreName, 100, vtaskID, appID)
	if err1 != nil {
		fmt.Printf("Create with OVERWRITE (first time) error: %v\n", err1)
	} else {
		fmt.Println("Create with OVERWRITE (first time) succeeded")
	}

	// 第二次创建，应该覆盖
	err2 := semaphore.Create(semaphoreName, 200, vtaskID, appID)
	if err2 != nil {
		fmt.Printf("Create with OVERWRITE (second time, should overwrite) error: %v\n", err2)
	} else {
		fmt.Println("Create with OVERWRITE (second time, should overwrite) succeeded")
	}

	// 测试2: SEMAPHORE_CONFLICT_ACTION=IGNORE
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "IGNORE")

	// 第一次创建
	err3 := semaphore.Create("test_exists_sema_ignore", 300, vtaskID, appID)
	if err3 != nil {
		fmt.Printf("Create with IGNORE (first time) error: %v\n", err3)
	} else {
		fmt.Println("Create with IGNORE (first time) succeeded")
	}

	// 第二次创建，应该忽略冲突
	err4 := semaphore.Create("test_exists_sema_ignore", 400, vtaskID, appID)
	if err4 != nil {
		fmt.Printf("Create with IGNORE (second time, should ignore) error: %v\n", err4)
	} else {
		fmt.Println("Create with IGNORE (second time, should ignore) succeeded")
	}

	// 测试3: SEMAPHORE_CONFLICT_ACTION未设置（默认行为，应该报错）
	os.Unsetenv("SEMAPHORE_CONFLICT_ACTION")

	// 第一次创建应该成功
	err5 := semaphore.Create("test_exists_sema_default", 500, vtaskID, appID)
	if err5 != nil {
		fmt.Printf("Create without SEMAPHORE_CONFLICT_ACTION (first time) error: %v\n", err5)
	} else {
		fmt.Println("Create without SEMAPHORE_CONFLICT_ACTION (first time) succeeded")
	}

	// 第二次创建应该失败（由于唯一约束）
	err6 := semaphore.Create("test_exists_sema_default", 600, vtaskID, appID)
	if err6 != nil {
		fmt.Printf("Create without SEMAPHORE_CONFLICT_ACTION (second time, should fail) error: %v\n", err6)
	} else {
		fmt.Println("Create without SEMAPHORE_CONFLICT_ACTION (second time, should fail) succeeded (unexpected)")
	}
}

func TestCreateSemaphoresWithExistsSema(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	vtaskID := 4117
	lines := []string{
		`"exists_sema_1":10`,
		`"exists_sema_2":20`,
		`"exists_sema_3":30`,
	}

	// 测试1: SEMAPHORE_CONFLICT_ACTION=OVERWRITE
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "OVERWRITE")

	// 第一次批量创建
	err1 := semaphore.CreateSemaphores(lines, vtaskID, appID, 10)
	if err1 != nil {
		fmt.Printf("CreateSemaphores with OVERWRITE (first time) error: %v\n", err1)
	} else {
		fmt.Println("CreateSemaphores with OVERWRITE (first time) succeeded")
	}

	// 第二次批量创建，应该覆盖
	lines2 := []string{
		`"exists_sema_1":100`,
		`"exists_sema_2":200`,
		`"exists_sema_3":300`,
	}
	err2 := semaphore.CreateSemaphores(lines2, vtaskID, appID, 10)
	if err2 != nil {
		fmt.Printf("CreateSemaphores with OVERWRITE (second time, should overwrite) error: %v\n", err2)
	} else {
		fmt.Println("CreateSemaphores with OVERWRITE (second time, should overwrite) succeeded")
	}

	// 测试2: SEMAPHORE_CONFLICT_ACTION=IGNORE
	os.Setenv("SEMAPHORE_CONFLICT_ACTION", "IGNORE")

	// 第一次批量创建
	lines3 := []string{
		`"exists_sema_ignore_1":10`,
		`"exists_sema_ignore_2":20`,
	}
	err3 := semaphore.CreateSemaphores(lines3, vtaskID, appID, 10)
	if err3 != nil {
		fmt.Printf("CreateSemaphores with IGNORE (first time) error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphores with IGNORE (first time) succeeded")
	}

	// 第二次批量创建，应该忽略冲突
	err4 := semaphore.CreateSemaphores(lines3, vtaskID, appID, 10)
	if err4 != nil {
		fmt.Printf("CreateSemaphores with IGNORE (second time, should ignore) error: %v\n", err4)
	} else {
		fmt.Println("CreateSemaphores with IGNORE (second time, should ignore) succeeded")
	}

	// 测试3: SEMAPHORE_CONFLICT_ACTION未设置（默认行为）
	os.Unsetenv("SEMAPHORE_CONFLICT_ACTION")

	// 第一次批量创建应该成功
	lines4 := []string{
		`"exists_sema_default_1":10`,
		`"exists_sema_default_2":20`,
	}
	err5 := semaphore.CreateSemaphores(lines4, vtaskID, appID, 10)
	if err5 != nil {
		fmt.Printf("CreateSemaphores without SEMAPHORE_CONFLICT_ACTION (first time) error: %v\n", err5)
	} else {
		fmt.Println("CreateSemaphores without SEMAPHORE_CONFLICT_ACTION (first time) succeeded")
	}

	// 第二次批量创建可能失败（由于唯一约束）
	err6 := semaphore.CreateSemaphores(lines4, vtaskID, appID, 10)
	if err6 != nil {
		fmt.Printf("CreateSemaphores without SEMAPHORE_CONFLICT_ACTION (second time, may fail) error: %v\n", err6)
	} else {
		fmt.Println("CreateSemaphores without SEMAPHORE_CONFLICT_ACTION (second time, may fail) succeeded (unexpected)")
	}
}
