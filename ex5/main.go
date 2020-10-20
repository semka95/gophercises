package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	link "github.com/semka95/gophercises/ex4"
)

var u = "https://www.calhoun.io"
var links map[string]struct{}

func main() {
	links = make(map[string]struct{})
	links[u] = struct{}{}
	queue := make([]string, 0)
	queue = append(queue, u)
	f(queue)

	for k := range links {
		fmt.Printf("link: %s\n", k)
	}
	fmt.Println(len(links))
}

func f(queue []string) {
	newQueue := make([]string, 0)
	for _, v := range queue {
		resp, err := http.Get(v)
		if err != nil {
			panic(err)
		}

		newQueue = append(newQueue, parseLinks(resp.Body)...)

		resp.Body.Close()
	}

	if len(newQueue) > 0 {
		f(newQueue)
	}
}

func parseLinks(r io.Reader) []string {
	res, err := link.ParseHTML(r)
	if err != nil {
		panic(err)
	}

	queue := make([]string, 0)

	for _, v := range res {
		if strings.HasPrefix(v.Href, u) {
			if _, ok := links[v.Href]; !ok {
				links[v.Href] = struct{}{}
				queue = append(queue, v.Href)
				continue
			}
		}

		if strings.HasPrefix(v.Href, "/") {
			if _, ok := links[u+v.Href]; !ok {
				links[u+v.Href] = struct{}{}
				queue = append(queue, u+v.Href)
			}
		}
	}

	return queue
}
