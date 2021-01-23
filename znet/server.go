package znet

import (
	"fmt"
	"net"
)

type Server struct {
	Name string
	IP string
	Port int
	IPVersion string
}

func NewServer(name string) *Server {
	return &Server{
		Name:      name,
		IP:        "127.0.0.1",
		Port:      8099,
		IPVersion: "tcp4",
	}
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s, Port:%d\n",s.IP,s.Port)

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
		// 3. 阻塞的等待客户端连接，处理客户端业务
		for true {
			conn, err := tcpListener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			fmt.Printf("%v", conn)
			fmt.Print("test")
			// 简单业务：回显最大512字节的数据
			go func (conn *net.TCPConn) {
				data:= make([]byte,512)
				for true {
					read, err := conn.Read(data)
					if err != nil {
						fmt.Println("TCP : read data error:",err)
						continue
					}
					if read<=0 {
						break
					}
					if _, err := conn.Write(data[0:read]); err!= nil {
						fmt.Println("TCP : write data error:",err)
						continue
					}
				}
			}(conn)
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
	select {
	}


}





