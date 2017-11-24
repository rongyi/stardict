package stardict

import (
	"github.com/jaytaylor/html2text"
)

var (
	option = &html2text.Options{
		PrettyTables: true,
	}
)

func Unhtml(raw []byte) (string, error) {
	r := string(raw)
	return html2text.FromString(r, *option)
}
