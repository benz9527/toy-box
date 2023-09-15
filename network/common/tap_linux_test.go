//go:build linux

package common

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestOpenTAPDev(t *testing.T) {
	ifaceMeta := net.IPNet{
		IP:   net.IPv4(10, 0, 11, 1),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
	tap, err := NewTAP("mytap0")
	assert.NoError(t, err)
	defer func() {
		_ = tap.CloseIface()
	}()
	err = SetupIface(ifaceMeta, tap.Name)
	assert.NoError(t, err)

	dataCh, errCh, release := tap.CreateDataChannels()
	defer func() { release() }()
	err = PingIface(ifaceMeta.IP)
	assert.NoError(t, err)
	tap.WaitOrReadData(dataCh, errCh)
}
