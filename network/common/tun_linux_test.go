//go:build linux

package common

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestOpenTUNDev(t *testing.T) {
	ifaceMeta := net.IPNet{
		IP:   net.IPv4(10, 0, 11, 1),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
	tun, err := NewTUN("mytun0")
	assert.NoError(t, err)
	defer func() {
		_ = tun.CloseIface()
	}()
	err = SetupIface(ifaceMeta, tun.Name)
	assert.NoError(t, err)

	dataCh, errCh, release := tun.CreateDataChannels()
	defer func() { release() }()
	err = PingIface(ifaceMeta.IP)
	assert.NoError(t, err)
	tun.WaitOrReadData(dataCh, errCh)
}
