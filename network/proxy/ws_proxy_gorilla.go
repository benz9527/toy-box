package proxy

import (
	"crypto/tls"
	"errors"
	Fws "github.com/gofiber/contrib/websocket"
	Fiber "github.com/gofiber/fiber/v2"
	Gws "github.com/gorilla/websocket"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
)

// https://datatracker.ietf.org/doc/html/rfc6455
// 5.2. Base Framing Protocol
//      0                   1                   2                   3
//      0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
//     +-+-+-+-+-------+-+-------------+-------------------------------+
//     |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
//     |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
//     |N|V|V|V|       |S|             |   (if payload len==126/127)   |
//     | |1|2|3|       |K|             |                               |
//     +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
//     |     Extended payload length continued, if payload len == 127  |
//     + - - - - - - - - - - - - - - - +-------------------------------+
//     |                               |Masking-key, if MASK set to 1  |
//     +-------------------------------+-------------------------------+
//     | Masking-key (continued)       |          Payload Data         |
//     +-------------------------------- - - - - - - - - - - - - - - - +
//     :                     Payload Data continued ...                :
//     + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
//     |                     Payload Data continued ...                |
//     +---------------------------------------------------------------+

// https://gorilla.github.io/
// benchmark test site
// https://github.com/crossbario/autobahn-testsuite

func wsProxyByGorilla(ctx *Fiber.Ctx) error {
	var (
		fwsHeader = ctx.GetRespHeaders()
		gwsHeader = make(http.Header)
		// 跳过 wss 协议中的 SSL 证书的客户端校验
		gwsDialer = Gws.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		gwsUrl    = url.URL{
			Scheme: "ws",
			Path:   "",
		}
	)
	// 头部信息传递处理
	if traceID, ok := fwsHeader["X-Request-ID"]; ok {
		gwsHeader["X-Request-ID"] = []string{traceID}
	}

	nextHopWS, nextHopResp, nextHopErr := gwsDialer.Dial(gwsUrl.String(), gwsHeader)
	if errors.Is(nextHopErr, io.ErrUnexpectedEOF) {
		// 尝试切换为 wss 的加密通讯。
		// 这种切换方式比较朴素，可以尝试改进
		gwsUrl = url.URL{
			Scheme: "wss",
			Path:   "",
		}
		nextHopWS, nextHopResp, nextHopErr = gwsDialer.Dial(gwsUrl.String(), gwsHeader)
	}
	if nextHopErr != nil {
		if nextHopResp != nil {
			ctx.Status(nextHopResp.StatusCode)
			return ctx.SendString(nextHopResp.Status)
		}
		return ctx.SendStatus(http.StatusServiceUnavailable)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	return Fws.New(func(conn *Fws.Conn) {
		defer func() {
			// 连接关闭是有顺序的
			slog.Info("closing gorilla websocket (next hop websocket server) connection...")
			_ = nextHopWS.Close()
			slog.Info("closing fiber websocket (proxy server) connection...")
			_ = conn.Close()
		}()
		// 实际使用需要去管理协程数量
		// 也可以使用 io.Copy ? 这里暂时应该是不行的
		go func() {
		WSRespLoop:
			for {
				msgT, msgB, err := nextHopWS.ReadMessage()
				if err != nil {
					wg.Done()
					break WSRespLoop
				}
				if err = conn.WriteMessage(msgT, msgB); err != nil {
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
				if err = nextHopWS.WriteMessage(msgT, msgB); err != nil {
					wg.Done()
					break WSRecvLoop
				}
			}
		}()
		wg.Wait()
	})(ctx)
}
