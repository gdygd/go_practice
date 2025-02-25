package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

func main() {
	// 소켓 생성
	listenFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal("socket creation failed:", err)
	}
	defer syscall.Close(listenFd)

	// 소켓에 bind
	addr := syscall.SockaddrInet4{Port: 23000}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	err = syscall.Bind(listenFd, &addr)
	if err != nil {
		log.Fatal("bind failed:", err)
	}

	// 소켓을 수신 대기 상태로 변경
	err = syscall.Listen(listenFd, syscall.SOMAXCONN)
	if err != nil {
		log.Fatal("listen failed:", err)
	}

	// socket option
	//----------------------------------------------------------------------------------------------
	errOpt := syscall.SetsockoptInt(listenFd, syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1)
	if errOpt != nil {
		fmt.Printf("SetsockoptInt keep Alive %v\n", errOpt)
	}

	errOpt = syscall.SetsockoptInt(listenFd, syscall.SOL_SOCKET, 0xf, 1)
	if errOpt != nil {
		fmt.Printf("SetsockoptInt SetReusedPort %v\n", errOpt)
	}

	errOpt = syscall.SetsockoptInt(listenFd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	if errOpt != nil {
		fmt.Printf("SetsockoptInt SetReusedAddr %v\n", errOpt)
	}

	errOpt = syscall.SetsockoptInt(listenFd, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
	if errOpt != nil {
		fmt.Printf("SetsockoptInt SetNodelay %v\n", errOpt)
	}
	//----------------------------------------------------------------------------------------------

	// epoll 생성
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal("epoll create failed:", err)
	}
	defer syscall.Close(epollFd)

	// epoll 이벤트 등록
	event := syscall.EpollEvent{
		Events: syscall.EPOLLIN,
		Fd:     int32(listenFd),
	}
	err = syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, listenFd, &event)
	if err != nil {
		log.Fatal("epoll ctl add failed:", err)
	}

	// epoll 이벤트 처리 루프
	events := make([]syscall.EpollEvent, 10)
	for {
		fmt.Printf("Epoll wait(1) \n")
		n, err := syscall.EpollWait(epollFd, events, 100)
		if err != nil {
			log.Fatal("epoll wait failed:", err)
		}

		fmt.Printf("Epoll wait(2) N : %d\n", n)

		for i := 0; i < n; i++ {
			if int(events[i].Fd) == listenFd {
				// 새로운 연결 요청이 들어옴
				connFd, _, err := syscall.Accept(listenFd)
				if err != nil {
					log.Println("accept failed:", err)
					continue
				}

				// 연결된 클라이언트 소켓도 epoll 이벤트 등록
				event := syscall.EpollEvent{
					Events: syscall.EPOLLIN,
					Fd:     int32(connFd),
				}
				err = syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_ADD, connFd, &event)
				if err != nil {
					log.Println("epoll ctl add failed:", err)
					continue
				}

				log.Println("New connection:", connFd)
			} else {
				// 클라이언트로부터 데이터 수신
				connFd := int(events[i].Fd)
				buffer := make([]byte, 1024)
				n, err := syscall.Read(connFd, buffer)
				if err != nil {
					log.Println("read failed:", err)
					syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_DEL, connFd, nil)
					syscall.Close(connFd)
					continue
				}

				data := buffer[:n]
				log.Println("Received data:", string(data))

				// 클라이언트에 데이터 전송
				_, err = syscall.Write(connFd, data)
				if err != nil {
					log.Println("write failed:", err)
					syscall.EpollCtl(epollFd, syscall.EPOLL_CTL_DEL, connFd, nil)
					syscall.Close(connFd)
					continue
				}
			}

			time.Sleep(time.Millisecond * 500)
		}
	}
}
