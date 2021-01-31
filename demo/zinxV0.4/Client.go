package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		go connectTCPServer(i)
	}
	time.Sleep(time.Second * 60)

}

func connectTCPServer(i int) {
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("[client] connect  to server error: ", err)
	} else {
		fmt.Printf("[client:%d] connect server success\n", i)
	}
	//reader := bufio.NewReader(os.Stdin)
	defer conn.Close()
	data := make([]byte, 128)

	for true {
		//msgByte, err := reader.ReadBytes('\n')
		msgByte := []byte("jkdsknb看过了你的反馈你看")

		_, err = conn.Write(msgByte)
		if err != nil {
			fmt.Println("send msg to server error:", err)
		}

		read, err := conn.Read(data)
		if err != nil {
			fmt.Println("read msg from server error:", err)
		} else {
			fmt.Printf("[client:%d] receve :%s\n", i, string(data[0:read]))
		}
		time.Sleep(time.Second)
	}

}
