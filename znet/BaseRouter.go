package znet

import "zinx/ziface"

// 定义一个baseRouter的基类，然后根据需要对这个基类的方法进行重写
// 方法没有具体的实现。有继承者 按需要 进行 重写
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {}

func (b *BaseRouter) Handle(request ziface.IRequest) {}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {}
