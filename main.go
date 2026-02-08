package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

var (
	listenAddr string
	targetAddr string
)

func init() {
	flag.StringVar(&listenAddr, "listen", ":8080", "Local listening address")
	flag.StringVar(&targetAddr, "target", "", "Remote target address (e.g. 1.1.1.1:80)")
	flag.Parse()
}

func main() {
	if targetAddr == "" {
		fmt.Println("Error: Target address is required!")
		flag.Usage()
		return
	}

	// 启动监听
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to listen on %s: %v", listenAddr, err)
	}
	log.Printf("[INFO] G-Tunnel Started! Listening on %s, Forwarding to %s", listenAddr, targetAddr)

	// 死循环接受连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] Accept failed: %v", err)
			continue
		}

		// 每来一个连接，开启一个 Goroutine 处理
		go handleConnection(clientConn)
	}
}

// core logic
func handleConnection(clientConn net.Conn) {
	// 函数结束时，关闭连接，防止资源泄露
	defer clientConn.Close()

	clientIP := clientConn.RemoteAddr().String()
	log.Printf("[INFO] New Connection from %s", clientIP)

	// 连接远程目标 (打通隧道)
	// 设置 5秒 超时，防止连不上一直卡住
	targetConn, err := net.DialTimeout("tcp", targetAddr, 5*time.Second)
	if err != nil {
		log.Printf("[ERROR] Failed to connect to target %s: %v", targetAddr, err)
		return
	}
	defer targetConn.Close()

	log.Printf("[DEBUG] Tunnel Established: %s <-> %s", clientIP, targetAddr)

	// 数据交换 (双向拷贝)
	// 需要两个 Goroutine：一个读左写右，一个读右写左
	var wg sync.WaitGroup
	wg.Add(2)

	// 方向 1: Client -> Target
	go func() {
		defer wg.Done()
		// io.Copy 是系统级优化调用 (splice/sendfile)，性能极高
		written, err := io.Copy(targetConn, clientConn)
		if err != nil {
			// 连接断开是正常的，不需要打印 ERROR
			log.Printf("[DEBUG] Client -> Target closed: %v", err)
		}
		log.Printf("[TRACE] Transferred %d bytes (Client -> Target)", written)
		// 一方断开，强制关闭另一方的写入
		targetConn.(*net.TCPConn).CloseWrite()
	}()

	// 方向 2: Target -> Client
	go func() {
		defer wg.Done()
		written, err := io.Copy(clientConn, targetConn)
		if err != nil {
			log.Printf("[DEBUG] Target -> Client closed: %v", err)
		}
		log.Printf("[TRACE] Transferred %d bytes (Target -> Client)", written)
		clientConn.(*net.TCPConn).CloseWrite()
	}()

	// 等待双方都传输完毕
	wg.Wait()
	log.Printf("[INFO] Connection Closed: %s", clientIP)
}
