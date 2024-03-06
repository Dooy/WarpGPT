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

	// 实现 HTTP 代理处理逻辑
	// 这里只是简单示例，实际应该解析 HTTP 请求，转发到目标服务器，并返回响应

	// 读取客户端的请求
	requestBuffer := make([]byte, 14096)
	n, err := conn.Read(requestBuffer)
	if err != nil {
		fmt.Println("Error reading HTTP request:", err)
		return
	}

	// 解析HTTP请求
	httpRequest := string(requestBuffer[:n])
	fmt.Println("Received HTTP request:\n", httpRequest)

	// // 这里可以添加逻辑，解析请求并决定要连接的目标服务器

	// // 与目标服务器建立连接
	// targetConn, err := net.Dial("tcp", "target-server.com:80")
	// if err != nil {
	// 	fmt.Println("Error connecting to target server:", err)
	// 	return
	// }
	// defer targetConn.Close()

	// // 转发HTTP请求到目标服务器
	// _, err = targetConn.Write([]byte(httpRequest))
	// if err != nil {
	// 	fmt.Println("Error forwarding HTTP request:", err)
	// 	return
	// }

	// // 从目标服务器读取响应并转发给客户端
	// responseBuffer := make([]byte, 4096)
	// n, err = targetConn.Read(responseBuffer)
	// if err != nil {
	// 	fmt.Println("Error reading response from target server:", err)
	// 	return
	// }

	// // 将目标服务器的响应转发给客户端
	// _, err = conn.Write(responseBuffer[:n])
	// if err != nil {
	// 	fmt.Println("Error forwarding response to client:", err)
	// 	return
	// }

	fmt.Println("HTTP request handled successfully.")

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
