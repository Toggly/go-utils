package utils

import (
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

//GetRealIPAddr gets real IP Addr from request's headers
func GetRealIPAddr(r *http.Request) string {
	forward := r.Header.Get("X-Forwarded-For")
	if forward != "" {
		ips := strings.Split(forward, ", ")
		// returns the first item is the list
		// it means that it the first proxies IPAddr
		return ips[0]
	}
	return r.RemoteAddr
}

// GenerateUUID gets UUID v4 string
func GenerateUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
