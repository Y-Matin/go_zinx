package znet

import (
	"fmt"
	"net"
	"time"
	"zinx/src/utils"
	"zinx/src/ziface"
)

type Server struct {
	Name        string                              // 服务器名称
	IP          string                              // ip
	Port        int                                 // port
	IPVersion   string                              // 服务器绑定的ip版本
	Routers     ziface.IMsgHandler                  // 多路由
	ConnManager ziface.IConnManager                 // 该server的连接管理器
	OnConnStart func(connection ziface.IConnection) // 创建连接之后调用
	OnConnStop  func(connection ziface.IConnection) // 销毁连接之前调用
}

func NewServer(name string) (server *Server) {
	defer func() {
		//补全 config对象
		utils.Config.TcpServer = server
	}()
	return &Server{
		Name:        name,
		IP:          utils.Config.Ip,
		Port:        utils.Config.Port,
		IPVersion:   utils.Config.IPVersion,
		Routers:     NewMsgHandler(),
		ConnManager: NewConnManager(),
		OnConnStart: nil,
		OnConnStop:  nil,
	}
}

func (s *Server) Start() {
	fmt.Println(utils.Config)
	fmt.Printf("[Start] Server Listenner at IP:%s, Port:%d\n", s.IP, s.Port)
	// 异步创建tcp监听
	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
		}
		// 2. 监听服务器的地址
		tcpListener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Printf("Listen %s err:%v\n", s.IPVersion, err)
		}
		fmt.Printf("[Start] Zinx success [%s] , Listenning....\n", s.Name)
		var connID uint32
		connID = 0
		// 3. 阻塞的等待客户端连接，处理客户端业务
		for true {
			// todo accept() 会阻塞程序
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			} else {
				fmt.Println("[server] accept conn")
			}
			//判断当前server的连接数是否大于最大连接限制
			if s.ConnManager.GetConnCount() > utils.Config.MaxConn {
				// todo 给客户端一个超过最大连接的错误提示
				conn.Close()
				fmt.Println("over MaxConnLimit")
				time.Sleep(time.Millisecond)
				continue
			}
			//  得到 connection 包装体
			connection := NewConnection(s, conn, connID, s.Routers)
			connID++
			// 同步到连接管理模块中
			s.ConnManager.AddConn(connection)
			// 启动连接
			go connection.Start()
		}
	}()
}

func (s *Server) Stop() {
	// todo  将一些服务器的资源、状态或者一些已经开辟的连接信息， 进行停止或回收
	s.ConnManager.Clear()
	fmt.Printf("[Stop] server:[%s]\n", s.Name)
}

func (s *Server) Serve() {
	// start 只负责 服务的启动，业务处理
	s.Start()

	// Serve 负责 阻塞
	// todo 做一些 启动服务之后的额外业务

	// 阻塞状态
	select {}

}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.Routers.Put(msgId, router)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(f func(connection ziface.IConnection)) {
	s.OnConnStart = f
}

func (s *Server) SetOnConnStop(f func(connection ziface.IConnection)) {
	s.OnConnStop = f
}

func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("----------->Call OnConnStart().......")
		s.OnConnStart(connection)
	}
}

func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("----------->Call OnConnStop().......")
		s.OnConnStop(connection)
	}
}
