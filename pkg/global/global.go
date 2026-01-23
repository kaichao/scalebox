package global

import (
	"database/sql"
	"errors"
	"fmt"

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
		return fmt.Errorf("unable to global-set,name=%s, value=%s, err: %w", name, value, err)
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
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("global %s not found: %w", name, err)
	}
	return "", fmt.Errorf("unable to global-get,name=%s, err: %w", name, err)
}
