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

	// 测试1: SEMAPHORE_AUTO_CREATE=yes，不存在的信号量自动创建
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")
	value1, err1 := semaphore.GetValue("test_getvalue_sema", testAppID)
	if err1 != nil {
		fmt.Printf("GetValue with auto-create error: %v\n", err1)
	} else {
		fmt.Printf("GetValue with auto-create succeeded, value=%d\n", value1)
	}

	// 测试2: SEMAPHORE_AUTO_CREATE未设置，不存在的信号量应报错
	os.Unsetenv("SEMAPHORE_AUTO_CREATE")
	value2, err2 := semaphore.GetValue("non_existent_sema", testAppID)
	if err2 != nil {
		fmt.Printf("GetValue without auto-create error (expected): %v\n", err2)
	} else {
		fmt.Printf("GetValue without auto-create succeeded (unexpected), value=%d\n", value2)
	}
}

func TestGetJSON(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 创建一些测试信号量
	semaphoreNames := []string{"test_json_a", "test_json_b", "test_json_c"}
	for i, name := range semaphoreNames {
		_ = semaphore.Create(name, (i+1)*10, testAppID)
	}

	// 测试1: 按前缀获取信号量
	jsonResult1, err1 := semaphore.GetJSON("test_json", testAppID)
	if err1 != nil {
		fmt.Printf("GetJSON with prefix error: %v\n", err1)
	} else {
		fmt.Printf("GetJSON with prefix succeeded, result=%s\n", jsonResult1)
	}

	// 测试2: 使用正则表达式获取
	jsonResult2, err2 := semaphore.GetJSON("test_json_.+", testAppID)
	if err2 != nil {
		fmt.Printf("GetJSON with regex error: %v\n", err2)
	} else {
		fmt.Printf("GetJSON with regex succeeded, result=%s\n", jsonResult2)
	}

	// 测试3: 获取不存在的信号量
	jsonResult3, err3 := semaphore.GetJSON("non_existent_prefix", testAppID)
	if err3 != nil {
		fmt.Printf("GetJSON non-existent error: %v\n", err3)
	} else {
		fmt.Printf("GetJSON non-existent succeeded, result=%s\n", jsonResult3)
	}
}

func TestGetValueNotFoundLogic(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	testCases := []struct {
		name        string
		autoCreate  string
		expectError bool
		expectValue int
	}{
		{"auto_create_yes", "yes", false, 0}, // 自动创建，返回0
		{"auto_create_no", "no", true, 0},    // 不自动创建，应该报错
		{"auto_create_unset", "", true, 0},   // 未设置，应该报错
	}

	for _, tc := range testCases {
		os.Setenv("SEMAPHORE_AUTO_CREATE", tc.autoCreate)

		value, err := semaphore.GetValue(tc.name, testAppID)

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
