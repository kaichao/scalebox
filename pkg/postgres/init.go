package postgres

import (
	"fmt"
	"os"
	"strings"

	"github.com/kaichao/scalebox/pkg/common"
	"github.com/sirupsen/logrus"
)

// getConnString ...
func getConnString() string {
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")
	if pgHost == "" {
		// in agent, set grpc server as default server
		grpcServer := os.Getenv("GRPC_SERVER")
		pgHost = strings.Split(grpcServer, ":")[0]
		logrus.Tracef("Set GRPC_SERVER %s as default db server.\n", grpcServer)
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

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		pgUser, pgPass, pgHost, pgPort, pgDB)
	logrus.Tracef("conn-string:%s\n", connString)
	return connString
}
