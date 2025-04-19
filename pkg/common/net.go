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
	if localIPIndexStr == "" {
		cmd = "hostname -i"
	} else if runtime.GOOS == "darwin" { // macos
		// 172.*.*.* are used for out-of-band Network/container network
		cmd = `ifconfig|grep "inet "|grep -v 127.0.0.1|grep -v "inet 172."|head -1|cut -d ' ' -f 2`
	} else { // linux
		cmd = fmt.Sprintf("hostname -I | awk '{print $%s}'", localIPIndexStr)
	}

	out, _ := exec.Command("/bin/bash", "-c", cmd).Output()
	localIP = strings.Split(string(out), " ")[0]
	// remove '\n'
	localIP = strings.Replace(localIP, "\n", "", -1)
	reIPv4 := regexp.MustCompile("^([0-9]+\\.){3}[0-9]+$")
	if !reIPv4.MatchString(localIP) {
		logrus.Warnf("error get_local_ip, localIP=%s\n", localIP)
	}
	return localIP
}
