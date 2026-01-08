package semagroup

import (
	"fmt"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"
)

// encodedSemaExpr 处理信号量表达式
func encodedSemaExpr(semaExpr string) string {
	// 首字母不是regex元字符，自动添加^
	if !common.IsRegexString(semaExpr[0:1]) {
		return "^" + semaExpr
	}
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
