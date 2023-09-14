//go:build linux

package common

import (
	"io"
	"net"
)

// TAP 工作在数据链路层(L2)，处理以太网帧

type TapIface struct {
	Name               string
	io.ReadWriteCloser // 实现从文件中读取信息
}

type FrameTag int

const (
	FrameNotTagged FrameTag = iota << 2
	FrameTagged
	FrameDoubleTagged
)

func (i *TapIface) ParseDstMACFromFrame(frame []byte) net.HardwareAddr {
	return frame[:6]
}

func (i *TapIface) ParseSrcMACFromFrame(frame []byte) net.HardwareAddr {
	return frame[6:12]
}

func (i *TapIface) GetFrameTag(frame []byte) FrameTag {
	switch {
	// 802.1Q 标签长度是 4 byte (2 byte TPID + 2 byte TCI:{Priority, CFI, VLAN ID})
	case frame[12] == 0x81 && frame[13] == 0x00:
		// 帧被单标记 TPID
		return FrameTagged
	case frame[12] == 0x88 && frame[13] == 0xA8:
		// 802.1ad TPID
		return FrameDoubleTagged
	}
	return FrameNotTagged
}

func (i *TapIface) ParseMACEtherType(frame []byte) EtherType {
	typePos := 12 + i.GetFrameTag(frame)
	return EtherType{frame[typePos], frame[typePos+1]}
}

func (i *TapIface) ParseMACPayload(frame []byte) []byte {
	return frame[12+i.GetFrameTag(frame)+2:]
}

func NewTAP(tunIfaceName string) (*TapIface, error) {
	return &TapIface{}, nil
}
