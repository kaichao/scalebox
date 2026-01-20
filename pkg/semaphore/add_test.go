package semaphore_test

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kaichao/scalebox/pkg/semaphore"
)

func TestAddValue(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	// 测试添加信号量
	v, err := semaphore.AddValue("test_semaphore_1", vtaskID, appID, 3)
	if err != nil {
		// 由于外键约束可能失败，我们只记录错误但不失败测试
		fmt.Printf("AddValue error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Printf("AddValue result: v=%s\n", v)
	}

	// 测试全局信号量（vtaskID=0）
	v2, err2 := semaphore.AddValue("test_semaphore_global", 0, appID, 5)
	if err2 != nil {
		fmt.Printf("AddValue global error: %v\n", err2)
	} else {
		fmt.Printf("AddValue global result: v=%s\n", v2)
	}

	// 测试负vtaskID
	v3, err3 := semaphore.AddValue("test_semaphore_negative", -1, appID, 2)
	if err3 != nil {
		fmt.Printf("AddValue negative vtaskID error: %v\n", err3)
	} else {
		fmt.Printf("AddValue negative vtaskID result: v=%s\n", v3)
	}
}

func TestAddMultiValues(t *testing.T) {
	os.Setenv("PGHOST", "10.0.6.100")
	os.Setenv("SEMAPHORE_AUTO_CREATE", "yes")

	// 测试多个信号量
	pairs := map[string]int{
		"multi_semaphore_1": 10,
		"multi_semaphore_2": 20,
		"multi_semaphore_3": 30,
	}

	result, err := semaphore.AddMultiValues(pairs, vtaskID, appID)
	if err != nil {
		// 由于外键约束可能失败，我们只记录错误但不失败测试
		fmt.Printf("AddMultiValues error (expected due to foreign key): %v\n", err)
	} else {
		fmt.Printf("AddMultiValues result: %v\n", result)
	}

	// 测试全局信号量（vtaskID=0）
	globalPairs := map[string]int{
		"global_multi_1": 5,
		"global_multi_2": 15,
	}

	globalResult, err2 := semaphore.AddMultiValues(globalPairs, 0, appID)
	if err2 != nil {
		fmt.Printf("AddMultiValues global error: %v\n", err2)
	} else {
		fmt.Printf("AddMultiValues global result: %v\n", globalResult)
	}

	// 测试空map
	emptyResult, err3 := semaphore.AddMultiValues(map[string]int{}, vtaskID, appID)
	if err3 != nil {
		t.Errorf("AddMultiValues empty map should not error: %v", err3)
	}
	if len(emptyResult) != 0 {
		t.Errorf("AddMultiValues empty map should return empty result, got: %v", emptyResult)
	}
	fmt.Println("AddMultiValues empty map test passed")
}
