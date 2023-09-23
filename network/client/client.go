package client

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func getFixDnsResolverClient() *http.Client {
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
				_dialer := net.Dialer{
					Timeout: 1 * time.Second,
				}
				// 可以使用本地的 dns 地址来加快访问 dns 的解析速度
				return _dialer.DialContext(ctx, "udp", "114.114.114.114:53")
			},
		},
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
		},
	}
}

// 自定义新的 Dialer

type MyDialer interface {
	Dial(network, addr string) (net.Conn, error)
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

type MyDialContext func(ctx context.Context, network, addr string) (net.Conn, error)

func (m MyDialContext) Dial(network, addr string) (net.Conn, error) {
	return m(context.Background(), network, addr)
}

func (m MyDialContext) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return m(ctx, network, addr)
}

type DomainType int

const (
	TypeA DomainType = iota
	TypeAAAA
)

type DnsCacheItem struct {
	duration int64          // random TTL
	dnsTable [][]net.IPAddr // CNAME -> 0:A, 1:AAAA
	count    []int          // Prepare for DNS Round-Robin
}

type MyResolver interface {
	LookupIpAddresses(ctx context.Context, hostDomainName string) ([][]net.IPAddr, error)
}

type MyDialerStats struct {
	CacheMiss       int64
	CacheHit        int64
	DNSQuery        int64
	SuccessDNSQuery int64
	ConnCount       int64
}

// dialer wrapper
type myDialer struct {
	dialer   MyDialer
	resolver MyResolver
	stats    MyDialerStats
	rand     *rand.Rand
	rwLock   *sync.RWMutex
	cache    map[string]*DnsCacheItem // Host domain name as key
}

func (m *myDialer) Dial(network, addr string) (net.Conn, error) {
	return m.DialContext(context.Background(), network, addr)
}

func (m *myDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	atomic.AddInt64(&m.stats.ConnCount, 1)
	if len(addr) > 0 {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		if len(strings.TrimSpace(host)) > 0 {
			ip, err := m.getIP(ctx, host)
			if err != nil {
				return nil, err
			}
			addr = net.JoinHostPort(ip.String(), port)
		}
	}
	return m.dialer.DialContext(ctx, network, addr)
}

func (m *myDialer) getIP(ctx context.Context, hostDomainName string) (net.IPAddr, error) {
	return net.IPAddr{}, nil
}

func (m *myDialer) getIPByDomainTypeFromCache(hostDomainName string) (net.IPAddr, error) {
	m.rwLock.RLock()
	_, _ = m.cache[hostDomainName]
	m.rwLock.RUnlock()

	// TODO Select A ip or AAAA ip
	// TODO If A ip list is empty then attempt to select ip from AAAA ip list
	// TODO If both A and AAAA ip list is empty, return error
	// TODO Select ip by policy, such as rand or round-robin
	// TODO Store the stats like cache hit, cache missing

	return net.IPAddr{}, nil
}

func (m *myDialer) exchangeFromUpstreamDNS(ctx context.Context, hostDomainName string) ([][]net.IPAddr, error) {
	// TODO Query DNS by resolver.
	// TODO Store the stats like cache query requests, success dns query requests
	// TODO New DNS cache item and set the cache expiry timestamp (duration)
	// TODO Add DNS cache item to time wheel
	return make([][]net.IPAddr, 0), nil
}

func (m *myDialer) Stats() MyDialerStats {
	return m.stats
}

func NewMyDialer() (MyDialer, error) {
	md := &myDialer{}

	if md.dialer == nil {
		md.dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	}

	return md, nil
}
