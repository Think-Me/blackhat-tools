package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"testing"
)

// 这里定义DNS协议数据的结构体 另外:dns协议规定两个字节的都用16位无符号整型表示
type DNSData struct {
	TransactionId uint16 // 属于header 客户端随机生成的一个无符号整数，范围是0~2^16(0~65536)。在响应头里面也会返回这个值作用是校验。如果值不相等，丢弃响应内容。
	Flags         uint16 // 属于header 16位标志位，如下所示QR、opcode、AA、TC、RD、RA、zero、rcode。一般查询flags为 00000001 00000000
	// QR(1bit) 查询应答标志，0表示这是查询报文，1表示这是应答报文。
	// opcode(4bit) 查询应答类型，0表示标准查询，1表示反向查询，2表示请求服务器状态。
	// AA(1bit) 表示权威回答( authoritative answer )，意味着当前查询结果是由域名的权威服务器给出的，仅由应答报文使用。
	// TC(1bit) 位表示截断( truncated )，使用 UDP 时，如果应答超过 512 字节，只返回前 512 个字节，仅当DNS报文使用UDP服务时使用。DNS 协议使用UDP服务，但也明确了 『当 DNS 查询被截断时，应该使用 TCP 协议进行重试』 这一规范。
	// RD(1bit) 表示递归查询标志 ( recursion desired )，在请求中设置，并在应答中返回。该位为 1 时，服务器必须处理这个请求：如果服务器没有授权回答，它必须替客户端请求其他 DNS 服务器，这也是所谓的 递归查询； 该位为 0 时，如果服务器没有授权回答，它就返回一个能够处理该查询的服务器列表给客户端，由客户端自己进行 迭代查询。
	// RA(1bit) 位表示可递归 ( recursion available )，如果服务器支持递归查询，就会在应答中设置该位，以告知客户端。仅由应答报文使用。
	// zero(3bit) 这三位未使用，固定为0。
	// rcode(4bit) 表示返回码（reply code），用来返回应答状态，常用返回码：0表示无错误，2表示格式错误，3表示域名不存在。
	Queries []dnsQuestion // 本身不属于header 表示查询请求记录内容数据，他的组数长度值为Questions(属于header)
	// 下面几个是应答记录中的内容，只有在应答消息中才会出现
	Answers       []dnsAnswer // 本身不属于header 应答资源记录数据（answer resource record, answer RR）此项只在DNS应答消息中存在，他的数组长度值为AnswerRRs(属于header)
	AuthorityRRs  uint16      // 属于header 授权资源记录数量（authority resource record, authority RR）此项只在DNS应答消息中存在
	AdditionalRRs uint16      // 属于header 附加资源记录数量（additional resource record, additional RR）此项只在DNS应答消息中存在
}

func (dDNSData *DNSData) SetFlag(QR uint16, Opcode uint16, AA uint16, TC uint16, RD uint16, RA uint16, Rcode uint16) {
	dDNSData.Flags = QR<<15 + Opcode<<11 + AA<<10 + TC<<9 + RD<<8 + RA<<7 + Rcode
}

// 定义请求数据的结构体
type dnsQuestion struct {
	QueriesName  string `net:"domain-name"` // 要查询的域名
	QueriesType  uint16 // 查询类型 1:A 2:NS 5:CNAME 6:SOA 12:PTR 15:MX 16:TXT 28:AAAA
	QueriesClass uint16 // 查询类 1:IN 2:CS 3:CH 4:HS 通常为1表示为TCP/IP互联网地址
}

// 定义响应数据的结构体
type dnsAnswer struct {
	AnswerName       uint16 // 同dnsQuestion.QueriesName 要查询的域名
	AnswerType       uint16 // 应答记录的类型 1:A 2:NS 5:CNAME 6:SOA 12:PTR 15:MX 16:TXT 28:AAAA
	AnswerClass      uint16 // 同dnsQuestion.QueriesClass
	AnswerTTL        uint32 // 32位生存时间(有效期)单位是秒
	AnswerDataLength uint16 // 16位无符号整数，表示应答资源记录中数据的长度
	AnswerCNAME      string `net:"domain-name"` // 别名
}

// 写入DNS协议头部数据
func (dDNSData *DNSData) WriteHeader() []byte {
	// DNS协议定义Header为12个字节的固定长度
	bs := make([]byte, 12)
	binary.BigEndian.PutUint16(bs[0:2], dDNSData.TransactionId)
	binary.BigEndian.PutUint16(bs[2:4], dDNSData.Flags)
	binary.BigEndian.PutUint16(bs[4:6], uint16(len(dDNSData.Queries))) // Queries的数组长度值为Questions
	binary.BigEndian.PutUint16(bs[6:8], uint16(len(dDNSData.Answers))) // Answers的数组长度值为AnswerRRs
	binary.BigEndian.PutUint16(bs[8:10], dDNSData.AuthorityRRs)
	binary.BigEndian.PutUint16(bs[10:12], dDNSData.AdditionalRRs)
	// 填充Question数据，要将域名进行转换。用.分割域名字符串
	ds := strings.Split(dDNSData.Queries[0].QueriesName, ".")
	// 循环遍历域名的每一部分，将其长度和内容写入到字节切片中。例如：list.eber.vip，写入的内容为4list4eber3vip0 末尾用0来表示结束。
	for _, d := range ds {
		bs = append(bs, byte(len(d)))
		bs = append(bs, []byte(d)...)
	}
	bs = append(bs, 0)
	// 添加查询类型和分类
	temp := make([]byte, 2)
	binary.BigEndian.PutUint16(temp, dDNSData.Queries[0].QueriesType)
	bs = append(bs, temp...)
	binary.BigEndian.PutUint16(temp, dDNSData.Queries[0].QueriesClass)
	bs = append(bs, temp...)
	return bs
}

// 测试发送DNS请求
func TestDNSDemo1(t *testing.T) {
	dnsServer := "114.114.114.114:53"
	dnsProtocol := "udp"
	dnsType := uint16(1)
	dnsClass := uint16(1)
	udpAddr, err := net.ResolveUDPAddr(dnsProtocol, dnsServer)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn, err := net.DialUDP(dnsProtocol, nil, udpAddr)

	question := dnsQuestion{"dns.eber.vip", dnsType, dnsClass}
	out := DNSData{}
	out.TransactionId = 2015
	out.SetFlag(0, 0, 0, 0, 1, 0, 0)
	out.Queries = append(out.Queries, question)
	header := out.WriteHeader()
	fmt.Println(header)
	_, err = conn.Write(header)
	var buf []byte
	buf = make([]byte, 512)
	n, err := conn.Read(buf[0:])
	fmt.Println(buf[0:n])
}
