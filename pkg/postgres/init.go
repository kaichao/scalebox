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
//	3. 证书认证       — $SCALEBOX_CERTS_DIR/client.crt + client.key 存在时自动启用
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
	pgHost, pgPort := resolveHostPort()
	pgUser := resolveUser()
	pgDB := resolveDB()

	// 3. 证书认证（仅当 SCALEBOX_CERTS_DIR 显式设置时启用）
	certDir := os.Getenv("SCALEBOX_CERTS_DIR")
	if certDir != "" {
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
	}

	// 4. 密码认证
	pgPass := os.Getenv("PGPASS")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgDB)
	if pgPass == "" {
		logrus.Warnf("PGPASS not set, connecting without password")
	}
	logrus.Debugf("Password auth: host=%s user=%s db=%s", pgHost, pgUser, pgDB)
	return connString
}

// ── 辅助函数 ──────────────────────────────────────────────────

func resolveHostPort() (string, string) {
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")

	// PGHOST 可能包含端口：host:port
	if pgHost != "" {
		if ss := strings.SplitN(pgHost, ":", 2); len(ss) == 2 {
			pgHost = ss[0]
			if pgPort == "" {
				pgPort = ss[1]
			}
		}
	} else {
		// agent 场景：用 gRPC server 地址推断
		if grpcServer := os.Getenv("GRPC_SERVER"); grpcServer != "" {
			pgHost = strings.Split(grpcServer, ":")[0]
			logrus.Debugf("PGHOST fallback from GRPC_SERVER: %s", pgHost)
		} else if localAddr := os.Getenv("LOCAL_ADDR"); localAddr != "" {
			pgHost = localAddr
		} else {
			pgHost = common.GetLocalIP()
			logrus.Debugf("PGHOST fallback to local IP: %s", pgHost)
		}
	}

	if pgPort == "" {
		pgPort = "5432"
	}
	return pgHost, pgPort
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
