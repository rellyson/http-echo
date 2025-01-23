package netutils

import (
	"errors"
	"net"
	"os"
)

// HostInfo contains the host information.
type HostInfo struct {
	Hostname string // application hostname
	IP       net.IP // application IP address
}

// Retrieves the host information.
func GetHostInfo() (*HostInfo, error) {
	hostname, err := os.Hostname()

	if err != nil {
		return &HostInfo{}, err
	}

	ip, err := localIP()

	if err != nil {
		return &HostInfo{}, err
	}

	return &HostInfo{
		Hostname: hostname,
		IP:       ip,
	}, nil
}

// Gets the host machine local IP address
func localIP() (net.IP, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip, nil
			}
		}
	}

	return nil, errors.New("no IP")
}

// Checks if the IP address is private
func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet

	for _, cidr := range []string{
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}
