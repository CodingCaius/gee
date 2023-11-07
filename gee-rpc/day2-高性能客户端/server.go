package geerpc

import (
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

// 用于配置RPC服务器的选项
type Option struct {
	MagicNumber int        // MagicNumber marks this's a geerpc request
	CodecType   codec.Type // 客户端可以选择不同的Codec来编码body
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server代表一个RPC服务器。
type Server struct{}

// NewServer 返回一个新服务器
func NewServer() *Server {
	return &Server{}
}

// DefaultServer 是 *Server 的默认实例
var DefaultServer = NewServer()

// RPC服务器在单个连接上提供服务的核心逻辑
// ServeConn 阻塞，为连接提供服务，直到客户端挂断
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	// 通过json.NewDecoder解码连接中的选项信息，并存储在opt变量中
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	// 检查连接中的MagicNumber是否与期望的值相匹配。MagicNumber用于标识这是一个geerpc请求
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	// 根据选项中的CodecType选择相应的编解码器函数
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	// 通过选择的编解码器函数创建一个具体的编解码器（f(conn)），然后调用serveCodec方法开始处理请求
	server.serveCodec(f(conn))
}

// invalidRequest 是发生错误时响应 argv 的占位符
var invalidRequest = struct{}{}

// serveCodec 负责在给定的编解码器上提供RPC服务
func (server *Server) serveCodec(cc codec.Codec) {
	// sending 是一个互斥锁，用于确保在发送完整的响应之前不会有其他响应被发送。这是因为在并发情况下，可能会有多个请求同时到达服务器
	sending := new(sync.Mutex) // make sure to send a complete response
	// wg 是一个等待组，用于等待所有的请求都被处理完毕。在每处理一个请求时，都会通过 wg.Add(1) 增加计数，处理完成时通过 wg.Done() 减少计数。最后，通过 wg.Wait() 等待所有请求的完成。
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // 无法恢复，所以关闭连接
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request stores all information of a call
type request struct {
	h            *codec.Header // 请求头
	argv, replyv reflect.Value // argv and replyv of request
}

func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	// TODO: now we don't know the type of request argv
	// day 1, just suppose it's string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO，应该调用注册的 rpc 方法来获得正确的回复v 第一天，只需打印 argv 并发送一条 hello 消息
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// Accept 接受侦听器上的连接并为每个传入连接提供处理方法
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		// 处理过程交给了 ServerConn 方法
		go server.ServeConn(conn)
	}
}

// Accept 接受侦听器上的连接并为每个传入连接提供处理方法
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }