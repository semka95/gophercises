package cyoa

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

type appEnv struct {
	outputCLI bool
	storyJSON string
}

func (app *appEnv) fromArgs(args []string) error {
	fl := flag.NewFlagSet("cyoa", flag.ContinueOnError)
	fl.StringVar(
		&app.storyJSON, "story", "./gopher.json", "Path to story file in json format",
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
		//cli := StoryCLI{
		//	Story: story,
		//}
		//err := cli.runCLI()
		//if err != nil {
		//	return err
		//}
		return nil
	}

	tmpl, err := app.parseTemplate()
	if err != nil {
		return err
	}

	storyServer := &StoryWebServer{
		Story:    story,
		Template: tmpl,
	}

	err = storyServer.startWeb()
	if err != nil {
		return err
	}

	return nil
}

func (app *appEnv) parseStory() (Story, error) {
	storyFile, err := os.Open(app.storyJSON)
	if err != nil {
		return nil, err
	}

	var story Story
	jd := json.NewDecoder(storyFile)
	err = jd.Decode(&story)
	if err != nil {
		return nil, err
	}

	return story, nil
}

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

func (s *StoryWebServer) startWeb() error {
	storyHandler := s.StoryHandler()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: storyHandler,
	}

	// Launch server
	go func() {
		log.Printf("Starting the server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	// Listen for interrupt signal to close http server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("Program interrupted")

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func (s *StoryWebServer) StoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
