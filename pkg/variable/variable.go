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
func Get(name string, appID int) (string, error) {
	sqlText := `
		WITH selected_rows AS (
			SELECT name,value
			FROM t_variable
			WHERE app=$2 AND (name ~ $1)
			ORDER BY 1
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM selected_rows
	`
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in get-variable(%s,%d), err-t=%T,err=%v",
			name, appID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	if common.IsRegexString(name) {
		return packed, nil
	}
	if packed == "{}" {
		// no semaphore
		errInfo := fmt.Sprintf("[ERROR]Variable(name:%s,app-id:%d) not-found", name, appID)
		logrus.Errorln(errInfo)
		return "", errors.New(errInfo)
	}
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
func Set(name string, value string, appID int) error {
	sqlText := `
		INSERT INTO t_variable(name,value,app)
		VALUES($1,$2,$3)
		ON CONFLICT (name,app)
		DO UPDATE SET value = EXCLUDED.value;
	`

	result, err := postgres.GetDB().Exec(sqlText, name, value, appID)
	if err != nil {
		logrus.Errorf("db-error in set-variable (name:%s,app-id:%d), err-t=%T,err=%v\n",
			name, appID, err, err)
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		logrus.Errorf(" Variable %s not-defined\n", name)
		return err
	}

	return nil
}
