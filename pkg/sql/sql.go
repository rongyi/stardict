package sql

import (
	"database/sql"
)

// Database is a rapper for db connection
type Database struct {
	db *sql.DB
}

// NewDatabase create a sqlite3 database
func NewDatabase(dbname string) (*Database, error) {
	ret := &Database{}

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	ret.db = db

	return ret, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func CreateLangdaoTable(dbname string) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
	create table words (id integer not null primary key autoincrement,
word text not null,
meaning text);
    `
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) InsertUniq(w, m string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into words(word, meaning) values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(w, m)

	err = tx.Commit()
	return err
}

func (d *Database) Insert(bulks [][]string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into words(word, meaning) values (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range bulks {
		stmt.Exec(record[0], record[1])
	}

	err = tx.Commit()
	return err
}
