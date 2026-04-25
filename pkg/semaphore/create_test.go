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

	jsonText := `{"sema-3":0,"sema-4":3}`

	err := semaphore.CreateJSONSemaphores(jsonText, testAppID, 10)
	if err != nil {
		t.Logf("CreateJSONSemaphores error: %v", err)
	}
}

func TestCreate(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 测试创建单个信号量
	err := semaphore.Create("test_create_single", 100, testAppID)
	if err != nil {
		fmt.Printf("Create single semaphore error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("Create single semaphore succeeded")
	}

	// 测试更新已存在的信号量
	err4 := semaphore.Create("test_create_single", 150, testAppID)
	if err4 != nil {
		fmt.Printf("Update existing semaphore error: %v\n", err4)
	} else {
		fmt.Println("Update existing semaphore succeeded")
	}

	// 测试 CONFLICT_ACTION=IGNORE
	os.Setenv("CONFLICT_ACTION", "IGNORE")
	err5 := semaphore.Create("test_create_single", 150, testAppID)
	if err5 != nil {
		fmt.Printf("Update existing semaphore with IGNORE error: %v\n", err5)
	} else {
		fmt.Println("Update existing semaphore with IGNORE succeeded")
	}

	// 测试 CONFLICT_ACTION=OVERWRITE
	os.Setenv("CONFLICT_ACTION", "OVERWRITE")
	err6 := semaphore.Create("test_create_single", 150, testAppID)
	if err6 != nil {
		fmt.Printf("Update existing semaphore with OVERWRITE error: %v\n", err6)
	} else {
		fmt.Println("Update existing semaphore with OVERWRITE succeeded")
	}
}

func TestCreateSemaphores(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	lines := []string{
		`"sema-1":10`,
		`"sema-2":20`,
		`"sema-3":30`,
	}

	// 测试批量创建信号量
	err := semaphore.CreateSemaphores(lines, testAppID, 10)
	if err != nil {
		fmt.Printf("CreateSemaphores error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Println("CreateSemaphores succeeded")
	}

	// 测试空列表
	emptyLines := []string{}
	err3 := semaphore.CreateSemaphores(emptyLines, testAppID, 10)
	if err3 != nil {
		fmt.Printf("CreateSemaphores empty lines error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphores empty lines succeeded")
	}

	// 测试无效格式的行
	invalidLines := []string{
		`"sema-1":10`,
		`invalid_format`,
		`"sema-3":30`,
	}
	err4 := semaphore.CreateSemaphores(invalidLines, testAppID, 10)
	if err4 != nil {
		fmt.Printf("CreateSemaphores with invalid lines error (expected): %v\n", err4)
	} else {
		fmt.Println("CreateSemaphores with invalid lines succeeded (partial success)")
	}
}

func TestCreateWithConflictAction(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	semaphoreName := "test_conflict_sema"

	// 测试1: CONFLICT_ACTION=OVERWRITE
	os.Setenv("CONFLICT_ACTION", "OVERWRITE")

	err1 := semaphore.Create(semaphoreName, 100, testAppID)
	if err1 != nil {
		fmt.Printf("Create with OVERWRITE (first time) error: %v\n", err1)
	} else {
		fmt.Println("Create with OVERWRITE (first time) succeeded")
	}

	err2 := semaphore.Create(semaphoreName, 200, testAppID)
	if err2 != nil {
		fmt.Printf("Create with OVERWRITE (second time) error: %v\n", err2)
	} else {
		fmt.Println("Create with OVERWRITE (second time) succeeded")
	}

	// 测试2: CONFLICT_ACTION=IGNORE
	os.Setenv("CONFLICT_ACTION", "IGNORE")

	err3 := semaphore.Create("test_conflict_ignore", 300, testAppID)
	if err3 != nil {
		fmt.Printf("Create with IGNORE (first time) error: %v\n", err3)
	} else {
		fmt.Println("Create with IGNORE (first time) succeeded")
	}

	err4 := semaphore.Create("test_conflict_ignore", 400, testAppID)
	if err4 != nil {
		fmt.Printf("Create with IGNORE (second time) error: %v\n", err4)
	} else {
		fmt.Println("Create with IGNORE (second time) succeeded")
	}

	// 测试3: CONFLICT_ACTION未设置（默认行为，应该报错）
	os.Unsetenv("CONFLICT_ACTION")

	err5 := semaphore.Create("test_conflict_default", 500, testAppID)
	if err5 != nil {
		fmt.Printf("Create without CONFLICT_ACTION (first time) error: %v\n", err5)
	} else {
		fmt.Println("Create without CONFLICT_ACTION (first time) succeeded")
	}

	err6 := semaphore.Create("test_conflict_default", 600, testAppID)
	if err6 != nil {
		fmt.Printf("Create without CONFLICT_ACTION (second time, should fail) error: %v\n", err6)
	} else {
		fmt.Println("Create without CONFLICT_ACTION (second time) succeeded (unexpected)")
	}
}

func TestCreateSemaphoresWithConflictAction(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	lines := []string{
		`"conflict_sema_1":10`,
		`"conflict_sema_2":20`,
		`"conflict_sema_3":30`,
	}

	// 测试1: CONFLICT_ACTION=OVERWRITE
	os.Setenv("CONFLICT_ACTION", "OVERWRITE")

	err1 := semaphore.CreateSemaphores(lines, testAppID, 10)
	if err1 != nil {
		fmt.Printf("CreateSemaphores with OVERWRITE (first time) error: %v\n", err1)
	} else {
		fmt.Println("CreateSemaphores with OVERWRITE (first time) succeeded")
	}

	lines2 := []string{
		`"conflict_sema_1":100`,
		`"conflict_sema_2":200`,
		`"conflict_sema_3":300`,
	}
	err2 := semaphore.CreateSemaphores(lines2, testAppID, 10)
	if err2 != nil {
		fmt.Printf("CreateSemaphores with OVERWRITE (second time) error: %v\n", err2)
	} else {
		fmt.Println("CreateSemaphores with OVERWRITE (second time) succeeded")
	}

	// 测试2: CONFLICT_ACTION=IGNORE
	os.Setenv("CONFLICT_ACTION", "IGNORE")

	lines3 := []string{
		`"conflict_ignore_1":10`,
		`"conflict_ignore_2":20`,
	}
	err3 := semaphore.CreateSemaphores(lines3, testAppID, 10)
	if err3 != nil {
		fmt.Printf("CreateSemaphores with IGNORE (first time) error: %v\n", err3)
	} else {
		fmt.Println("CreateSemaphores with IGNORE (first time) succeeded")
	}

	err4 := semaphore.CreateSemaphores(lines3, testAppID, 10)
	if err4 != nil {
		fmt.Printf("CreateSemaphores with IGNORE (second time) error: %v\n", err4)
	} else {
		fmt.Println("CreateSemaphores with IGNORE (second time) succeeded")
	}

	// 测试3: CONFLICT_ACTION未设置（默认行为）
	os.Unsetenv("CONFLICT_ACTION")

	lines4 := []string{
		`"conflict_default_1":10`,
		`"conflict_default_2":20`,
	}
	err5 := semaphore.CreateSemaphores(lines4, testAppID, 10)
	if err5 != nil {
		fmt.Printf("CreateSemaphores without CONFLICT_ACTION (first time) error: %v\n", err5)
	} else {
		fmt.Println("CreateSemaphores without CONFLICT_ACTION (first time) succeeded")
	}

	err6 := semaphore.CreateSemaphores(lines4, testAppID, 10)
	if err6 != nil {
		fmt.Printf("CreateSemaphores without CONFLICT_ACTION (second time, may fail) error: %v\n", err6)
	} else {
		fmt.Println("CreateSemaphores without CONFLICT_ACTION (second time) succeeded (unexpected)")
	}
}
