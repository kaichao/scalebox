package semagroup

import (
	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
)

// DiffMax 计算信号量组最大值与当前信号量值的差值
// semaExpr:正则表达式，其形式为(group-prefix)expr，其中圆括号中group-prefix标识为信号量组的前缀
func DiffMax(semaExpr string, appID int) (int, error) {
	// 解析group表达式
	groupExpr, expr, err := parseGroupExpr(semaExpr)
	if err != nil {
		return 0, err
	}

	// 在事务中执行查询
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return 0, errors.WrapE(err, "begin transaction")
	}
	defer tx.Rollback()

	// 查询当前信号量值
	var currentValue int
	err = tx.QueryRow(`
		SELECT value 
		FROM t_semaphore 
		WHERE name = $1 AND app = $2
	`,
		expr, appID).Scan(&currentValue)
	if err != nil {
		return 0, errors.WrapE(err, "query current semaphore",
			"app-id", appID, "sema-expr", expr)
	}

	// 查询组最大值
	var maxValue int
	err = tx.QueryRow(`
		SELECT MAX(value)
		FROM t_semaphore
		WHERE name ~ $1 AND app = $2
	`,
		groupExpr, appID).Scan(&maxValue)
	if err != nil {
		return 0, errors.WrapE(err, "query max semaphore",
			"app-id", appID, "group-expr", groupExpr)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return 0, errors.WrapE(err, "commit transaction")
	}

	return maxValue - currentValue, nil
}
