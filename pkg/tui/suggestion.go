package tui

import ()

type SuggestionInterface interface {
	Get(keyword string) []string
}

type Suggestion struct {
}

func NewSuggestion() *Suggestion {
	return &Suggestion{}
}

func (s *Suggestion) Get(keyword string) []string {
	// TODO: connect to db
	return []string{"hello", "hell"}
}
