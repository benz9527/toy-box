//go:build linux

package common

import (
	"golang.org/x/sys/unix"
	"io"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// TUN 虚拟网卡，主要处理 IP 报文，处于网络层 (L3)
// TUN - IPv4 Packet:
//  +---------------------------------------------------------------------------------------------------------------+
//  |       | Octet |           0           |           1           |           2           |           3           |
//  | Octet |  Bit  |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|
//  +---------------------------------------------------------------------------------------------------------------+
//  |   0   |   0   |  Version  |    IHL    |      DSCP       | ECN |                 Total  Length                 |
//  +---------------------------------------------------------------------------------------------------------------+
//  |   4   |  32   |                Identification                 | Flags  |           Fragment Offset            |
//  +---------------------------------------------------------------------------------------------------------------+
//  |   8   |  64   |     Time To Live      |       Protocol        |                Header Checksum                |
//  +---------------------------------------------------------------------------------------------------------------+
//  |  12   |  96   |                                       Source IP Address                                       |
//  +---------------------------------------------------------------------------------------------------------------+
//  |  16   |  128  |                                    Destination IP Address                                     |
//  +---------------------------------------------------------------------------------------------------------------+
//  |  20   |  160  |                                     Options (if IHL > 5)                                      |
//  +---------------------------------------------------------------------------------------------------------------+
//  |  24   |  192  |                                                                                               |
//  |  30   |  224  |                                            Payload                                            |
//  |  ...  |  ...  |                                                                                               |
//  +---------------------------------------------------------------------------------------------------------------+

type TunIface struct {
	Name               string
	io.ReadWriteCloser // 实现从文件中读取信息
}

func (i *TunIface) CloseIface() error {
	if i.ReadWriteCloser != nil {
		if err := i.ReadWriteCloser.Close(); err != nil {
			return err
		}
	}
	return TearDownIface(i.Name)
}

func (i *TunIface) CreateDataChannels() (<-chan []byte, <-chan error) {
	bufSize := 1522
	dataCh, errCh := make(chan []byte), make(chan error)
	go func() {
		for {
			buf := make([]byte, bufSize)
			n, err := i.Read(buf)
			if err != nil {
				errCh <- err
			} else {
				dataCh <- buf[:n:n]
			}
		}
	}()
	return dataCh, errCh
}

func (i *TunIface) WaitOrReadData(
	dataCh <-chan []byte,
	errCh <-chan error,
) {
	timeoutTimer := time.NewTimer(6 * time.Second).C
	for {
		select {
		case <-timeoutTimer:
			log.Fatal("waiting for broadcast packet timeout")
		case err := <-errCh:
			log.Fatalf("read packet error: %v", err)
		case buf := <-dataCh:
			// 0x60, 0x0,  0x0,  0x0,  0x0,  0x8, 0x3a, 0xff,
			// 0xfe, 0x80, 0x0,  0x0,  0x0,  0x0, 0x0,  0x0,
			// 0x8b, 0x5b, 0x1f, 0x3b, 0xd9, 0xd, 0x3d, 0xa7,
			// 0xff, 0x2,  0x0,  0x0,  0x0,  0x0, 0x0,  0x0,
			// 0x0,  0x0,  0x0,  0x0,  0x0,  0x0, 0x0,  0x2,
			// 0x85, 0x0,  0xbb, 0xeb, 0x0,  0x0, 0x0,  0x0
			log.Printf("broadcast frame: %#v\n", buf)
		}
	}
}

func NewTUN(tunIfaceName string) (*TunIface, error) {
	tunDev, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	// 创建接口
	ifr, err := unix.NewIfreq(tunIfaceName)
	if err != nil {
		return nil, err
	}
	// IFF_TUN 指定 TUN 设备，但是报文不包含 Ethernet header
	// IFF_NO_PI 置位后表示让 linux core 不提供报文信息，只需要纯 IP 报文;
	// 否则 linux core 会在报文开始的地方添加 2 byte 标识和 2 byte 协议
	ifr.SetUint16(unix.IFF_TUN | unix.IFF_NO_PI) // 这里还有一个 multiple queue: unix.IFF_MULTI_QUEUE
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, tunDev.Fd(), syscall.TUNSETIFF, uintptr(unsafe.Pointer(ifr)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	// 配置 TUN dev 用户属性，先做默认配置
	anyUserAsOwner, anyGroup := -1, -1
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunDev.Fd(), syscall.TUNSETOWNER, uintptr(unsafe.Pointer(&anyUserAsOwner)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunDev.Fd(), syscall.TUNSETGROUP, uintptr(unsafe.Pointer(&anyGroup)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	// 配置 TUN dev 的持久化属性，先做默认配置
	persist := 0
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunDev.Fd(), syscall.TUNSETPERSIST, uintptr(unsafe.Pointer(&persist)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}

	return &TunIface{
		Name:            tunIfaceName,
		ReadWriteCloser: tunDev,
	}, nil
}
