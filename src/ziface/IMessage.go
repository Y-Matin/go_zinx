package ziface

/**
对 request中的data进行封装为 message
*/
type IMessage interface {
	// 得到消息id
	GetMsgId() uint32
	// 得到消息长度
	GetMsgLength() uint32
	// 得到消息内容
	GetMsgData() []byte

	SetMsgId(id uint32)

	SetMsgLength(length uint32)

	SetMsgData(data []byte)
}
