package main

import (
	"fmt"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
)

//var shminst shm.ShmHnd

type ShmMem struct {
	Name [10]byte
	Num1 int
	Num2 int
	Addr [10]byte
}

var SharedMem *ShmMem

const skey = 0x1234

func main() {
	shminst := shmlinux.NewLinuxShm()

	shminst.InitShm(skey, 1024)

	err := shminst.CreateShm()
	if err != nil {
		fmt.Println("CreateShm err : ", err)
	}

	err = shminst.AttachShm()
	if err != nil {
		fmt.Println("AttachShm err : ", err)
	}
	SharedMem = (*ShmMem)(unsafe.Pointer(shminst.Addr))

	copy(SharedMem.Addr[:], []byte(string("1234aa")))
	copy(SharedMem.Name[:], []byte(string("YGD")))
	SharedMem.Num1 = 11
	SharedMem.Num2 = 99

	var n int
	fmt.Scanf("%d", &n)

	fmt.Printf("\nmem#1 : %v\n", SharedMem)

	err = shminst.DeleteShm()
	if err != nil {
		fmt.Println("DeleteShm err:", err)
	}
}
