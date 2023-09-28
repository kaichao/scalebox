package misc

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

// GetLocalIP ...
// hostname -i / -I
func GetLocalIP() string {
	// 0=1, 2, 3, .. ,n
	localIPIndexStr := os.Getenv("LOCAL_IP_INDEX")
	var cmd string
	if localIPIndexStr == "" {
		cmd = "hostname -i"
	} else {
		cmd = fmt.Sprintf("hostname -I | awk '{print $%s}'", localIPIndexStr)
	}

	out, _ := exec.Command("/bin/bash", "-c", cmd).Output()
	localIP := strings.Split(string(out), " ")[0]
	// remove '\n'
	localIP = strings.Replace(localIP, "\n", "", -1)
	if localIP == "" {
		log.Printf("get_local_ip error\n")
	}
	return localIP
}

// GetLocalIPv4Addr ...
// @Deprecated
func getLocalIPv4Addr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
	}
	var ip string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return ip
}

// GetLocalIPv4AddrByInterface ...
// @Deprecated
func getLocalIPv4AddrByInterface(interfaceName string) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	var ip string
	for _, interf := range interfaces {
		if interf.Name == "en7" {
			addrs, err := interf.Addrs()
			if err != nil {
				panic(err)
			}
			for _, add := range addrs {
				if ipnet, ok := add.(*net.IPNet); ok {
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()
					}
				}
			}
		}
	}
	return ip
}
