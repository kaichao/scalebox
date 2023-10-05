package main

import (
	"database/sql"
	"os"

	scalebox "github.com/kaichao/scalebox/golang/misc"
	"github.com/sirupsen/logrus"
)

// DataSet ...
type DataSet struct {
	DatasetID string

	KeyGroupRegex string
	KeyGroupIndex string

	RootDir string
	SinkJob string

	// "H" / "V"
	GroupType string

	// for type "H", x-coord
	HorizontalWidth int

	// for type "V", y-coord
	VerticalStart  int
	VerticalLength int
	// vertical group length
	GroupLength int
	// vertical interleaved
	Interleaved bool
}

// Entity ...
type Entity struct {
	ID        int
	name      string
	datasetID string
	x         int
	y         int
	flag      string
}

var (
	db         *sql.DB
	mapDataset = make(map[string]*DataSet)

	logger *logrus.Logger

	messageFile string
	datasetFile string
	sqliteFile  string
)

func init() {
	var err error
	logrus.SetReportCaller(true)

	messageFile = os.Getenv("MESSAGE_FILE")
	if messageFile == "" {
		messageFile = "/tmp/messages.txt"
	}
	datasetFile = os.Getenv("DATASET_FILE")
	if datasetFile == "" {
		datasetFile = "/tmp/datasets.txt"
	}
	sqliteFile = os.Getenv("SQLITE_FILE")
	if sqliteFile == "" {
		sqliteFile = "/tmp/my.db"
	}

	// set database connection
	if db, err = sql.Open("sqlite3", sqliteFile); err != nil {
		logrus.Fatalln("Unable to open sqlite3 database:", err)
	}
	sqlText := `
		CREATE TABLE IF NOT EXISTS t_entity (
			id INTEGER PRIMARY KEY autoincrement,
			name TEXT,
			dataset_id TEXT,
			x INTEGER,
			y INTEGER,
			flag TEXT
		);
		CREATE UNIQUE INDEX IF NOT EXISTS i_entity_0 ON t_entity(name,dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_1 ON t_entity(dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_2 ON t_entity(dataset_id,x);
	`

	if _, err = db.Exec(sqlText); err != nil {
		logrus.Fatal(err)
	}

	if lines, err := scalebox.GetTextFileLines(datasetFile); err == nil {
		for _, line := range lines {
			dataset := parseDataSet(line)
			mapDataset[dataset.DatasetID] = dataset
		}
	}
}
