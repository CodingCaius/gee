package main

import (
	"encoding/json"
	"fmt"
	"geerpc"

	"geerpc/codec"
	"log"
	"net"
	"time"
)

// startServer 函数启动一个RPC服务器，并监听端口 8080。它将服务器的地址发送到一个通道 addr 中，以便在 main 函数中获取
func startServer(addr chan string) {
	// pick a free port
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}

func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)

	// 事实上，下面的代码就像一个简单的geerpc客户端

	// 通过 <-addr 从通道中获取服务器地址，并使用 net.Dial 连接到RPC服务器
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// 通过 json.NewEncoder 发送默认的选项（geerpc.DefaultOption）到服务器
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// send request & receive response
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
	// 客户端首先发送 Option 进行协议交换，接下来发送消息头 h := &codec.Header{}，和消息体 geerpc req ${h.Seq}。
	// 最后解析服务端的响应 reply，并打印出来。
}
