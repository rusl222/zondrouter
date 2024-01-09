package gateway

import (
	"log"
	"net"

	"github.com/rusl222/zondrouter/internal/config"
)

type Gateway struct {
	Client  config.Direction
	Servers []config.Direction
}

var client *connection
var servers []*connection

type connection struct {
	readBuffer []byte
	connUdp    *net.UDPConn
	connTcp    *net.TCPConn
}

func (s Gateway) Run() error {
	for _, dir := range append(s.Servers, s.Client) {
		if !dir.Self.IsValid() {
			return &net.AddrError{Err: "Address not valid", Addr: dir.Self.String()}
		}
		if !dir.Remote.IsValid() {
			return &net.AddrError{Err: "Address not valid", Addr: dir.Remote.String()}
		}
	}

	con1, err := s.connect(s.Client)
	if err == nil {
		client = con1
	}

	for _, dir := range s.Servers {
		con1, err := s.connect(dir)
		if err == nil {
			servers = append(servers, con1)
		}
	}

	for _, con1 := range servers {
		go s.transport(con1, []*connection{client})
	}

	s.transport(client, servers)

	return nil
}

func (s Gateway) transport(src *connection, dst []*connection) {
	for {
		var n int
		if src.connUdp != nil {
			n, _, _ = src.connUdp.ReadFromUDP(src.readBuffer)
		}
		if src.connTcp != nil {
			n, _ = src.connTcp.Read(src.readBuffer)
		}
		for _, c := range dst {
			if c.connUdp != nil {
				c.connUdp.Write(src.readBuffer[:n])
			}
			if c.connTcp != nil {
				c.connTcp.Write(src.readBuffer[:n])
			}
		}
	}
}

func (s Gateway) connect(dir config.Direction) (*connection, error) {
	var err error
	switch dir.Net {
	case config.Udp:
		log.Printf("Подключение %s - %s\n", dir.Self, dir.Remote)
		conn1, err := net.DialUDP("udp", net.UDPAddrFromAddrPort(dir.Self), net.UDPAddrFromAddrPort(dir.Remote))
		if err != nil {
			log.Fatalf("Не удалось подключится! %s - %s", dir.Self, dir.Remote)
		} else {
			return &connection{
				readBuffer: make([]byte, 300),
				connUdp:    conn1}, nil
		}
	case config.Tcp:
		log.Printf("Подключение %s - %s\n", dir.Self, dir.Remote)
		conn2, err := net.DialTCP("tcp", net.TCPAddrFromAddrPort(dir.Self), net.TCPAddrFromAddrPort(dir.Remote))
		if err != nil {
			log.Fatalf("Не удалось подключится! %s - %s", dir.Self, dir.Remote)
		} else {
			return &connection{
				readBuffer: make([]byte, 300),
				connTcp:    conn2}, nil
		}
	}
	return nil, err
}
