package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"zinx/src/ziface"
	"zinx/src/znet"
)

// 对 Zinx 的可用性进行验证测试
//读写分离
func main() {
	// 得到一个server,使用server的api
	server := znet.NewServer("zinxV0.9")
	// 为message设置对应的路由器，处理逻辑
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &RandomRouter{})
	//为 conn 创建后，销毁前的额外处理逻辑
	server.SetOnConnStart(createConn)
	server.SetOnConnStop(closeConn)
	// 运行server
	server.Serve()
}

type PingRouter struct {
	router znet.BaseRouter
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	fmt.Printf("Data:[%s]\n", request.GetData())
	request.GetConnection().SendMsg(1, []byte("ping... ping... ping..."))
}

func (p *PingRouter) PostHandle(request ziface.IRequest) {
}

type RandomRouter struct {
	znet.BaseRouter
}

func (r *RandomRouter) Handle(request ziface.IRequest) {
	fmt.Printf("Data:[%s]\n", request.GetData())
	request.GetConnection().SendMsg(1, []byte("Roudom:["+strconv.Itoa(rand.Intn(100))+"]"))
}

func createConn(connection ziface.IConnection) {
	fmt.Printf("HOOK函数==>>创建Conn：[%d]\n", connection.GetConnID())
}
func closeConn(connection ziface.IConnection) {
	fmt.Printf("HOOK函数==>>销毁Conn：[%d]\n", connection.GetConnID())
}
