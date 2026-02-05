package global

import (
	"database/sql"

	"github.com/kaichao/gopkg/errors"
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
	logrus.Tracef("In global.Set(),global-name:%s,global-value:%s,err:%v\n",
		name, value, err)
	if err != nil {
		return errors.WrapE(err, "global-set", "name", name, "value", value)
	}
	return nil
}

// Get ...
func Get(name string) (string, error) {
	sqlText := `SELECT value FROM t_global WHERE name=$1`
	var value string
	err := postgres.GetDB().QueryRow(sqlText, name).Scan(&value)
	logrus.Tracef("In global.Get(),global-name:%s,global-value:%s,err:%v\n",
		name, value, err)
	if err == nil {
		return value, nil
	}
	if err == sql.ErrNoRows {
		return "", errors.WrapE(err, "global not found", "name", name)
	}
	return "", errors.WrapE(err, "global get", "name", name)
}
