package ziface

type Iserver interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()
	//路由功能：给当前的服务注册一个路有方法，功能客户端的连接处理使用
	AddRouter(msgId uint32, router IRouter)
}
