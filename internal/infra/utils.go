package infra

// Este arquivo contém trechos de código licenciados sob a Licença MIT,
// originalmente criados por devfullcycle. Veja o arquivo LICENSE para detalhes.

import (
	"net"
	"net/http"
)

// getClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := net.ParseIP(xff)
		if ips != nil {
			return ips.String()
		}
	}

	// Extract from RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
