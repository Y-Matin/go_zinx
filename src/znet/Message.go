package znet

import "zinx/src/ziface"

type Message struct {
	Id     uint32 // id
	Length uint32 // length
	Data   []byte // data
}

func NewMessage(id uint32, data []byte) ziface.IMessage {
	return &Message{
		Id:     id,
		Length: uint32(len(data)),
		Data:   data,
	}
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLength() uint32 {
	return m.Length
}

func (m *Message) GetMsgData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetMsgLength(length uint32) {
	m.Length = length
}

func (m *Message) SetMsgData(data []byte) {
	m.Data = data
}
