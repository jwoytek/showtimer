package main

import (
	"errors"
	"net"
	"time"
)

func findMyIPs() ([]string, error) {
	ifaces, err := net.Interfaces()
	ips := make([]string, 0, 10)
	if err != nil {
		return ips, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		//if iface.Flags&net.FlagLoopback != 0 {
		//	continue // loopback interface
		//}
		addrs, err := iface.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil { //|| ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ips = append(ips, ip.String())
			//return ip.String(), nil
		}
	}
	if len(ips) == 0 {
		return ips, errors.New("are you connected to the network?")
	}
	return ips, nil
}

func delayUntilNextSecond() {
	now := time.Now()
	// Compute the next whole second.
	// If we are already at the top of the second, then time.Until returns 0.
	next := now.Truncate(time.Second).Add(time.Second)
	time.Sleep(time.Until(next))
}
