// 请求一个网站，返回结果
// 第一次请求，通过form传递认证信息
// 以后多次请求时，通过cookies实现已登录状态页面访问

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

func main() {
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	authUrl := "http://v.yupoo.com/account/login/"
	value := url.Values{
		"username":   {"<username>"},
		"password":   {"<password>"},
		"grant_type": {"password"},
	}
	req, err := http.NewRequest("POST", authUrl, strings.NewReader(value.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(req.Header)
	//fmt.Println(resp.Header)
	defer resp.Body.Close()
	//cookies := resp.Cookies()
	//fmt.Println(cookies)
	url := "http://v.yupoo.com/account/layout/"
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(req.Header)
	//fmt.Println(resp.Header)
	defer resp.Body.Close()
	fmt.Println(string(data))
}
