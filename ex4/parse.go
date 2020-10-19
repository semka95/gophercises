package link

import (
	"golang.org/x/net/html"
	"io"
)

// Link represents HTML link tag
type Link struct {
	Href string
	Text string
}

// ParseHTML parses given html file and returns all links
func ParseHTML(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)

	var parseNode func(node *html.Node)
	parseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			link := Link{
				Href: parseHref(n.Attr),
				Text: parseLinkText(n),
			}

			links = append(links, link)

		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseNode(c)
		}
	}

	parseNode(doc)

	return links, nil
}

// parseLinkText extracts text from link tag
func parseLinkText(n *html.Node) (text string) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			text += c.Data
		}

		if c.Type == html.ElementNode && c.FirstChild != nil {
			text += c.FirstChild.Data
		}
	}

	return
}

// parseHref extracts href attribute from link tag
func parseHref(attrs []html.Attribute) (href string) {
	for _, a := range attrs {
		if a.Key == "href" {
			href = a.Val
			break
		}
	}

	return
}
