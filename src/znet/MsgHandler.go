package znet

import (
	"fmt"
	"zinx/src/ziface"
)

type MsgHandler struct {
	Api map[uint32]ziface.IRouter
}

// 初始化方法
func NewMsgHandler() ziface.IMsgHandler {
	return &MsgHandler{
		Api: make(map[uint32]ziface.IRouter),
	}
}

func (m *MsgHandler) Put(id uint32, router ziface.IRouter) {
	//	 是否存在该id对应的router
	_, ok := m.Api[id]
	if ok {
		return
	}
	//	添加
	m.Api[id] = router
}

func (m *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	id := request.GetMsgId()
	router, ok := m.Api[id]
	if !ok {
		fmt.Printf("api msgid:[%v] is not found the router ,need register!\n", request.GetMsgId())
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
	fmt.Printf("api msgID:[%v],exec  router finish!\n", request.GetMsgId())
}
