package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Story map[string]Chapter

type Chapter struct {
	Title      string    `json:"title"`
	Paragraphs []string  `json:"story"`
	Options    []Options `json:"options"`
}

type Options struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Srv struct {
	Story    Story
	Template *template.Template
}

func main() {
	var story Story

	storyFile, err := os.Open("./gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	jd := json.NewDecoder(storyFile)
	err = jd.Decode(&story)
	if err != nil {
		log.Fatal(err)
	}

	t, err := ioutil.ReadFile("./template.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("Story").Parse(string(t))
	if err != nil {
		log.Fatal(err)
	}

	srv := &Srv{
		Story:    story,
		Template: tmpl,
	}

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", srv)
}

func (s *Srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if v, ok := s.Story[path]; ok {
		err := s.Template.Execute(w, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "http://localhost:8080/intro", http.StatusFound)
}
