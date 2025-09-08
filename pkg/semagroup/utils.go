package semagroup

import (
	"fmt"
	"regexp"
)

// processSemaExpr 处理信号量表达式
func processSemaExpr(semaExpr string) string {
	// 不再自动添加通配符
	return semaExpr
}

// parseGroupExpr 解析(group-prefix)expr格式的表达式
// 返回: groupExpr, fullExpr, error
func parseGroupExpr(semaExpr string) (string, string, error) {
	// 提取group-prefix

	// 正则分组：第1组括号里的内容，第2组后面的内容
	matches := regexp.MustCompile(`^\(([^)]*)\)(.*)$`).FindStringSubmatch(semaExpr)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("not valid semagroup expression ")
	}
	part1 := matches[1]
	part2 := matches[2]

	return part1 + ".*", part1 + part2, nil
}
