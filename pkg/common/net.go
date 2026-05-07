package common

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// GetLocalIP ...
// hostname -i / -I
func GetLocalIP() string {
	localIP := os.Getenv("LOCAL_IP")
	if localIP != "" {
		return localIP
	}
	// 0=1, 2, 3, .. ,n
	localIPIndexStr := os.Getenv("LOCAL_IP_INDEX")
	var cmd string

	if runtime.GOOS == "darwin" { // macos
		// macOS: use ifconfig to get IP address
		// 172.*.*.* are used for out-of-band Network/container network
		cmd = `ifconfig | grep "inet " | grep -v 127.0.0.1 | grep -v "inet 172\." | head -1 | awk '{print $2}'`
	} else if localIPIndexStr == "" { // linux, default case
		cmd = "hostname -i"
	} else { // linux with specific index
		cmd = fmt.Sprintf("hostname -I | awk '{print $%s}'", localIPIndexStr)
	}

	out, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		logrus.Warnf("error executing command to get local IP: %v, cmd=%s", err, cmd)
		// Fallback to simple hostname -i
		if out, err = exec.Command("hostname", "-i").Output(); err != nil {
			logrus.Warnf("error with hostname -i fallback: %v", err)
			return "127.0.0.1"
		}
	}

	localIP = strings.TrimSpace(string(out))
	// Take only the first IP if multiple are returned
	localIP = strings.Split(localIP, " ")[0]
	localIP = strings.Split(localIP, "\n")[0]

	reIPv4 := regexp.MustCompile(`^([0-9]+\.){3}[0-9]+$`)
	if !reIPv4.MatchString(localIP) {
		logrus.Warnf("invalid IP address format, localIP=%s", localIP)
		// For Linux default case, fallback to hostname -I
		if runtime.GOOS != "darwin" && localIPIndexStr == "" {
			if out, err = exec.Command("/bin/bash", "-c", "hostname -I | awk '{print $1}'").Output(); err == nil {
				fallbackIP := strings.TrimSpace(string(out))
				fallbackIP = strings.Split(fallbackIP, " ")[0]
				fallbackIP = strings.Split(fallbackIP, "\n")[0]
				if reIPv4.MatchString(fallbackIP) {
					return fallbackIP
				}
			}
		}
		// Return loopback as last resort
		return "127.0.0.1"
	}
	return localIP
}
