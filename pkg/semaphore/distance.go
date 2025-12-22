package semaphore

import (
	"database/sql"
	"errors"

	"github.com/kaichao/scalebox/pkg/postgres"
	"github.com/sirupsen/logrus"
)

// GroupDistance ...
// Deprecated
func GroupDistance(name string, appID int) (int, error) {
	sqlText := `
        WITH current_host AS (
            SELECT group_id, cluster
            FROM t_host 
            WHERE hostname = split_part($1,':',3)
        ),
        group_hosts AS (
            SELECT hostname 
            FROM t_host
            JOIN current_host ON t_host.cluster = current_host.cluster
            AND t_host.group_id IS NOT DISTINCT FROM current_host.group_id
        ),
        group_min AS (
            SELECT MIN(value) AS value
            FROM t_semaphore
            WHERE app = $2
                AND split_part(name,':',2) = split_part($1,':',2)
                AND split_part(name,':',3) IN (SELECT hostname FROM group_hosts)
        )
        SELECT t_semaphore.value - group_min.value
        FROM t_semaphore
        CROSS JOIN group_min
        WHERE name = $1 AND app = $2
	`
	return distance(name, appID, sqlText)
}

// GlobalDistance ...
// Deprecated
func GlobalDistance(name string, appID int) (int, error) {
	sqlText := `
        WITH current_host AS (
            SELECT group_id, cluster
            FROM t_host 
            WHERE hostname = split_part($1,':',3)
        ),
        global_hosts AS (
            SELECT hostname 
            FROM t_host
            JOIN current_host ON t_host.cluster = current_host.cluster
            AND t_host.group_id IS NOT NULL
        ),
        global_min AS (
            SELECT MIN(value) AS value
            FROM t_semaphore
            WHERE app = $2
                AND split_part(name,':',2) = split_part($1,':',2)
                AND split_part(name,':',3) IN (SELECT hostname FROM global_hosts)
        )
        SELECT t_semaphore.value - global_min.value
        FROM t_semaphore
        CROSS JOIN global_min
        WHERE name = $1 AND app = $2
	`
	return distance(name, appID, sqlText)
}

func distance(name string, appID int, sqlText string) (int, error) {
	n := 0
	if err := postgres.GetDB().QueryRow(sqlText, name, appID).Scan(&n); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// create semaphore first time, initial value=0
			err = Create(name, 0, appID)
			if err != nil {
				logrus.Errorf(" Semaphore (name:%s,app-id:%d), create error,err-info:%v\n",
					name, appID, err)
				return -1, err
			}
		} else {
			logrus.Warnf("db-error in semaphore-distance(name:%s,app-id:%d), err-t=%T,err=%v\n",
				name, appID, err, err)
			return -1, err
		}
	}
	return n, nil
}
