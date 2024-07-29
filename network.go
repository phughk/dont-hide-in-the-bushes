package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"time"
)

type Network struct {
	peers []*ConnHandler
}

func (n *Network) Connect(host string, port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	ch := newConnHandler(conn, n.receiveMessage)
	n.peers = append(n.peers, ch)
	return nil
}

func (n *Network) Bind(req_port int) (string, int, error) {
	addresses, err := getLocalAddresses()
	if err != nil {
		return "", 0, err
	}
	fmt.Printf("local addresses: %v\n", addresses)
	host := addresses[0]
	fmt.Printf("picked: %s\n", host)
	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, req_port))
	if err != nil {
		return "", 0, err
	}
	addr := conn.Addr().String()
	host, portStr, err := net.SplitHostPort(addr)
	fmt.Printf("Host is %s and port is %s\n", host, portStr)
	if err != nil {
		return "", 0, err
	}
	port, err := strconv.Atoi(portStr)
	// Now get the public IP via UPnP
	ctx := context.Background()
	externalPort := port
	var port_u16 uint16 = 0
	host, port_u16, err = GetIPAndForwardPort(ctx, uint16(externalPort), host, uint16(port))
	port = int(port_u16)
	if err != nil {
		return "", 0, err
	} else {
		// start a function to renew request
		go func() {
			for {
				// every 30s renew the port
				time.Sleep(time.Duration(30) * time.Second)
				_, _, err := GetIPAndForwardPort(ctx, uint16(externalPort), host, uint16(port))
				if err != nil {
					logrus.Errorf("Failed to renew upnp: %e", err)
					break
				}
			}
		}()
	}
	return host, port, nil
}

func getLocalAddresses() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error getting interfaces: %e", err)
	}

	ret_addrs := make([]string, 0)
	for _, iface := range interfaces {
		// Skip down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Get all addresses associated with this interface
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("error getting addresses: %e", err)
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Check if the IP address is a valid unicast address (IPv4 or IPv6)
			if ip == nil || ip.IsLoopback() || !ip.IsGlobalUnicast() {
				continue
			}

			ret_addrs = append(ret_addrs, ip.String())
			//fmt.Printf("Interface: %s, IP: %s\n", iface.Name, ip.String())
		}
	}
	return ret_addrs, nil
}

func (n *Network) receiveMessage(ch *ConnHandler, message *AnyMessage) {
	fmt.Println("Received message:", message)
}

func (n *Network) Close() []error {
	errors := make([]error, 0)
	for _, conn := range n.peers {
		err := conn.Close()
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}
