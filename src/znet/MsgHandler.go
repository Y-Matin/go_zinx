package znet

import (
	"fmt"
	"sync"
	"zinx/src/ziface"
)

type MsgHandler struct {
	// 路由集合
	Api map[uint32]ziface.IRouter

	// worker数量 池大小
	WorkerPoolSize uint32
	// 消息队列集合
	TaskQueues  map[uint32]chan ziface.IRequest
	QueueLength uint32
	// 下一个存放任务的序号
	nextNumber uint32
}

var msgHandler *MsgHandler
var mutex sync.Mutex
var wg = &sync.WaitGroup{}

func init() {
	workSize := 10
	msgHandler = &MsgHandler{
		Api:            make(map[uint32]ziface.IRouter),
		WorkerPoolSize: 10,
		TaskQueues:     make(map[uint32]chan ziface.IRequest),
		QueueLength:    100,
		nextNumber:     0,
	}
	// 创建 固定数量的worker-goroutine
	wg.Add(workSize)
	for i := 0; i < workSize; i++ {
		// 此处一定要传参进去，否则 i 的值始终是一个。 不是原子性，添加入参，使用i的副本，保证原子性
		go func(i uint32) {
			channel := make(chan ziface.IRequest, msgHandler.QueueLength)
			// 将channel放入等待队列集合中
			mutex.Lock()
			msgHandler.TaskQueues[i] = channel
			mutex.Unlock()
			wg.Done()
			for true {
				select {
				//一直读取channel中的request
				case req := <-channel:
					msgHandler.DoMsgHandler(req)
				}
			}
		}(uint32(i))
	}
	wg.Wait()
	fmt.Printf("初始化Worker池，worker数量：[%d],等待队列长度:[%d]\n", msgHandler.WorkerPoolSize, msgHandler.QueueLength)
}

// 初始化方法: 单例
func NewMsgHandler() ziface.IMsgHandler {
	return msgHandler
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

// 默认使用轮询算法分配task
func (m *MsgHandler) AddTask(req ziface.IRequest) {
	index := m.nextNumber % uint32(len(m.TaskQueues))
	m.nextNumber++
	m.TaskQueues[index] <- req
	fmt.Println("Add Task to NO.[", index, "] worker")
}
