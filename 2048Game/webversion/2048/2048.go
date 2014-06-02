package main

import (
	"2048/game"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Page struct {
	Win   string
	Tile  [][]game.Tile
	Uzero uint32
	Size  int
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
		if r.Method == "GET" && r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		renderTemplate(w, templ, page)
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	router(w, r, "2048", nil)
}
func gridHandler(w http.ResponseWriter, r *http.Request) {
	game.Reset()
	p.Tile = game.InitGrid()
	p.Size = game.GridSize()
	p.Win = game.Win()
	router(w, r, "grid", p)
}

func moveHandler(w http.ResponseWriter, r *http.Request) {
	direction, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if strings.Contains(r.URL.Path, ".") {
		mux.ServeHTTP(w, r)
	} else {
		if r.Method == "GET" && r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		p.Tile = game.Move(string(direction))
		p.Win = game.Win()
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		err := enc.Encode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
