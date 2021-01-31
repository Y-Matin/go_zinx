package znet

import (
	"fmt"
	"net"
	"zinx/src/utils"
	"zinx/src/ziface"
)

type Connection struct {
	// 当前连接的套接字
	Conn *net.TCPConn
	// 连接的id
	ID uint32
	// 连接的状态
	Closed bool
	// 当前连接的处理方法
	HandleAPI ziface.HandleFunc

	// 连接关闭 的标志 channel
	ExitChan chan bool
	// 当前连接注册的router 方法
	Router ziface.IRouter
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, id uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ID:       id,
		Closed:   false,
		Router:   router,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn[%d] Start().... \n", c.ID)
	defer c.Stop()
	//  todo 启动 当前连接的写业务
	data := make([]byte, utils.Config.MaxPackageSize)
	for {
		// todo read() 会阻塞程序
		length, err := c.Conn.Read(data)
		if err != nil {
			fmt.Printf("read Conn[%d] eror : %v\n", c.ID, err)
			continue
		}

		// 使用该Conn绑定的路由方法
		req := Request{
			conn: c,
			data: data[0:length],
		}
		// 执行路由中注册的处理方法
		go func(req ziface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&req)

	}

}

func (c *Connection) Stop() {
	fmt.Printf("Conn[%d] Stop().... \n", c.ID)
	if c.Closed {
		return
	}
	c.Closed = true
	// 关闭连接
	_ = c.Conn.Close()
	// 回收资源 ，关闭管道
	close(c.ExitChan)
	fmt.Printf("Conn[%d] Stoped \n", c.ID)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ID
}

func (c *Connection) IsClosed() bool {
	return c.Closed
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(bytes []byte) error {
	_, err := c.Conn.Write(bytes)
	if err != nil {
		fmt.Println("[connection]: send data error:", err)
	}
	return err
}
