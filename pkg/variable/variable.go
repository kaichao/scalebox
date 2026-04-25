package variable

import (
	// "errors"

	"regexp"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetValue ...
func GetValue(name string, appID int) (string, error) {
	sqlText := `
		SELECT value
		FROM t_variable
		WHERE app=$2 AND (name = $1) AND vtask IS NULL
	`
	var v string
	err := postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v)
	return v, errors.WrapE(err, "get variable",
		"app-id", appID, "var-name", name)
}

// GetJSON ...
func GetJSON(name string, appID int) (string, error) {
	sqlText := `
		WITH selected_rows AS (
			SELECT name,value
			FROM t_variable
			WHERE app=$2 AND (name ~ $1) AND vtask IS NULL
			ORDER BY 1
		)
		SELECT COALESCE(JSON_OBJECT_AGG(name, value), '{}') AS aggregated_values
		FROM selected_rows
	`
	var v string
	err := postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&v)
	packed := regexp.MustCompile(`\s+`).ReplaceAllString(v, "")
	return packed, errors.WrapE(err, "get variable-json",
		"app-id", appID, "var-name", name)
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
		return errors.WrapE(err, "set-variable",
			"app-id", appID, "var-name", name, "var-value", value)
	}
	logrus.Tracef("In variable.Set(),name=%s,value=%s,app-id:%d,err:%v\n",
		name, value, appID, err)

	if n, _ := result.RowsAffected(); n == 0 {
		return errors.E("variable not defined", "app-id", appID, "var-name", name)
	}

	return nil
}
