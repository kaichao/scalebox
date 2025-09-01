package semagroup

import (
	"fmt"
	"strings"
)

// processSemaExpr 处理信号量表达式
func processSemaExpr(semaExpr string) string {
	// 不再自动添加通配符
	return semaExpr
}

// parseGroupExpr 解析(group-prefix):expr格式的表达式
// 返回: groupExpr, fullExpr, error
func parseGroupExpr(semaExpr string) (string, string, error) {
	// 提取group-prefix
	prefixStart := strings.Index(semaExpr, "(")
	prefixEnd := strings.Index(semaExpr, ")")
	if prefixStart == -1 || prefixEnd == -1 || prefixStart >= prefixEnd {
		return "", "", fmt.Errorf("invalid semaExpr format, expected (group-prefix):expr")
	}
	groupPrefix := semaExpr[prefixStart+1 : prefixEnd]
	expr := strings.TrimLeft(semaExpr[prefixEnd+1:], ":")

	if expr == "" {
		return "", "", fmt.Errorf("empty expression after group prefix")
	}

	// 组合完整信号量名称
	fullExpr := groupPrefix + ":" + expr
	groupExpr := groupPrefix + ".*"

	return groupExpr, fullExpr, nil
}
