package utils

import (
	"fmt"
	TPing "github.com/digineo/go-ping"
	UPing "github.com/prometheus-community/pro-bing"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"time"
)

// PingByTCP
// https://manpages.ubuntu.com/manpages/lunar/en/man7/packet.7.html
// https://squidarth.com/networking/systems/rc/2018/05/28/using-raw-sockets.html
// sudo setcap cap_net_raw,cap_net_admin=eip ./your_exeutable
func PingByTCP(ip net.IP, timeout time.Duration, count int) error {
	pinger, err := TPing.New("0.0.0.0", "::")
	if err != nil {
		return err
	}

	_, err = pinger.PingAttempts(&net.IPAddr{IP: ip}, timeout, count)
	if err != nil {
		return err
	}
	return nil
}

// PingByUDP
// NonPrivileged must set 'sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"'
func PingByUDP(ip net.IP, timeout time.Duration, count int) error {
	pinger, err := UPing.NewPinger(ip.String())
	if err != nil {
		return err
	}
	pinger.Timeout = timeout
	pinger.Count = count
	pinger.SetPrivileged(true)

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			pinger.Stop()
		}
	}()

	pinger.OnRecv = func(pkt *UPing.Packet) {
		slog.Info(fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt))
	}

	pinger.OnDuplicateRecv = func(pkt *UPing.Packet) {
		slog.Info(fmt.Sprintf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL))
	}

	pinger.OnFinish = func(stats *UPing.Statistics) {
		slog.Info(fmt.Sprintf("\n--- %s ping statistics ---\n", stats.Addr))
		slog.Info(fmt.Sprintf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss))
		slog.Info(fmt.Sprintf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt))
	}

	slog.Info("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	return pinger.Run()
}
