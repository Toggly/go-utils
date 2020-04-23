package utils

import (
	"net"
	"net/http"
	"strings"

	uuid "github.com/toggly/go.uuid"
)

//GetRealIPAddr gets real IP Addr from request's headers
func GetRealIPAddr(r *http.Request) string {

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	userIP := net.ParseIP(ip).String()

	forward := r.Header.Get("X-Forwarded-For")
	if forward != "" {
		ips := strings.Split(forward, ", ")
		// returns the first item is the list
		// it means that it the first proxies IPAddr
		userIP = ips[0]
	}
	return userIP
}

// GenerateUUID gets UUID v4 string
func GenerateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
