package test

import (
	"fmt"
	"net"
)

func handleSocks5(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Socks5 word")
	// 实现 SOCKS5 协议处理逻辑
	// 这里只是简单示例，实际应该根据协议规范来处理请求
}

func handleHTTP(conn net.Conn) {
	defer conn.Close()
	fmt.Println("handleHTTP hello word")
	// 实现 HTTP 代理处理逻辑
	// 这里只是简单示例，实际应该解析 HTTP 请求，转发到目标服务器，并返回响应
}

func MyProxy() {
	server := "0.0.0.0:6050"
	// 启动监听在本地的端口
	listener, err := net.Listen("tcp", server)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Proxy server listening on ", server)

	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// 根据请求协议类型分发连接到不同的处理函数
		go func(conn net.Conn) {
			// 读取请求的头部信息，以确定连接类型
			buffer := make([]byte, 256)
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Error reading request:", err)
				conn.Close()
				return
			}

			// 判断请求的类型
			if buffer[0] == 0x05 {
				// 如果是 SOCKS5 协议
				handleSocks5(conn)
			} else {
				// 其他情况认为是 HTTP 请求
				handleHTTP(conn)
			}
		}(conn)
	}
}
