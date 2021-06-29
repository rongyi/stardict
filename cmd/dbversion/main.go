package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/tui"
)

func main() {
	eng := tui.NewEngine("/home/coder/go/src/github.com/rongyi/stardict/cmd/dbversion/langdao-ec.db")

	err := eng.Run()
	if err != nil {
		log.Fatal(err)
	}
}
