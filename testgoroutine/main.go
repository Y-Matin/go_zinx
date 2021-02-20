package main

import (
	"fmt"
	"sync"
)

/**
 目的：一次输出 cat、fish、dog 100次
使用channel ：有缓冲区、无缓冲区的有能实现。
有缓冲区的channel 包容性比较强。并发性能
无缓冲的channel需要额外小心逻辑。避免死锁。
	注意点：在dog的goroutine中，在最后一轮，不能忘catChan里面放。
          因为 cat这个goroutine在最后一轮结束，没有从catChan中取数据，因此在dog-goroutine中年放 catCHan里面放数据，会dog-goroutine阻塞程序，导致wg没有释放，阻塞主程序。
*/
var wg sync.WaitGroup

func main() {
	catChan := make(chan struct{}, 1)
	fishChan := make(chan struct{}, 1)
	dogChan := make(chan struct{}, 1)
	wg.Add(3)

	go cat(catChan, fishChan, &wg)
	go fish(fishChan, dogChan, &wg)
	go dog(dogChan, catChan, &wg)
	catChan <- struct{}{}

	wg.Wait()

}

func cat(catChan <-chan struct{}, fishChan chan<- struct{}, wg *sync.WaitGroup) {
	for count := 0; count < 100; count++ {
		<-catChan
		fmt.Println("cat")
		fishChan <- struct{}{}
	}
	wg.Done()
	return
}

func fish(fishChan <-chan struct{}, dogChan chan<- struct{}, wg *sync.WaitGroup) {

	for count := 0; count < 100; count++ {
		<-fishChan
		fmt.Println("fish")
		dogChan <- struct{}{}
	}
	wg.Done()
	return
}
func dog(dogChan <-chan struct{}, catChan chan<- struct{}, wg *sync.WaitGroup) {
	for count := 0; count < 100; count++ {
		<-dogChan
		fmt.Println("dog")
		if count != 99 {
			catChan <- struct{}{}
		}
	}
	wg.Done()
	return
}
