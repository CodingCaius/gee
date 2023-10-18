// 提供了消息编码和解码的功能
package codec

import (
	"io"
)

// 消息的头部信息
type Header struct {
	ServiceMethod string // 格式为 "Service.Method"
	Seq           uint64 // 客户端选择的序列号
	Error         string
}

// 抽象出对消息体进行编解码的接口 Codec，抽象出接口是为了实现不同的 Codec 实例
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

// 抽象出 Codec 的构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

// 定义了两个实例
const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

// 声明了一个名为 NewCodecFuncMap 的映射，将 Type 与 NewCodecFunc（构造函数）关联起来
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}