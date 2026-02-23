package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMax 获取信号量组最大值
func GetMax(semaExpr string, appID int) (name string, value int, err error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	err = postgres.GetDB().QueryRow(`
		SELECT name, value
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value DESC, name
		LIMIT 1`,
		semaExpr, appID).Scan(&name, &value)
	if err != nil {
		return "", 0, errors.WrapE(err, "query max semaphore value failed",
			"app-id", appID, "sema-expr", semaExpr)
	}

	return name, value, nil
}
