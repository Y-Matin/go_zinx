package main

import (
	"fmt"
	"time"
)

/**
验证 爷goroutine阻塞，父goroutine结束,由父创建的子goroutine是否结束
结果： 子goroutine不会结束
*/
func main() {
	go func() {
		defer fmt.Println("12312321312")
		go func() {
			for true {
				fmt.Println("go No.1")
				time.Sleep(time.Millisecond * 500)

			}
		}()
		go func() {
			for {
				fmt.Println("go No.2")
				time.Sleep(time.Millisecond * 500)
			}

		}()
	}()

	select {}
}
