package semagroup

import (
	"fmt"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// GetMax 获取信号量组最大值
func GetMax(semaExpr string, appID int) (string, error) {
	// 处理信号量表达式
	semaExpr = encodedSemaExpr(semaExpr)

	// 从t_semaphore表中，做postgres的sql查询
	var name string
	var maxValue int
	err := postgres.GetDB().QueryRow(`
		SELECT name, value
		FROM t_semaphore 
		WHERE name ~ $1 AND app = $2
		ORDER BY value DESC, name ASC
		LIMIT 1`,
		semaExpr, appID).Scan(&name, &maxValue)
	if err != nil {
		return "", errors.WrapE(err, "query max semaphore value failed",
			"app-id", appID, "sema-expr", semaExpr)
	}

	return fmt.Sprintf(`"%s":%d`, name, maxValue), nil
}
