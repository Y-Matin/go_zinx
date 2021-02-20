package ziface

/**
封包、拆包 模块
直接面向TCP连接的数据流，用于处理TCP粘包问题
*/
type IDataPackage interface {
	// 得到数据包的长度
	GetHeadLen() uint32
	// 封包
	Pack(msg IMessage) ([]byte, error)

	// 拆包
	Unpack([]byte) (IMessage, error)
}
