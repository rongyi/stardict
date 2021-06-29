package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rongyi/stardict/pkg/parser"
	"github.com/rongyi/stardict/pkg/sql"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "tosql",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "droot",
				Usage: "dictionary root",
			},
			&cli.StringFlag{
				Name:  "dbname",
				Usage: "db name",
			},
		},
		Usage:  "convert stardict format to sql database(sqlite file)",
		Action: rootCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func find(sdroot string) (string, string, string) {
	files, err := ioutil.ReadDir(sdroot)
	if err != nil {
		panic(err)
	}
	var ifoFile, idxFile, dictFile string

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if strings.HasSuffix(f.Name(), ".ifo") {
			ifoFile = filepath.Join(sdroot, f.Name())
		} else if strings.HasSuffix(f.Name(), ".idx") {
			idxFile = filepath.Join(sdroot, f.Name())
		} else if strings.HasSuffix(f.Name(), "dict.dz") {
			dictFile = filepath.Join(sdroot, f.Name())
		}
	}

	return ifoFile, idxFile, dictFile
}

func rootCmd(c *cli.Context) error {
	sdRoot := c.String("droot")
	dbName := c.String("dbname")
	if sdRoot == "" || dbName == "" {
		return errors.New("need to specify the stardict file and to sql file")
	}
	// assume the existed file is a db we created
	// create db if needed
	if _, err := os.Stat(dbName); err != nil {
		if err := sql.CreateDatabase(dbName); err != nil {
			log.Fatal(err)
		}
	}
	// create db
	db, err := sql.NewDatabase(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ifoFile, idxFile, dictFile := find(sdRoot)
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

	return d.ParseDB(db)
}
