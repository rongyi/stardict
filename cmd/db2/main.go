package main

import (
	"flag"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/tui"
)

var (
	dict string
)

func argparse() {
	flag.StringVar(&dict, "dict", "./langdao-ec.db", "specify the dict file")
	flag.Parse()
}

func main() {
	argparse()
	if dict == "" {
		log.Fatal("need a dict db")
	}
	e := tui.NewEngine(dict)
	defer e.Stop()
	if err := e.Run(); err != nil {
		log.Fatal(err)
	}
}
