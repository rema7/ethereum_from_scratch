package main

import (
	"bytes"
	"crypto/ecdsa"
	"ethereum_from_scratch/rlp"
	"fmt"
	"net"

	"time"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
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
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Error(err.Error())
	}
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, remoteAddr, err := 0, new(net.UDPAddr), error(nil)
	for err == nil {
		n, remoteAddr, err = conn.ReadFromUDP(buffer)
		fmt.Println("from", remoteAddr)
		receive <- string(buffer[:n])
	}
}

//func (node *Node) wrapPacket(to EndPoint) {
//
//hash || signature || packet-type || packet-data
//hash: sha3(signature || packet-type || packet-data)
//signature: sign(privkey, sha3(packet-type || packet-data))
//packet-type: single byte < 2**7 // valid values are [1,4]
//packet-data: RLP encoded list. Packet properties are serialized in the order in which they're defined. See packet-data below.
//}

func (node *Node) SendPing(to EndPoint, key ecdsa.PrivateKey) {
	ping := PingNode{
		EndpointFrom: node.endPoint,
		EndpointTo:   to,
	}

	var b bytes.Buffer
	rlp.Encode(&b, ping.pack())
	payload := bytes.Join([][]byte{{ping.packetType}, b.Bytes()}, []byte{})

	sha3Payload := Keccak256(payload)
	signature, _ := secp256k1.Sign(sha3Payload, key.D.Bytes())

	payload = bytes.Join([][]byte{signature, payload}, []byte{})
	payloadHash := Keccak256(payload)

	result := bytes.Join([][]byte{payloadHash, payload}, []byte{})
	node.Send(result, to)
}

func (node *Node) Send(payload []byte, to EndPoint) {
	localAddr, _ := net.ResolveUDPAddr("udp", node.endPoint.UdpAddress())
	destinationAddress, _ := net.ResolveUDPAddr("udp", to.UdpAddress())
	connection, err := net.DialUDP("udp", localAddr, destinationAddress)
	if err != nil {
		log.Error(err.Error())
	}
	defer connection.Close()

	//for {
	connection.Write(payload)
	//}
}

func main() {
	ep1 := EndPoint{
		Ip:      "192.168.0.198",
		UdpPort: 30303,
		TcpPort: 30303,
	}
	//
	ep2 := EndPoint{
		Ip:      "127.0.0.1",
		UdpPort: 30303,
		TcpPort: 30303,
	}
	wallet := InitWallet()
	newNode := NewServer(ep1)
	receive, _ := newNode.Start()

	go func() {
		for {
			select {
			case msg := <-receive:
				fmt.Printf("received message %s", msg)
			}
		}
	}()

	time.Sleep(2 * time.Second)

	newNode.SendPing(ep2, wallet.PrivateKey)

	select {}
}
