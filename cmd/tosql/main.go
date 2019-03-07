package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/parser"
	"github.com/rongyi/stardict/pkg/sql"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	dbName = "langdao.db"
)

func main() {
	// assume the existed file is a db we created
	// create db if needed
	if _, err := os.Stat(dbName); err != nil {
		if err := sql.CreateLangdaoTable("langdao.db"); err != nil {
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
	f1, err := os.Open("../../pkg/parser/testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo")
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f1.Close()
	f2, err := os.Open("../../pkg/parser/testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		log.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f2.Close()
	f3, err := os.Open("../../pkg/parser/testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
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
