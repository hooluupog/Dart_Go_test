package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

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
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	data <- -1
	return written, err
}
func downloadFromUrl(url string, fileLength, data chan int64) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	if response.StatusCode == 404 {
		fmt.Println("file NotFound: ", url)
		return
	}
	fileLength <- response.ContentLength
	_, err = Copy(output, response.Body, data)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
}

func progress(length chan int64, data chan int64) {
	fileLength := <-length
	fmt.Printf("Total Size: %.2fKB\n", float64(fileLength)/1024)
	for currentSize := int64(0); currentSize != -1; currentSize = <-data {
		present := float64(currentSize) / float64(fileLength)
		i := int(present * 50)
		h := strings.Repeat("=", i) + strings.Repeat(" ", 50-i)
		fmt.Printf("\r%.2f%%[%s] %.2fKB", present*100, h, float64(currentSize)/1024)
	}
}

func main() {
	data := make(chan int64)
	fileLength := make(chan int64)
	//url := "http://www.baidu.com/img/bdlogo.gif"
	url := "http://down.sandai.net/thunder7/Thunder_dl_7.9.19.4736.exe"
	go downloadFromUrl(url, fileLength, data)
	progress(fileLength, data)
	fmt.Print("\nDownload finished.")
}
