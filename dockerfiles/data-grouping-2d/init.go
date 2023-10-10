package main

import (
	"database/sql"
	"fmt"
	"os"

	scalebox "github.com/kaichao/scalebox/golang/misc"
	"github.com/sirupsen/logrus"
)

// DataSet ...
type DataSet struct {
	// prefix ':' type ':' sub-id
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
	VerticalHeight int
	// vertical group length
	GroupSize int
	// GroupStep   int
	// vertical interleaved
	Interleaved bool
}

// Entity ...
type Entity struct {
	ID        int
	name      string
	datasetID string
	x         string
	y         string
	flag      string
}

var (
	db         *sql.DB
	mapDataset = make(map[string]*DataSet)

	logger *logrus.Logger

	messageFile string
	datasetFile string
	sqliteFile  string

	datasetPrefix  string
	isIntegerCoord bool
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

	datasetPrefix = os.Getenv("DATASET_PREFIX")
	isIntegerCoord = os.Getenv("COORD_TYPE") == "integer"

	// set database connection
	if db, err = sql.Open("sqlite3", sqliteFile); err != nil {
		logrus.Fatalln("Unable to open sqlite3 database:", err)
	}
	sqlTextInteger := `
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

	sqlTextString := `
		CREATE TABLE IF NOT EXISTS t_entity (
			id INTEGER PRIMARY KEY autoincrement,
			name TEXT,
			dataset_id TEXT,
			x TEXT,
			y TEXT,
			flag TEXT
		);
		CREATE UNIQUE INDEX IF NOT EXISTS i_entity_0 ON t_entity(name,dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_1 ON t_entity(dataset_id);
		CREATE INDEX IF NOT EXISTS i_entity_2 ON t_entity(dataset_id,x);
	`

	sqlText := sqlTextString
	if isIntegerCoord {
		sqlText = sqlTextInteger
	}
	if _, err = db.Exec(sqlText); err != nil {
		logrus.Fatal(err)
	}

	if lines, err := scalebox.GetTextFileLines(datasetFile); err == nil {
		for _, line := range lines {
			dataset := parseDataSet(line)
			fmt.Println("loaded dataset-id:", dataset.DatasetID, "dataset:", dataset)
			mapDataset[dataset.DatasetID] = dataset
		}
	}
	fmt.Println("mapDataset:", mapDataset)
}
