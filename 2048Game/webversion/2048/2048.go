package main

import (
	"2048/game"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Page struct {
	Win   string
	Tile  [][]uint32
	Uzero uint32
}

var mux = http.NewServeMux()
var p *Page = new(Page)
var templates = template.Must(template.ParseFiles("2048.html", "grid.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func router(w http.ResponseWriter, r *http.Request, templ string, page *Page) {
	if strings.Contains(r.URL.Path, ".") {
		mux.ServeHTTP(w, r)
	} else {
		w.Header().Set("Content-Type", "text/html")
		renderTemplate(w, templ, page)
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	router(w, r, "2048", nil)
}
func gridHandler(w http.ResponseWriter, r *http.Request) {
	p.Tile = game.InitGrid()
	router(w, r, "grid", p)
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	direction, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.Tile = game.Move(string(direction))
	p.Win = game.Win()
	router(w, r, "grid", p)
}

func main() {
	mux.Handle("/", http.FileServer(http.Dir("./")))
	mux.Handle("/grid", http.FileServer(http.Dir("./")))
	mux.Handle("/up", http.FileServer(http.Dir("./")))
	mux.Handle("/down", http.FileServer(http.Dir("./")))
	mux.Handle("/left", http.FileServer(http.Dir("./")))
	mux.Handle("/right", http.FileServer(http.Dir("./")))
	http.HandleFunc("/", initHandler)
	http.HandleFunc("/grid", gridHandler)
	http.HandleFunc("/up", moveHandler)
	http.HandleFunc("/down", moveHandler)
	http.HandleFunc("/left", moveHandler)
	http.HandleFunc("/right", moveHandler)
	http.ListenAndServe(":8080", nil)
}
