package postgres

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kaichao/scalebox/pkg/common"
)

// GetPgxConnString ...
func GetPgxConnString() string {
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	if pgHost == "" {
		// in agent, set grpc server as default server
		grpcServer := os.Getenv("GRPC_SERVER")
		pgHost = strings.Split(grpcServer, ":")[0]
		fmt.Printf("[INFO] %s Set GRPC_SERVER %s as default db server.\n",
			time.Now().Format("15:04:05.000"), grpcServer)
	}
	if pgHost == "" {
		pgHost = os.Getenv("LOCAL_ADDR")
	}
	if pgHost == "" {
		localIP := common.GetLocalIP()
		pgHost = localIP
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

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgDB)

}
