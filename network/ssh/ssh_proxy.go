package ssh

import (
	"io"
	"log/slog"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	host string
	port int16
}

func NewSSHEndpoint(host string, port int16) *Endpoint {
	return &Endpoint{
		host: host,
		port: port,
	}
}

func (ep *Endpoint) String() string {
	_port := strconv.Itoa(int(ep.port))
	return ep.host + ":" + _port
}

type Tunnel struct {
	Local  *Endpoint
	Remote *Endpoint
	Server *Endpoint
	Config *ssh.ClientConfig
}

func (t *Tunnel) Start() error {
	lis, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}
	defer func() {
		if lis != nil {
			_ = lis.Close()
		}
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go t.forward(conn)
	}
}

func (t *Tunnel) forward(localConn net.Conn) {
	conn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		slog.Error("ssh dial failed", "error", err)
		return
	}

	remoteConn, err := conn.Dial("tcp", t.Remote.String())
	if err != nil {
		slog.Error("remote ssh server dial failed", "error", err)
		return
	}
	connCopy := func(writer, reader net.Conn) {
		defer func() {
			if writer != nil {
				_ = writer.Close()
			}
			if reader != nil {
				_ = reader.Close()
			}
		}()
		if _, err := io.Copy(writer, reader); err != nil {
			slog.Error("io copy failed", "error", err)
		}
	}
	go connCopy(localConn, remoteConn)
	go connCopy(remoteConn, localConn)
}
