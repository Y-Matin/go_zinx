package znet

import (
	"zinx/src/ziface"
)

type Request struct {
	conn    ziface.IConnection
	message ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.message.GetMsgData()
}

func (r *Request) GetMsgId() uint32 {
	return r.message.GetMsgId()
}

func NewRequest(conn ziface.IConnection, message ziface.IMessage) *Request {
	return &Request{
		conn:    conn,
		message: message,
	}
}
