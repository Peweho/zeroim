package libnet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	maxBodySize = 1 << 12
	//协议字段
	packSize      = 4
	headerSize    = 2
	verSize       = 1
	statusSize    = 1
	serviceIdsize = 2
	cmdSize       = 2
	seqSize       = 4
	//body

	rawHeaderSize = verSize + statusSize + serviceIdsize + cmdSize + seqSize
	maxPackSize   = maxBodySize + rawHeaderSize + headerSize + packSize
	//offset
	headerOffset    = 0
	verOffset       = headerOffset + headerSize
	statusOffset    = verOffset + verSize
	serviceIdOffset = statusOffset + statusSize
	cmdOffset       = serviceIdOffset + serviceIdsize
	seqOffset       = cmdOffset + cmdSize
	bodyOffset      = seqOffset + seqSize
)

var (
	ErrRawPackLen   = errors.New("")
	ErrRawHeaderLen = errors.New("")
)

type Header struct {
	Version   uint8
	Status    uint8
	ServiceId uint16
	Cmd       uint16
	Seq       uint32
}

type Message struct {
	Header
	Body []byte
}

func (m *Message) Fromat() string {
	return fmt.Sprintf("Version:%d, Status:%d, ServiceId:%d, Cmd:%d, Seq:%d, Body:%s",
		m.Version, m.Status, m.ServiceId, m.Cmd, m.Seq, string(m.Body))
}

type imCodec struct {
	conn net.Conn
}

// 读取内容转换为数字
func (c *imCodec) readPackSize() (uint32, error) {
	return c.readUint32BE()
}

// 从conn读取packSize大小的内容，并转换为大端存储的数字
func (c *imCodec) readUint32BE() (uint32, error) {
	b := make([]byte, packSize)
	_, err := io.ReadFull(c.conn, b)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b), nil
}

// 从conn读取msgSize大小的内容返回
func (c *imCodec) readPacket(msgSize uint32) ([]byte, error) {
	b := make([]byte, msgSize)
	_, err := io.ReadFull(c.conn, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 接收数据包，并解析为Message
func (c *imCodec) Receive() (*Message, error) {
	//读取包的总长度
	packLen, err := c.readPackSize()
	if err != nil {
		return nil, err
	}
	//如果大于最大包长，返回错误
	if packLen > maxPackSize {
		return nil, ErrRawPackLen
	}
	//读取包的全部内容
	buf, err := c.readPacket(packLen)
	if err != nil {
		return nil, err
	}

	// 定义消息体，设置各字段内容
	msg := &Message{}
	msg.Version = buf[verOffset]
	msg.Status = buf[statusOffset]
	msg.ServiceId = binary.BigEndian.Uint16(buf[serviceIdOffset:cmdOffset])
	msg.Cmd = binary.BigEndian.Uint16(buf[cmdOffset:seqOffset])
	msg.Seq = binary.BigEndian.Uint32(buf[seqOffset:bodyOffset])

	headerLen := binary.BigEndian.Uint16(buf[headerOffset:verOffset])
	// 判断头部长度是否正确
	if headerLen != rawHeaderSize {
		return nil, ErrRawHeaderLen
	}
	// 判断是否存在消息体
	// 如果没有，包的长度等于头部长度
	if packLen > uint32(headerLen) {
		msg.Body = buf[bodyOffset:packLen]
	}
	return msg, nil
}

// Message转换为字节流，输入到conn中
func (c *imCodec) Send(msg Message) error {
	//设置包的长度，数值为 头部大小数据+头部数据+消息体
	packLen := headerSize + rawHeaderSize + len(msg.Body)
	packLenBuf := make([]byte, packSize)
	// 将packLen转换为大端存储的数字写入到packLenBuf
	binary.BigEndian.PutUint32(packLenBuf[:packSize], uint32(packLen))

	buf := make([]byte, packLen)
	// header
	binary.BigEndian.PutUint16(buf[headerOffset:], uint16(rawHeaderSize))
	buf[verOffset] = msg.Version
	buf[statusOffset] = msg.Status
	binary.BigEndian.PutUint16(buf[serviceIdOffset:], msg.ServiceId)
	binary.BigEndian.PutUint16(buf[cmdOffset:], msg.Cmd)
	binary.BigEndian.PutUint32(buf[seqOffset:], msg.Seq)

	// body
	copy(buf[headerSize+rawHeaderSize:], msg.Body)
	// 合并包长度和包内容
	allBuf := append(packLenBuf, buf...)
	n, err := c.conn.Write(allBuf)
	if err != nil {
		return err
	}
	if n != len(allBuf) {
		return fmt.Errorf("n:%d, len(buf):%d", n, len(buf))
	}
	return nil
}

// 设置读超时时间
func (c *imCodec) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

// 设置写超时时间
func (c *imCodec) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *imCodec) Close() error {
	return c.conn.Close()
}

type Protocol interface {
	NewCodec(conn net.Conn) Codec
}

type IMProtocol struct{}

func NewIMProtocol() Protocol {
	return &IMProtocol{}
}

func (p *IMProtocol) NewCodec(conn net.Conn) Codec {
	return &imCodec{conn: conn}
}
