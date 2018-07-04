package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type EndPoint struct {
	Ip      string
	UdpPort uint16
	TcpPort uint16
}

func (ep *EndPoint) pack() []byte {

	ip16 := binary.BigEndian.Uint32(net.ParseIP(ep.Ip).To4())
	fmt.Printf("%x\n", ip16)
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, ip16)
	fmt.Printf("%x\n", ip)

	udpPort := make([]byte, 2)
	binary.BigEndian.PutUint16(udpPort, ep.UdpPort)
	fmt.Printf("%x\n", binary.BigEndian.Uint16(udpPort))

	tcpPort := make([]byte, 2)
	binary.BigEndian.PutUint16(tcpPort, ep.TcpPort)
	fmt.Printf("%x\n", binary.BigEndian.Uint16(tcpPort))

	pack := bytes.Join([][]byte{ip, udpPort, tcpPort}, []byte{})

	return pack[:]
}

func (ep *EndPoint) UdpAddress() string {
	return fmt.Sprintf("%s:%d", ep.Ip, ep.UdpPort)
}

func (ep *EndPoint) TcpAddress() string {
	return fmt.Sprintf("%s:%d", ep.Ip, ep.TcpPort)
}

type PingNode struct {
	packetType   uint8
	version      uint8
	EndpointFrom EndPoint
	EndpointTo   EndPoint
	timestamp    uint32
}

func (pn *PingNode) pack() []byte {
	ts := make([]byte, 4)
	binary.BigEndian.PutUint32(ts, 0)

	pn.version = 1

	return bytes.Join([][]byte{
		{pn.version},
		pn.EndpointFrom.pack(),
		pn.EndpointTo.pack(),
		ts,
	}, []byte{})

}
