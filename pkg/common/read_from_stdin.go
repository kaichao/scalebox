package common

import (
	"bufio"
	"os"

	"github.com/kaichao/gopkg/errors"
)

// ReadLinesFromStdin ...
func ReadLinesFromStdin() ([]string, error) {
	// Read lines from stdin
	// Check if stdin is coming from a pipe or file
	fi, err := os.Stdin.Stat()
	if err != nil {
		return []string{}, errors.WrapE(err, "failed to get stdin info")
	}
	// If stdin is a character device (terminal), no redirection or pipe was used
	if fi.Mode()&os.ModeCharDevice != 0 {
		return []string{},
			errors.E("no standard input detected")
	}

	var lines []string
	// Standard input is available; read and print line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return []string{}, errors.WrapE(err, "failed to read standard input")
	}
	return lines, nil
}
