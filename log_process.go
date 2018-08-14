package main

import (
	"strings"
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"bytes"
	"regexp"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(rc chan string)
}

type ReadFromFile struct {
	path string //读取文件路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influx data source
}

func (r *ReadFromFile) Read(rc chan []byte) {
	// 读取模块
	// 打开模块

	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}

	//从文件末尾开始逐行读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
		} else if err != nil{
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}
		line = bytes.Trim(line,"\r\n")
		if len(line) > 0 {
			rc <- line
		}

	}

}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块

	for v := range wc {
		fmt.Println(v)
	}
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

func (l *LogProcess) Process() {
	//解析模块

	r := regexp.MustCompile(`([\d\.]+)\s+([^\[]+)\s+([^\[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}

}

func main() {

	r := &ReadFromFile{
		path: "./access.log",
	}

	w := &WriteToInfluxDB{
		influxDBDsn: "username&password..",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(200*time.Second)
}
