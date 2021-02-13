package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"zinx/src/ziface"
	"zinx/src/znet"
)

// 对 Zinx 的可用性进行验证测试
func main() {
	// 得到一个server,使用server的api
	server := znet.NewServer("zinxV0.6")
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &RandomRouter{})
	// 运行server
	server.Serve()
}

type PingRouter struct {
	router znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("before handle")
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("do handle")
	fmt.Printf("Data:[%s]\n", request.GetData())

	dp := znet.NewDataPackage()
	message := znet.NewMessage(request.GetMsgId(), []byte("ping... ping... ping..."))
	bytes, err := dp.Pack(message)
	if err != nil {
		fmt.Println("Pack error:", err)
	}
	_, err = request.GetConnection().GetTCPConnection().Write(bytes)
	if err != nil {
		fmt.Println("Call back handle error:", err)
	}
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("after handle ")
}

type RandomRouter struct {
	znet.BaseRouter
}

func (r *RandomRouter) Handle(request ziface.IRequest) {
	fmt.Println("do handle")
	fmt.Printf("Data:[%s]\n", request.GetData())

	dp := znet.NewDataPackage()
	message := znet.NewMessage(request.GetMsgId(), []byte("Roudom:[%v]"+strconv.Itoa(rand.Intn(100))))
	bytes, err := dp.Pack(message)
	if err != nil {
		fmt.Println("Pack error:", err)
	}
	_, err = request.GetConnection().GetTCPConnection().Write(bytes)
	if err != nil {
		fmt.Println("Call back handle error:", err)
	}
}
