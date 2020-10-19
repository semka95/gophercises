package cyoa

import "html/template"

// Story represents a Choose Your Own Adventure story
type Story map[string]Chapter

// Chapter represents a CYOA story chapter
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option represents a choice offered at the end of a story
// chapter
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

// StoryWebServer contains data needed to run CYOA web server
type StoryWebServer struct {
	Story        Story
	Template     *template.Template
	IntroChapter string
}

// StoryCLI contains data needed to run CYOA cli
type StoryCLI struct {
	Story        Story
	IntroChapter string
}
