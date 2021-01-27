package znet

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IP        string
	Port      int
	IPVersion string
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IP:        "0.0.0.0",
		Port:      8999,
		IPVersion: "tcp4",
	}
}

// 业务处理方法
func CallbackMsgToClient(conn *net.TCPConn, data []byte, length int) error {
	fmt.Println("[Conn Handle] CallbackMsgToClient()... ")
	_, err := conn.Write(data[0:length])
	if err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallbackMsgToClient() error")
	}
	return nil

}

func (s *Server) Start() {
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
			fmt.Print("test")
			//  得到 connection 包装体
			connection := NewConnection(conn, connID, CallbackMsgToClient)
			connID++
			// 启动连接
			go connection.Start()
		}
	}()

}

func (s *Server) Stop() {
	// todo  将一些服务器的资源、状态或者一些已经开辟的连接信息， 进行停止或回收

}

func (s *Server) Serve() {
	// start 只负责 服务的启动，业务处理
	s.Start()

	// Serve 负责 阻塞
	// todo 做一些 启动服务之后的额外业务

	// 阻塞状态
	select {}

}
