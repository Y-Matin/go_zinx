package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/src/ziface"
)

type Connection struct {
	// 当前连接所属的server
	TcpServer ziface.Iserver
	// 当前连接的套接字
	Conn *net.TCPConn
	// 连接的id
	ID uint32
	// 连接的状态
	Closed bool

	// 连接关闭 的标志 channel
	ExitChan chan bool

	// 写消息管道(无缓冲channel)
	MsgChan chan []byte

	//多路由
	Routers ziface.IMsgHandler
}

// 初始化连接模块的方法
func NewConnection(server ziface.Iserver, conn *net.TCPConn, id uint32, routers ziface.IMsgHandler) *Connection {
	return &Connection{
		TcpServer: server,
		Conn:      conn,
		ID:        id,
		Closed:    false,
		Routers:   routers,
		ExitChan:  make(chan bool, 1),
		MsgChan:   make(chan []byte),
	}
}

func (c *Connection) Start() {
	fmt.Printf("conn[%d] Start().... \n", c.ID)
	go c.startReader()
	go c.startWriter()

}

func (c *Connection) startReader() {
	fmt.Printf("conn[%d] Reader Goroutine is running\n", c.ID)
	defer fmt.Printf("conn[%d] Reader Goroutine is closed\n", c.ID)
	defer c.Stop()
	// TLV： 读取head头部数据
	dp := NewDataPackage()
	data := make([]byte, dp.GetHeadLen())
	for {
		// todo read() 会阻塞程序
		_, err := c.Conn.Read(data)
		if err != nil {
			fmt.Printf("read conn[%d] eror : %v\n", c.ID, err)
			break
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
			break
		}
		message.SetMsgData(body)
		req := NewRequest(c, message)
		// 执行路由中注册的处理方法
		//将任务 放入到任务池中，等待被worker执行
		c.Routers.AddTask(req)
	}
}

func (c *Connection) startWriter() {
	fmt.Printf("conn[%d] Writer Goroutine is running\n", c.ID)
	defer fmt.Printf("conn[%d] Write Goroutine is closed \n", c.ID)
	for true {
		select {
		case bytes := <-c.MsgChan:
			// 有数据要写回客户端
			_, err := c.Conn.Write(bytes)
			if err != nil {
				fmt.Println("Write() error:", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}

	}

}

func (c *Connection) Stop() {
	fmt.Printf("conn[%d] Stop().... \n", c.ID)
	if c.Closed {
		return
	}
	c.Closed = true
	c.ExitChan <- true
	// 关闭连接
	_ = c.Conn.Close()
	// 同步移除连接管理器中的连接
	c.TcpServer.GetConnManager().RemoveConn(c.ID)
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
	// 往写通道里面放入数据
	c.MsgChan <- bytes
	return nil
}
