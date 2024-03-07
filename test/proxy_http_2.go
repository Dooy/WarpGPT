package test

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	// // 复制请求头
	for {
		line, err := reader.ReadString('\n')
		if err != nil || line == "\r\n" {
			//if err != nil {
			break
		}
		log.Println(line)
	}

	// // 读取请求行

	// requestLine, err := reader.ReadString('\n')
	// if err != nil {
	// 	log.Println("Error reading request:", err)
	// 	return
	// }
	// parts := strings.Fields(requestLine)
	// if len(parts) != 3 {
	// 	log.Println("Malformed request line:", requestLine)
	// 	return
	// }
	// method, target := parts[0], parts[1]

	// // 解析目标地址
	// targetURL, err := url.Parse(target)
	// log.Println("request:", method, " target=", target)
	// if err != nil {
	// 	log.Println("Error parsing target URL:", err)
	// 	return
	// }

	// // 创建与目标服务器的连接
	// targetConn, err := net.Dial("tcp", targetURL.Host)
	// if err != nil {
	// 	log.Println("Error connecting to target:", err)
	// 	return
	// }
	// defer targetConn.Close()

	// // 如果是CONNECT请求，建立隧道
	// if method == "CONNECT" {
	// 	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n\r\n")
	// } else {
	// 	// 将请求发送到目标服务器
	// 	fmt.Fprint(targetConn, method, " ", targetURL.RequestURI(), " HTTP/1.1\r\n")
	// }

	// // 复制请求头
	// for {
	// 	line, err := reader.ReadString('\n')
	// 	if err != nil || line == "\r\n" {
	// 		break
	// 	}
	// 	if method != "CONNECT" {
	// 		fmt.Fprint(targetConn, line)
	// 	}
	// }

	// // 将目标服务器的响应返回给客户端
	// go func() {
	// 	io.Copy(conn, targetConn)
	// }()
	// io.Copy(targetConn, conn)
}

func StartHttpV2(server string) {
	listener, err := net.Listen("tcp", server)
	if err != nil {
		log.Fatal("Error starting proxy server:", err)
	}
	defer listener.Close()

	fmt.Println("Proxy server listening on port  ", server)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn)
	}
}
