package general

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type OLog2 struct {
	name    string
	rootDir string
	stamp   string
	level   int
	fp      *os.File
	fileErr error
	path    string
}

var lgMutex = &sync.Mutex{}

func InitLogEnv(rootPath string, name string, level int) *OLog2 {
	lg := OLog2{}
	lg.rootDir = rootPath
	lg.name = name
	lg.level = level

	lg.path = lg.getFilePath()

	lg.fp, lg.fileErr = lg.fileopen(lg.path)

	return &lg
}

func (lg *OLog2) fileopen(path string) (*os.File, error) {

	// 이전 FD체크
	if lg.fp != nil {
		//fmt.Printf("file close..\n")
		lg.Fileclose()
		//fmt.Printf("file close.. [%v][%v]\n", *lg.fp, lg.fileErr)
	}

	lgMutex.Lock()
	// open
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	fmt.Printf(">>>fd(1) : %v", file)
	if err != nil {
		lg.fp = nil
		log.Println("Failed to open info log file : ", err)
		lg.fileErr = err
	} else {
		log.Println("file open : ", path)
	}

	lgMutex.Unlock()
	return file, err
}

func (lg *OLog2) Fileclose() {
	lgMutex.Lock()
	//fmt.Printf("file close..(1) [%s]\n", lg.path)
	lg.fp.Close()
	//lg.fp = nil
	//fmt.Printf("file close..(2) [%s]\n", lg.path)
	lgMutex.Unlock()
}

func (lg *OLog2) SetLevel(level int) {
	lg.level = level
}

func (lg *OLog2) GetLevel() int {
	return (lg.level)
}

func (lg OLog2) getFilePath() string {
	// log 패키지에서 시간 찍음
	year, month, day := time.Now().Date()

	filePath := fmt.Sprintf("%02d%02d%02d", year, int(month), day)
	dirPath := lg.rootDir + "/" + filePath

	// log directory 존재 확인
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir(dirPath, os.ModePerm)
	}

	path := lg.rootDir + "/" + filePath + "/" + lg.name + "-" + filePath + ".log"

	//return filePath
	return path
}

func (lg OLog2) checkFilePath(path string) bool {
	var isok bool = true
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("File does not exist[%s]", path)
		isok = false
	}

	return isok
}

func (lg OLog2) checkLevel(level int) bool {
	if lg.level > level {
		return false
	}

	return true
}

func (lg *OLog2) logging2(lv int, stamp string, format string, v ...interface{}) {
	// check level
	if !lg.checkLevel(lv) {
		return
	}

	lg.path = lg.getFilePath()
	//path := lg.rootDir + "/" + filePath + "/" + lg.name + "-" + filePath + ".log"
	// path 존재확인
	// 없으면 생성

	if !lg.checkFilePath(lg.path) {
		lg.fp, lg.fileErr = lg.fileopen(lg.path)
	}

	//fmt.Printf("fileopen info [%v][%v] \n", lg.fp, lg.fileErr)

	if lg.fileErr != nil {
		log.Printf("file fd error..%v \n", lg.fileErr)
		return
	}

	//fmt.Printf("logging2 :[%v] [%v]\n", lg.fp, lg.fileErr)

	var logger *log.Logger
	logger = log.New(io.MultiWriter(lg.fp, os.Stderr),
		stamp,
		log.Ldate|log.Ltime|log.Lmicroseconds)

	format += "\n"

	logger.Printf(format, v...)

}

func (lg OLog2) logging(lv int, stamp string, format string, v ...interface{}) {
	// check level
	if !lg.checkLevel(lv) {
		return
	}

	filePath := lg.getFilePath()
	path := lg.rootDir + "/" + filePath + "/" + lg.name + "-" + filePath + ".log"
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open info log file : ", err)
		return
	}

	defer func() {
		file.Close()
	}()

	var logger *log.Logger
	// logger = log.New(io.MultiWriter(file,os.Stderr),
	// 	stamp,
	// 	log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger = log.New(io.MultiWriter(file, os.Stderr),
		stamp,
		log.Ldate|log.Ltime|log.Lmicroseconds)

	format += "\n"

	logger.Printf(format, v...)
}

func (lg *OLog2) Always(format string, v ...interface{}) {
	stamp := "ALWAYS : "

	lg.logging2(99, stamp, format, v...)
}

func (lg *OLog2) Info(format string, v ...interface{}) {
	stamp := "INFO   : "
	lg.logging2(INFO, stamp, format, v...)
}

func (lg *OLog2) Print(level int, format string, v ...interface{}) {

	stamp := "PRINT   : "
	lg.logging2(level, stamp, format, v...)
}

func (lg *OLog2) Debug(format string, v ...interface{}) {
	stamp := "DEBUG  : "
	lg.logging2(DEBUG, stamp, format, v...)
}

func (lg *OLog2) DebugDump(level int, format string, v ...interface{}) {
	stamp := "DUMP  : "
	lg.logging2(level, stamp, format, v...)
}

func (lg *OLog2) Warn(format string, v ...interface{}) {
	stamp := "WARN   : "
	lg.logging2(WARN, stamp, format, v...)
}

func (lg *OLog2) Error(format string, v ...interface{}) {
	stamp := "ERROR  : "
	lg.logging2(ERROR, stamp, format, v...)
}

func (lg *OLog2) Dump(level int, stamp string, bytes []byte, length int) {
	header_msg := fmt.Sprintf("  %s  [%d]\n", stamp, length)
	var message string = ""

	for idx := 0; idx < length; idx++ {
		if idx%20 == 0 {
			if idx != 0 {
				message = message + "\n"
			}
			message += "	"
		}
		message += fmt.Sprintf(" %02X", bytes[idx])
	}

	message += "\n"

	logmsg := header_msg + message

	lg.DebugDump(level, logmsg)
}
