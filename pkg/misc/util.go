package misc

import (
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// AddTimeStamp ...
//
//	add-timestamp in agent
func AddTimeStamp(label string) {
	fileName := os.Getenv("WORK_DIR") + "/timestamps.txt"
	timeStamp := time.Now().Format("2006-01-02T15:04:05.000000Z07:00")
	// fmt.Printf("timestamp:%s\n", timeStamp)
	AppendToFile(fileName, timeStamp+","+label)
}

// SplitCommaWithEscapeSupport ..
func SplitCommaWithEscapeSupport(s string) []string {
	var ret []string
	i0 := 0
	i := 0
	for ; i < len(s); i++ {
		if s[i] == ',' && i > 0 && s[i-1] != '\\' {
			ret = append(ret, strings.ReplaceAll(s[i0:i], "\\,", ","))
			i0 = i + 1
		}
	}
	if i0 < i {
		ret = append(ret, strings.ReplaceAll(s[i0:i], "\\,", ","))
	}
	return ret
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

// IsRunnable ...
func IsRunnable(runFile string) bool {
	stat, err := os.Stat(runFile)
	if err != nil {
		return false
	}
	return stat.Mode()&0111 != 0
}
