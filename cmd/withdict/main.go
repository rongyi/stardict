package main

import (
	"log"
	"bytes"
	"stardict"
)

func main() {
	a, err := Asset("langdao-ec-gb.ifo")
	if err != nil {
		log.Fatal(err)
	}
	b, err := Asset("langdao-ec-gb.idx")
	if err != nil {
		log.Fatal(err)
	}
	c, err := Asset("langdao-ec-gb.dict.dz")


	r1 := bytes.NewReader(a)
	r2 := bytes.NewReader(b)
	r3 := bytes.NewReader(c)

	engine, err := stardict.NewEngine(r1, r2, r3)
	if err != nil {
		log.Fatal(err)
	}
	engine.RunWithOutput()
}
