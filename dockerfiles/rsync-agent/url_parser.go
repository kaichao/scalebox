package main

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
	"strings"
)

/*
 * url parser
 *   LOCAL : local directory
 *   SSH   : user@host:port/root-dir
 *   RSYNC : rsync://user@host:port/root-dir
 *   FTP   : ftp://user@host:port/root-dir
 *
 *  LOCAL_ROOT / REMOTE_ROOT should have "/local" prefix
 *		for mapping to paths outside the container
 */

func main() {
	url := os.Args[1]
	serverContainerized := "yes" == os.Getenv("SERVER_CONTAINERIZED")

	if strings.HasPrefix(url, "/") {
		// MODE | LOCAL_ROOT
		url = "/local" + url
		fmt.Printf("%s %s", "LOCAL", url)
		return
	}

	if strings.HasPrefix(url, "ftp://") {
		// MODE | FTP_URL | REMOTE_ROOT
		reg := regexp.MustCompile("(ftp://([^:]+:[^@]+@)?[^/:]+(:[^/]+)?)(/.*)")
		ss := reg.FindStringSubmatch(url)
		if ss == nil {
			fmt.Printf("ERROR_FORMAT, url:%s", url)
		} else {
			ftpURL := ss[1]
			remoteRoot := ss[4]
			fmt.Printf("%s %s %s", "FTP", ftpURL, remoteRoot)
		}
		return
	}

	var mode string
	if strings.HasPrefix(url, "rsync://") {
		mode = "RSYNC"
		url = url[8:]
	} else {
		mode = "SSH"
	}
	reg := regexp.MustCompile("^(([^@]+)@)?([^:/]+)(:([0-9]+))?(/.*)$")
	ss := reg.FindStringSubmatch(url)
	if ss == nil {
		fmt.Printf("ERROR_FORMAT, url:%s", url)
	}
	uname := ss[2]
	host := ss[3]
	port := ss[5]
	path := ss[6]
	if uname == "" {
		if mode == "SSH" {
			if !serverContainerized {
				if u, err := user.Current(); err == nil {
					uname = u.Username
				}
			}
			if uname == "" {
				uname = "root"
			}
		} else { // "RSYNC"
			if serverContainerized {
				// rsyncd image's default user name
				uname = "user"
			}
		}
	}

	if port == "" {
		if mode == "SSH" {
			if serverContainerized {
				port = "2222"
			} else {
				port = "22"
			}
		} else {
			port = "873"
		}
	}
	if host == "" || path == "" {
		fmt.Fprintf(os.Stderr, "REMOTE_HOST or REMOTE_ROOT is null, url=#%s#", url)
	}
	if serverContainerized {
		path = "/local" + path
	}
	// for anonymous rsync && not containerized, REMOTE_USER is null
	// MODE | REMOTE_HOST | REMOTE_PORT | REMOTE_ROOT | REMOTE_USER
	fmt.Printf("%s %s %s %s %s", mode, host, port, path, uname)
}
