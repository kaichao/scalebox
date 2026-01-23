package postgres

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	// Register pgx driver with database/sql
	_ "github.com/jackc/pgx/v5/stdlib"
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

	connString := getConnString()
	s := os.Getenv("PG_MAX_IDLE_CONNS")
	maxIdles, _ := strconv.Atoi(s)
	if maxIdles <= 0 {
		maxIdles = 1
	}
	s = os.Getenv("PG_MAX_OPEN_CONNS")
	maxOpens, _ := strconv.Atoi(s)
	if maxOpens <= 0 {
		maxOpens = 4
	}
	// set database connection
	var err error
	if db, err = sql.Open("pgx", connString); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	// 设置较短间隔，主要不用作连接池
	db.SetConnMaxLifetime(1 * time.Second)
	db.SetMaxIdleConns(maxIdles)
	db.SetMaxOpenConns(maxOpens)
	// db.Stats()
	return db
}
