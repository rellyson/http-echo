package netutils

import (
	"net"
	"testing"
)

func TestGetHostInfo(t *testing.T) {
	hostInfo, err := GetHostInfo()

	if err != nil {
		t.Errorf("GetHostInfo() error = %v", err)
		return
	}

	if hostInfo.Hostname == "" {
		t.Error("GetHostInfo() returned empty hostname")
	}

	if hostInfo.IP == nil {
		t.Error("GetHostInfo() returned nil IP")
	}
}

func TestLocalIP(t *testing.T) {
	ip, err := localIP()

	if err != nil {
		t.Errorf("localIP() error = %v", err)
		return
	}

	if ip == nil {
		t.Error("localIP() returned nil IP")
	}

	if !isPrivateIP(ip) {
		t.Error("localIP() returned non-private IP")
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"loopback v4", "127.0.0.1", true},
		{"private class A", "10.0.0.1", true},
		{"private class B", "172.16.0.1", true},
		{"private class C", "192.168.0.1", true},
		{"public IP", "8.8.8.8", false},
		{"IPv6 loopback", "::1", true},
		{"IPv6 link-local", "fe80::1", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			if got := isPrivateIP(ip); got != tt.expected {
				t.Errorf("isPrivateIP() = %v, want %v", got, tt.expected)
			}
		})
	}
}
