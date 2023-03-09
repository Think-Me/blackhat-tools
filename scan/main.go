package main

import (
	"fmt"
	"net"
)

// Scan0 笨笨的单端口扫描
func Scan0() {
	_, err := net.Dial("tcp", "www.eber.vip:80")
	if err == nil {
		fmt.Printf("Connection to www.eber.vip port :80 [tcp] succeeded! ")
	} else {
		fmt.Println(err)
	}
}

// Scan1 笨笨的多端口for循环扫描 从1-443端口
func Scan1() {
	for i := 1; i < 443; i++ {
		dial, err := net.Dial("tcp", fmt.Sprintf("www.eber.vip:%d", i))
		if err == nil {
			// 监听到端口open后  关闭连接
			err = dial.Close()
			fmt.Printf("Discovered open port %d/tcp on www.eber.vip \n", i)
		}
	}

}
