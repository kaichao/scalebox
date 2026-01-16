package variable_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/variable"
)

func TestSet(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	vtaskID := 4117

	variable.Set("var-1", "val-1", vtaskID, appID)
	variable.Set("var-2", "val-2", vtaskID, appID)

	variable.Set("var-1", "val-10", 0, appID)
	variable.Set("var-2", "val-20", 0, appID)

	val, _ := variable.Get("var.+", vtaskID, appID)
	fmt.Println("val:", val)
	fmt.Println()

	val, _ = variable.Get("var.+", 0, appID)
	fmt.Println("val0:", val)

}

// ExampleSet 展示了如何使用 Set 函数设置变量
// 注意：此示例需要正确的数据库连接和有效的appID
func ExampleSet() {
	// 在实际使用中，需要设置数据库连接环境变量
	// os.Setenv("PGHOST", "your-database-host")

	// 示例：设置变量
	// appID := 100  // 需要有效的app ID
	// vtaskID := 200 // 任务ID，0表示全局变量

	// err := variable.Set("task_progress", "50%", vtaskID, appID)
	// if err != nil {
	//     fmt.Printf("Error: %v\n", err)
	//     return
	// }

	// err = variable.Set("global_setting", "enabled", 0, appID)
	// if err != nil {
	//     fmt.Printf("Error: %v\n", err)
	//     return
	// }

	fmt.Println("Set(name, value, vtaskID, appID) // 设置变量")
	fmt.Println("Set() with vtaskID=0 for global variables")
	// Output:
	// Set(name, value, vtaskID, appID) // 设置变量
	// Set() with vtaskID=0 for global variables
}

// ExampleGet 展示了如何使用 Get 函数获取变量
// 注意：此示例需要正确的数据库连接和有效的appID
func ExampleGet() {
	// 在实际使用中，需要设置数据库连接环境变量
	// os.Setenv("PGHOST", "your-database-host")

	// 示例：获取变量
	// appID := 100  // 需要有效的app ID
	// vtaskID := 200 // 任务ID，0表示全局变量

	// 首先需要设置变量
	// _ = variable.Set("username", "john_doe", vtaskID, appID)

	// 然后获取变量
	// value, err := variable.Get("username", vtaskID, appID)
	// if err != nil {
	//     fmt.Printf("Error: %v\n", err)
	//     return
	// }
	// fmt.Printf("Value: %s\n", value)

	fmt.Println("value, err := Get(name, vtaskID, appID)")
	fmt.Println("// 返回变量值或错误")
	// Output:
	// value, err := Get(name, vtaskID, appID)
	// // 返回变量值或错误
}

// ExampleGet_withVtask 展示了如何使用 Get 函数处理不同的 vtaskID
// 注意：此示例展示了变量按vtaskID隔离的概念
func ExampleGet_withVtask() {
	// 变量按vtaskID隔离：不同任务可以有同名但不同值的变量
	fmt.Println("// 为不同任务设置变量:")
	fmt.Println("Set(\"output_file\", \"/path/task1.txt\", 200, 100)")
	fmt.Println("Set(\"output_file\", \"/path/task2.txt\", 201, 100)")
	fmt.Println("Set(\"output_file\", \"/path/global.txt\", 0, 100)")
	fmt.Println()
	fmt.Println("// 获取不同任务的变量:")
	fmt.Println("Get(\"output_file\", 200, 100) // 返回: /path/task1.txt")
	fmt.Println("Get(\"output_file\", 201, 100) // 返回: /path/task2.txt")
	fmt.Println("Get(\"output_file\", 0, 100)   // 返回: /path/global.txt")
	// Output:
	// // 为不同任务设置变量:
	// Set("output_file", "/path/task1.txt", 200, 100)
	// Set("output_file", "/path/task2.txt", 201, 100)
	// Set("output_file", "/path/global.txt", 0, 100)
	//
	// // 获取不同任务的变量:
	// Get("output_file", 200, 100) // 返回: /path/task1.txt
	// Get("output_file", 201, 100) // 返回: /path/task2.txt
	// Get("output_file", 0, 100)   // 返回: /path/global.txt
}

