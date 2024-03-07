package test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
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

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()

	clientReader := bufio.NewReader(client)

	method, requestAddress, protocol, headers, headerLines, err := decodeHeader(clientReader)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("method:%s", method)
	nFirstLine := method + " " + requestAddress + " " + protocol
	var serverAddress, oldHost string
	if method == "CONNECT" {
		serverAddress = requestAddress
	} else {
		hostPortURL, err := url.Parse(requestAddress)
		if err != nil {
			log.Println("url解析错误: " + nFirstLine)
			log.Println(err)
			return
		}
		oldHost = hostPortURL.Host
		if !strings.Contains(oldHost, ":") {
			serverAddress = oldHost + ":80"
		} else {
			serverAddress = oldHost
		}
	}

	log.Println(nFirstLine + " " + serverAddress + "\n")
	server, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	if method == "CONNECT" {
		//fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
		go io.Copy(server, clientReader)
	} else {
		needDecodeHeader := false
		go func() {
			for {
				if needDecodeHeader {
					method, requestAddress, protocol, headers, headerLines, err = decodeHeader(clientReader)
					if err != nil {
						log.Println(err)
						return
					}
				} else {
					needDecodeHeader = true
				}
				requestPath := append(strings.Split(requestAddress, oldHost), "/")[1]
				server.Write([]byte(method + " " + requestPath + " " + protocol + "\r\n"))
				for _, line := range headerLines {
					server.Write([]byte(line))
				}
				server.Write([]byte("\r\n"))

				length64, err := strconv.ParseInt(headers["content-length"], 10, 64)
				if err == nil {
					if length64 == -1 {
						io.Copy(server, clientReader)
						return
					}
					limitedReader := io.LimitReader(clientReader, length64)
					io.Copy(server, limitedReader)
					limitedReader = io.LimitReader(clientReader, 2)
					io.Copy(server, limitedReader)
				}
			}
		}()
	}
	io.Copy(client, server)
}

func decodeHeader(render *bufio.Reader) (string, string, string, map[string]string, []string, error) {
	var method, requestAddress, protocol string
	var headers = map[string]string{}
	var headerLines = []string{}
	lineData, err := render.ReadBytes('\n')
	if err != nil {
		return method, requestAddress, protocol, headers, headerLines, err
	}

	line := string(lineData)
	fmt.Sscanf(line, "%s%s%s", &method, &requestAddress, &protocol)
	if line != method+" "+requestAddress+" "+protocol+"\r\n" {
		log.Println("解析错误: " + line)
	}
	for {
		lineData, err := render.ReadBytes('\n')
		if err != nil {
			return method, requestAddress, protocol, headers, headerLines, err
		}
		if len(lineData) == 2 {
			break
		}
		line := string(lineData)
		index := strings.Index(line, ":")
		keyLower := strings.ToLower(strings.Trim(line[:index], "\r\n "))
		value := line[index+1:]
		if strings.HasPrefix(keyLower, "proxy-") {
			log.Println(line)
		}
		headers[keyLower] = strings.Trim(value, "\r\n ")
		if keyLower == "proxy-connection" {
			headerLines = append(headerLines, "Connection:"+value)
		} else {
			headerLines = append(headerLines, line)
		}
	}
	return method, requestAddress, protocol, headers, headerLines, err
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
		//go handleConnection(conn)
		go handleClientRequest(conn)
	}
}
