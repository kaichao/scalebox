package misc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// AppendToFile ...
func AppendToFile(fileName string, line string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file %s error,err-info:%v\n", fileName, err)
		fmt.Fprintln(os.Stderr, os.Args)
		os.Exit(3)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(line + "\n")
	writer.Flush()
}

// GetTextFileLines ...
func GetTextFileLines(textFile string) ([]string, error) {
	if _, err := os.Stat(textFile); err != nil {
		_, ok := err.(*os.PathError)
		if ok && strings.Contains(err.Error(), "no such file or directory") {
			// file not exists
			return []string{}, nil
		}
		return []string{}, err
	}

	fileContent, err := ioutil.ReadFile(textFile)
	if err != nil {
		return []string{}, fmt.Errorf("Read file error, filename:%s, err:%v", textFile, err)
	}
	var lines []string
	for _, line := range strings.Split(string(fileContent), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