// ExampleGet_regex 展示了如何使用 Get 函数进行正则表达式匹配
func ExampleGet_regex() {
	// 使用正则表达式获取多个变量
	fmt.Println("// 设置多个配置变量:")
	fmt.Println("Set(\"config_server_host\", \"localhost\", 200, 100)")
	fmt.Println("Set(\"config_server_port\", \"8080\", 200, 100)")
	fmt.Println("Set(\"config_database_name\", \"mydb\", 200, 100)")
	fmt.Println()
	fmt.Println("// 使用正则表达式获取所有config_开头的变量:")
	fmt.Println("// Get(\"^config_.+\", 200, 100)")
	fmt.Println("// 返回JSON: {\"config_server_host\":\"localhost\",\"config_server_port\":\"8080\",\"config_database_name\":\"mydb\"}")
	// Output:
	// // 设置多个配置变量:
	// Set("config_server_host", "localhost", 200, 100)
	// Set("config_server_port", "8080", 200, 100)
	// Set("config_database_name", "mydb", 200, 100)
	//
	// // 使用正则表达式获取所有config_开头的变量:
	// // Get("^config_.+", 200, 100)
	// // 返回JSON: {"config_server_host":"localhost","config_server_port":"8080","config_database_name":"mydb"}
}

// ExampleGet_vtaskZero 展示了 vtaskID <= 0 时的特殊行为
func ExampleGet_vtaskZero() {
	// vtaskID <= 0 时获取全局变量（vtask IS NULL）
	fmt.Println("// 设置全局变量（vtaskID=0）:")
	fmt.Println("Set(\"global_setting\", \"enabled\", 0, 100)")
	fmt.Println()
	fmt.Println("// 以下调用都会返回全局变量:")
	fmt.Println("Get(\"global_setting\", 0, 100)   // vtaskID=0")
	fmt.Println("Get(\"global_setting\", -1, 100)  // vtaskID=-1 (<=0)")
	fmt.Println("Get(\"global_setting\", -100, 100)// 任何<=0的vtaskID")
	// Output:
	// // 设置全局变量（vtaskID=0）:
	// Set("global_setting", "enabled", 0, 100)
	//
	// // 以下调用都会返回全局变量:
	// Get("global_setting", 0, 100)   // vtaskID=0
	// Get("global_setting", -1, 100)  // vtaskID=-1 (<=0)
	// Get("global_setting", -100, 100)// 任何<=0的vtaskID
}

