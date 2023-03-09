package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestScan0(t *testing.T) {
	Scan0("tcp", "130.162.128.154:80")
}
func TestScan1(t *testing.T) {
	Scan1()
}
func TestScan2(t *testing.T) {
	Scan2()
}

// 后面都用test写，方便运行

// goroutine并行操作
func TestScan3(t *testing.T) {
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

func TestScan4(t *testing.T) {
	ports := make(chan int, 10)
	var wg sync.WaitGroup
	// 获取非空长度
	fmt.Println(len(ports))
	// 管道的左指针到末尾的空间容量
	fmt.Println(cap(ports))
	println("---------------")
	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg)
	}
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		ports <- i
	}
	wg.Wait()
	close(ports)
}
