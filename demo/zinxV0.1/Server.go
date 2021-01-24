package main

import "zinx/znet"

func main() {
	// 得到一个server,使用server的api
	server := znet.NewServer("zinxV0.1")
	// 运行server
	server.Serve()
}
