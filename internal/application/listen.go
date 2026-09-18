package application

import (
	"fmt"
	"net"
)

func listenLoopback(host string, startPort, attempts int) (net.Listener, int, error) {
	var last error
	for i := 0; i < attempts; i++ {
		port := startPort + i
		ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			last = err
			continue
		}
		return ln, port, nil
	}
	endPort := startPort + attempts - 1
	if last == nil {
		return nil, 0, fmt.Errorf("unable to bind %s:%d-%d", host, startPort, endPort)
	}
	return nil, 0, fmt.Errorf("unable to bind %s:%d-%d after %d attempts: %w", host, startPort, endPort, attempts, last)
}
