package core

import (
	"fmt"
	"net"
)

// GetInterfaces returns a formatted string of local network interfaces and their IPs
func GetInterfaces() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "Error retrieving interfaces"
	}

	result := "Available Network Interfaces:\n\n"
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// Only show IPv4 for simplicity, like typical necromancy
			if ip.To4() != nil {
				result += fmt.Sprintf(" - %s: %s\n", i.Name, ip.String())
			}
		}
	}
	return result
}
