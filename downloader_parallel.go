// 2014/4/20
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

type fileInfo struct {
	fileName   string
	fileLength int64
}

var (
	currentSize int64
	lastSize    int64
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
		fmt.Println("")
		fmt.Println("Trying ", i+1, " times")
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		response, err := client.Do(req)
		if err != nil {
			continue
		}
		fmt.Println("\nConnected.Restart downloading.")
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
		fmt.Println("")
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
		fmt.Println("\nConnected.Resume downloading.")
		return response
	}
	return nil
}

// File downloading scheduler.
func transfer(dst io.Writer, src io.Reader, url string, fileName string, fileLength int64) (written int64, err error) {
	buf := int64(32 * 1024)
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
				if resp.StatusCode != http.StatusPartialContent { //206
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
		w, er := io.CopyN(dst, src, buf)
		written += w
		if er == io.EOF {
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
func downloadFromUrl(url string, fInfo *fileInfo, quit chan int) {
	var rUrl string
	fInfo.fileName, rUrl = parseUrl(url)
	fmt.Print("File: ", fInfo.fileName, " ")
	tempName := fInfo.fileName + ".download"
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
	fInfo.fileLength = response.ContentLength
	totalSize := sizeFormat(float64(fInfo.fileLength))
	fmt.Printf("Total Size: %s\n", totalSize)
	written, err := transfer(output, response.Body, rUrl, tempName, fInfo.fileLength)
	do(err, "Error while downloading ", rUrl)
	output.Close()
	if written == fInfo.fileLength {
		err := os.Rename(tempName, fInfo.fileName)
		do(err, "Error while renaming file", tempName)
	}
	quit <- 1
}

// Format total seconds into xx:xx:xx.(hour:minute:second)
func timeFormat(sec int64) string {
	fmtTime := ""
	if float64(sec)/3600 > 24 {
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

func progress(fInfo *fileInfo, elapsed int64) { // Real-time displaying rate of progress.
	var size int64
	var speed float64
	if fInfo.fileName == "" {
		return
	}
	stat, err := os.Stat(fInfo.fileName + ".download")
	if err != nil {
		stat, err = os.Stat(fInfo.fileName)
	}
	do(err, "", fInfo.fileName)
	currentSize = stat.Size()
	totalSize := sizeFormat(float64(fInfo.fileLength))
	present := float64(currentSize) / float64(fInfo.fileLength)
	//i := int(present * 50)
	//h := strings.Repeat("=", i) + strings.Repeat(" ", 50-i)
	if currentSize < lastSize { // A new dowloading begins.
		lastSize = currentSize
	}
	size = currentSize - lastSize
	lastSize = currentSize
	speed = float64(size) / 1024
	cSize := sizeFormat(float64(currentSize))
	remained := int64(3600*24 + 1)
	if int64(speed) != 0 {
		remained = (fInfo.fileLength - currentSize) / 1024 / int64(speed)
	}
	if currentSize == fInfo.fileLength {
		remained = 0
	}
	rTime := timeFormat(remained)
	eTime := timeFormat(elapsed)
	if fInfo.fileLength != 0 {
		fmt.Printf("\r%4.1f%%  %8s/%-8s  %4.0f KB/S  Elapsed[%8s] Remain[%8s]", present*100, cSize, totalSize, speed, eTime, rTime)
	}
}

func main() {
	fInfo := new(fileInfo)
	elapsed := int64(0)
	quit := make(chan int)
	//url := "http://www.baidu.com/img/bdlogo.gif"
	url := "http://down.sandai.net/thunder7/Thunder_dl_7.9.20.4754.exe"
	//url := "http://www.ubuntukylin.com/downloads/download.php?id=25"
	//url := "http://cdimage.ubuntu.com/daily-live/current/trusty-desktop-i386.iso"
	//url := "https://d3g7pb956c5s5x.cloudfront.net/ubuntukylin-14.04-desktop-i386.iso"
	//url := "https://codeload.github.com/gabrielecirulli/2048/zip/master"
	//url := "https://gitcafe.com/riku/Markdown-Syntax-CN/tarball/master"
	//url := "http://download.skycn.com/hao123-soft-online-bcs/soft/P/2013-12-31_PowerWord.100@7728@_sky4.exe"
	//runtime.GOMAXPROCS(runtime.NumCPU())
	tick := time.Tick(1000 * time.Millisecond)
	go downloadFromUrl(url, fInfo, quit)
	for {
		select {
		case <-tick:
			elapsed++
			progress(fInfo, elapsed)
		case <-quit:
			progress(fInfo, elapsed)
			fmt.Print("\nDownload finished.")
			return
		}
	}
}
