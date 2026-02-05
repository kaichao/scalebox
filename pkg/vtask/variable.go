package vtask

import (
	"fmt"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetVariable ...
func GetVariable(name string, vtaskID int64, appID int) (string, error) {
	sqlFmt := `
		SELECT value
		FROM t_variable
		WHERE app=$2 AND (name = $1) AND %s 
	`
	var v string
	var err error

	if vtaskID <= 0 {
		// vtaskID <= 0 时，查询vtask IS NULL的记录
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID).Scan(&v)
	} else {
		// vtaskID > 0 时，需要匹配vtask参数
		vtaskExpr := "vtask=$3"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, vtaskID).Scan(&v)
	}
	if err != nil {
		return "", errors.WrapE(err, "get variable",
			"app-id", appID, "vtask-id", vtaskID, "var-name", name)
	}
	return v, nil
}

// SetVariable ...
func SetVariable(name string, value string, vtaskID int64, appID int) error {
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
	if err != nil {
		return errors.WrapE(err, "set-variable",
			"app-id", appID, "vtask-id", vtaskID, "var-name", name, "var-value", value)
	}
	logrus.Tracef("In variable.Set(),name=%s,value=%s,vtask-id:%d,app-id:%d,err:%v\n",
		name, value, vtaskID, appID, err)

	if n, _ := result.RowsAffected(); n == 0 {
		return errors.E("variable not defined",
			"app-id", appID, "vtask-id", vtaskID, "var-name", name)
	}

	return nil
}
