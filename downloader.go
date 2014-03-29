package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	buf := make([]byte, 32*1024)
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
func downloadFromUrl(url string, data chan int64) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	output, err := os.Create(fileName)
	do(err, "Error while creating ", fileName)
	defer output.Close()
	response, err := http.Get(url)
	do(err, "Error while downloading ", url)
	if response.StatusCode == 404 {
		log.Fatal("file NotFound: ", url)
	}
	defer response.Body.Close()
	data <- response.ContentLength
	_, err = Copy(output, response.Body, data)
	do(err, "Error while downloading ", url)
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
		fmt.Printf("\r%3.1f%%[%s] %4.0f KB/S", present*100, h, speed)
	}
}

func main() {
	data := make(chan int64)
	//url := "http://www.baidu.com/img/bdlogo.gif"
	url := "http://releases.ubuntu.com/trusty/ubuntu-14.04-beta2-desktop-i386.iso"
	go downloadFromUrl(url, data)
	progress(data)
	fmt.Print("\nDownload finished.")
}
