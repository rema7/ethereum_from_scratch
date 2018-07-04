package main

import (
	"fmt"
	"net"
)

type Node struct {
	endPoint EndPoint
}

func NewServer(endpoint EndPoint) *Node {
	newNode := &Node{
		endPoint: endpoint,
	}

	return newNode
}

func (node *Node) Start() (<-chan string, <-chan string) {
	receive := make(chan string, 10)
	send := make(chan string, 10)
	go node.listen(receive)
	return receive, send
}

func (node *Node) listen(receive chan string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", node.endPoint.UdpAddress())
	conn, _ := net.ListenUDP("udp", udpAddr)
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for err == nil {
		n, remoteAddr, err = conn.ReadFromUDP(buffer)
		fmt.Println("from", remoteAddr)
		receive <- string(buffer[:n])
	}
}

func (node *Node) Send(payload chan string, endpoint EndPoint) {
	localAddr, _ := net.ResolveUDPAddr("udp", node.endPoint.UdpAddress())
	destinationAddress, _ := net.ResolveUDPAddr("udp", endpoint.UdpAddress())
	connection, _ := net.DialUDP("udp", localAddr, destinationAddress)
	defer connection.Close()

	for {
		message := <-payload
		connection.Write([]byte(message))
	}
}

func main() {
	ep1 := EndPoint{
		Ip:      "127.0.0.1",
		UdpPort: 30303,
		TcpPort: 30303,
	}
	//
	//ep2 := EndPoint{
	//	Ip:"13.93.211.84",
	//	UdpPort: 30303,
	//	TcpPort: 30303,
	//213.208.177.250}
	//wallet := InitWallet()
	newNode := NewServer(ep1)
	receive, _ := newNode.Start()

	for {
		select {
		case msg := <-receive:
			fmt.Printf("received message %s", msg)
		}
	}
}
