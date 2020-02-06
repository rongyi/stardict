package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rongyi/stardict/pkg/tui"
)

var (
	dictionaryDir string

	ifoFile  string
	idxFile  string
	dictFile string
)

func bindDictFile() {
	flag.StringVar(&dictionaryDir, "d", "./", "specify the dictionary dir")
	flag.Parse()

	files, err := ioutil.ReadDir(dictionaryDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".ifo") {
			ifoFile = filepath.Join(dictionaryDir, f.Name())
		} else if strings.HasSuffix(f.Name(), ".idx") {
			idxFile = filepath.Join(dictionaryDir, f.Name())
		} else if strings.HasSuffix(f.Name(), "dict.dz") {
			dictFile = filepath.Join(dictionaryDir, f.Name())
		}
	}
	if ifoFile == "" || idxFile == "" || dictFile == "" {
		log.Fatal("empty dictionary file")
	}
}

func main() {
	bindDictFile()

	f1, err := os.Open(ifoFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(idxFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	f3, err := os.Open(dictFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f3.Close()

	engine, err := tui.NewEngine(f1, f2, f3)
	if err != nil {
		log.Fatal(err)
	}
	engine.RunWithOutput()
}
