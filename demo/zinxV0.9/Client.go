package main

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
	"zinx/src/znet"
)

func main() {
	group := sync.WaitGroup{}
	group.Add(1)
	// 测试 server的最大连接数限制是否有效（10）
	for i := 0; i < 15; i++ {
		go connectTCPServer(i)
	}
	group.Wait()
}

func connectTCPServer(k int) {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("[client] connect  to server error: ", err)
	} else {
		fmt.Printf("[client:%d] connect server success\n", k)
	}
	//reader := bufio.NewReader(os.Stdin)
	var i uint32
	go func() {
		for ; i < 1; i++ {
			//	封包
			dp := znet.NewDataPackage()
			i2 := rand.Intn(3)
			content := []byte("今天星期?" + strconv.Itoa(int(i2)))
			message := znet.NewMessage(uint32(i2), content)
			pack, err := dp.Pack(message)
			if err != nil {
				fmt.Println("pack error:", err)
			}
			_, err = conn.Write(pack)
			if err != nil {
				fmt.Println("write error:", err)
			}
			time.Sleep(time.Millisecond * 500)

		}
	}()
	go func() {
		for true {
			dataPackage := znet.NewDataPackage()
			headBytes := make([]byte, dataPackage.GetHeadLen())
			_, err := conn.Read(headBytes)
			if err != nil {
				fmt.Println("read error:", err)
			}

			msg, err := dataPackage.Unpack(headBytes)
			if err != nil {
				fmt.Println("Unpack error:", err)
			}
			bodyBytes := make([]byte, msg.GetMsgLength())
			_, err = conn.Read(bodyBytes)
			if err != nil {
				fmt.Println("Resd Body error:", err)
			} else {
				fmt.Printf("[Receive]:[%s]\n", string(bodyBytes))
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	// 防止当前goroutine 在创建两个goroutine后，运行结束，因此加上了select阻塞当前goroutine
	//select {}

}
