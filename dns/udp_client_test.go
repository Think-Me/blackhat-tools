package dns

import (
	"fmt"
	"net"
	"testing"
)

func TestDNSClient(t *testing.T) {

	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9888,
	})
	if err != nil {
		fmt.Println("net.DialUDP error : ", err)
		return
	}
	defer socket.Close()
	sendData := []byte("hello，I am client!!")
	println(socket.LocalAddr())
	// 发送数据
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("socket.Write error : ", err)
		return
	}
	data := make([]byte, 2048)
	// 接收数据
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("socket.ReadFromUDP error : ", err)
		return
	}
	fmt.Printf("data == %v  , addr == %v , count == %v\n", string(data[:n]), remoteAddr, n)
}
