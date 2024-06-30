package server

import (
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"sync/atomic"
	"zeroim/common/libnet"
)

type Session struct {
	tcpServer *TCPServer
	codec     libnet.Codec
	closeChan chan int
	closeFlag int32
	sendChan  chan libnet.Message
}

const (
	loginCmd  = 1
	LogoutCmd = 2
)

var (
	SessionClosedError  = errors.New("session closed")
	SessionBlockedError = errors.New("session blocked")
)

func NewSession(tcpServer *TCPServer, codec libnet.Codec, sendChanSize int) *Session {
	sn := &Session{
		tcpServer: tcpServer,
		codec:     codec,
		closeChan: make(chan int),
	}
	if sendChanSize > 0 {
		sn.sendChan = make(chan libnet.Message, sendChanSize)
		go sn.sendLoop()
	}
	return sn
}

// 执行cmd命令
func (s *Session) HandlePacket(msg *libnet.Message) error {
	// 心跳处理
	// 登出
	if msg.Cmd == LogoutCmd {
		s.Close()
	}
	// 消息转发，调用rpc进行处理
	err := s.codec.Send(*msg)
	if err != nil {
		logx.Errorf("[HandlePacket] s.codec.Send msg: %v error: %v", msg, err)
	}

	return nil
}

// 将消息放入发送管道中，起到缓冲作用
func (s *Session) Send(msg libnet.Message) error {
	if s.IsClosed() {
		return SessionClosedError
	}
	if s.sendChan == nil {
		return s.codec.Send(msg)
	}
	select {
	case s.sendChan <- msg:
		return nil
	default:
		return SessionBlockedError
	}
}

// 取出管道中的消息，使用codec发送
func (s *Session) sendLoop() {
	defer s.Close()
	for {
		select {
		case msg := <-s.sendChan:
			err := s.codec.Send(msg)
			if err != nil {
				logx.Errorf("[sendLoop] s.codec.Send msg: %v error: %v", msg, err)
			}
		//	管道关闭后取出零值
		case <-s.closeChan:
			return
		}
	}
}

func (s *Session) Close() error {
	if atomic.CompareAndSwapInt32(&s.closeFlag, 0, 1) {
		err := s.codec.Close()
		close(s.closeChan)
		return err
	}
	return SessionClosedError
}

// 检查发送数据管道是否关闭
func (s *Session) IsClosed() bool {
	return atomic.LoadInt32(&s.closeFlag) == 1
}

func (s *Session) Receive() (*libnet.Message, error) {
	return s.codec.Receive()
}
