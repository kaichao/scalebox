package misc

import (
	"os"
	"strings"
	"time"

	"github.com/kaichao/gopkg/common"
)

// AddTimeStamp ...
//
//	add-timestamp in agent
func AddTimeStamp(label string) {
	fileName := os.Getenv("WORK_DIR") + "/timestamps.txt"
	timeStamp := time.Now().Format("2006-01-02T15:04:05.000000Z07:00")
	// fmt.Printf("timestamp:%s\n", timeStamp)
	common.AppendToFile(fileName, timeStamp+","+label)
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
