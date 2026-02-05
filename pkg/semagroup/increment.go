package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// Increment 对信号量组，其值最小的信号量加一，返回最小信号量及其值，用冒号连接
// 信号量表达式定义与GetMax/GetMin相同
func Increment(semaExpr string, appID int) (string, int, error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 在事务中执行查询和更新
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return "", 0, errors.WrapE(err, "begin transaction")
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
		return "", 0, errors.WrapE(err, "query min semaphore",
			"app-id", appID, "sema-expr", semaExpr)
	}

	// 更新信号量值
	_, err = tx.Exec(`
		UPDATE t_semaphore 
		SET value = value + 1 
		WHERE name = $1 AND app = $2`,
		name, appID)
	if err != nil {
		return "", 0, errors.WrapE(err, "increment semagroup",
			"app-id", appID, "name", name)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return "", 0, errors.WrapE(err, "commit transaction")
	}

	return name, value + 1, nil
}
