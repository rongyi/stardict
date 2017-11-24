package main
import (
	"stardict"
	"log"
)

func main() {
	engine, err := stardict.NewEngine("../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo",
		"../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx",
	"../testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
	if err != nil {
		log.Fatal(err)
	}
	engine.RunWithOutput()
}
