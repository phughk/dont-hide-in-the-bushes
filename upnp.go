package main

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/huin/goupnp/dcps/internetgateway2"
	"github.com/huin/goupnp/soap"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type RouterClient interface {
	AddPortMapping(
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
		NewInternalPort uint16,
		NewInternalClient string,
		NewEnabled bool,
		NewPortMappingDescription string,
		NewLeaseDuration uint32,
	) (err error)

	GetExternalIPAddress() (
		NewExternalIPAddress string,
		err error,
	)
}

func PickRouterClient(ctx context.Context) (RouterClient, error) {
	tasks, _ := errgroup.WithContext(ctx)
	// Request each type of client in parallel, and return what is found.
	var ip1Clients []*internetgateway2.WANIPConnection1
	tasks.Go(func() error {
		var err error
		ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
		return err
	})
	var ip2Clients []*internetgateway2.WANIPConnection2
	tasks.Go(func() error {
		var err error
		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
		return err
	})
	var ppp1Clients []*internetgateway2.WANPPPConnection1
	tasks.Go(func() error {
		var err error
		ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
		return err
	})

	if err := tasks.Wait(); err != nil {
		return nil, err
	}

	// Trivial handling for where we find exactly one device to talk to, you
	// might want to provide more flexible handling than this if multiple
	// devices are found.
	fmt.Printf("There are [%d, %d, %d] clients\n", len(ip2Clients), len(ip1Clients), len(ppp1Clients))
	switch {
	case len(ip2Clients) == 1:
		fmt.Println("Picking router client 1")
		return ip2Clients[0], nil
	case len(ip1Clients) == 1:
		fmt.Println("Picking router client 2")
		return ip1Clients[0], nil
	case len(ppp1Clients) == 1:
		fmt.Println("Picking router client 3")
		return ppp1Clients[0], nil
	default:
		return nil, errors.New("multiple or no services found")
	}
}

const SOAP_ERR_MAPPED = 718

func GetIPAndForwardPort(client RouterClient, externalPort uint16, internalHost string, internalPort uint16, lease_duration_seconds uint32) (string, uint16, error) {
	externalIP, err := client.GetExternalIPAddress()
	if err != nil {
		return "", 0, err
	}
	fmt.Println("Our external IP address is: ", externalIP)

	err = client.AddPortMapping(
		"",
		// External port number to expose to Internet:
		externalPort,
		// Forward TCP (this could be "UDP" if we wanted that instead).
		"TCP",
		// Internal port number on the LAN to forward to.
		// Some routers might not support this being different to the external
		// port number.
		internalPort,
		// Internal address on the LAN we want to forward to.
		internalHost,
		// Enabled:
		true,
		// Informational description for the client requesting the port forwarding.
		"DontHideInTheBushes",
		// How long should the port forward last for in seconds.
		// If you want to keep it open for longer and potentially across router
		// resets, you might want to periodically request before this elapses.
		lease_duration_seconds,
	)
	if err != nil {
		if soapErr, ok := err.(*soap.SOAPFaultError); ok {
			if soapErr.Detail.UPnPError.Errorcode == SOAP_ERR_MAPPED {
				logrus.Info("Failed to renew lease, port already mapped")
			} else {
				xmlErr, serErr := xml.Marshal(soapErr)
				if serErr != nil {
					return "", 0, fmt.Errorf("unhandled soap error: %+v", soapErr)
				} else {
					return "", 0, fmt.Errorf("unhandled soap error: %+v", string(xmlErr))
				}
			}
		} else {
			return "", 0, fmt.Errorf("not a soap error: %+v", err)
		}
	}
	return externalIP, externalPort, nil
}
