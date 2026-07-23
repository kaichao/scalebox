package postgres

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	pool         *pgxpool.Pool
	lastPoolStr  string
	poolMutex    sync.Mutex
)

// GetPgxPool 返回 *pgxpool.Pool。
// 当连接参数变化时自动重建连接池。
func GetPgxPool() *pgxpool.Pool {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	connString := getConnString()

	if pool != nil {
		if lastPoolStr == connString {
			return pool
		}
		logrus.Debugf("Connection config changed, closing old pool")
		pool.Close()
		pool = nil
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logrus.Errorf("Failed to parse connection string: %v", err)
		return nil
	}

	// 连接池大小
	s := os.Getenv("PG_MAX_CONNS")
	maxConns, err := strconv.Atoi(s)
	if err != nil || maxConns <= 0 {
		maxConns = 20
	}
	s = os.Getenv("PG_MIN_CONNS")
	minConns, err := strconv.Atoi(s)
	if err != nil || minConns <= 0 {
		minConns = 5
	}

	// 连接生命周期
	s = os.Getenv("PG_MAX_CONN_LIFETIME_MIN")
	maxConnLifetime, err := strconv.Atoi(s)
	if err != nil || maxConnLifetime <= 0 {
		maxConnLifetime = 30
	}
	s = os.Getenv("PG_MAX_CONN_IDLE_TIME_MIN")
	maxConnIdleTime, err := strconv.Atoi(s)
	if err != nil || maxConnIdleTime <= 0 {
		maxConnIdleTime = 5
	}

	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnLifetime = time.Duration(maxConnLifetime) * time.Minute
	config.MaxConnIdleTime = time.Duration(maxConnIdleTime) * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	ctx := context.Background()
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logrus.Errorf("Failed to create connection pool: %v", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		logrus.Errorf("Unable to ping database: %v", err)
	}

	lastPoolStr = connString
	logrus.Debugf("Connection pool created")
	return pool
}
