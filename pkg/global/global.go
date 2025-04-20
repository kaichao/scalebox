package global

import (
	"database/sql"

	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// Set ...
func Set(name string, value string) error {
	sqlText := `
		INSERT INTO t_global (name, value)
		VALUES ($1, $2)
		ON CONFLICT (name)
		DO UPDATE SET value = EXCLUDED.value;
	`
	_, err := postgres.GetDB().Exec(sqlText, name, value)
	if err != nil {
		logrus.Errorf("Unable to global-set,name=%s, value=%s, err=%T, err=%v\n",
			name, value, err, err)
		return err
	}
	return nil
}

// Get ...
func Get(name string) (string, error) {
	sqlText := `SELECT value FROM t_global WHERE name=$1`
	var value string
	err := postgres.GetDB().QueryRow(sqlText, name).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("global %s not-found\n", name)
		} else {
			logrus.Errorf("Unable to global-get,name=%s,  err=%T, err=%v\n",
				name, err, err)
		}
		return "", err
	}

	return value, nil
}