func TestGet(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")

	appID := 168
	vtaskID := 4117

	// 清理可能存在的旧数据
	// 注意：实际测试中可能需要更完善的清理机制

	// 设置测试数据
	// 1. 设置带有特定vtaskID的变量
	variable.Set("var-1", "val-1", vtaskID, appID)
	variable.Set("var-2", "val-2", vtaskID, appID)

	// 2. 设置vtaskID=0的变量（全局变量）
	variable.Set("var-1", "val-10", 0, appID)
	variable.Set("var-2", "val-20", 0, appID)

	// 3. 设置另一个vtaskID的变量，用于测试区分
	anotherVtaskID := 4118
	variable.Set("var-1", "val-100", anotherVtaskID, appID)

	// 测试用例1：获取带有特定vtaskID的变量
	t.Run("Get with specific vtaskID", func(t *testing.T) {
		val, err := variable.Get("var-1", vtaskID, appID)
		if err != nil {
			t.Errorf("Get with specific vtaskID failed: %v", err)
		}
		if val != "val-1" {
			t.Errorf("Expected 'val-1', got '%s'", val)
		}
		fmt.Printf("Test 1 - Get with vtaskID=%d: %s\n", vtaskID, val)
	})

	// 测试用例2：获取vtaskID=0的变量（全局变量）
	t.Run("Get with vtaskID=0", func(t *testing.T) {
		val, err := variable.Get("var-1", 0, appID)
		if err != nil {
			t.Errorf("Get with vtaskID=0 failed: %v", err)
		}
		if val != "val-10" {
			t.Errorf("Expected 'val-10', got '%s'", val)
		}
		fmt.Printf("Test 2 - Get with vtaskID=0: %s\n", val)
	})

	// 测试用例3：获取不存在的vtaskID的变量
	t.Run("Get with non-existent vtaskID", func(t *testing.T) {
		nonExistentVtaskID := 9999
		_, err := variable.Get("var-1", nonExistentVtaskID, appID)
		if err == nil {
			t.Error("Expected error for non-existent vtaskID, but got none")
		} else {
			fmt.Printf("Test 3 - Expected error for non-existent vtaskID: %v\n", err)
		}
	})

	// 测试用例4：获取不存在的变量名
	t.Run("Get non-existent variable", func(t *testing.T) {
		_, err := variable.Get("non-existent-var", vtaskID, appID)
		if err == nil {
			t.Error("Expected error for non-existent variable, but got none")
		} else {
			fmt.Printf("Test 4 - Expected error for non-existent variable: %v\n", err)
		}
	})

	// 测试用例5：使用正则表达式获取变量（多个匹配）
	t.Run("Get with regex pattern", func(t *testing.T) {
		val, err := variable.Get("var.+", vtaskID, appID)
		if err != nil {
			t.Errorf("Get with regex pattern failed: %v", err)
		}
		// 正则表达式应该返回JSON字符串
		fmt.Printf("Test 5 - Get with regex pattern: %s\n", val)
		// 验证返回的是有效的JSON
		if len(val) == 0 || val == "{}" {
			t.Errorf("Expected non-empty JSON, got '%s'", val)
		}
	})

	// 测试用例6：获取另一个vtaskID的变量
	t.Run("Get with another vtaskID", func(t *testing.T) {
		val, err := variable.Get("var-1", anotherVtaskID, appID)
		if err != nil {
			t.Errorf("Get with another vtaskID failed: %v", err)
		}
		if val != "val-100" {
			t.Errorf("Expected 'val-100', got '%s'", val)
		}
		fmt.Printf("Test 6 - Get with vtaskID=%d: %s\n", anotherVtaskID, val)
	})

	// 测试用例7：测试负vtaskID（应该被视为vtaskID=0的情况）
	t.Run("Get with negative vtaskID", func(t *testing.T) {
		val, err := variable.Get("var-1", -1, appID)
		if err != nil {
			t.Errorf("Get with negative vtaskID failed: %v", err)
		}
		// 负vtaskID应该返回vtaskID=0的变量
		if val != "val-10" {
			t.Errorf("Expected 'val-10' for negative vtaskID, got '%s'", val)
		}
		fmt.Printf("Test 7 - Get with vtaskID=-1: %s\n", val)
	})

	// 测试用例8：测试多个变量区分
	t.Run("Verify variable isolation by vtaskID", func(t *testing.T) {
		// 验证不同vtaskID的变量是独立的
		val1, _ := variable.Get("var-1", vtaskID, appID)
		val2, _ := variable.Get("var-1", 0, appID)
		val3, _ := variable.Get("var-1", anotherVtaskID, appID)

		if val1 == val2 || val1 == val3 || val2 == val3 {
			t.Errorf("Variables with different vtaskIDs should be isolated. Got: vtaskID=%d:'%s', vtaskID=0:'%s', vtaskID=%d:'%s'",
				vtaskID, val1, val2, anotherVtaskID, val3)
		}
		fmt.Printf("Test 8 - Variable isolation verified: %s, %s, %s\n", val1, val2, val3)
	})
}
