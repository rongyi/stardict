package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/tui"
)

func main() {
	ea := &tui.EngineAttribute{
		DefaultQuery: "",
		Monochrome:   false,
	}

	e, err := tui.NewEngine("/home/ry/go/src/github.com/rongyi/stardict/cmd/dbversion/langdao-ec.db", ea)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(run(e))

}

func run(e tui.EngineInterface) int {
	result := e.Run()
	if result.GetError() != nil {
		return 2
	}
	fmt.Printf("%s\n", result.GetQueryString())
	fmt.Printf("%s\n", result.GetContent())
	return 0
}
