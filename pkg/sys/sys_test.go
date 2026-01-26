package sys_test

import (
	"strings"
	"testing"

	"github.com/kaichao/scalebox/pkg/sys"
)

// Named function for nested call
func nestedFunc() string {
	return sys.GetCurrentGoroutineStack()
}

// Named function for deep recursive call
func deepStack(depth int) string {
	if depth == 0 {
		return sys.GetCurrentGoroutineStack()
	}
	return deepStack(depth - 1)
}

func TestGetCurrentGoroutineStack(t *testing.T) {
	// Test 1: Basic call, check for current function name
	stack := sys.GetCurrentGoroutineStack()
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
