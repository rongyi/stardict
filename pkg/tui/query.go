package tui

import (
	// "log"
	"strings"

	"github.com/mattn/go-runewidth"
)

type QueryInterface interface {
	Get() []rune
	Set(query []rune) []rune
	Insert(query []rune, idx int) []rune
	Add(query []rune) []rune
	Delete(i int) []rune
	Clear() []rune
	Length() int
	IndexOffset(int) int
	GetChar(int) rune
	StringGet() string
	StringSet(query string) string
	StringInsert(query string, idx int) string
	StringAdd(query string) string
}

type Query struct {
	query    *[]rune
	complete *[]rune
}

func NewQuery(query []rune) *Query {
	q := &Query{
		query:    &[]rune{},
		complete: &[]rune{},
	}
	_ = q.Set(query)
	return q
}
func NewQueryWithString(query string) *Query {
	return NewQuery([]rune(query))
}

func (q *Query) Get() []rune {
	return *q.query
}

func (q *Query) GetChar(idx int) rune {
	var r rune = 0
	qq := q.Get()
	if l := len(qq); l > idx && idx >= 0 {
		r = qq[idx]
	}
	return r
}

func (q *Query) Length() int {
	return len(q.Get())
}

func (q *Query) IndexOffset(i int) int {
	o := 0
	if l := q.Length(); i >= l {
		o = runewidth.StringWidth(q.StringGet())
	} else if i >= 0 && i < l {
		o = runewidth.StringWidth(string(q.Get()[:i]))
	}
	return o
}

func (q *Query) Set(query []rune) []rune {
	str := validate(query)
	a := []rune(str)
	q.query = &a

	return q.Get()
}

func (q *Query) Insert(query []rune, idx int) []rune {
	qq := q.Get()
	if idx == 0 {
		qq = append(query, qq...)
	} else if idx > 0 && len(qq) >= idx {
		_q := make([]rune, idx)
		copy(_q, qq[:idx])
		qq = append(append(_q, query...), qq[idx:]...)
	}
	return q.Set(qq)
}

func (q *Query) StringInsert(query string, idx int) string {
	return string(q.Insert([]rune(query), idx))
}

func (q *Query) Add(query []rune) []rune {
	return q.Set(append(q.Get(), query...))
}

func (q *Query) Delete(i int) []rune {
	d := []rune{}
	qq := q.Get()
	lastIdx := len(qq)
	if i < 0 {
		if lastIdx+i >= 0 {
			d = qq[lastIdx+i:]
			qq = qq[0 : lastIdx+i]
		} else {
			d = qq
			qq = qq[0:0]
		}
	} else if i >= 0 && i < lastIdx {
		d = []rune{qq[i]}
		qq = append(qq[:i], qq[i+1:]...)
	}
	_ = q.Set(qq)
	return d
}

func (q *Query) Clear() []rune {
	return q.Set([]rune(""))
}

func (q *Query) StringGet() string {
	return string(q.Get())
}

func (q *Query) StringSet(query string) string {
	return string(q.Set([]rune(query)))
}

func (q *Query) StringAdd(query string) string {
	return string(q.Add([]rune(query)))
}

func validate(r []rune) string {
	s := string(r)
	if s == "" {
		return ""
	}
	sec := strings.Fields(s)

	if len(sec) == 0 {
		return ""
	}

	ret := strings.Join(sec, " ")
	if s[len(s) - 1] == ' ' {
		ret += " "
	}
	return ret
}
