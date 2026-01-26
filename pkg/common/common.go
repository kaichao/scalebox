package common

import (
	"os"
	"strings"
	"time"
)

// IsRegexString 判断字符串是否包含正则表达式元字符
func IsRegexString(s string) bool {
	// metachars := `.*+?^$[]{}()|\` // List of regex metacharacters
	metachars := `*+?^$[]{}()|\` // List of regex metacharacters
	escaped := false             // Flag to track if the previous character was a backslash

	for _, r := range s {
		if escaped {
			// If the previous character was a backslash, this character is escaped
			escaped = false // Reset the flag and skip this character
			continue
		}
		if r == '\\' {
			// Found a backslash, mark the next character as escaped
			escaped = true
		} else if strings.ContainsRune(metachars, r) {
			// Found an unescaped metacharacter
			return true
		}
	}
	return false // No unescaped metacharacters found
}

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
