package postgres

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var pool *pgxpool.Pool

func initPool() {
	// 数据库连接字符串
	connString := getConnString() + "?sslmode=disable"

	// 创建连接池配置
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logrus.Fatalf("无法解析连接字符串: %v", err)
	}

	// 配置最大和最小连接数
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

	// 配置连接存活时间和空闲时间
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

	// 设置连接池参数
	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnLifetime = time.Duration(maxConnLifetime) * time.Minute
	config.MaxConnIdleTime = time.Duration(maxConnIdleTime) * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	// 创建连接池
	ctx := context.Background()
	pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logrus.Fatalf("无法创建连接池: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		logrus.Fatalf("无法连接到数据库: %v", err)
	}
}

// GetPgxPool ...
func GetPgxPool() *pgxpool.Pool {
	if pool == nil {
		initPool()
	}
	return pool
}
