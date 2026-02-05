package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMin 获取信号量组最小值
func GetMin(semaExpr string, appID int) (name string, value int, err error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	err = postgres.GetDB().QueryRow(`
		SELECT name, value
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value, name
		LIMIT 1`,
		semaExpr, appID).Scan(&value)
	if err != nil {
		return "", 0, errors.WrapE(err, "query min semaphore value failed",
			"app-id", appID, "sema-expr", semaExpr)
	}
	return name, value, nil
}
