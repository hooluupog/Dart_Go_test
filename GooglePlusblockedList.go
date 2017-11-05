package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	res, err := ioutil.ReadFile("list.html")
	if err != nil {
		fmt.Println(err)
	}
	exp := regexp.MustCompile("已屏蔽.*?</div>")
	match := exp.FindAllString(string(res), -1)
	fmt.Println("已屏蔽 ", len(match), " 人：")
	for _, v := range match {
		r := []rune(v)
		fmt.Println(string(r[3 : len(r)-6]))
	}
}
