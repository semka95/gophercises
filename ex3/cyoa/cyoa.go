// cyoa is a package for building Choose Your Own Adventure
// stories that can be rendered via the resulting http.Handler
// or command line interface
package cyoa

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

// CLI runs the go-cyoa command line app and returns its exit status.
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
	outputCLI bool
	storyJSON string
	intro     string
}

// fromArgs parses command line arguments into appEnv struct
func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("cyoa", flag.ContinueOnError)
	fl.StringVar(
		&app.storyJSON, "story", "./gopher.json", "Path to story file in json format",
	)
	fl.StringVar(
		&app.intro, "intro", "intro", "Intro chapter name",
	)
	outputType := fl.String(
		"o", "web", "Print output in format: web/cli",
	)
	if err := fl.Parse(args); err != nil {
		return err
	}
	if *outputType != "web" && *outputType != "cli" {
		fmt.Fprintf(os.Stderr, "got bad output type: %q\n", *outputType)
		fl.Usage()
		return flag.ErrHelp
	}
	app.outputCLI = *outputType == "cli"
	return nil
}

func (app *appEnv) run() error {
	story, err := app.parseStory()
	if err != nil {
		return err
	}

	if app.outputCLI {
		cli := StoryCLI{
			Story:        story,
			IntroChapter: app.intro,
		}
		err := cli.runCLI()
		if err != nil {
			return err
		}
		return nil
	}

	tmpl, err := app.parseTemplate()
	if err != nil {
		return err
	}

	storyServer := &StoryWebServer{
		Story:        story,
		Template:     tmpl,
		IntroChapter: app.intro,
	}

	err = storyServer.runWeb()
	if err != nil {
		return err
	}

	return nil
}

// parseStory parses story from json file
func (app *appEnv) parseStory() (Story, error) {
	storyFile, err := os.Open(app.storyJSON)
	if err != nil {
		return nil, err
	}
	defer storyFile.Close()

	var story Story
	jd := json.NewDecoder(storyFile)
	err = jd.Decode(&story)
	if err != nil {
		return nil, err
	}

	return story, nil
}

// parseTemplate parses html template to be used in web server
func (app *appEnv) parseTemplate() (*template.Template, error) {
	tf, err := ioutil.ReadFile("./template.html")
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("Story").Parse(string(tf))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
