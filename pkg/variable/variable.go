package variable

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// Get ...
func Get(name string, vtaskID int64, appID int) (string, error) {
	// 构建SQL查询，考虑vtaskID参数
	sqlText := ""
	var v string
	var err error

	if vtaskID <= 0 {
		// vtaskID <= 0 时，查询vtask IS NULL的记录
		sqlText = `
			WITH selected_rows AS (
				SELECT name,value
				FROM t_variable
				WHERE app=$2 AND (name ~ $1) AND vtask IS NULL
				ORDER BY 1
			)
			SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
			FROM selected_rows
		`
		err = postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v)
	} else {
		// vtaskID > 0 时，需要匹配vtask参数
		sqlText = `
			WITH selected_rows AS (
				SELECT name,value
				FROM t_variable
				WHERE app=$2 AND (name ~ $1) AND vtask=$3
				ORDER BY 1
			)
			SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
			FROM selected_rows
		`
		err = postgres.GetDB().QueryRow(sqlText, name, appID, vtaskID).Scan(&v)
	}
	logrus.Debugf("In variable.Get(),name=%s,value=%s,vtask-id:%d,app-id:%d,err:%v\n",
		name, v, vtaskID, appID, err)

	if err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in get-variable(%s,%d,vtask:%d), err-t=%T,err=%v",
			name, appID, vtaskID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if common.IsRegexString(name) {
		return packed, nil
	}
	if packed == "{}" {
		// no variable found
		errInfo := fmt.Sprintf("[ERROR]Variable(name:%s,app-id:%d,vtask:%d) not-found", name, appID, vtaskID)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	// 提取变量值的字符串
	re := regexp.MustCompile(`{".+":"(.+)"}`)
	ss := re.FindStringSubmatch(packed)
	if len(ss) == 0 {
		errInfo := fmt.Sprintf("[ERROR]Invalid JSON string, value=%s", packed)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
	return ss[1], nil
}

// Set ...
func Set(name string, value string, vtaskID int64, appID int) error {
	sqlText := `
		INSERT INTO t_variable(name,value,vtask,app)
		VALUES($1,$2,$3,$4)
		ON CONFLICT (name,vtask,app)
		DO UPDATE SET value = EXCLUDED.value;
	`

	pVtaskID := &vtaskID
	if vtaskID <= 0 {
		pVtaskID = nil
	}
	result, err := postgres.GetDB().Exec(sqlText, name, value, pVtaskID, appID)
	logrus.Debugf("In variable.Set(),name=%s,value=%s,vtask-id:%d,app-id:%d,err:%v\n",
		name, value, vtaskID, appID, err)
	if err != nil {
		logrus.Errorf("db-error in set-variable (name:%s,app-id:%d,vtask:%d), err-t=%T,err=%v\n",
			name, appID, vtaskID, err, err)
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		logrus.Errorf(" Variable %s not-defined\n", name)
		return err
	}

	return nil
}
