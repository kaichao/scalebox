package semagroup

import (
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMin 获取信号量组最小值
func GetMin(semaExpr string, appID int) (int, error) {
	// 处理信号量表达式
	semaExpr = processSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	var minValue int
	err := postgres.GetDB().QueryRow(`
		SELECT MIN(value) 
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2`,
		semaExpr, appID).Scan(&minValue)
	if err != nil {
		return 0, fmt.Errorf("failed to query min semaphore value: %w", err)
	}

	return minValue, nil
}
