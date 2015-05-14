package utils

import (
		"net"
		//"os"
		"errors"
		"github.com/apenella/messageOutput"
)

func IsLocalIPAddress(ip string) error {

	if ip == "0.0.0.0" {return nil}

	if addrs, err := net.InterfaceAddrs(); err == nil {
		message.WriteDebug("(utils::IsLocalIPAddress) validation host IP "+ip)
		
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok {
				if ipnet.IP.String() == ip{
					return nil
				}
			}
		} 
	}else {
		return err
	}

	return errors.New("Take care, the desired IP does not belong to this server")
}