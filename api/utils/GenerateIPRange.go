package utils

import (
	"encoding/binary"
	"errors"
	"net"
)

func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}

func GenerateIPRange(start, end string) ([]string, error) {
	startIP := net.ParseIP(start)
	endIP := net.ParseIP(end)
	if startIP == nil || endIP == nil {
		return nil, errors.New("invalid IP format")
	}

	startInt := ipToUint32(startIP)
	endInt := ipToUint32(endIP)

	if endInt < startInt {
		return nil, errors.New("end IP is smaller than start IP")
	}

	ips := make([]string, 0, endInt-startInt+1)
	for i := startInt; i <= endInt; i++ {
		ips = append(ips, uint32ToIP(i).String())
	}
	return ips, nil
}
