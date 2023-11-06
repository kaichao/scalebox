package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	input := os.Args[1]
	if !checkFormat(input) {
		os.Exit(1)
	}
	fmt.Println(getJumpServers(input))
}

func split(input string) []string {
	idx := []int{0}
	status := 1
	for i, c := range input {
		switch status {
		case 1:
			if c == ',' {
				idx = append(idx, i+1)
			} else if c == '[' {
				status = 2
			}
		case 2:
			if c == ']' {
				status = 1
			}
		default:
		}
	}
	idx = append(idx, len(input)+1)

	var ret []string
	for i := 0; i < len(idx)-1; i++ {
		ret = append(ret, input[idx[i]:idx[i+1]-1])
	}
	return ret
}

func checkFormat(input string) bool {
	re1 := `([a-zA-Z0-9_\-]+)@([0-9\.]+)`
	re2 := fmt.Sprintf(`(\[((%s,[0-9]+)/)+(%s,[0-9]+)\])`, re1, re1)
	re3 := fmt.Sprintf(`^(%s|%s)$`, re1, re2)
	for _, s := range split(input) {
		matched, _ := regexp.MatchString(re3, s)
		if !matched {
			fmt.Fprintf(os.Stderr,
				"'%s' is not valid jump-servers ,and '%s' not valid.\n", input, s)
			return false
		}
	}
	return true
}

func getJumpServers(input string) string {
	var servers []string
	// parts := strings.FieldsFunc(input, func(c rune) bool { return c == ',' })
	for _, s := range split(input) {
		servers = append(servers, getJumpServer(s))
	}
	return strings.Join(servers, ",")
}

func getJumpServer(input string) string {
	var (
		servers []string
		factors []float32
		sum     float32
	)
	s := regexp.MustCompile(`\[(.+)\]`).FindStringSubmatch(input)
	if s == nil {
		return input
	}
	str := s[1]
	sParts := strings.FieldsFunc(str, func(c rune) bool { return c == '/' })
	for _, s := range sParts {
		parts := strings.FieldsFunc(s, func(c rune) bool { return c == ',' })
		servers = append(servers, parts[0])
		if len(parts) > 1 {
			f, _ := strconv.ParseFloat(parts[1], 32)
			factors = append(factors, float32(f))
		} else {
			factors = append(factors, 1.0)
		}
	}

	for _, f := range factors {
		sum += f
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Float32()

	var index int
	for i, f := range factors {
		r -= f / sum
		if r <= 0 {
			index = i
			break
		}
	}
	return servers[index]
}
