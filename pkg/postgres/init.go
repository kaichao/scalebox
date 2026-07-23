package postgres

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/sirupsen/logrus"
)

// getConnString 按优先级构建数据库连接串：
//
//	1. DATABASE_URL — 完整连接串（业界标准）
//	2. PGURL         — 同上（存量兼容）
//	3. 证书认证       — $PG_CERT_DIR/client.crt + client.key 存在时自动启用
//	4. 密码认证       — PGPASS 环境变量
func getConnString() string {
	// 1. DATABASE_URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		logrus.Debugf("Using DATABASE_URL")
		return dbURL
	}

	// 2. PGURL（存量兼容）
	if pgURL := os.Getenv("PGURL"); pgURL != "" {
		logrus.Debugf("Using PGURL")
		return pgURL
	}

	// 解析连接参数
	pgHost := resolveHost()
	pgPort := resolvePort(pgHost)
	pgUser := resolveUser()
	pgDB := resolveDB()

	// 3. 证书认证
	certDir := os.Getenv("PG_CERT_DIR")
	if certDir == "" {
		certDir = "./certs"
	}
	certFile := filepath.Join(certDir, "client.crt")
	keyFile := filepath.Join(certDir, "client.key")
	caFile := filepath.Join(certDir, "ca.crt")

	if fileExists(certFile) && fileExists(keyFile) {
		sslMode := "verify-full"
		if !fileExists(caFile) {
			logrus.Warnf("CA cert not found at %s, falling back to sslmode=require", caFile)
			sslMode = "require"
		}
		connString := fmt.Sprintf(
			"postgres://%s@%s:%s/%s?sslmode=%s&sslcert=%s&sslkey=%s&sslrootcert=%s",
			pgUser, pgHost, pgPort, pgDB, sslMode, certFile, keyFile, caFile,
		)
		logrus.Debugf("Cert auth: sslmode=%s cert=%s", sslMode, certFile)
		return connString
	}

	// 4. 密码认证
	pgPass := os.Getenv("PGPASS")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		pgUser, pgPass, pgHost, pgPort, pgDB)
	if pgPass == "" {
		logrus.Warnf("PGPASS not set, connecting without password")
	}
	logrus.Debugf("Password auth: host=%s user=%s db=%s", pgHost, pgUser, pgDB)
	return connString
}

// ── 辅助函数 ──────────────────────────────────────────────────

func resolveHost() string {
	pgHost := os.Getenv("PGHOST")
	if pgHost != "" {
		return pgHost
	}

	// agent 场景：用 gRPC server 地址推断
	if grpcServer := os.Getenv("GRPC_SERVER"); grpcServer != "" {
		host := strings.Split(grpcServer, ":")[0]
		logrus.Debugf("PGHOST fallback from GRPC_SERVER: %s", host)
		return host
	}

	if localAddr := os.Getenv("LOCAL_ADDR"); localAddr != "" {
		return localAddr
	}

	localIP := common.GetLocalIP()
	logrus.Debugf("PGHOST fallback to local IP: %s", localIP)
	return localIP
}

func resolvePort(pgHost string) string {
	// 支持 PGHOST=host:port 格式
	if ss := strings.Split(pgHost, ":"); len(ss) == 2 {
		return ss[1]
	}
	if pgPort := os.Getenv("PGPORT"); pgPort != "" {
		return pgPort
	}
	return "5432"
}

func resolveUser() string {
	if u := os.Getenv("PGUSER"); u != "" {
		return u
	}
	return "scalebox"
}

func resolveDB() string {
	if d := os.Getenv("PGDB"); d != "" {
		return d
	}
	return "scalebox"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
