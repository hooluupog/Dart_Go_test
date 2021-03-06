// 2014/4/7
// File downloader
// Lee
package main

import (
	"errors"
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
func newDownload(url string) *http.Response {
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
func resumeDownload(url string, fileSize, fileLength int64) *http.Response {
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
		//fmt.Println(response.Header)
		//fmt.Println(response.Header.Get("Location"))
		//fmt.Println(response.Header.Get("Content-Length"), fileLength)
		//fmt.Println(response.Header.Get("Content-Range"))
		//fmt.Println(response.StatusCode)
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
				fmt.Println("File broken.Will restart downloading.")
				defer output.Close()
				dst = output
				resp = newDownload(url)
				src = resp.Body
				written = 0
			} else { // Continue dowloading.
				fileSize := stat.Size()
				resp = resumeDownload(url, fileSize, fileLength)
				if resp.StatusCode != http.StatusPartialContent { // 206
					fmt.Println("Do not support partial download.Restart downloading.")
					output, err := os.Create(fileName)
					do(err, "Can not create ", fileName)
					defer output.Close()
					dst = output
					src = resp.Body
					written = 0
				} else {
					src = resp.Body
				}
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

//Check whether it is redirect url. Limit the redirect num <= 10.
func redirectChecker(location string) (string, string, error) {
	name, Location := location, location
	for redirect := 0; redirect < 10; redirect++ {
		req, err := http.NewRequest("GET", Location, nil)
		do(err, "bad redirect link.", Location)
		trans := &http.Transport{}
		res, err := trans.RoundTrip(req)
		do(err, "Error while return response", Location)
		contentDisposition := res.Header.Get("Content-Disposition")
		Location = res.Header.Get("Location")
		if Location == "" {
			if contentDisposition != "" {
				name = contentDisposition
				return name, location, nil
			} else {
				return name, location, nil
			}
		}
		if contentDisposition != "" {
			name = contentDisposition
		} else {
			name, location = Location, Location
		}
	}
	err := errors.New("> 10 times redirection.")
	return name, location, err
}

//Get fileName and real url from original URL.
func parseUrl(url string) (name string, rUrl string) {
	req, err := http.NewRequest("GET", url, nil)
	do(err, "Error while requesting ", url)
	// do request without redirect support.
	trans := &http.Transport{}
	res, err := trans.RoundTrip(req)
	do(err, "Error while return response", url)
	contentDisposition := res.Header.Get("Content-Disposition")
	Location := res.Header.Get("Location")
	if Location == "" { // Non redirection.
		if contentDisposition == "" {
			// Empty Content-Disposition value,
			// parse fileName directly from url.
			tokens := strings.Split(url, "/")
			name = tokens[len(tokens)-1]
			rUrl = url
			return name, rUrl
		}
		name = strings.SplitAfter(contentDisposition, "filename=")[1]
		name = strings.Trim(name, "\"")
		rUrl = url
		return name, rUrl
	}
	name, Location, err = redirectChecker(Location)
	do(err, "Too many redirect link.", url)
	tokens := strings.Split(name, "/")
	name = tokens[len(tokens)-1]
	if strings.Contains(name, "fn=") {
		name = strings.SplitAfter(name, "fn=")[1]
		name = strings.Trim(name, "\"")
	}
	if strings.Contains(name, "filename=") {
		name = strings.SplitAfter(name, "filename=")[1]
		name = strings.Trim(name, "\"")
	}
	rUrl = Location
	return name, rUrl
}

// Start a dowloading task.
func downloadFromUrl(url string, data chan int64) {
	fileName, rUrl := parseUrl(url)
	fmt.Print("File: ", fileName, " ")
	tempName := fileName + ".download"
	output, err := os.Create(tempName)
	do(err, "Error while creating ", tempName)
	client := &http.Client{}
	req, err := http.NewRequest("GET", rUrl, nil)
	do(err, "Error while requesting ", rUrl)
	response, err := client.Do(req)
	do(err, "Error while downloading ", rUrl)
	if response.StatusCode == http.StatusNotFound { //404
		log.Fatal("file NotFound: ", rUrl)
	}
	defer response.Body.Close()
	fileLength := response.ContentLength
	data <- fileLength
	written, err := transfer(output, response.Body, rUrl, tempName, fileLength, data)
	do(err, "Error while downloading ", rUrl)
	output.Close()
	if written == fileLength {
		err := os.Rename(tempName, fileName)
		do(err, "Error while renaming file", tempName)
	}
	// Wakeup progress().
	<-data
}

// Format total seconds into xx:xx:xx.(hour:minute:second)
func timeFormat(sec int64) string {
	fmtTime := ""
	if sec == -1 || (sec/3600 > 24) {
		fmtTime = "> 1 day"
		return fmtTime
	}
	hour := sec / 3600
	minute := sec % 3600 / 60
	second := sec % 3600 % 60
	switch {
	case hour < 10:
		fmtTime = "0" + strconv.FormatInt(hour, 10) + ":"
	default:
		fmtTime = strconv.FormatInt(hour, 10) + ":"
	}
	switch {
	case minute < 10:
		fmtTime += "0" + strconv.FormatInt(minute, 10) + ":"
	default:
		fmtTime += strconv.FormatInt(minute, 10) + ":"
	}
	switch {
	case second < 10:
		fmtTime += "0" + strconv.FormatInt(second, 10)
	default:
		fmtTime += strconv.FormatInt(second, 10)
	}
	return fmtTime
}

//Format file size Bytes show GB,MB,KB,B.
func sizeFormat(size float64) string {
	fmtSize := ""
	switch {
	case size/1024 <= 1:
		fmtSize = strconv.FormatFloat(size, 'f', 0, 64) + "B"
	case size/1024/1024 <= 1:
		fmtSize = strconv.FormatFloat(size/1024, 'f', 1, 64) + "KB"
	case size/1024/1024/1024 <= 1:
		fmtSize = strconv.FormatFloat(size/1024/1024, 'f', 1, 64) + "MB"
	default:
		fmtSize = strconv.FormatFloat(size/1024/1024/1024, 'f', 1, 64) + "GB"
	}
	return fmtSize
}

func progress(data chan int64) { // Real-time displaying rate of progress.
	fileLength := <-data
	start := time.Now()
	earlest := start
	var size int64
	var speed, updateSpeed float64
	totalSize := sizeFormat(float64(fileLength))
	fmt.Printf("Total Size: %s\n", totalSize)
	for lastSize, currentSize := int64(0), int64(0); currentSize != -1; currentSize = <-data {
		present := float64(currentSize) / float64(fileLength)
		//i := int(present * 50)
		//h := strings.Repeat("=", i) + strings.Repeat(" ", 50-i)
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
		cSize := sizeFormat(float64(currentSize))
		remained := int64(-1)
		if int64(speed) != 0 {
			remained = (fileLength - currentSize) / 1024 / int64(speed)
		}
		if currentSize == fileLength {
			remained = 0
		}

		elapsed := int64(time.Now().Sub(earlest).Seconds())
		rTime := timeFormat(remained)
		eTime := timeFormat(elapsed)
		fmt.Printf("\r%4.1f%%  %8s/%-8s  %4.0f KB/S  Elapsed[%8s] Remain[%8s]", present*100, cSize, totalSize, speed, eTime, rTime)
	}
	//Wait downloadfromurl() finish executing.
	data <- 1
	close(data)
}

func main() {
	data := make(chan int64)
	//url := "http://www.baidu.com/img/bdlogo.gif"
	//url := "http://down.sandai.net/thunder7/Thunder_dl_7.9.20.4754.exe"
	//url := "http://www.ubuntukylin.com/downloads/download.php?id=25"
	//url := "http://cdimage.ubuntu.com/daily-live/current/trusty-desktop-i386.iso"
	url := "http://nj.baidupcs.com/file/1fcbccdbdf4e64c9320d503e42ba1c39?xcode=9e6be5366124c95d5eb3f719eabe654969ca792c15612f02&fid=1746018321-250528-75500252898357&time=1397108108&sign=FDTAXER-DCb740ccc5511e5e8fedcff06b081203-lEUebeGcrMY0WduRWWsz%2BGyRURg%3D&to=nb&fm=N,B,T,t&newver=1&expires=1397127564&rt=sh&r=308399442&logid=337143703&sh=1&vuk=1746018321&fn=trusty-desktop-i386.iso"
	//url := "https://codeload.github.com/gabrielecirulli/2048/zip/master"
	//url := "https://gitcafe.com/riku/Markdown-Syntax-CN/tarball/master"
	//url := "http://download.skycn.com/hao123-soft-online-bcs/soft/P/2013-12-31_PowerWord.100@7728@_sky4.exe"
	go downloadFromUrl(url, data)
	progress(data)
	fmt.Print("\nDownload finished.")
}
