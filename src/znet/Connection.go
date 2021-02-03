package znet

import (
	"errors"
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
	fmt.Printf("conn[%d] Start().... \n", c.ID)
	defer c.Stop()
	// TLV： 读取head头部数据
	data := make([]byte, utils.Config.MaxPackageSize)
	for {
		// todo read() 会阻塞程序
		_, err := c.Conn.Read(data)
		if err != nil {
			fmt.Printf("read conn[%d] eror : %v\n", c.ID, err)
			continue
		}
		// TLV：拆包
		dataPackage := DataPackage{}
		message, err := dataPackage.Unpack(data)
		if err != nil {
			fmt.Println("[TLV] read  head get  body length error:", err)
			break
		}
		body := make([]byte, message.GetMsgLength())
		_, err = c.Conn.Read(body)
		if err != nil {
			fmt.Println("[TLV] read body error ：", err)
		}
		message.SetMsgData(body)
		req := NewRequest(c, message)

		// 执行路由中注册的处理方法
		go func(req ziface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)

	}

}

func (c *Connection) Stop() {
	fmt.Printf("conn[%d] Stop().... \n", c.ID)
	if c.Closed {
		return
	}
	c.Closed = true
	// 关闭连接
	_ = c.Conn.Close()
	// 回收资源 ，关闭管道
	close(c.ExitChan)
	fmt.Printf("conn[%d] Stoped \n", c.ID)
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

// 将数据 封装为msg，在将msg 封包得到最终要发送的TLV格式的[]byte
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.Closed {
		return errors.New("conn is closed when send msg")
	}
	msg := NewMessage(msgId, data)
	dp := DataPackage{}
	bytes, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("[Pack] err :", err)
		return errors.New("Pack data error")
	}
	_, err = c.Conn.Write(bytes)
	if err != nil {
		fmt.Println("[connection]: send data error:", err)
	}
	return err
}
