package main

import (
	"Test_Pig/controller"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var chttp = http.NewServeMux()

type Page struct {
	Image []string
	Num   int
}

func loadImage() *Page {
	dir, err := os.Open("./")
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	fi, err := dir.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileNames := make([]string, 0)
	if fi.IsDir() {
		fis, err := dir.Readdir(-1)
		if err != nil {
			log.Fatal(err)
		}
		for _, fileInfo := range fis {
			if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".png" {
				fileNames = append(fileNames, fileInfo.Name())
			}
		}
	}
	return &Page{Image: fileNames}
}

var templates = template.Must(template.ParseFiles("view.html", "edit.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	num, err := strconv.Atoi(string(body))
	if num < 0 {
		err = errors.New("Invalid number.")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	controller.Run(num)
	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		p := loadImage()
		p.Num = num
		renderTemplate(w, "view", p)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, ".") {
		chttp.ServeHTTP(w, r)
	} else {
		w.Header().Set("Content-Type", "text/html")
		renderTemplate(w, "edit", nil)
	}
}

func main() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))
	chttp.Handle("/view", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", editHandler)
	http.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":8080", nil)
}
