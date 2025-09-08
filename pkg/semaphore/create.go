package semaphore

import (
	"bytes"
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
	var jsonData struct {
		RawData json.RawMessage `json:"semaphores"`
	}
	if err := json.Unmarshal([]byte(jsonText), &jsonData); err != nil {
		logrus.Errorf("Invalid json format, err-info:%v, json-text:%s\n", err, jsonText)
		return err
	}

	var ordered []*Sema
	decoder := json.NewDecoder(bytes.NewReader(jsonData.RawData))
	decoder.Token()
	for decoder.More() {
		k, _ := decoder.Token()
		var v int
		if err := decoder.Decode(&v); err != nil {
			logrus.Errorf("Error decoding semaphore value, err-info:%v, json-text:%s\n", err, jsonText)
			return err
		}
		ordered = append(ordered, &Sema{Name: k.(string), Value: v})
	}
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
