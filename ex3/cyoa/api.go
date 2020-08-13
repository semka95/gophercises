package cyoa

import "html/template"

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

type StoryWebServer struct {
	Story    Story
	Template *template.Template
}

type StoryCLI struct {
	Story Story
}
