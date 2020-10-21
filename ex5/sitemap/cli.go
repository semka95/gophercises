package sitemap

import (
	"flag"
	"fmt"
	"os"
)

// appEnv represents parsed command line arguments
type appEnv struct {
	rootLink   string
	outputFile string
	depth      int
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

// fromArgs parses command line arguments into appEnv struct
func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("go-sitemap", flag.ContinueOnError)

	l := fl.String(
		"rootLink", "", "Link to the website to build sitemap",
	)
	fl.StringVar(
		&app.outputFile, "output", "", "Path to output .xml file. By default outputs to stdout",
	)
	fl.IntVar(&app.depth, "depth", -1, "Maximum number of links to follow when building a sitemap. By default depth is not set")

	if err := fl.Parse(args); err != nil {
		return err
	}

	if *l == "" {
		fmt.Fprintf(os.Stderr, "got bad rootLink: %q\n", *l)
		fl.Usage()
		return flag.ErrHelp
	}
	app.rootLink = *l

	return nil
}

func (app *appEnv) run() error {
	queue := []string{app.rootLink}
	sm := NewSitemap(app.rootLink)
	err := sm.BuildSitemap(queue, app.depth)
	if err != nil {
		return err
	}

	w := os.Stdout
	if app.outputFile != "" {
		w, err = os.OpenFile(app.outputFile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
	}

	err = sm.buildXML(w)
	if err != nil {
		return err
	}

	return nil
}
