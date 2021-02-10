package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"zinx/src/utils"
	"zinx/src/ziface"
)

/**
封包、拆包的具体实现对象
*/

type DataPackage struct {
}

func NewDataPackage() *DataPackage {
	return &DataPackage{}
}

func (d *DataPackage) GetHeadLen() uint32 {
	//	len:unint32(4字节) + ID:uint32(4字节)
	return 8
}

func (d *DataPackage) Pack(msg ziface.IMessage) ([]byte, error) {
	//  创建一个 存放一个bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLength()); err != nil {
		return nil, err
	}

	// 将msgId写进dataBuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data写进dataBuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

/**
解析消息的头部，得到 消息内容的长度
*/
func (d *DataPackage) Unpack(datas []byte) (ziface.IMessage, error) {
	// 1. 读取 head （8个字节）
	//headByte:= bytes[0:8]
	//length := int32(headByte[0:4])

	msg := Message{}

	reader := bytes.NewReader(datas)
	if err := binary.Read(reader, binary.LittleEndian, &msg.Length); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断 length 是否已经超过最大限制
	if msg.Length >= utils.Config.MaxPackageSize {

		return nil, errors.New("the msg is too large :[" + strconv.Itoa(int(msg.Length)) + "]")
	}
	return &msg, nil
}
