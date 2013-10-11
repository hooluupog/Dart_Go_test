//统计出一个项目中每种编程语言所使用的数量。
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
        "os"
)

var c_count int = 0
var cpp_count int = 0
var h_count int = 0
var java_count int = 0
var go_count int = 0
var dart_count int = 0

func visit(path string, fileInfo os.FileInfo, err error) error {
	if !fileInfo.IsDir() {
		counter(path,fileInfo)
	}
	return err
}

func counter(path string,file os.FileInfo) {
	extension := filepath.Ext(file.Name())
	switch extension {
	case ".c":
		c_count++
	case ".cpp" :
		cpp_count++
        case ".cc" :
		cpp_count++
	case ".java":
		java_count++
	case ".go":
		go_count++
	case ".dart":
		dart_count++
	}
}

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, visit)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("c : %d\ncpp : %d\njava : %d\ngo : %d\ndart : %d\n", c_count, cpp_count, java_count, go_count, dart_count)
}
