package common_test

import (
	"strings"
	"testing"

	"github.com/kaichao/gopkg/common"
)

// Named function for nested call
func nestedFunc() string {
	return common.GetCurrentGoroutineStack()
}

// Named function for deep recursive call
func deepStack(depth int) string {
	if depth == 0 {
		return common.GetCurrentGoroutineStack()
	}
	return deepStack(depth - 1)
}

func TestGetCurrentGoroutineStack(t *testing.T) {
	// Test 1: Basic call, check for current function name
	stack := common.GetCurrentGoroutineStack()
	if !strings.Contains(stack, "TestGetCurrentGoroutineStack") {
		t.Errorf("Expected stack to contain 'TestGetCurrentGoroutineStack', got:\n%s", stack)
	}

	// Test 2: Nested call, verify multi-layer stack
	stack = nestedFunc()
	if !strings.Contains(stack, "nestedFunc") || !strings.Contains(stack, "TestGetCurrentGoroutineStack") {
		t.Errorf("Expected stack to contain 'nestedFunc' and 'TestGetCurrentGoroutineStack', got:\n%s", stack)
	}

	// Test 3: Deep call, simulate complex stack before error
	stack = deepStack(5) // Simulate 5 layers of recursion
	if !strings.Contains(stack, "deepStack") {
		t.Errorf("Expected stack to contain 'deepStack', got:\n%s", stack)
	}
}

// TestIsRegexString 测试 IsRegexString 函数
func TestIsRegexString(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"hello", false},       // 普通字符串
		{"12345", false},       // 纯数字
		{"hello.*", true},      // 含元字符
		{"a+b", true},          // 含元字符
		{".*+?", true},         // 多个元字符
		{"\\.", false},         // 转义字符，当前视为普通字符
		{"", false},            // 空字符串
		{"*", true},            // 仅元字符
		{"?", true},            // 仅元字符
		{"[a-z]", true},        // 范围元字符
		{"normal text", false}, // 含空格的普通字符串
	}

	for _, test := range tests {
		result := common.IsRegexString(test.input)
		if result != test.expected {
			t.Errorf("For input %q, expected %v, got %v", test.input, test.expected, result)
		}
	}
}
