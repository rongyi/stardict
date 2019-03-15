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

func (d *Database) Prefix(key string) ([]string, error) {
	stmt, err := d.db.Prepare("select word from words where word like ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(key + `%`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []string{}

	for rows.Next() {
		var m string
		err = rows.Scan(&m)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (d *Database) Exact(key string) (string, error) {
	stmt, err := d.db.Prepare("select meaning from words where word = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var ret string
	row := stmt.QueryRow(key)

	err = row.Scan(&ret)
	if err != nil {
		return ret, err
	}

	return ret, nil
}
