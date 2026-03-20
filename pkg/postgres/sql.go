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
	db         *sql.DB
	lastPGHost string
	dbMutex    sync.RWMutex
)

// GetDB ...
func GetDB() *sql.DB {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	currentPGHost := os.Getenv("PGHOST")
	if db != nil {
		if lastPGHost == currentPGHost {
			// 如果连接已存在且PGHOST未变化，直接返回
			return db
		}
		// PGHOST发生变化，需要创建新连接
		db.Close()
		db = nil
		logrus.Debugf("PGHOST changed from %s to %s, closing old connection", lastPGHost, currentPGHost)
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
	db, err = sql.Open("pgx", connString)
	if err != nil {
		logrus.Errorf("Unable to connect to database:%v\n", err)
		// 即使连接失败，也返回db对象（可能为nil或无效连接）
		// 不更新lastPGHost，这样下次会重试
	} else {
		// 更新记录的PGHOST值
		lastPGHost = currentPGHost
		logrus.Debugf("Created new database connection with PGHOST=%s", currentPGHost)
	}

	// 如果db不为nil，设置连接参数
	if db != nil {
		// 设置较短间隔，主要不用作连接池
		db.SetConnMaxLifetime(1 * time.Second)
		db.SetMaxIdleConns(maxIdles)
		db.SetMaxOpenConns(maxOpens)
		// db.Stats()
	}
	return db
}
