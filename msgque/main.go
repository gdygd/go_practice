package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	KEY       = 1234
	MSG_TYPE  = 1
	MSG_FLAGS = 0
	MSG_SIZE  = 128
	IPC_CREAT = 01000 // 추가된 상수

	sysMsgGet = 68
	sysMsgSnd = 69
	sysMsgRcv = 70
	sysMsgCtl = 71
)

type msgbuf struct {
	Mtype int64
	Mtext [MSG_SIZE]byte
}

func main() {
	// 메시지 큐를 생성합니다.
	msqid, _, errno := syscall.Syscall(sysMsgGet, uintptr(KEY), uintptr(IPC_CREAT|0666), uintptr(int32(0)))

	if errno != 0 {
		fmt.Printf("Msgget: %v\n", errno)
		return
	}
	fmt.Printf("Message queue ID: %d\n", msqid)

	// 메시지 구조체를 정의합니다.
	msg := msgbuf{Mtype: MSG_TYPE}
	copy(msg.Mtext[:], "Hello, World!")
	_, _, errno = syscall.Syscall(sysMsgSnd, msqid, uintptr(unsafe.Pointer(&msg)), uintptr(int32(MSG_FLAGS)))
	if errno != 0 {
		fmt.Printf("Msgsnd: %v\n", errno)
		return
	}
	fmt.Println("Message sent")

	// 메시지를 받습니다.
	var rcvmsg msgbuf
	//_, _, errno = syscall.Syscall(sysMsgRcv, msqid, uintptr(unsafe.Pointer(&rcvmsg)), MSG_SIZE, MSG_TYPE|MSG_FLAGS)
	_, _, errno = syscall.Syscall6(sysMsgRcv, msqid, uintptr(unsafe.Pointer(&rcvmsg)), MSG_SIZE, MSG_TYPE, MSG_FLAGS, 0)
	if errno != 0 {
		fmt.Printf("Msgrcv: %v\n", errno)
		return
	}
	fmt.Printf("Received message: %s\n", rcvmsg.Mtext)
}
