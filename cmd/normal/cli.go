package main

import (
	"log"
	"os"
	"stardict"
)

func main() {
	f1, err := os.Open("../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo")
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open("../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	f3, err := os.Open("../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
	if err != nil {
		log.Fatal(err)
	}
	defer f3.Close()

	engine, err := stardict.NewEngine(f1, f2, f3)
	if err != nil {
		log.Fatal(err)
	}
	engine.RunWithOutput()
}
