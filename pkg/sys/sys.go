package sys

import (
	"reflect"
	"runtime"
	"strings"
)

// GetCurrentGoroutineStack ...
func GetCurrentGoroutineStack() string {
	buf := make([]byte, 4096) // Initial small capacity
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		// Buffer too small, double its size and retry
		buf = make([]byte, len(buf)*2)
	}
}

// GetFunctionName ...
func GetFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
