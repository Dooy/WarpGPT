package test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

func handleConnection(conn net.Conn) {

	defer conn.Close()
	request, err := http.ReadRequest(bufio.NewReader(conn)) // 读取请求
	if err != nil {
		fmt.Printf("Failed to read request: %s", err)
		return
	}
	//fmt.Printf("Proxying request to: %s", request.URL.String())
	log.Println(request.URL.String())
	client, err := net.Dial("tcp", request.Host) // 建立连接
	if err != nil {
		fmt.Printf("Failed to dial server: %s", err)
		return

	}

	defer client.Close()

	err = request.Write(client) // 替换Host头并发送请求
	if err != nil {
		fmt.Printf("Failed to write request: %s", err)
		return
	}

	response, err := http.ReadResponse(bufio.NewReader(client), request) // 读取响应

	if err != nil {
		fmt.Printf("Failed to read response: %s", err)
		return
	}
	defer response.Body.Close()
	for k, v := range response.Header { // 将响应头写回客户端

		for _, vv := range v {
			conn.Write([]byte(k + ": " + vv + ""))
		}
	}
	conn.Write([]byte(""))       // 写入响应行与响应头的分隔符
	io.Copy(conn, response.Body) // 将响应体写回客户端
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
		//go handleRequest(conn)
		go handleConnection(conn)
	}
}
