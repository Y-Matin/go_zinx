package ziface

type IMsgHandler interface {
	//调度 执行对应 Router的消息处理方法
	DoMsgHandler(request IRequest)
	// 添加路由
	Put(id uint32, router IRouter)
}
