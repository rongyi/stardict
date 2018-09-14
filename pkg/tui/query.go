package stardict

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type Query struct {
	query *[]rune
}

func NewQuery(query []rune) *Query {
	q := &Query{}
	_ = q.Set(query)
	return q
}

func (q *Query) Set(query []rune) []rune {
	q.query = &query

	return q.Get()
}

func (q *Query) Get() []rune {
	return *q.query
}

func (q *Query) Length() int {
	return len(q.Get())
}

func (q *Query) StringGet() string {
	rawStr := string(q.Get())
	return strings.ToLower(rawStr)
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

func (q *Query) Insert(query []rune, idx int) []rune {
	qq := q.Get()
	if idx == 0 {
		qq = append(query, qq...)
	} else if idx > 0 && len(qq) >= idx {
		_q := make([]rune, idx+len(query)-1)
		copy(_q, qq[:idx])
		qq = append(append(_q, query...), qq[idx:]...)
	}
	return q.Set(qq)
}

func (q *Query) Delete(i int) []rune {
	var d []rune
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
	} else if i == 0 {
		d = []rune{}
		qq = qq[1:]
	} else if i > 0 && i < lastIdx {
		d = []rune{qq[i]}
		qq = append(qq[:i], qq[i+1:]...)
	}
	_ = q.Set(qq)
	return d
}
