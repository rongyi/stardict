package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/parser"
	"github.com/rongyi/stardict/pkg/sql"
	log "github.com/sirupsen/logrus"
)

var (
	dbName        string
	dictionaryDir string

	ifoFile  string
	idxFile  string
	dictFile string
)

func bindDictFile() {
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
}

func init() {
	flag.StringVar(&dbName, "db", "", "specify db name")
	flag.StringVar(&dictionaryDir, "dict", "", "specify dictionary dir")
	flag.Parse()
}

func main() {
	if dbName == "" || dictionaryDir == "" {
		log.Fatal("need db name and dictionary dir")
	}
	bindDictFile()
	// assume the existed file is a db we created
	// create db if needed
	if _, err := os.Stat(dbName); err != nil {
		if err := sql.CreateLangdaoTable(dbName); err != nil {
			log.Fatal(err)
		}
	}
	// create db
	db, err := sql.NewDatabase(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// create dictionary
	f1, err := os.Open(ifoFile)
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f1.Close()
	f2, err := os.Open(idxFile)
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f2.Close()
	f3, err := os.Open(dictFile)
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f3.Close()

	d, err := parser.NewDictionary(f1, f2, f3)
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}

	// dump word to dictionary
	if err := d.DumpLangdao(db); err != nil {
		log.Fatal(err)
	} else {
		log.Infoln("dump to db success!")
	}
}
