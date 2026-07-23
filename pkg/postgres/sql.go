package postgres

import (
	"database/sql"
	"os"
	"strconv"
	"sync"
	"time"

	// Register pgx driver with database/sql
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
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
	db          *sql.DB
	lastConnStr string
	dbMutex     sync.Mutex
)

// GetDB 返回 *sql.DB。
// 当连接参数变化（DATABASE_URL/PGURL/PGHOST/PGPASS/PG_CERT_DIR 等）时自动重建连接。
func GetDB() *sql.DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	connString := getConnString()

	if db != nil {
		if lastConnStr == connString {
			return db
		}
		logrus.Debugf("Connection config changed, closing old connection")
		db.Close()
		db = nil
	}

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

	var err error
	db, err = sql.Open("pgx", connString)
	if err != nil {
		logrus.Errorf("Unable to connect to database: %v", err)
		return db
	}

	lastConnStr = connString

	db.SetConnMaxLifetime(1 * time.Second)
	db.SetMaxIdleConns(maxIdles)
	db.SetMaxOpenConns(maxOpens)

	logrus.Debugf("Database connection established")
	return db
}
