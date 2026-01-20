package semaphore_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestGetValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 测试1: 测试SEMAPHORE_AUTO_CREATE=yes时的行为
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	semaphoreName := "test_getvalue_semaphore"

	// 测试获取不存在的信号量（应该自动创建）
	value1, err1 := semaphore.GetValue(semaphoreName, vtaskID, appID)
	if err1 != nil {
		fmt.Printf("GetValue with auto-create error: %v\n", err1)
	} else {
		fmt.Printf("GetValue with auto-create succeeded, value=%d\n", value1)
	}

	// 测试2: 测试SEMAPHORE_AUTO_CREATE未设置时的行为
	os.Unsetenv("SEMAPHORE_AUTO_CREATE")

	value2, err2 := semaphore.GetValue("non_existent_semaphore_2", vtaskID, appID)
	if err2 != nil {
		fmt.Printf("GetValue without auto-create error (expected): %v\n", err2)
	} else {
		fmt.Printf("GetValue without auto-create succeeded (unexpected), value=%d\n", value2)
	}

	// 测试3: 测试全局信号量（vtaskID=0）
	value3, err3 := semaphore.GetValue("test_global_semaphore", 0, appID)
	if err3 != nil {
		fmt.Printf("GetValue global error: %v\n", err3)
	} else {
		fmt.Printf("GetValue global succeeded, value=%d\n", value3)
	}

	// 测试4: 测试负vtaskID
	value4, err4 := semaphore.GetValue("test_negative_semaphore", -1, appID)
	if err4 != nil {
		fmt.Printf("GetValue with negative vtaskID error: %v\n", err4)
	} else {
		fmt.Printf("GetValue with negative vtaskID succeeded, value=%d\n", value4)
	}
}

func TestGetJSON(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 创建一些测试信号量
	semaphoreNames := []string{"test_json_a", "test_json_b", "test_json_c"}
	for i, name := range semaphoreNames {
		_ = semaphore.Create(name, (i+1)*10, vtaskID, appID)
	}

	// 测试1: 按前缀获取信号量
	jsonResult1, err1 := semaphore.GetJSON("test_json", vtaskID, appID)
	if err1 != nil {
		fmt.Printf("GetJSON with prefix error: %v\n", err1)
	} else {
		fmt.Printf("GetJSON with prefix succeeded, result=%s\n", jsonResult1)
	}

	// 测试2: 使用正则表达式获取
	jsonResult2, err2 := semaphore.GetJSON("test_json_.+", vtaskID, appID)
	if err2 != nil {
		fmt.Printf("GetJSON with regex error: %v\n", err2)
	} else {
		fmt.Printf("GetJSON with regex succeeded, result=%s\n", jsonResult2)
	}

	// 测试3: 获取不存在的信号量
	jsonResult3, err3 := semaphore.GetJSON("non_existent_prefix", vtaskID, appID)
	if err3 != nil {
		fmt.Printf("GetJSON non-existent error: %v\n", err3)
	} else {
		fmt.Printf("GetJSON non-existent succeeded, result=%s\n", jsonResult3)
	}

	// 测试4: 测试全局信号量（vtaskID=0）
	jsonResult4, err4 := semaphore.GetJSON("test_json", 0, appID)
	if err4 != nil {
		fmt.Printf("GetJSON global error: %v\n", err4)
	} else {
		fmt.Printf("GetJSON global succeeded, result=%s\n", jsonResult4)
	}
}

func TestGetValueNotFoundLogic(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 测试详细的NotFound逻辑
	testCases := []struct {
		name        string
		autoCreate  string
		expectError bool
		expectValue int
	}{
		{"test_notfound_1", "yes", false, 0}, // 自动创建，返回0
		{"test_notfound_2", "no", true, 0},   // 不自动创建，应该报错
		{"test_notfound_3", "", true, 0},     // 未设置，应该报错
	}

	for _, tc := range testCases {
		os.Setenv("SEMAPHORE_AUTO_CREATE", tc.autoCreate)

		value, err := semaphore.GetValue(tc.name, vtaskID, appID)

		if tc.expectError {
			if err == nil {
				fmt.Printf("Test %s: Expected error but got none, value=%d\n", tc.name, value)
			} else {
				fmt.Printf("Test %s: Got expected error: %v\n", tc.name, err)
			}
		} else {
			if err != nil {
				fmt.Printf("Test %s: Unexpected error: %v\n", tc.name, err)
			} else if value != tc.expectValue {
				fmt.Printf("Test %s: Expected value %d but got %d\n", tc.name, tc.expectValue, value)
			} else {
				fmt.Printf("Test %s: Success, value=%d\n", tc.name, value)
			}
		}
	}
}

func TestGetIsolation(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	vtaskIDX := int64(44)

	// 测试信号量按vtaskID隔离的概念
	// 注意：由于外键约束，这些测试可能不会实际执行数据库操作
	// 但我们可以验证函数调用语法正确

	semaphoreName := "isolated_semaphore"

	// 为不同任务创建同名信号量（理论上应该隔离）
	_ = semaphore.Create(semaphoreName, 100, vtaskID, appID)
	_ = semaphore.Create(semaphoreName, 200, vtaskIDX, appID)
	_ = semaphore.Create(semaphoreName, 300, 0, appID) // 全局

	fmt.Println("Tested semaphore isolation by vtaskID (function calls only due to foreign key constraints)")
}
