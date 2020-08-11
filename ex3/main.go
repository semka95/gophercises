package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Arc struct {
	Title   string
	Story   []string
	Options []Options
}

type Options struct {
	Text string
	Arc  string
}

type Srv struct {
	Arcs     map[string]Arc
	Template *template.Template
}

func main() {
	arcs := make(map[string]Arc)

	story, err := ioutil.ReadFile("./gopher.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(story, &arcs)
	if err != nil {
		log.Fatal(err)
	}

	t, err := ioutil.ReadFile("./template.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("Arc").Parse(string(t))
	if err != nil {
		log.Fatal(err)
	}

	srv := &Srv{
		Arcs:     arcs,
		Template: tmpl,
	}

	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", srv)
}

func (s *Srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	if v, ok := s.Arcs[path]; ok {
		err := s.Template.Execute(w, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "http://localhost:8080/intro", http.StatusFound)
}
