package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// PublicIPInfo represents public IP information
type PublicIPInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Provider string `json:"org"`
}

// GetPublicIP gets the public IP address using multiple services
func GetPublicIP() (string, error) {
	// Try multiple services for redundancy
	services := []string{
		"https://api.ipify.org?format=text",
		"https://checkip.amazonaws.com",
		"https://ifconfig.me",
		"https://icanhazip.com",
	}
	
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	for _, service := range services {
		resp, err := client.Get(service)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			
			ip := strings.TrimSpace(string(body))
			if net.ParseIP(ip) != nil {
				return ip, nil
			}
		}
	}
	
	return "", fmt.Errorf("unable to get public IP from any service")
}

// GetPublicIPInfo gets detailed public IP information
func GetPublicIPInfo() (*PublicIPInfo, error) {
	// Try ipinfo.io first (most reliable)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	
	resp, err := client.Get("https://ipinfo.io/json")
	if err != nil {
		// Fallback to basic IP
		ip, err := GetPublicIP()
		if err != nil {
			return nil, err
		}
		return &PublicIPInfo{IP: ip}, nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ipinfo.io returned status %d", resp.StatusCode)
	}
	
	var info PublicIPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	
	return &info, nil
}

// GetLocalIP gets the local IP address
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	
	return "", fmt.Errorf("no local IP found")
}

// GetNetworkInfo gets comprehensive network information
func GetNetworkInfo() map[string]string {
	info := make(map[string]string)
	
	// Get local IP
	if localIP, err := GetLocalIP(); err == nil {
		info["local_ip"] = localIP
	} else {
		info["local_ip"] = "127.0.0.1"
	}
	
	// Get public IP
	if publicIP, err := GetPublicIP(); err == nil {
		info["public_ip"] = publicIP
	} else {
		info["public_ip"] = "Unknown"
	}
	
	// Get detailed info if available
	if ipInfo, err := GetPublicIPInfo(); err == nil && ipInfo.IP != "" {
		info["city"] = ipInfo.City
		info["region"] = ipInfo.Region
		info["country"] = ipInfo.Country
		info["provider"] = ipInfo.Provider
	}
	
	return info
}

// FormatNetworkInfo formats network information for display
func FormatNetworkInfo(info map[string]string) string {
	var result strings.Builder
	
	result.WriteString(fmt.Sprintf("[cyan]Local IP:[white] %s\n", info["local_ip"]))
	result.WriteString(fmt.Sprintf("[cyan]Public IP:[white] %s\n", info["public_ip"]))
	
	if city, ok := info["city"]; ok && city != "" {
		result.WriteString(fmt.Sprintf("[cyan]Location:[white] %s, %s, %s\n", city, info["region"], info["country"]))
	}
	
	if provider, ok := info["provider"]; ok && provider != "" {
		result.WriteString(fmt.Sprintf("[cyan]Provider:[white] %s\n", provider))
	}
	
	return result.String()
}

// IsPrivateIP checks if an IP address is private
func IsPrivateIP(ip string) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
	}
	
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	
	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(parsedIP) {
			return true
		}
	}
	
	return false
}