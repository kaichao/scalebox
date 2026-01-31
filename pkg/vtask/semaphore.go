package vtask

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// CreateSemaphore ...
func CreateSemaphore(name string, value int, vtaskID int64, appID int) error {
	pVtask := &vtaskID
	if vtaskID <= 0 {
		pVtask = nil
	}

	// 根据环境变量 CONFLICT_ACTION 决定冲突处理逻辑
	conflictAction := os.Getenv("CONFLICT_ACTION")
	var sqlText string

	switch conflictAction {
	case "OVERWRITE":
		// 覆盖现有值
		sqlText = `
			INSERT INTO t_semaphore(name,value,value0,vtask,app)
			VALUES($1,$2,$2,$3,$4)
			ON CONFLICT (name, vtask, app)
				DO UPDATE SET
				value  = EXCLUDED.value,
				value0 = EXCLUDED.value0
		`
	case "IGNORE":
		// 忽略冲突，不报错
		sqlText = `
			INSERT INTO t_semaphore(name,value,value0,vtask,app)
			VALUES($1,$2,$2,$3,$4)
			ON CONFLICT (name, vtask, app) DO NOTHING
		`
	default:
		// 默认行为：报错（不使用 ON CONFLICT 子句）
		sqlText = `
			INSERT INTO t_semaphore(name,value,value0,vtask,app)
			VALUES($1,$2,$2,$3,$4)
		`
	}
	logrus.Tracef("In semaphore.Create(),sqlText:%s\n", sqlText)

	if _, err := postgres.GetDB().Exec(sqlText, name, value, pVtask, appID); err != nil {
		return errors.WrapE(err, "semaphore-create failed",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "value", value, "conflict-action", conflictAction)
	}
	logrus.Tracef("semaphore-create: name=%s,value=%d,vtask-id=%d,app-id=%d,conflict-action=%s\n",
		name, value, vtaskID, appID, conflictAction)

	return nil
}

// GetSemaphore ...
func GetSemaphore(name string, vtaskID int64, appID int) (value int, err error) {
	sqlFmt := `
		SELECT value
		FROM t_semaphore
		WHERE app=$2 AND name=$1 AND %s
	`

	if vtaskID > 0 {
		// vtaskID > 0 时，需要匹配vtask参数
		vtaskExpr := "vtask = $3"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, vtaskID).Scan(&value)

	} else {
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID).Scan(&value)
	}
	logrus.Tracef("In semaphore.GetValue(),name=%s,vtask-id:%d,app-id:%d,value:%d,err:%v\n",
		name, vtaskID, appID, value, err)

	if err == nil {
		return value, nil
	}

	if err != sql.ErrNoRows {
		return -1, errors.WrapE(err, "get semaphore failed",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}
	// not-defined semaphore
	if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
		// create semaphore first time
		if err := CreateSemaphore(name, 0, vtaskID, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore failed",
				"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
}

// AddSemaphore ...
func AddSemaphore(name string, delta int, vtaskID int64, appID int) (v int, err error) {
	sqlFmt := `
		UPDATE t_semaphore
		SET value = value + $3
		WHERE name = $1 AND app = $2 AND %s
		RETURNING value
	`
	if vtaskID > 0 {
		// vtaskID > 0 时，需匹配vtask参数
		vtaskExpr := "vtask = $4"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, delta, vtaskID).Scan(&v)
	} else {
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, delta).Scan(&v)
	}
	logrus.Tracef("In semaphore.AddValue(),name=%s,vtask-id:%d,app-id:%d,delta:%d,ret-value:%d,err:%v\n",
		name, vtaskID, appID, delta, v, err)

	if err == nil {
		return v, nil
	}

	if err != sql.ErrNoRows {
		return v, errors.WrapE(err, "update semaphore failed",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
	}
	// not-defined semaphore
	if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
		// create semaphore first time
		if err := CreateSemaphore(name, delta, vtaskID, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore failed",
				"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "delta", delta)
}
