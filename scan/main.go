package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 1; i <= 81; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("130.162.128.154:%d", j)
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
			fmt.Printf("%d open\n", j)
		}(i)
	}
	wg.Wait()
}

// Scan0 笨笨的单端口扫描
func Scan0(protocol string, address string) {
	_, err := net.Dial(protocol, address)
	if err == nil {
		fmt.Printf("Connection to %s [%s] succeeded! ", protocol, address)
	} else {
		fmt.Println(err)
	}
}

// Scan1 笨笨的多端口for循环扫描 从1-443端口
func Scan1() {
	fmt.Println(time.Now())
	for i := 1; i <= 443; i++ {
		dial, err := net.Dial("tcp", fmt.Sprintf("www.eber.vip:%d", i))
		if err == nil {
			// 监听到端口open后  关闭连接
			err = dial.Close()
			fmt.Printf("Discovered open port %d/tcp on www.eber.vip \n", i)
		}
	}
	fmt.Println(time.Now())
}

// Scan2 错误的协程示范
func Scan2() {
	fmt.Println(time.Now())
	for i := 1; i <= 443; i++ {
		go func(j int) {
			dial, err := net.Dial("tcp", fmt.Sprintf("www.eber.vip:%d", j))
			if err == nil {
				// 监听到端口open后  关闭连接
				err = dial.Close()
				fmt.Printf("Discovered open port %d/tcp on www.eber.vip \n", j)
			}
		}(i)
	}
	fmt.Println(time.Now())
}

func worker(ports chan int, wg *sync.WaitGroup) {
	for port := range ports {
		fmt.Println(port)
		wg.Done()
	}
}
