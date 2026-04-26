package vtask

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GetSemaphore ...
func GetSemaphore(name string, vtaskID int64, appID int) (value int, err error) {
	sqlFmt := `
		SELECT value
		FROM t_semaphore
		WHERE app=$2 AND name=$1 AND %s
	`

	if vtaskID > 0 {
		// vtaskID > 0 时，需要匹配vtask参数
		vtaskExpr := "vtask = $3"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID, vtaskID).Scan(&value)

	} else {
		vtaskExpr := "vtask IS NULL"
		err = postgres.GetDB().QueryRow(fmt.Sprintf(sqlFmt, vtaskExpr),
			name, appID).Scan(&value)
	}
	logrus.Tracef("In semaphore.GetValue(),name=%s,vtask-id:%d,app-id:%d,value:%d,err:%v\n",
		name, vtaskID, appID, value, err)

	if err == nil {
		return value, nil
	}

	if err != sql.ErrNoRows {
		return -1, errors.WrapE(err, "get semaphore",
			"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
	}
	// not-defined semaphore
	if os.Getenv("SEMAPHORE_AUTO_CREATE") == "yes" {
		// create semaphore first time
		if err := CreateSemaphore(name, 0, vtaskID, appID); err != nil {
			return -1, errors.WrapE(err, "create semaphore",
				"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
		}
		return 0, nil
	}
	return -1, errors.WrapE(err, "semaphore not found",
		"app-id", appID, "vtask-id", vtaskID, "sema-name", name)
}
