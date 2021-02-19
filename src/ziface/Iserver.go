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
	//返回当前server的连接管理器
	GetConnManager() IConnManager
	// 设置建立连接之后的处理
	SetOnConnStart(func(connection IConnection))
	// 设置销毁连接之前的处理
	SetOnConnStop(func(connection IConnection))
	//调用onConStart 钩子函数的方法
	CallOnConnStart(connection IConnection)
	//调用onCon了Stop 钩子函数的方法
	CallOnConnStop(connection IConnection)
}
