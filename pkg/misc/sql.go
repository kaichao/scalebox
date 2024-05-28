package misc

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
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
	if db != nil {
		return db
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if len(databaseURL) == 0 {
		pgHost := os.Getenv("PGHOST")
		pgPort := os.Getenv("PGPORT")
		if pgHost == "" {
			pgHost = os.Getenv("LOCAL_ADDR")
			if pgHost == "" {
				pgHost = GetLocalIP()
			}
		}
		// ${PGHOST}:${PGPORT}
		ss := strings.Split(pgHost, ":")
		if len(ss) == 2 {
			pgHost = ss[0]
			pgPort = ss[1]
		}
		if pgPort == "" {
			pgPort = "5432"
		}
		pgUser := os.Getenv("PGUSER")
		if len(pgUser) == 0 {
			pgUser = "scalebox"
		}
		pgPass := os.Getenv("PGPASS")
		if len(pgPass) == 0 {
			pgPass = "changeme"
		}
		pgDB := os.Getenv("PGDB")
		if len(pgDB) == 0 {
			pgDB = "scalebox"
		}

		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			pgUser, pgPass, pgHost, pgPort, pgDB)
	}
	s := os.Getenv("PG_MAX_IDLE_CONNS")
	maxIdles, _ := strconv.Atoi(s)
	if maxIdles <= 0 {
		maxIdles = 50
	}
	s = os.Getenv("PG_MAX_OPEN_CONNS")
	maxOpens, _ := strconv.Atoi(s)
	if maxOpens <= 0 {
		maxOpens = 20
	}
	// set database connection
	var err error
	if db, err = sql.Open("pgx", databaseURL); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	db.SetConnMaxLifetime(500)
	db.SetMaxIdleConns(maxIdles)
	db.SetMaxOpenConns(maxOpens)
	// db.Stats()
	return db
}
