package main

import (
	"fmt"
	"zinx/src/ziface"
	"zinx/src/znet"
)

// 对 Zinx 的可用性进行验证测试
func main() {
	// 得到一个server,使用server的api
	server := znet.NewServer("zinxV0.4")
	server.AddRouter(&PingRouter{})
	// 运行server
	server.Serve()
}

type PingRouter struct {
	router znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("before handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("befor ping..."))
	if err != nil {
		fmt.Println("Call back before error:", err)
	}
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("do handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping... ping... ping..."))
	if err != nil {
		fmt.Println("Call back handle error:", err)
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("after handle ")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("Call back after error:", err)
	}
}
