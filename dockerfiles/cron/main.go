package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type cronItem struct {
	cronText string
	modName  string
}

func main() {
	var cronItems []cronItem

	file, err := os.Open("/cron.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open file error:%v", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if len(s) > 0 && s[0] != '#' { // not comment
			ss := strings.Split(s, ",")
			if len(ss) != 2 {
				fmt.Fprintf(os.Stderr, "format error: %s\n", s)
				continue
			}
			cronItems = append(cronItems, cronItem{cronText: ss[0], modName: ss[1]})
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scan file error:%v", err)
		os.Exit(2)
	}

	c := cron.New()

	for _, item := range cronItems {
		// CAUTION: local variable needed
		modName := item.modName
		c.AddFunc(item.cronText, func() {
			current := time.Now().Format("2006-01-02T15:04:05")
			cmd := exec.Command("send-job-message", modName, current)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "send-job-message error:%v", err)
			}
			fmt.Println(string(output))
		})
	}

	c.Start()
	// block main goroutine
	ch := make(chan struct{})
	<-ch
}
