package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/src/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	conLock     sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}
func (c *ConnManager) AddConn(conn ziface.IConnection) {
	//写操作，加写锁
	c.conLock.Lock()
	defer c.conLock.Unlock()
	c.connections[conn.GetConnID()] = conn
	fmt.Printf("Add ConnId:[%d] to ConnManager success!\n", conn.GetConnID())
}

func (c *ConnManager) RemoveConn(id uint32) {
	//写操作，加写锁
	c.conLock.Lock()
	defer c.conLock.Unlock()
	delete(c.connections, id)
	fmt.Printf("Remove ConnId:[%d] to ConnManager success!\n", id)
}

func (c *ConnManager) GetConn(id uint32) (ziface.IConnection, error) {
	//	读操作，加读锁
	c.conLock.RLock()
	defer c.conLock.RUnlock()
	connection, ok := c.connections[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("GetConn failed, connId:[%d] is not exist", id))
	}
	return connection, nil
}

func (c *ConnManager) GetConnCount() int {
	return len(c.connections)
}

func (c *ConnManager) Clear() {
	//	清除操作，加写锁
	c.conLock.Lock()
	defer c.conLock.Unlock()
	for k, conn := range c.connections {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, k)
	}
	fmt.Println("Clear Connections success!")
}
