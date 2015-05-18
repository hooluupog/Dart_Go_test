// Login to v2ex and do something once successing login.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
)

func main() {
	var username, password string
	fmt.Scan(&username, &password)

	// Cookie manager
	client := &http.Client{}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	// Simulate browser to login.
	url := "https://v2ex.com/signin"
	value := neturl.Values{
		"u": {username},
		"p": {password},
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(value.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	//fmt.Println(req.Header)
	//fmt.Println(resp.Header)
	//fmt.Println(req.Cookies())
	//fmt.Println(resp.Cookies())

	// Do something after suceesing login.
	url = "https://v2ex.com/member/" + username
	resp, err = client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(req.Header)
	//fmt.Println(resp.Header)
	//fmt.Println(req.Cookies())
	//fmt.Println(resp.Cookies())
	fmt.Println(string(data))
}
