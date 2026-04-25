package semaphore

import (
	"database/sql"
	"os"
	"regexp"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetJSON ...
//
//	按regex获取semaphore列表name/value的json格式
func GetJSON(name string, appID int) (v string, err error) {
	sqlText := `
		WITH selected_rows AS (
			SELECT name,value
			FROM t_semaphore
			WHERE app=$2 AND (name ~ $1) AND vtask IS NULL
			ORDER BY 1
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM selected_rows
	`

	if !common.IsRegexString(name[0:1]) {
		// 首字母不是regex元字符，自动添加^
		name = "^" + name
	}

	err = postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v)
	if err != nil {
		v = "{}"
	} else {
		// 删除结果的空字符
		v = regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	}
	logrus.Tracef("In semaphore.GetValue(),name=%s,app-id:%d,json-value:%s,err:%v\n",
		name, appID, v, err)
	return v, errors.WrapE(err, "get-semaphore",
		"app-id", appID, "sema-name", name)
}

// GetValue ...
func GetValue(name string, appID int) (value int, err error) {
	sqlText := `
		SELECT value
		FROM t_semaphore
		WHERE app=$2 AND name=$1 AND vtask IS NULL
	`
	err = postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&value)
	logrus.Tracef("In semaphore.GetValue(),name=%s,app-id:%d,value:%d,err:%v\n",
		name, appID, value, err)
	if err == nil {
		return value, nil
	}
	if err != sql.ErrNoRows {
		return -1, errors.WrapE(err, "get semaphore",
			"app-id", appID, "sema-name", name)
	}
	// not-defined semaphore
	if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
		// create semaphore first time
		if err := Create(name, 0, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore",
				"app-id", appID, "sema-name", name)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "sema-name", name)
}
