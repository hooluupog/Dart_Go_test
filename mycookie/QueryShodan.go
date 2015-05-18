package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mycookie"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"regexp"
	"strings"
)

// My personal cookieJar implemention to allow non stardard cookie.
type MyJar struct {
	jar map[string][]*http.Cookie
}

var client = &http.Client{}

func (j *MyJar) SetCookies(u *neturl.URL, cookies []*http.Cookie) {
	j.jar[u.Host] = cookies
}

func (j *MyJar) Cookies(u *neturl.URL) []*http.Cookie {
	return j.jar[u.Host]
}

// Add MyJar's cookies into request header.
func SendCookies(req *http.Request, j *MyJar) {
	cookies := j.Cookies(req.URL)
	if len(cookies) > 0 {
		for _, cookie := range cookies {
			s := cookie.Name + "=" + cookie.Value
			if c := req.Header.Get("Cookie"); c != "" {
				req.Header.Add("Cookie", s)
			} else {
				req.Header.Set("Cookie", s)
			}
		}
	}
}

//Store non stardard cookies into MyJar.
func ReceiveCookies(u *neturl.URL, resp *http.Response, j *MyJar) {
	var cookies []*http.Cookie
	c := mycookie.Cookies(resp.Header)
	for _, cookie := range c {
		_, success := mycookie.ParseCookieValue(cookie.Value, true)
		if !success {
			cookies = append(cookies, cookie)
		}
	}
	j.SetCookies(u, cookies)
}

func Query(method, url string, value io.Reader, output bool, myjar *MyJar) []byte {
	req, err := http.NewRequest(method, url, value)
	if err != nil {
		log.Fatal(err)
	}
	SendCookies(req, myjar)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.152 Safari/537.36")
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ReceiveCookies(req.URL, resp, myjar)
	if output {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return data
	}
	return nil
}

// Parse data and print out wanted result.
func printList(data []byte) {
	var pList [][]byte
	re := regexp.MustCompile(`(?s)<div class="search-result-summary">.*?<span[^>]*>(.*?)</span>.*?</div>`)
	list := re.FindAllSubmatch(data, -1)
	// list is [][][]byte
	for _, v := range list {
		pList = append(pList, v[1])
	}
	for _, v := range pList {
		fmt.Println(string(v))
	}
}

func main() {
	var username, password string
	fmt.Scan(&username, &password)

	// Cookie manager
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	myjar := &MyJar{}
	myjar.jar = make(map[string][]*http.Cookie)

	url := "https://account.shodan.io/login"
	Query("GET", url, nil, false, myjar)
	url = "https://account.shodan.io/login"
	value := neturl.Values{
		"username":   {username},
		"password":   {password},
		"grant_type": {"password"},
	}
	Query("POST", url, strings.NewReader(value.Encode()), false, myjar)

	// For instance,search google ip, keyword="gws"
	var keyword string
	fmt.Println("Please input search keyword:")
	fmt.Scan(&keyword)
	url = "https://www.shodan.io/search?query=" + keyword + "&page=1"
	data := Query("GET", url, nil, true, myjar)
	printList(data)
	url = "https://www.shodan.io/search?query=gws&page=2"
	data = Query("GET", url, nil, true, myjar)
	printList(data)
}
