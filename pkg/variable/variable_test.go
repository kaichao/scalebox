package variable_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/variable"
)

const testAppID = 46

func TestSet(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	err := variable.Set("var-1", "val-1", testAppID)
	if err != nil {
		t.Logf("Set var-1 error: %v", err)
	}
	err = variable.Set("var-2", "val-2", testAppID)
	if err != nil {
		t.Logf("Set var-2 error: %v", err)
	}

	val, _ := variable.GetJSON("var.+", testAppID)
	fmt.Println("val:", val)
}

func TestGetValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 设置测试数据
	variable.Set("var-1", "val-1", testAppID)
	variable.Set("var-2", "val-2", testAppID)

	// 测试用例1：获取存在的变量
	t.Run("existing variable", func(t *testing.T) {
		val, err := variable.GetValue("var-1", testAppID)
		if err != nil {
			t.Errorf("GetValue failed: %v", err)
		}
		if val != "val-1" {
			t.Errorf("Expected 'val-1', got '%s'", val)
		}
		fmt.Printf("GetValue var-1: %s\n", val)
	})

	// 测试用例2：获取不存在的变量
	t.Run("non-existent variable", func(t *testing.T) {
		_, err := variable.GetValue("non-existent-var", testAppID)
		if err == nil {
			t.Error("Expected error for non-existent variable, but got none")
		} else {
			fmt.Printf("Expected error for non-existent variable: %v\n", err)
		}
	})
}

func TestGetJSON(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	// 设置测试数据
	variable.Set("var-1", "val-1", testAppID)
	variable.Set("var-2", "val-2", testAppID)

	// 测试用例1：使用正则表达式获取变量（多个匹配）
	t.Run("regex pattern", func(t *testing.T) {
		val, err := variable.GetJSON("var.+", testAppID)
		if err != nil {
			t.Errorf("GetJSON with regex pattern failed: %v", err)
		}
		fmt.Printf("GetJSON with regex pattern: %s\n", val)
		if len(val) == 0 || val == "{}" {
			t.Errorf("Expected non-empty JSON, got '%s'", val)
		}
	})

	// 测试用例2：获取不存在的正则表达式匹配
	t.Run("non-matching regex", func(t *testing.T) {
		val, err := variable.GetJSON("non-matching-.+", testAppID)
		if err != nil {
			t.Errorf("GetJSON with non-matching regex failed: %v", err)
		}
		if val != "{}" {
			t.Errorf("Expected '{}' for non-matching regex, got '%s'", val)
		}
		fmt.Printf("GetJSON with non-matching regex: %s\n", val)
	})

	// 测试用例3：精确匹配
	t.Run("exact match", func(t *testing.T) {
		val, err := variable.GetJSON("var-1", testAppID)
		if err != nil {
			t.Errorf("GetJSON with exact match failed: %v", err)
		}
		fmt.Printf("GetJSON with exact match: %s\n", val)
	})
}

// ExampleSet 展示了如何使用 Set 函数设置变量
func ExampleSet() {
	fmt.Println("err := Set(name, value, appID)")
	// Output:
	// err := Set(name, value, appID)
}

// ExampleGetValue 展示了如何使用 GetValue 函数获取单个变量
func ExampleGetValue() {
	fmt.Println("value, err := GetValue(name, appID)")
	// Output:
	// value, err := GetValue(name, appID)
}

// ExampleGetJSON 展示了如何使用 GetJSON 函数进行正则表达式匹配
func ExampleGetJSON() {
	fmt.Println("// 使用正则表达式获取多个变量:")
	fmt.Println("json, err := GetJSON(\"^config_.+\", appID)")
	fmt.Println("// 返回JSON: {\"key1\":\"value1\",\"key2\":\"value2\"}")
	// Output:
	// // 使用正则表达式获取多个变量:
	// json, err := GetJSON("^config_.+", appID)
	// // 返回JSON: {"key1":"value1","key2":"value2"}
}
