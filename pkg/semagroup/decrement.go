package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// Decrement 对信号量组，其值最大的信号量减一，返回最大信号量及其值，用冒号连接
// 信号量表达式定义与GetMax/GetMin相同
func Decrement(semaExpr string, appID int) (string, int, error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 在事务中执行查询和更新
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return "", 0, errors.WrapE(err, "begin transaction")
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
		return "", 0, errors.WrapE(err, "query max semaphore",
			"app-id", appID, "sema-expr", semaExpr)
	}

	// 更新信号量值
	_, err = tx.Exec(`
		UPDATE t_semaphore 
		SET value = value - 1 
		WHERE name = $1 AND app = $2`,
		name, appID)
	if err != nil {
		return "", 0, errors.WrapE(err, "decrement semagroup",
			"app-id", appID, "name", name)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return "", 0, errors.WrapE(err, "commit transaction")
	}

	return name, value - 1, nil
}
