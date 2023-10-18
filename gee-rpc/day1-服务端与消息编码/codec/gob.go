// 实现了 Codec 接口的 Gob 编解码器 GobCodec
package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser // 读写流
	buf  *bufio.Writer // 缓冲写入器
	dec  *gob.Decoder // 解码器
	enc  *gob.Encoder // 编码器
}

// 检查 *GobCodec 类型是否实现了 Codec 接口
var _ Codec = (*GobCodec)(nil)

// 创建并返回一个新的 GobCodec 实例
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// 从流中解码消息头部信息并填充到给定的 Header 结构体中
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// 从流中解码消息主体信息并填充到给定的接口类型中
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// 将消息头和消息体编码并写入到流中。它使用 Gob 编码器来实现这一过程
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

// 关闭连接
func (c *GobCodec) Close() error {
	return c.conn.Close()
}