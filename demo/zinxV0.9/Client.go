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
	// main Goroutine一致阻塞
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
		defer fmt.Printf("conn:[%d] closed\n", k)
		for ; i < 100; i++ {
			//	封包
			dp := znet.NewDataPackage()
			i2 := rand.Intn(3)
			content := []byte("今天星期?" + strconv.Itoa(int(i2)))
			message := znet.NewMessage(uint32(i2), content)
			pack, err := dp.Pack(message)
			if err != nil {
				fmt.Println("pack error:", err)
				return
			}
			_, err = conn.Write(pack)
			if err != nil {
				fmt.Println("write error:", err)
				return
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	go func() {
		defer fmt.Printf("conn:[%d] closed\n", k)
		for true {
			dataPackage := znet.NewDataPackage()
			headBytes := make([]byte, dataPackage.GetHeadLen())
			_, err := conn.Read(headBytes)
			if err != nil {
				fmt.Println("read error:", err)
				return
			}

			msg, err := dataPackage.Unpack(headBytes)
			if err != nil {
				fmt.Println("Unpack error:", err)
			}
			bodyBytes := make([]byte, msg.GetMsgLength())
			_, err = conn.Read(bodyBytes)
			if err != nil {
				fmt.Println("Resd Body error:", err)
				return
			} else {
				fmt.Printf("[Receive]:[%s]\n", string(bodyBytes))
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
}
