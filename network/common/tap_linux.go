//go:build linux

package common

import (
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net"
	"os"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

// TAP (以太网设备) 工作在数据链路层(L2)，处理以太网帧

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

func (i *TapIface) CloseIface() error {
	if i.ReadWriteCloser != nil {
		if err := i.ReadWriteCloser.Close(); err != nil {
			return err
		}
	}
	return TearDownIface(i.Name)
}

func (i *TapIface) CreateDataChannels() (<-chan []byte, <-chan error, func()) {
	bufSize := 1522
	dataCh, errCh := make(chan []byte), make(chan error)
	isClosed := atomic.Bool{} // 如果不使用这种方式，就得 panic-recover 的重操作
	go func() {
	WriteFrame:
		for {
			buf := make([]byte, bufSize)
			n, err := i.Read(buf)
			if isClosed.Load() {
				break WriteFrame
			}
			if err != nil {
				errCh <- err
			} else {
				dataCh <- buf[:n:n]
			}
		}
	}()
	// 返回只读通道和通道释放方法
	return dataCh, errCh, func() {
		// 关闭通道
		close(dataCh)
		close(errCh)
		isClosed.Store(true)
	}
}

func (i *TapIface) WaitOrReadData(
	dataROCh <-chan []byte,
	errROCh <-chan error,
) {
	timeoutTimer := time.NewTimer(2 * time.Second).C
ReadFrame:
	for {
		select {
		case <-timeoutTimer:
			log.Println("waiting for broadcast packet timeout")
			break ReadFrame
		case err, isOpen := <-errROCh: // for-range 遍历读通道可以高效检查是否关闭
			if !isOpen {
				break ReadFrame
			}
			log.Printf("read packet error: %v\n", err)
			break ReadFrame
		case buf, isOpen := <-dataROCh: // for-range 遍历读通道可以高效检查是否关闭
			if !isOpen {
				break ReadFrame
			}
			log.Printf("tap parse broadcast frame dst mac: %#v, src mac: %#v\n", i.ParseDstMACFromFrame(buf), i.ParseSrcMACFromFrame(buf))
			log.Printf("tap parse broadcast frame ethertype: %#v", i.ParseMACEtherType(buf))
			log.Printf("tap parse broadcast frame payload: %#v\n", i.ParseMACPayload(buf))
		}
	}
}

func NewTAP(tapIfaceName string) (*TapIface, error) {
	// 运行需要 CAP_NET_ADMIN 权限的账户
	tunAsTapDev, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	// 创建接口
	ifr, err := unix.NewIfreq(tapIfaceName)
	if err != nil {
		return nil, err
	}
	// IFF_TAP 指定 TAP 设备，报文包含 Ethernet header
	// IFF_NO_PI 置位后表示让 linux core 不提供报文信息，只需要纯 IP 报文;
	// 否则 linux core 会在报文开始的地方添加 2 byte 标识和 2 byte 协议
	ifr.SetUint16(unix.IFF_TAP | unix.IFF_NO_PI) // 这里还有一个 multiple queue: unix.IFF_MULTI_QUEUE
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, tunAsTapDev.Fd(), syscall.TUNSETIFF, uintptr(unsafe.Pointer(ifr)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	// 配置 TAP dev 用户属性，先做默认配置
	anyUserAsOwner, anyGroup := -1, -1
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunAsTapDev.Fd(), syscall.TUNSETOWNER, uintptr(unsafe.Pointer(&anyUserAsOwner)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunAsTapDev.Fd(), syscall.TUNSETGROUP, uintptr(unsafe.Pointer(&anyGroup)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}
	// 配置 TAP dev 的持久化属性，先做默认配置
	persist := 0
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, tunAsTapDev.Fd(), syscall.TUNSETPERSIST, uintptr(unsafe.Pointer(&persist)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl", errno)
	}

	return &TapIface{
		Name:            tapIfaceName,
		ReadWriteCloser: tunAsTapDev,
	}, nil
}
