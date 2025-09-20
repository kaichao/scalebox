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
func Create(name string, value int, appID int) error {
	// overwrite existing
	sqlText := `
		INSERT INTO t_semaphore(name,value,value0,app)
		VAlUES($1,$2,$2,$3)
		ON CONFLICT (name, app)
			DO UPDATE SET
    		value  = EXCLUDED.value,
    		value0 = EXCLUDED.value0
	`

	if _, err := postgres.GetDB().Exec(sqlText, name, value, appID); err != nil {
		errInfo := fmt.Sprintf("semaphore-create: name=%s,value=%d,app-id=%d,err=%v",
			name, value, appID, err)
		logrus.Errorln(errInfo)
		return err
	}
	return nil
}

// CreateFileSemaphores ...
func CreateFileSemaphores(fileName string, appID int, batchSize int) error {
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
	return CreateSemaphores(semas, appID, batchSize)
}

// CreateJSONSemaphores ...
func CreateJSONSemaphores(jsonText string, appID int, batchSize int) error {
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
	return CreateSemaphores(ordered, appID, batchSize)
}

// Sema ...
type Sema struct {
	Name  string
	Value int
}

// CreateSemaphores ...
func CreateSemaphores(ordered []*Sema, appID int, batchSize int) error {
	// start transaction
	tx, err := postgres.GetDB().Begin()
	if err != nil {
		logrus.Errorf("err:%v\n", err)
		return err
	}
	defer tx.Rollback()

	// ignore existing
	sqlText := `
		INSERT INTO t_semaphore(name,value,value0,app) 
		VALUES($1,$2,$2,$3)
		ON CONFLICT (name, app) DO NOTHING; 
	`

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
			if _, err := stmt.Exec(v.Name, v.Value, appID); err != nil {
				logrus.Errorf("err:%v\n", err)
				return err
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
