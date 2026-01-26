package common_test

import (
	"testing"

	"github.com/kaichao/scalebox/pkg/common"
)

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
