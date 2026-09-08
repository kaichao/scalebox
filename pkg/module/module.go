// Package module provides shared entry-point logic for scalebox module binaries:
// CLI arg parsing, JSON header deserialization, action dispatch, header helpers,
// and structured error exit codes. A scalebox module main() is typically just
// os.Exit(module.Run(handlers, "action")).
package module

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/kaichao/gopkg/errors"
	"github.com/kaichao/gopkg/logger"
	"github.com/sirupsen/logrus"
)

// Exit codes returned by Run. errors.TracedError 的业务错误码直接透传。
const (
	ExitOK         = 0
	ExitBadArgs    = 1 // <task-body> <task-headers> 参数数量不足，或 handler 返回未包装的普通 error
	ExitBadHeaders = 2 // task headers JSON 解析失败
	ExitNoHandler  = 4 // 无匹配的 handler
)

// HandlerFunc is the unified handler signature for scalebox module actions.
// body is the task body (os.Args[1]), headers are parsed from the JSON task headers.
// Returns nil on success, or a *errors.TracedError with an exit code on failure.
type HandlerFunc func(body string, headers map[string]string) error

// LogEntry is the package-level logrus entry used by logger.LogError.
var LogEntry *logrus.Entry

// InitLogging configures logrus from LOG_LEVEL env var.
func InitLogging() {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{DisableQuote: true}
	logrus.SetFormatter(formatter)

	log := logrus.New()
	log.SetLevel(level)
	log.SetFormatter(formatter)
	LogEntry = logrus.NewEntry(log)
}

// Run is the shared entry point for scalebox modules.
// It parses <task-body> <task-headers-json> from os.Args,
// dispatches to the handler matching the route key in headers,
// and returns the exit code: ExitOK on success, TracedError.Code on failure.
//
// routeKey selects which header field to use for dispatch.
// Default is "from_module"; pass "action" for tape-mgr / tape-io style modules.
func Run(handlers map[string]HandlerFunc, routeKey ...string) int {
	InitLogging()

	key := "from_module"
	if len(routeKey) > 0 && routeKey[0] != "" {
		key = routeKey[0]
	}

	if len(os.Args) < 3 {
		logrus.Errorf("usage: %s <task-body> <task-headers>\nparameters expect=2, actual=%d\n",
			os.Args[0], len(os.Args)-1)
		return ExitBadArgs
	}

	headers := make(map[string]string)
	if err := json.Unmarshal([]byte(os.Args[2]), &headers); err != nil {
		logrus.Errorf("parse headers: %v\n", err)
		return ExitBadHeaders
	}

	handler := handlers[headers[key]]
	if handler == nil {
		logrus.Warnf("handler not found: %s=%s, task-body=%s\n",
			key, headers[key], os.Args[1])
		return ExitNoHandler
	}

	err := handler(os.Args[1], headers)
	if err == nil {
		return ExitOK
	}

	logger.LogError(err, LogEntry)
	if te, ok := err.(*errors.TracedError); ok {
		return te.Code
	}
	// handler 返回了未包装 TracedError 的普通 error，无业务码可用
	return ExitBadArgs
}

// GetString extracts a string header value, returning fallback if missing or empty.
func GetString(headers map[string]string, key, fallback string) string {
	if v, ok := headers[key]; ok && v != "" {
		return v
	}
	return fallback
}

// GetInt extracts an int header value, returning fallback if missing or unparseable.
func GetInt(headers map[string]string, key string, fallback int) int {
	if v, ok := headers[key]; ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// GetStringSlice unmarshals a JSON-encoded string array from headers.
func GetStringSlice(headers map[string]string, key string) []string {
	v := headers[key]
	if v == "" {
		return nil
	}
	var arr []string
	if err := json.Unmarshal([]byte(v), &arr); err != nil {
		return nil
	}
	return arr
}

// GetStringMap unmarshals a JSON-encoded string→string map from headers.
func GetStringMap(headers map[string]string, key string) map[string]string {
	v := headers[key]
	if v == "" {
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(v), &m); err != nil {
		return nil
	}
	return m
}

// PrintJSON JSON-marshals v and prints it to stdout.
func PrintJSON(v any) {
	b, _ := json.Marshal(v)
	fmt.Println(string(b))
}
