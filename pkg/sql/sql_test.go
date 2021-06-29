package sql

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestPrefix(t *testing.T) {
	a := require.New(t)
	db, err := NewDatabase("./langdao-ec.db")
	a.Nil(err, "fail to connect to db")
	defer db.Close()
	prefix, err := db.Prefix("hell")
	a.Nil(err, "fail")
	t.Log(prefix)
}

func TestExact(t *testing.T) {
	a := require.New(t)
	db, err := NewDatabase("./langdao-ec.db")
	a.Nil(err, "fail to connect to db")
	defer db.Close()
	prefix, err := db.Exact("hello program")
	a.Nil(err, "fail")
	t.Log(prefix)
}
