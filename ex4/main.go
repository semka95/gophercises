package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
)

func main() {
	r, err := os.Open("./ex4.html")
	if err != nil {
		log.Fatal(err)
	}

	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					fmt.Print("link: ", a.Val)
					break
				}
			}
			fmt.Print(" data: ")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					fmt.Print(c.Data)
				}
			}
			fmt.Println()
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
}
