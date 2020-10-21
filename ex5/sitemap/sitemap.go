// sitemap package builds a map of all of the pages within
// a specific domain and provides it in specific xml format -
// https://www.sitemaps.org/schemas/sitemap/0.9/
package sitemap

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	link "github.com/semka95/gophercises/ex4"
)

const Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

// Urlset represents sitemap xml struct https://www.sitemaps.org/schemas/sitemap/0.9/
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Links   []Loc    `xml:"url"`
}

// Loc represents loc field in sitemap xml struct
type Loc struct {
	URL string `xml:"loc"`
}

// CLI runs the go-sitemap command line app and returns its exit status.
func CLI(args []string) int {
	var app appEnv
	err := app.fromArgs(args)
	if err != nil {
		return 2
	}
	if err = app.run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}
	return 0
}

// appEnv represents parsed command line arguments
type appEnv struct {
	link         string
	outputFile   string
	depth        int
	visitedLinks map[string]struct{}
}

// fromArgs parses command line arguments into appEnv struct
func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-sitemap", flag.ContinueOnError)

	l := fl.String(
		"link", "", "Link to the website to build sitemap",
	)
	fl.StringVar(
		&app.outputFile, "output", "", "Path to output .xml file. By default outputs to stdout",
	)
	fl.IntVar(&app.depth, "depth", -1, "Maximum number of links to follow when building a sitemap. By default depth is not set")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if *l == "" {
		fmt.Fprintf(os.Stderr, "got bad link: %q\n", *l)
		fl.Usage()
		return flag.ErrHelp
	}
	app.link = *l

	app.visitedLinks = make(map[string]struct{})
	app.visitedLinks[app.link] = struct{}{}

	return nil
}

func (app *appEnv) run() error {
	queue := []string{app.link}
	output := os.Stdout

	err := app.BuildSitemap(queue, app.depth)
	if err != nil {
		return err
	}

	if app.outputFile != "" {
		output, err = os.OpenFile(app.outputFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer output.Close()
	}

	err = app.buildXML(output)
	if err != nil {
		return err
	}

	return nil
}

// BuildSitemap visits links in queue and parses links in other queue and
// recursively visits them. If depth is greater than zero it defines number
// of links to follow. All visited links recorded to appEnv.visitedLinks
func (app *appEnv) BuildSitemap(queue []string, depth int) error {
	if depth == 0 {
		return nil
	}

	discoveredLinks := make([]string, 0)

	for _, v := range queue {
		resp, err := http.Get(v)
		if err != nil {
			return err
		}

		l, err := app.parseLinks(resp.Body)
		if err != nil {
			return err
		}
		discoveredLinks = append(discoveredLinks, l...)

		resp.Body.Close()
	}

	if len(discoveredLinks) > 0 {
		depth--
		return app.BuildSitemap(discoveredLinks, depth)
	}

	return nil
}

// parseLinks reads from io.Reader and creates array of links.
// It only parses links with appEnv.link domain
func (app *appEnv) parseLinks(r io.Reader) ([]string, error) {
	res, err := link.ParseHTML(r)
	if err != nil {
		return nil, err
	}

	links := make([]string, 0)

	for _, v := range res {
		href := v.Href
		if strings.HasPrefix(href, app.link) {
			if _, ok := app.visitedLinks[href]; !ok {
				app.visitedLinks[href] = struct{}{}
				links = append(links, href)
				continue
			}
		}

		if strings.HasPrefix(href, "/") {
			href = app.link + href
			if _, ok := app.visitedLinks[href]; !ok {
				app.visitedLinks[href] = struct{}{}
				links = append(links, href)
			}
		}
	}

	return links, nil
}

// buildXML builds xml using appEnv.visitedLinks and writes it
// to io.Writer
func (app *appEnv) buildXML(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	links := make([]Loc, 0, len(app.visitedLinks))
	for k := range app.visitedLinks {
		links = append(links, Loc{URL: k})
	}

	urlset := &Urlset{
		Xmlns: Xmlns,
		Links: links,
	}

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")

	if err := enc.Encode(urlset); err != nil {
		return err
	}

	return nil
}
