package core

import (
	"fmt"
	"net"
	"strings"
)

// GetListeningPorts returns the current listening ports from configuration
func GetListeningPorts() []string {
	if GlobalConfig == nil || GlobalConfig.Ports == "" {
		return []string{"4444"} // Default port
	}
	
	// Split ports by comma and trim whitespace
	ports := strings.Split(GlobalConfig.Ports, ",")
	for i := range ports {
		ports[i] = strings.TrimSpace(ports[i])
	}
	
	return ports
}

// GetListeningAddresses returns the full listening addresses
func GetListeningAddresses() []string {
	if GlobalConfig == nil {
		return []string{"0.0.0.0:4444"}
	}
	
	ports := GetListeningPorts()
	interfaceAddr := GlobalConfig.Interface
	if interfaceAddr == "" {
		interfaceAddr = "0.0.0.0"
	}
	
	var addresses []string
	for _, port := range ports {
		addresses = append(addresses, fmt.Sprintf("%s:%s", interfaceAddr, port))
	}
	
	return addresses
}

// IsPortListening checks if a port is currently listening
func IsPortListening(port string) bool {
	if GlobalConfig == nil {
		return false
	}
	
	interfaceAddr := GlobalConfig.Interface
	if interfaceAddr == "" {
		interfaceAddr = "0.0.0.0"
	}
	
	address := fmt.Sprintf("%s:%s", interfaceAddr, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false // Port is already in use
	}
	listener.Close()
	return true // Port is available
}

// GetCurrentIP gets the current IP address for payload generation
func GetCurrentIP() string {
	if GlobalConfig == nil || GlobalConfig.Interface == "" || GlobalConfig.Interface == "0.0.0.0" {
		// Try to get the primary network interface IP
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			return "YOUR_IP"
		}
		defer conn.Close()
		
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}
	
	// Use the configured interface
	if GlobalConfig.Interface == "127.0.0.1" || GlobalConfig.Interface == "localhost" {
		return "YOUR_IP"
	}
	
	return GlobalConfig.Interface
}