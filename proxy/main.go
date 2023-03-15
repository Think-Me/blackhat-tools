package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":81")
	if err != nil {
		log.Fatalln(err)
	}
	accept, err := listen.Accept()
	if err != nil {
		log.Fatalln(err)
	}
	go handle(accept)
}

func handle(dst net.Conn) {
	src, err := net.Dial("tcp", "10.0.1.5:15672")
	if err != nil {
		log.Fatalln("连接源站失败")
	}
	defer src.Close()
	go func() {
		// 将源站输出复制dst
		if _, err := io.Copy(src, dst); err != nil {
			log.Fatalln(err)
		}
	}()
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalln(err)
	}
}
