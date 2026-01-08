package semagroup

import (
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMin 获取信号量组最小值
func GetMin(semaExpr string, appID int) (string, error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	var name string
	var minValue int
	err := postgres.GetDB().QueryRow(`
		SELECT name, value
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value, name
		LIMIT 1`,
		semaExpr, appID).Scan(&minValue)
	if err != nil {
		return "", fmt.Errorf("failed to query min semaphore value: %w", err)
	}
	return fmt.Sprintf(`"%s":%d`, name, minValue), nil
}
