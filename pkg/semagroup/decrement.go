package semagroup

import (
	"fmt"

	"github.com/kaichao/scalebox/pkg/postgres"
)

// Decrement 对信号量组，其值最大的信号量减一，返回最大信号量及其值，用冒号连接
// 信号量表达式定义与GetMax/GetMin相同
func Decrement(semaExpr string, appID int) (string, error) {
	// 处理信号量表达式
	semaExpr = processSemaExpr(semaExpr)

	// 在事务中执行查询和更新
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 查询最大信号量及其值
	var name string
	var value int
	err = tx.QueryRow(`
		SELECT name, value 
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value DESC
		LIMIT 1`,
		semaExpr, appID).Scan(&name, &value)
	if err != nil {
		return "", fmt.Errorf("failed to query max semaphore: %w", err)
	}

	// 更新信号量值
	_, err = tx.Exec(`
		UPDATE t_semaphore 
		SET value = value - 1 
		WHERE name = $1 AND app = $2`,
		name, appID)
	if err != nil {
		return "", fmt.Errorf("failed to decrement semaphore: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return fmt.Sprintf("%s:%d", name, value-1), nil
}
