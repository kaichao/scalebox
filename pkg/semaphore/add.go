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

// AddValue ...
func AddValue(name string, appID int, delta int) (string, error) {
	sqlFmt := `
		WITH updated_rows AS (
    		UPDATE t_semaphore
    		SET value = value + $3
    		WHERE (name %s $1) AND app = $2
    		RETURNING name,value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM updated_rows
	`
	op := "="
	if !common.IsRegexString(name) {
		op = "~"
	}
	sqlText := fmt.Sprintf(sqlFmt, op)
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, name, appID, delta).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%s,%d), err-t=%T,err=%v",
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
			if err := Create(name, delta, appID); err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d), create error,err-info:%v\n",
					name, appID, err)
				return "", err
			}
			return fmt.Sprintf("%d", delta), nil
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

// AddListValue ...
func AddListValue(names []string, appID int, delta int) (string, error) {
	sqlText := `
		WITH updated_rows AS (
    		UPDATE t_semaphore
    		SET value = value + $3
    		WHERE name = ANY($1) AND app = $2
    		RETURNING name,value
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM updated_rows
	`
	v := ""
	if err := postgres.GetDB().QueryRow(sqlText, names, appID, delta).Scan(&v); err != nil {
		errInfo := fmt.Sprintf("[ERROR]db-error in semaphore-op (%s,%d), err-t=%T,err=%v",
			names, appID, err, err)
		logrus.Errorln(errInfo)
		return "", err
	}

	return regexp.MustCompile(`\s+`).ReplaceAllString(v, ""), nil
}
