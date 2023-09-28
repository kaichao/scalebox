package misc

import (
	"bufio"
	"fmt"
	"os"
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
