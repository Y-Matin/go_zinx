package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
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
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, id uint32, callbackApi ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:      conn,
		ID:        id,
		Closed:    false,
		HandleAPI: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
}

func (c *Connection) Start() {
	fmt.Printf("Conn[%d] Start().... \n", c.ID)
	defer c.Stop()
	//  todo 启动 当前连接的写业务
	data := make([]byte, 512)
	for {
		length, err := c.Conn.Read(data)
		if err != nil {
			fmt.Printf("read Conn[%d] eror : %v\n", c.ID, err)
			continue
		}
		//  调用绑定的业务处理方法
		err = c.HandleAPI(c.Conn, data, length)
		if err != nil {
			fmt.Printf("Conn[%d] handlerAPI exec error:\n", c.ID, err)
			break
		}

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
