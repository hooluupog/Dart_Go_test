package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func do(err error, info ...interface{}) {
	if err != nil {
		log.Fatal(info, "\n", err)
	}
}

func Copy(dst io.Writer, src io.Reader, data chan int64) (written int64, err error) {
	buf := make([]byte, 8*1024*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				data <- written
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			data <- -1
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
func downloadFromFile(path string, data chan int64) {
	tokens := strings.Split(path, "/")
	fileName := tokens[len(tokens)-1]
	output, err := os.Create(fileName)
	do(err, "Error while creating ", fileName)
	defer output.Close()
	stat, err := os.Stat(path)
	do(err, "Error while getting file info ", fileName)
	data <- stat.Size()
	src, err := os.Open(path)
	do(err, "Error while reading file ", fileName)
	defer src.Close()
	_, err = Copy(output, src, data)
	do(err, "Error while downloading ", path)
}

func progress(data chan int64) {
	fileLength := <-data
	start := time.Now()
	var size int64
	var speed, updateSpeed float64
	totalSize := float64(fileLength)
	switch {
	case totalSize/1024 <= 1:
		fmt.Printf("Total Size: %fB\n", totalSize)
	case totalSize/1024/1024 <= 1:
		fmt.Printf("Total Size: %.0fkB\n", totalSize/1024)
	case totalSize/1024/1024/1024 <= 1:
		fmt.Printf("Total Size: %.0fMB\n", totalSize/1024/1024)
	default:
		fmt.Printf("Total Size: %.1fGB\n", totalSize/1024/1024/1024)
	}
	for lastSize, currentSize := int64(0), int64(0); currentSize != -1; currentSize = <-data {
		present := float64(currentSize) / float64(fileLength)
		i := int(present * 50)
		h := strings.Repeat("=", i) + strings.Repeat(" ", 50-i)
		duration := time.Since(start).Seconds()
		size = size + currentSize - lastSize
		speed = updateSpeed
		lastSize = currentSize
		if duration > 1 {
			start = time.Now()
			updateSpeed = float64(size) / 1024 / duration
			size = 0
		}
		switch {
		case speed/1024 < 1:
			fmt.Printf("\r%3.1f%%[%s] %4.0f KB/S", present*100, h, speed)
		default:
			fmt.Printf("\r%3.1f%%[%s] %4.0f MB/S", present*100, h, speed/1024)
		}
	}
}

func main() {
	data := make(chan int64)
	//url := "http://www.baidu.com/img/bdlogo.gif"
	path := "D:/备份/时空幻境 Braid 终生珍藏版 V1.0.zip"
	go downloadFromFile(path, data)
	progress(data)
	fmt.Print("\nDownload finished.")
}
