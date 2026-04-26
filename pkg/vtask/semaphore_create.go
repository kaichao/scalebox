package vtask

import (
	"fmt"
	"os"
	"regexp"

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
		return errors.WrapE(err, "semaphore-create",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name, "value", value, "conflict-action", conflictAction)
	}
	logrus.Tracef("semaphore-create: name=%s,value=%d,vtask-id=%d,app-id=%d,conflict-action=%s\n",
		name, value, vtaskID, appID, conflictAction)

	return nil
}

// CreateSemaphores ...
func CreateSemaphores(lines []string, vtaskID int64, appID int, batchSize int) error {
	var semas []*Sema
	re := regexp.MustCompile(`"([^"]+)":(\d+)`)
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			key := matches[1] // 提取到的 key
			var value int
			fmt.Sscanf(matches[2], "%d", &value) // 将 value 转为整数
			semas = append(semas, &Sema{Name: key, Value: value})
		} else {
			logrus.Warnf("Not matched semaphore :%s,\n", line)
		}
	}
	return createSemaphores(semas, vtaskID, appID, batchSize)
}

// Sema ...
type Sema struct {
	Name  string
	Value int
}

func createSemaphores(ordered []*Sema, vtaskID int64, appID int, batchSize int) error {
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

	// start transaction
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		return errors.WrapE(err, "begin transaction")
	}
	defer tx.Rollback()

	for i := 0; i < len(ordered); i += batchSize {
		stmt, err := tx.Prepare(sqlText)
		if err != nil {
			logrus.Warnf("err:%v\n", err)
			continue
		}
		defer stmt.Close()

		end := i + batchSize
		if end > len(ordered) {
			end = len(ordered)
		}

		for _, v := range ordered[i:end] {
			if _, err := stmt.Exec(v.Name, v.Value, pVtask, appID); err != nil {
				// 如果存在冲突且不是IGNORE模式，记录错误但不一定失败
				if conflictAction != "IGNORE" {
					return errors.WrapE(err, "create prepared-semaphore",
						"app-id", appID, "vtask-id", vtaskID, "sema-name", v.Name, "sema-value", v.Value)
				}
				// IGNORE模式下，冲突错误被忽略
			}
		}
		if err = tx.Commit(); err != nil {
			return errors.WrapE(err, "commit transaction")
		}
		logrus.Infof("[%d..%d], %d row(s) inserted.\n", i, end-1, end-i)

		// start next batch
		if tx, err = postgres.GetDB().Begin(); err != nil {
			return errors.WrapE(err, "begin transaction")
		}
	}

	return nil
}
