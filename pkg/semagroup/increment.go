package semagroup

import (
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
)

// Increment 对信号量组，其值最小的信号量加一，返回最小信号量及其值，用冒号连接
// 信号量表达式定义与GetMax/GetMin相同
func Increment(semaExpr string, appID int) (string, error) {
	// 处理信号量表达式
	semaExpr = processSemaExpr(semaExpr)

	// 在事务中执行查询和更新
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 查询最小信号量及其值
	var name string
	var value int
	err = tx.QueryRow(`
		SELECT name, value 
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value ASC
		LIMIT 1`,
		semaExpr, appID).Scan(&name, &value)
	if err != nil {
		return "", fmt.Errorf("failed to query min semaphore: %w", err)
	}

	// 更新信号量值
	_, err = tx.Exec(`
		UPDATE t_semaphore 
		SET value = value + 1 
		WHERE name = $1 AND app = $2`,
		name, appID)
	if err != nil {
		return "", fmt.Errorf("failed to increment semaphore: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return fmt.Sprintf("%s:%d", name, value+1), nil
}
