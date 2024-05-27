package main

import (
	"database/sql"
	"fmt"
	"os"

	misc "github.com/kaichao/scalebox/pkg/misc"

	"github.com/sirupsen/logrus"
)

var (
	db         *sql.DB
	mapDataset = make(map[string]*DataSet)

	logger *logrus.Logger

	datasetFile   string
	workDir       string
	datasetPrefix string
)

func init() {
	var err error
	logrus.SetReportCaller(true)

	workDir = os.Getenv("WORD_DIR")
	if workDir == "" {
		workDir = "/work"
	}
	datasetFile = workDir + "/.scalebox/datasets.txt"

	datasetPrefix = os.Getenv("DATASET_PREFIX")

	// set database connection
	if db, err = sql.Open("sqlite3", workDir+"/.scalebox/sqlite.db"); err != nil {
		logrus.Fatalln("Unable to open sqlite3 database:", err)
	}
	sqlTextFmt := `
		CREATE TABLE IF NOT EXISTS t_entity (
			id INTEGER PRIMARY KEY autoincrement,
			name TEXT,
			dataset_id TEXT,
			x %s,
			y %s,
			flag TEXT
		);
		CREATE UNIQUE INDEX IF NOT EXISTS i_entity_0 ON t_entity(name,dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_1 ON t_entity(dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_2 ON t_entity(dataset_id,x);
	`

	sqlText := fmt.Sprintf(sqlTextFmt, "TEXT", "TEXT")
	if os.Getenv("COORD_TYPE") == "integer" {
		sqlText = fmt.Sprintf(sqlTextFmt, "INTEGER", "INTEGER")
	}

	if _, err = db.Exec(sqlText); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	if lines, err := misc.GetTextFileLines(datasetFile); err == nil {
		for _, line := range lines {
			dataset := parseDataSet(line)
			fmt.Println("loaded dataset-id:", dataset.DatasetID, "dataset:", dataset)
			mapDataset[dataset.DatasetID] = dataset
		}
	}
	fmt.Println("mapDataset:", mapDataset)
}
