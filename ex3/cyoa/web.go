package cyoa

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func (s *StoryWebServer) runWeb() error {
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

// StoryHandler
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

		// redirect to intro if chapter was not found
		http.Redirect(w, r, "http://localhost:8080/"+s.IntroChapter, http.StatusFound)
	}
}
