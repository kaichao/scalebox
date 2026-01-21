package semaphore

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetJSON ...
//
//	按前缀获取semaphore列表name/value的json格式
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
	logrus.Debugf("In semaphore.GetValue(),name=%s,vtask-id:%d,app-id:%d,json-value:%s,err:%v\n",
		name, vtaskID, appID, v, err)

	if err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in get-semaphore (%s,%d,vtask:%d), err-t=%T,err=%v",
			name, appID, vtaskID, err, err)
		logrus.Errorln(errInfo)
		return "{}", err
	}
	// 删除结果的空字符
	v = regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	return v, nil
	// if common.IsRegexString(name) {
	// 	// regex name, return json-value
	// 	return packed, nil
	// }
	// // non-regex, not null, return int-value
	// re := regexp.MustCompile(`{".+":(-?[0-9]+)}`)
	// ss := re.FindStringSubmatch(packed)
	// if len(ss) == 0 {
	// 	errInfo := fmt.Sprintf("[ERROR]Invalid JSON string, value=%s", packed)
	// 	logrus.Errorln(errInfo)
	// 	return "", errors.New(errInfo)
	// }
	// return ss[1], nil
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
	logrus.Debugf("In semaphore.GetValue(),name=%s,vtask-id:%d,app-id:%d,value:%d,err:%v\n",
		name, vtaskID, appID, value, err)

	if err == nil {
		return value, nil
	}

	// 检查是否为"未找到"错误
	// 使用 errors.Is 检查 sql.ErrNoRows，通过 postgres 包间接引用
	if errors.Is(err, sql.ErrNoRows) {
		// not-defined semaphore
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			// create semaphore first time
			if createErr := Create(name, 0, vtaskID, appID); createErr != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d,vtask:%d), create error,err-info:%v\n",
					name, appID, vtaskID, createErr)
				return -1, createErr
			}
			return 0, nil
		}
		errInfo := fmt.Sprintf("[ERROR]Semaphore(name:%s,app-id:%d,vtask:%d) not-found", name, appID, vtaskID)
		logrus.Errorln(errInfo)
		return 0, errors.New(errInfo)
	}
	// 其他数据库错误
	errInfo := fmt.Sprintf("[ERROR]db-error in get-semaphore (%s,%d,vtask:%d), err-t=%T,err=%v",
		name, appID, vtaskID, err, err)
	logrus.Errorln(errInfo)
	return -1, err
}
