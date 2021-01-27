package ziface

// 路由的抽象接口
// 路由的数据都是IRequest
type IRouter interface {
	// 处理 conn业务方法之前
	PreHandle(request IRequest)
	// 处理 conn业务方法
	Handle(request IRequest)
	// 处理 conn业务方法之后
	PostHandle(request IRequest)
}
