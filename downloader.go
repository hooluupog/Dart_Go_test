package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func do(err error, info ...interface{}) {
	if err != nil {
		log.Fatal("\n", info, "\n", err)
	}
}

func isConnected(url string) bool {
	_, err := http.Head(url)
	if err != nil {
		return false
	}
	return true
}
func new_download(url string) io.Reader {
	client := &http.Client{}
	for i := 0; i < 10; i++ {
		time.Sleep(10000 * time.Millisecond)
		fmt.Println("Trying ", i+1, " times")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		response, err := client.Do(req)
		if err != nil {
			continue
		}
		fmt.Println("Connected.Restart downloading.")
		return response.Body
	}
	return nil
}
func resume_download(url string, fileSize, fileLength int64) io.Reader {
	client := &http.Client{}
	for i := 0; i < 10; i++ {
		time.Sleep(10000 * time.Millisecond)
		fmt.Println("Trying ", i+1, " times")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		contentRange := "bytes=" + strconv.FormatInt(fileSize, 10) + "-" + strconv.FormatInt(fileLength-1, 10)
		req.Header.Add("Range", contentRange)
		response, err := client.Do(req)
		if err != nil {
			continue
		}
		fmt.Println("Connected.Resume downloading.")
		return response.Body
	}
	return nil
}
func transfer(dst io.Writer, src io.Reader, url string, fileName string, fileLength int64, data chan int64) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		if !isConnected(url) {
			fmt.Println("\nConnection closed.Retry connecting...")
			stat, err := os.Stat(fileName)
			if err != nil {
				output, err := os.Create(fileName)
				do(err, "Can not create ", fileName)
				fmt.Println("File broken.Will restart downloaing.")
				defer output.Close()
				dst = output
				src = new_download(url)
				written = 0
			} else {
				fileSize := stat.Size()
				src = resume_download(url, fileSize, fileLength)
			}
			if src == nil {
				fmt.Println("download failed.")
				os.Exit(1)
			}
		}
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
			break
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
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	do(err, "Error while requesting ", url)
	response, err := client.Do(req)
	do(err, "Error while downloading ", url)
	if response.StatusCode == 404 {
		log.Fatal("file NotFound: ", url)
	}
	defer response.Body.Close()
	fileLength := response.ContentLength
	data <- fileLength
	_, err = transfer(output, response.Body, url, fileName, fileLength, data)
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
		if currentSize < lastSize {
			lastSize = currentSize
		}
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
	url := "http://releases.ubuntu.com/14.04/ubuntu-14.04-beta2-desktop-i386.iso"
	go downloadFromUrl(url, data)
	progress(data)
	fmt.Print("\nDownload finished.")
}
