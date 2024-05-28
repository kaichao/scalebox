package misc

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

// NewSQLNullString ...
func NewSQLNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

var (
	db *sql.DB
)

// GetDB ...
func GetDB() *sql.DB {
	if db == nil {
		var databaseURL string
		if databaseURL = os.Getenv("DATABASE_URL"); len(databaseURL) == 0 {

			dbHost := os.Getenv("PGHOST")
			dbPort := os.Getenv("PGPORT")
			if dbHost == "" {
				dbHost = os.Getenv("LOCAL_ADDR")
				if dbHost == "" {
					dbHost = GetLocalIP()
				}
			}
			ss := strings.Split(dbHost, ":")
			if len(ss) == 2 {
				dbHost = ss[0]
				dbPort = ss[1]
			}
			if dbPort == "" {
				dbPort = "5432"
			}
			databaseURL = fmt.Sprintf("postgres://scalebox:changeme@%s:%s/scalebox", dbHost, dbPort)
		}
		// set database connection
		var err error
		if db, err = sql.Open("pgx", databaseURL); err != nil {
			log.Fatal("Unable to connect to database:", err)
		}
		db.SetConnMaxLifetime(500)
		db.SetMaxIdleConns(50)
		db.SetMaxOpenConns(20)
		// db.Stats()
	}
	return db
}
