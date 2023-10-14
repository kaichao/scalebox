package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	regex *regexp.Regexp
	idx   []int

	fmtDataSet = ` {
		"datasetID":"%s",
		"horizontalWidth":%d,
		"verticalStart": %d,
		"verticalHeight": %d
	}
	`
	dataset string
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "number of args for get_2d_meta should be 3!")
		os.Exit(1)
	}
	regex = regexp.MustCompile(os.Args[1])
	for _, s := range strings.Split(os.Args[2], ",") {
		if i, err := strconv.Atoi(s); err == nil {
			idx = append(idx, i)
		}
	}

	if len(idx) != 3 {
		fmt.Fprintln(os.Stderr, "number of indexes for get_2d_meta should be 3!")
		os.Exit(2)
	}
	var (
		minX = math.MaxInt32
		minY = math.MaxInt32
		maxX = -1
		maxY = -1
	)
	var n int
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		x, y := getXY(line)
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
		n++
	}
	width := maxX - minX + 1
	height := maxY - minY + 1
	if width*height != n {
		fmt.Fprintf(os.Stderr, "%s is not 2d-dataset, numX=%d,numY=%d,count=%d!\n",
			dataset, width, height, n)
		os.Exit(3)
	}

	format := regexp.MustCompile("\\s+").ReplaceAllString(fmtDataSet, "")
	fmt.Printf(format, dataset, width, minY, height)
	os.Exit(0)
}

func getXY(s string) (int, int) {
	if !regex.MatchString(s) {
		return 0, 0
	}
	ss := regex.FindStringSubmatch(s)
	dataset = ss[idx[0]]
	s1 := ss[idx[1]]
	s2 := ss[idx[2]]
	x, _ := strconv.Atoi(s1)
	y, _ := strconv.Atoi(s2)
	return x, y
}
