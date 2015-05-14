package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func main() {
	keyword := "gws"
	var iplist [][]byte
	re := regexp.MustCompile(`(?s)<div class="search-result-summary">.*?<span[^>]*>(.*?)</span>.*?</div>`)
	client := &http.Client{}
	url := "https://www.shodan.io/search?query=" + keyword
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	list := re.FindAllSubmatch(data, -1)
	// list is [][][]byte
	for _, v := range list {
		iplist = append(iplist, v[1])
	}
	//time.Sleep(5 * time.Second)
	resp.Body.Close()
	for _, v := range iplist {
		fmt.Println(string(v))
	}
}
