package semaphore

import (
	"database/sql"
	"fmt"
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
func GetJSON(name string, vtaskID int64, appID int) (v string, err error) {
	// 构建SQL查询，考虑vtaskID参数
	sqlFmt := `
		WITH selected_rows AS (
			SELECT name,value
			FROM t_semaphore
			WHERE app=$2 AND (name ~ $1) AND %s
			ORDER BY 1
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM selected_rows
	`

	if !common.IsRegexString(name[0:1]) {
		// 首字母不是regex元字符，自动添加^
		name = "^" + name
	}

	if vtaskID > 0 {
		// vtaskID > 0 时，需匹配vtask参数
		vtaskExpr := "vtask = $3"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, vtaskID).Scan(&v)
	} else {
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID).Scan(&v)
	}
	if err != nil {
		return "{}", errors.WrapE(err, "get-semaphore failed",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}
	logrus.Tracef("In semaphore.GetValue(),name=%s,vtask-id:%d,app-id:%d,json-value:%s,err:%v\n",
		name, vtaskID, appID, v, err)

	// 删除结果的空字符
	v = regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	return v, nil
}

// GetValue ...
func GetValue(name string, vtaskID int64, appID int) (value int, err error) {
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
		if err := Create(name, 0, vtaskID, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore failed",
				"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
}
