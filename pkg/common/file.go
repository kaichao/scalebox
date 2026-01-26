package common

import (
	"bufio"
	"os"
	"strings"

	"github.com/kaichao/gopkg/errors"
)

// AppendToFile ...
func AppendToFile(fileName string, line string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.WrapE(err, "open file failed", "filename", fileName)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(line + "\n")
	writer.Flush()
	return nil
}

// GetTextFileLines ...
func GetTextFileLines(textFile string) ([]string, error) {
	if _, err := os.Stat(textFile); err != nil {
		if _, ok := err.(*os.PathError); ok {
			// file not exists
			return []string{}, errors.WrapE(err, "file not found", "filename", textFile)
		}
		return []string{}, errors.WrapE(err, "file open failed", "filename", textFile)
	}

	fileContent, err := os.ReadFile(textFile)
	if err != nil {
		return []string{}, errors.WrapE(err, "file read failed", "filename", textFile)
	}
	var lines []string
	for _, line := range strings.Split(string(fileContent), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// IsRunnable ...
func IsRunnable(runFile string) bool {
	stat, err := os.Stat(runFile)
	if err != nil {
		return false
	}
	return stat.Mode()&0111 != 0
}
