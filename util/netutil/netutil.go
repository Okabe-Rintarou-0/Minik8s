package netutil

import (
	"fmt"
	"net"
	"sync"
)

var lock sync.Mutex

// GetAvailablePort returns a random available TCP port
func GetAvailablePort() (int, error) {
	// We must use lock here, because many goroutines may call this function
	lock.Lock()
	defer lock.Unlock()
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

// IsPortAvailable judges whether given port is available
func IsPortAvailable(port int) bool {
	address := fmt.Sprintf("%s:%d", "0.0.0.0", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("port %s is taken: %s\n", address, err)
		return false
	}

	defer listener.Close()
	return true
}

func GetHostIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, address := range addrs {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP.String())
			}
		}
	}
}
