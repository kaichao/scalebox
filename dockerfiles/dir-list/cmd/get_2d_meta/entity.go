package main

import (
	"database/sql"
	"os"

	scalebox "github.com/kaichao/scalebox/golang/misc"
	"github.com/sirupsen/logrus"
)

var (
	db                  *sql.DB
	workDir, sqliteFile string
)

func init() {
	workDir = os.Getenv("WORK_DIR")
	if workDir == "" {
		workDir = "/work"
	}
	sqliteFile := workDir + "/my.db"
	scalebox.ExecShellCommand("rm -f " + sqliteFile)
}

func initialize() {
	sqlText := `
		CREATE TABLE IF NOT EXISTS t_entity (
			id INTEGER PRIMARY KEY autoincrement,
			dataset_id TEXT,
			x INTEGER,
			y INTEGER
		);
		CREATE UNIQUE INDEX IF NOT EXISTS i_entity_0 ON t_entity(x,y);
	`
	if _, err := db.Exec(sqlText); err != nil {
		logrus.Errorln(err)
	}

}

func saveEntity(x, y int) bool {
	sqlText := `
		INSERT INTO t_entity(x,y) VALUES($1,$2)
	`
	_, err := db.Exec(sqlText, x, y)
	if err != nil {
		logrus.Errorf("add entity, err=%v\n", err)
		return false
	}
	return true
}

// getStat
//
//	input:
//		num :	Total number of entities
//	output:
//		width:	The width of the 2d-dataset in x-direction; if not 1d-dataset, then -1
//		minY:	The first element of 2d-dataset in y-direction;	 if not 2d-dataset, then -1
//		height:	The height of the 2d-dataset in y-direction; if not 2d-dataset, then -1
func getStat(num int) (int, int, int) {
	sqlText := `
		SELECT DISTINCT cnt
		FROM
			(SELECT y,count(*) cnt FROM t_entity GROUP BY 1) t
	`
	var width, minY, maxY int
	err := db.QueryRow(sqlText).Scan(&width)
	if err != nil {
		logrus.Warnf("Error getting width, err:%v\n", err)
		// not 1d-dataset
		return -1, 0, 0
	}

	sqlText = `SELECT min(y),max(y) FROM t_entity`
	err = db.QueryRow(sqlText).Scan(&minY, &maxY)
	if err != nil {
		logrus.Warnf("Error getting height, err:%v\n", err)
		// not 2d-dataset
		return width, -1, -1
	}
	if num != width*(maxY-minY+1) {
		// not continuous 2d-dataset
		return width, -1, -1
	}

	return width, minY, maxY - minY + 1
}
