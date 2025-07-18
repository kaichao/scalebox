package semagrp

import (
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMax 获取信号量组最大值
func GetMax(semaExpr string, appID int) (int, error) {
	// 处理信号量表达式
	semaExpr = processSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	var maxValue int
	err := postgres.GetDB().QueryRow(`
		SELECT MAX(value) 
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2`,
		semaExpr, appID).Scan(&maxValue)
	if err != nil {
		return 0, fmt.Errorf("failed to query max semaphore value: %w", err)
	}

	return maxValue, nil
}
