package proxy

import (
	"fmt"
	Gws "github.com/gobwas/ws"
	Gwsflate "github.com/gobwas/ws/wsflate"
	Fws "github.com/gofiber/contrib/websocket"
	Fiber "github.com/gofiber/fiber/v2"
	"io"
	"log/slog"
	"net"
	"net/http"
	"runtime"
	"sync"
)

// https://github.com/gobwas/ws

func wsProxyByGobwas(ctx *Fiber.Ctx) error {
	var (
		fwsHeader = ctx.GetRespHeaders()
		gwsExt    = Gwsflate.Extension{
			Parameters: Gwsflate.DefaultParameters,
		}
		nextHopDialer = net.Dialer{
			FallbackDelay: -1, // disable dual stack
		}
		upgrader = Gws.Upgrader{
			Negotiate: gwsExt.Negotiate,
			OnHost: func(host []byte) error {
				if string(host) != "" {
					return nil
				}
				return Gws.RejectConnectionError(
					Gws.RejectionStatus(http.StatusForbidden),
					Gws.RejectionHeader(Gws.HandshakeHeaderString("X-None-Empty-Host: true\r\n")),
				)
			},
			OnHeader: func(key, val []byte) error {
				// 字符串安全转义打印
				slog.Info(fmt.Sprintf("gobwas websocket header %q=%q", key, val))
				if string(key) != "" && len(val) > 0 {
					return nil
				}
				return Gws.RejectConnectionError(
					Gws.RejectionReason("empty request headers"),
					Gws.RejectionStatus(http.StatusBadRequest),
				)
			},
			OnBeforeUpgrade: func() (Gws.HandshakeHeader, error) {
				// 头部信息传递处理
				header := Gws.HandshakeHeaderHTTP{
					"X-Go-Client-Version": []string{runtime.Version()},
				}
				if traceID, ok := fwsHeader["X-Request-ID"]; ok {
					header["X-Request-ID"] = []string{traceID}
				}
				return header, nil
			},
		}
	)

	nextHopWS, err := nextHopDialer.DialContext(ctx.UserContext(), "tcp", "")
	if err != nil {
		return err
	}
	// zero-copy upgrade
	_, err = upgrader.Upgrade(nextHopWS)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	return Fws.New(func(conn *Fws.Conn) {
		defer func() {
			// 连接关闭是有顺序的
			slog.Info("closing gobwas websocket (next hop websocket server) connection...")
			_ = nextHopWS.Close()
			slog.Info("closing fiber websocket (proxy server) connection...")
			_ = conn.Close()
		}()
		// 实际使用需要去管理协程数量
		go func() {
		WSRespLoop:
			for {
				hdr, err := Gws.ReadHeader(nextHopWS)
				if err != nil {
					wg.Done()
					break WSRespLoop
				}

				// 获取响应
				msg := make([]byte, hdr.Length)
				_, err = io.ReadFull(nextHopWS, msg)
				if err != nil {
					wg.Done()
					break WSRespLoop
				}
				if hdr.Masked {
					Gws.Cipher(msg, hdr.Mask, 0)
				}

				if err = conn.WriteMessage(int(hdr.OpCode), msg); err != nil {
					wg.Done()
					break WSRespLoop
				}
			}
		}()
		go func() {
		WSRecvLoop:
			for {
				msgT, msgB, err := conn.ReadMessage()
				if err != nil {
					wg.Done()
					break WSRecvLoop
				}

				// 数据转换
				var frame Gws.Frame
				switch msgT {
				case Fws.BinaryMessage:
					frame = Gws.NewBinaryFrame(msgB)
				case Fws.TextMessage:
					frame = Gws.NewTextFrame(msgB)
				case Fws.PingMessage:
					frame = Gws.NewPingFrame(msgB)
				case Fws.PongMessage:
					frame = Gws.NewPongFrame(msgB)
				default:
					frame = Gws.NewCloseFrame(msgB)
				}

				if err = Gws.WriteFrame(nextHopWS, frame); err != nil {
					wg.Done()
					break WSRecvLoop
				}
			}
		}()
		wg.Wait()
	})(ctx)
}
