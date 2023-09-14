//go:build linux

package common

import (
	"net"
	"os/exec"
)

func SetupIface(ip net.IPNet, devName string) error {
	if err := exec.Command("ip", "link", "set", devName, "up").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "addr", "add", ip.String(), "dev", devName).Run(); err != nil {
		return err
	}
	return nil
}

func TearDownIface(devName string) error {
	if err := exec.Command("ip", "link", "set", devName, "down").Run(); err != nil {
		return err
	}
	if err := exec.Command("ip", "link", "delete", devName).Run(); err != nil {
		return err
	}
	return nil
}

func PingIface(dstIP net.IP) error {
	if err := exec.Command("ping", "-c", "4", dstIP.String()).Run(); err != nil {
		return err
	}
	return nil
}

func IsBroadcastAddr(addr net.HardwareAddr) bool {
	yes := true
	for i := 0; i < 6; i++ {
		if addr[i] != 0xFF {
			yes = false
			break
		}
	}
	return yes
}

func IsIPv4MulticastAddr(addr net.HardwareAddr) bool {
	return addr[0] == 0x01 && addr[1] == 0x00 && addr[2] == 0x5E
}

type EtherType [2]byte

// http://en.wikipedia.org/wiki/Ethertype
var (
	IPv4                = EtherType{0x08, 0x00}
	ARP                 = EtherType{0x08, 0x06}
	WakeOnLAN           = EtherType{0x08, 0x42}
	TRILL               = EtherType{0x22, 0xF3}
	DECnetPhase4        = EtherType{0x60, 0x03}
	RARP                = EtherType{0x80, 0x35}
	AppleTalk           = EtherType{0x80, 0x9B}
	AARP                = EtherType{0x80, 0xF3}
	IPX1                = EtherType{0x81, 0x37}
	IPX2                = EtherType{0x81, 0x38}
	QNXQnet             = EtherType{0x82, 0x04}
	IPv6                = EtherType{0x86, 0xDD}
	EthernetFlowControl = EtherType{0x88, 0x08}
	IEEE802_3           = EtherType{0x88, 0x09}
	CobraNet            = EtherType{0x88, 0x19}
	MPLSUnicast         = EtherType{0x88, 0x47}
	MPLSMulticast       = EtherType{0x88, 0x48}
	PPPoEDiscovery      = EtherType{0x88, 0x63}
	PPPoESession        = EtherType{0x88, 0x64}
	JumboFrames         = EtherType{0x88, 0x70}
	HomePlug1_0MME      = EtherType{0x88, 0x7B}
	IEEE802_1X          = EtherType{0x88, 0x8E}
	PROFINET            = EtherType{0x88, 0x92}
	HyperSCSI           = EtherType{0x88, 0x9A}
	AoE                 = EtherType{0x88, 0xA2}
	EtherCAT            = EtherType{0x88, 0xA4}
	EthernetPowerlink   = EtherType{0x88, 0xAB}
	LLDP                = EtherType{0x88, 0xCC}
	SERCOS3             = EtherType{0x88, 0xCD}
	HomePlugAVMME       = EtherType{0x88, 0xE1}
	MRP                 = EtherType{0x88, 0xE3}
	IEEE802_1AE         = EtherType{0x88, 0xE5}
	IEEE1588            = EtherType{0x88, 0xF7}
	IEEE802_1ag         = EtherType{0x89, 0x02}
	FCoE                = EtherType{0x89, 0x06}
	FCoEInit            = EtherType{0x89, 0x14}
	RoCE                = EtherType{0x89, 0x15}
	CTP                 = EtherType{0x90, 0x00}
	VeritasLLT          = EtherType{0xCA, 0xFE}
)
