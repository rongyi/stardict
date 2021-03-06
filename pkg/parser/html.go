package parser

import (
	"github.com/jaytaylor/html2text"
)

var (
	option = &html2text.Options{
		PrettyTables: true,
	}
)

// Unhtml strip html tag to basic txt
func Unhtml(raw []byte) (string, error) {
	r := string(raw)
	return html2text.FromString(r, *option)
}
