package utils

import (
	"net"
	"testing"
	"time"
)

func TestPingByTCP(t *testing.T) {
	type args struct {
		ip      net.IP
		timeout time.Duration
		count   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ip:      net.ParseIP("192.168.157.1"),
				timeout: 2 * time.Second,
				count:   4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PingByTCP(tt.args.ip, tt.args.timeout, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("PingByTCP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPingByUDP(t *testing.T) {
	type args struct {
		ip      net.IP
		timeout time.Duration
		count   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ip:      net.ParseIP("192.168.157.1"),
				timeout: 2 * time.Second,
				count:   4,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PingByUDP(tt.args.ip, tt.args.timeout, tt.args.count); (err != nil) != tt.wantErr {
				t.Errorf("PingByUDP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
