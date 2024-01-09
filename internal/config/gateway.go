package config

import (
	"net/netip"
)

type GatewayConfig interface {
	Masters() []Line
}

type gatewayConfig struct {
	Lines []Line
}

type Network int

const (
	Udp Network = iota
	Tcp
)

func (d Network) String() string {
	return [...]string{"udp", "tcp"}[d]
}

type Direction struct {
	Net    Network
	Self   netip.AddrPort
	Remote netip.AddrPort
}

var lines []Line

type Line struct {
	Description string
	Master      Direction
	Slave       []Direction
}

func (_ gatewayConfig) Masters() []Line {
	return lines
}

func NewGatewayConfig() (GatewayConfig, error) {
	return &gatewayConfig{Lines: lines}, nil

}
