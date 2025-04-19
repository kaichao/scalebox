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
func Get(name string, appID int) (string, error) {
	sqlText := `
		WITH selected_rows AS (
			SELECT name,value
			FROM t_semaphore
			WHERE app=$2 AND (name ~ $1)
			ORDER BY 1
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM selected_rows
	`
	if !common.IsRegexString(name) {
		sqlText = `
			WITH selected_rows AS (
				SELECT name,value
				FROM t_semaphore
				WHERE app=$2 AND (name = $1)
				ORDER BY 1
			)
			SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
			FROM selected_rows
		`
	}
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in get-semaphore (%s,%d), err-t=%T,err=%v",
			name, appID, err, err)
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
			if err := Create(name, 0, appID); err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d), create error,err-info:%v\n",
					name, appID, err)
				return "", err
			}
			return "0", nil
		}
		errInfo := fmt.Sprintf("[ERROR]Semaphore(name:%s,app-id:%d) not-found", name, appID)
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
