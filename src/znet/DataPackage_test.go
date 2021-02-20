package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

/**
负责 测试消息的封包、拆包的单元测试
*/
func TestDataPackage_Pack(t *testing.T) {
	/**
	模拟服务器  tcpconnection
	*/
	listen, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("get addr err :", err)
	}

	go func() {
		conn, err1 := listen.Accept()
		if err1 != nil {
			fmt.Println("accept error:", err1)
		}
		//  todo  此处 新建了goroutine去单独处理conn，是因为内部 有一个for循环，一直读取conn的数据，导致上面的Accept() 无法执行。造成该server只能处理一个client请求。所以，创建了goroutine去单独处理connection
		go func() {
			for true {
				// readLength head
				d := NewDataPackage()
				data := make([]byte, d.GetHeadLen())
				readLength, err := conn.Read(data)
				if err != nil {
					fmt.Println("readLength error:", err)
				}
				// 拆包 ；将 字节切片解析为message对象得到具体的消息内容
				dataPackage := NewDataPackage()
				message, err := dataPackage.Unpack(data[0:readLength])
				if err != nil {
					fmt.Println("Unpack error:", err)
				}
				m := message.(*Message)
				m.Data = make([]byte, m.Length)
				if _, err := conn.Read(m.Data); err != nil {
					fmt.Println("read body error:", err)
				}
				fmt.Println("[Server] receive:", string(m.Data))
			}

		}()
	}()

	time.Sleep(time.Second)

	/**
	模拟 客户端
	*/

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("conn server error:", err)
	}
	message1 := Message{
		Id:     0,
		Length: 0,
		Data:   []byte("白日依山尽"),
	}
	message2 := Message{
		Id:     0,
		Length: 0,
		Data:   []byte("黄河入海流"),
	}
	d2 := NewDataPackage()
	message1.Id = uint32(1)
	message1.Length = uint32(len(message1.Data))
	data, err := d2.Pack(&message1)
	message2.Id = uint32(2)
	message2.Length = uint32(len(message2.Data))
	data1, err := d2.Pack(&message2)
	// 模拟粘包的情况  把两个 消息 []byte 合在一起写入conn发送给server
	data = append(data, data1...)

	if err != nil {
		fmt.Println("Pack error:", err)
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Write error:", err)
	} else {
		fmt.Println("[Client] send:", string(data[8:]))
	}

	time.Sleep(time.Second)

}
