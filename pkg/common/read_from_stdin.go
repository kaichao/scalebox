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
		// logrus.Errorf("Failed to get stdin info:%v\n", err)
		return []string{}, errors.WrapE(err, "Failed to get stdin info")
	}
	// If stdin is a character device (terminal), no redirection or pipe was used
	if fi.Mode()&os.ModeCharDevice != 0 {
		// logrus.Warnln("No standard input detected; please provide data via pipe or redirection.")
		return []string{},
			errors.E("No standard input detected; please provide data via pipe or redirection.")
	}

	var lines []string
	// Standard input is available; read and print line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		// logrus.Warnf("Error reading standard input:%v\n", err)
		return []string{}, errors.WrapE(err, "Error reading standard input")
	}
	return lines, nil
}
