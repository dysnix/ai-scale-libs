package configs

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GetRootRepositoryPath() string {
	out, err := exec.Command("sh", "-c", "git rev-parse --show-toplevel").Output()
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(string(out), "\n")
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

// HumanDuration returns a human-readable approximation of a duration
// (eg. "About a minute", "4 hours ago", etc.).
func HumanDuration(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func IsUrl(str string) bool {
	urlStr, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}

	address := net.ParseIP(urlStr.Host)

	if address == nil {
		return strings.Contains(urlStr.Host, ".")
	}

	return true
}

func ConvertDurationToStr(d time.Duration) (result string) {
	if int64(d/time.Second) <= 60 {
		result = fmt.Sprintf("%ds", int64(d/time.Second))
	} else if int64(d/time.Minute) <= 60 {
		result = fmt.Sprintf("%dm", int64(d/time.Minute))
	} else if int64(d/time.Hour) <= 60 {
		result = fmt.Sprintf("%dh", int64(d/time.Hour))
	} else if int64(d/(time.Hour)) >= 24 {
		result = fmt.Sprintf("%dd", int64(d/(time.Hour*24)))
	}

	return result
}
