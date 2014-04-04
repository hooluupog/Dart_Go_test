// 2014/3/31
// File downloader

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

// Error Handler function
func do(err error, info ...interface{}) {
	if err != nil {
		log.Fatal("\n", info, "\n", err)
	}
}

//Decide whether connection is alive
func isConnected(url string) bool {
	_, err := http.Head(url)
	if err != nil {
		return false
	}
	return true
}

// Start a new dowloading.
func new_download(url string) *http.Response {
	client := &http.Client{}
	// Try connecting 10 times
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
		return response
	}
	return nil
}

// Continue downloading if sucessing to connect server host.
func resume_download(url string, fileSize, fileLength int64) *http.Response {
	// Try connecting 10 times
	client := &http.Client{}
	for i := 0; i < 10; i++ {
		time.Sleep(10000 * time.Millisecond)
		fmt.Println("Trying ", i+1, " times")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		// Partial content transmit. Http 206 protocol.
		contentRange := "bytes=" + strconv.FormatInt(fileSize, 10) + "-" + strconv.FormatInt(fileLength-1, 10)
		req.Header.Add("Range", contentRange)
		response, err := client.Do(req)
		if err != nil {
			continue
		}
		fmt.Println("Connected.Resume downloading.")
		return response
	}
	return nil
}

// File downloading scheduler.
func transfer(dst io.Writer, src io.Reader, url string, fileName string, fileLength int64, data chan int64) (written int64, err error) {
	buf := make([]byte, 32*1024)
	var resp *http.Response
	for {
		if !isConnected(url) {
			fmt.Println("\nConnection closed.Retry connecting...")
			stat, err := os.Stat(fileName)
			if err != nil { // File broken and recreating a new file.
				output, err := os.Create(fileName)
				do(err, "Can not create ", fileName)
				fmt.Println("File broken.Will restart downloaing.")
				defer output.Close()
				dst = output
				resp = new_download(url)
				src = resp.Body
				written = 0
			} else { // Continue dowloading.
				fileSize := stat.Size()
				resp = resume_download(url, fileSize, fileLength)
				src = resp.Body
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
			data <- -1 // Send dowloading task finished signal.
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	if resp != nil {
		resp.Body.Close()
	}
	return written, err
}

//Get fileName from URL
func getFileName(url string) (name string) {
	header, err := http.Head(url)
	do(err, "Error while getting header", url)
	contentDisposition := header.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		// Empty Content-Disposition value,
		// parse fileName directly from url.
		tokens := strings.Split(url, "/")
		name := tokens[len(tokens)-1]
		return name
	}

	name = strings.SplitAfter(contentDisposition, "filename=")[1]
	return name
}

// Start a dowloading task.
func downloadFromUrl(url string, data chan int64) {
	fileName := getFileName(url)
	fmt.Print("File: ", fileName, " ")
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

func progress(data chan int64) { // Real-time displaying rate of progress.
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
		if currentSize < lastSize { // A new dowloading begins.
			lastSize = currentSize
		}
		size = size + currentSize - lastSize
		speed = updateSpeed
		lastSize = currentSize
		if duration > 1 { // Update progress every 1 seconds or so.
			start = time.Now()
			updateSpeed = float64(size) / 1024 / duration
			size = 0
		}
		fmt.Printf("\r%3.1f%%[%s] %4.0f KB/S", present*100, h, speed)
	}
}

func main() {
	data := make(chan int64)
	url := "http://www.baidu.com/img/bdlogo.gif"
	//url := "http://releases.ubuntu.com/14.04/ubuntu-14.04-beta2-desktop-i386.iso"
	//url := "https://codeload.github.com/gabrielecirulli/2048/zip/master"
	go downloadFromUrl(url, data)
	progress(data)
	fmt.Print("\nDownload finished.")
}
