package ziface

type IConnManager interface {
	//添加连接
	AddConn(conn IConnection)
	//删除连接
	RemoveConn(id uint32)
	//根据id获取conn
	GetConn(id uint32) (IConnection, error)
	//获取中连接数
	GetConnCount() int
	//清空所有连接
	Clear()
}
