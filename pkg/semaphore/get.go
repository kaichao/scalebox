package semaphore

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// Get ...
func Get(name string, vtaskID int, appID int) (string, error) {
	// 构建SQL查询，考虑vtaskID参数
	sqlFmt := ""
	sqlFmt = `
			WITH selected_rows AS (
				SELECT name,value
				FROM t_semaphore
				WHERE app=$2 AND (name %s $1) AND %s
				ORDER BY 1
			)
			SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
			FROM selected_rows
		`

	vtaskExpr := "vtask IS NULL"
	if vtaskID > 0 {
		// vtaskID > 0 时，需要匹配vtask参数
		vtaskExpr = "vtask = $3"
	}

	op := "="
	if common.IsRegexString(name) {
		op = "~"
		if !common.IsRegexString(name[0:1]) {
			// 首字母不是regex元字符，自动添加^
			name = "^" + name
		}
	}
	sqlText := fmt.Sprintf(sqlFmt, op, vtaskExpr)

	var v string
	var err error
	if vtaskID <= 0 {
		err = postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v)
	} else {
		err = postgres.GetDB().QueryRow(sqlText, name, appID, vtaskID).Scan(&v)
	}

	if err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in get-semaphore (%s,%d,vtask:%d), err-t=%T,err=%v",
			name, appID, vtaskID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if common.IsRegexString(name) {
		// regex name, return json-value
		return packed, nil
	}
	if packed == "{}" {
		// not-defined semaphore
		if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
			// create semaphore first time
			if err := Create(name, 0, vtaskID, appID); err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d,vtask:%d), create error,err-info:%v\n",
					name, appID, vtaskID, err)
				return "", err
			}
			return "0", nil
		}
		errInfo := fmt.Sprintf("[ERROR]Semaphore(name:%s,app-id:%d,vtask:%d) not-found", name, appID, vtaskID)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	// non-regex, not null, return int-value
	re := regexp.MustCompile(`{".+":(-?[0-9]+)}`)
	ss := re.FindStringSubmatch(packed)
	if len(ss) == 0 {
		errInfo := fmt.Sprintf("[ERROR]Invalid JSON string, value=%s", packed)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	return ss[1], nil
}
