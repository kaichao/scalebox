package semaphore

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"

	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// Create ...
func Create(name string, value int, vtaskID int64, appID int) error {
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

	if _, err := postgres.GetDB().Exec(sqlText, name, value, pVtask, appID); err != nil {
		errInfo := fmt.Sprintf("semaphore-create: name=%s,value=%d,vtask-id=%d,app-id=%d,conflict-action=%s,err=%v",
			name, value, vtaskID, appID, conflictAction, err)
		logrus.Errorln(errInfo)
		return err
	}
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
			logrus.Errorf("Not matched semaphore :%s,\n", line)
		}
	}
	return createSemaphores(semas, vtaskID, appID, batchSize)
}

// CreateFileSemaphores ...
func CreateFileSemaphores(fileName string, vtaskID int64, appID int, batchSize int) error {
	lines, err := common.GetTextFileLines(fileName)
	if err != nil {
		logrus.Errorf("file-name:%s, err-info:%v\n", fileName, err)
		return err
	}
	if len(lines) == 0 {
		logrus.Warnf("Null file, filename:%s\n", fileName)
	}
	var semas []*Sema
	re := regexp.MustCompile(`"([^"]+)":(\d+)`)
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); len(matches) == 3 {
			key := matches[1] // 提取到的 key
			var value int
			fmt.Sscanf(matches[2], "%d", &value) // 将 value 转为整数
			semas = append(semas, &Sema{Name: key, Value: value})
		} else {
			logrus.Errorf("Not matched semaphore :%s,\n", line)
		}
	}
	return createSemaphores(semas, vtaskID, appID, batchSize)
}

// CreateJSONSemaphores ...
func CreateJSONSemaphores(jsonText string, vtaskID int64, appID int, batchSize int) error {
	// Define a struct type for the semaphores
	type semaItem struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	var items []semaItem

	// Try unmarshalling as top-level array of objects
	if err := json.Unmarshal([]byte(jsonText), &items); err == nil {
		// Successfully parsed as array
	} else {
		// If that fails, try as object with key-value pairs
		var kvPairs map[string]int
		if err := json.Unmarshal([]byte(jsonText), &kvPairs); err != nil {
			// Try alternative pattern expecting top-level "semaphores" property
			var jsonData struct {
				Semaphores map[string]int `json:"semaphores"`
			}
			if err := json.Unmarshal([]byte(jsonText), &jsonData); err != nil {
				logrus.Errorf("Invalid JSON format, err-info:%v, json-text:%s\n", err, jsonText)
				return err
			}
			kvPairs = jsonData.Semaphores
		}

		// Convert key-value map to array of semaItem
		items = make([]semaItem, 0, len(kvPairs))
		for name, value := range kvPairs {
			items = append(items, semaItem{Name: name, Value: value})
		}
	}

	// Create semaphore slice in input order
	ordered := make([]*Sema, 0, len(items))
	for _, item := range items {
		ordered = append(ordered, &Sema{Name: item.Name, Value: item.Value})
	}

	logrus.Debugf("Unmarshalled %d semaphores from JSON text", len(ordered))
	return createSemaphores(ordered, vtaskID, appID, batchSize)
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
		logrus.Errorf("err:%v\n", err)
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(ordered); i += batchSize {
		stmt, err := tx.Prepare(sqlText)
		if err != nil {
			logrus.Errorf("err:%v\n", err)
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
					logrus.Errorf("err:%v\n", err)
					return err
				}
				// IGNORE模式下，冲突错误被忽略
			}
		}
		if err = tx.Commit(); err != nil {
			logrus.Errorf("Commit, err-info:%v\n", err)
			return err
		}
		fmt.Fprintf(os.Stderr, "[%d..%d], %d row(s) inserted.\n", i, end-1, end-i)

		// start next batch
		if tx, err = postgres.GetDB().Begin(); err != nil {
			logrus.Errorf("err:%v\n", err)
			return err
		}
	}

	return nil
}
