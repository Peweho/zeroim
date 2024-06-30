package server

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net"
	"time"
	"zeroim/common/libnet"
	"zeroim/edge/internal/svc"
	"zeroim/imrpc/imrpc"
)

type TCPServer struct {
	Listener net.Listener
	svcCtx   *svc.ServiceContext
	Protocol libnet.Protocol
}

func NewTCPServer(svcCtx *svc.ServiceContext, address string) (*TCPServer, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &TCPServer{
		svcCtx:   svcCtx,
		Listener: listener,
		Protocol: libnet.NewIMProtocol(),
	}, nil
}

func (srv *TCPServer) HandleRequest() {
	var tmpDelay time.Duration

	for {
		//阻塞等待连接
		conn, err := srv.Listener.Accept()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				if tmpDelay == 0 {
					tmpDelay = 10 * time.Microsecond
				} else {
					tmpDelay *= 2
				}
				if max := time.Second; tmpDelay > max {
					tmpDelay = max
				}
				time.Sleep(tmpDelay)
				continue
			}
		}
		//异步处理连接请求
		threading.GoSafe(func() {
			srv.handleRequest(conn)
		})
	}
}

// 请求处理方法
func (srv *TCPServer) handleRequest(conn net.Conn) {
	codec := srv.Protocol.NewCodec(conn)
	session := NewSession(srv, codec, srv.svcCtx.Config.SendChanSize)
	msg, err := session.Receive()
	if err != nil {
		logx.Errorf("[HandleRequest] session.Receive error: %v", err)
		_ = session.Close()
		return
	}

	logx.Infof("[HandleRequest] session.Receive msg: %s", msg.Fromat())

	// 登录校验
	if err := srv.Login(session, msg); err != nil {
		logx.Errorf("[HandleRequest] srv.Login error: %v", err)
		_ = session.Close()
		return
	}

	for {
		//收到数据包，解析出消息
		message, err := session.Receive()
		if err != nil {
			logx.Errorf("[HandleRequest] session.Receive error: %v", err)
			_ = session.Close()
			break
		}

		logx.Infof("[HandleRequest] session.Receive message: %s", message.Fromat())
		// 进行消息处理
		if err = session.HandlePacket(message); err != nil {
			logx.Errorf("[HandleRequest] session.HandlePacket message: %v error: %v", message, err)
		}
	}
}

func (srv *TCPServer) Close() error {
	return srv.Listener.Close()
}

// 调用rpc方法进行登录校验
func (srv *TCPServer) Login(session *Session, msg *libnet.Message) error {
	_, err := srv.svcCtx.IMRpc.Login(context.Background(), &imrpc.LoginRequest{})
	if err != nil {

	}
	// 登录成功，将消息发送到conn中
	_ = session.Send(*msg)

	return nil
}

// 调用rpc方法进行登出
func (srv *TCPServer) Logout(session *Session, msg *libnet.Message) error {
	_, err := srv.svcCtx.IMRpc.Logout(context.Background(), &imrpc.LogoutRequest{})
	if err != nil {

	}
	_ = session.Close()

	return nil
}
