package tui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {
	var assert = assert.New(t)

	assert.True(validate([]rune("hello")) == "hello")
	assert.True(validate([]rune("HELLO")) == "HELLO")
	assert.True(validate([]rune("123")) == "123", "fail")
}

func TestNewQuery(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("name")
	q := NewQuery(v)

	assert.Equal(*q.query, []rune("name"))
	assert.Equal(*q.complete, []rune(""))
}

func TestNewQueryWithInvalidQuery(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("  ")
	q := NewQuery(v)

	assert.Equal(*q.query, []rune(""))
	assert.Equal(*q.complete, []rune(""))
}

func TestNewQueryWithString(t *testing.T) {
	var assert = assert.New(t)

	q := NewQueryWithString("a   b")

	assert.Equal(*q.query, []rune("a b"))

	q = NewQueryWithString("a ")
	assert.Equal(*q.query, []rune("a "))

	q = NewQueryWithString("a b a ")
	assert.Equal(*q.query, []rune("a b a "))
}

func TestNewQueryWithStringWithInvalidQuery(t *testing.T) {
	var assert = assert.New(t)

	q := NewQueryWithString("   ")

	assert.Equal(*q.query, []rune(""))
	assert.Equal(*q.complete, []rune(""))
}

func TestQueryGet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("test")
	q := NewQuery(v)

	assert.Equal(q.Get(), []rune("test"))
}


func TestQueryGetChar(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("test")
	q := NewQuery(v)

	assert.Equal('e', q.GetChar(1))
	assert.Equal('t', q.GetChar(3))
	assert.Equal(rune(0), q.GetChar(-1))
	assert.Equal(rune(0), q.GetChar(6))

}

func TestQuerySet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("hello")
	q := NewQuery(v)

	assert.Equal([]rune("world"), q.Set([]rune("world")))
	assert.Equal("", string(q.Set([]rune(""))))
}


func TestQueryAdd(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("hello")
	q := NewQuery(v)

	assert.Equal(q.Add([]rune("world")), []rune("helloworld"))
}
func TestQueryInsert(t *testing.T) {
	var assert = assert.New(t)
	v := []rune("hello")
	q := NewQuery(v)

	// assert.Equal([]rune("whello"), q.Insert([]rune("w"), 0))
	assert.Equal([]rune("helloxxx"), q.Insert([]rune("xxx"), 5))
	assert.Equal([]rune("hxxelloxxx"), q.Insert([]rune("xx"), 1))
	assert.Equal([]rune("xxhxxelloxxx"), q.Insert([]rune("xx"), 0))
	// assert.Equal([]rune("wwhello"), q.Insert([]rune("w"), 1))
	// assert.Equal([]rune(".whello.world"), q.Insert([]rune("w"), 1))
	// assert.Equal([]rune(".wwhello.world"), q.Insert([]rune("w"), 1))
	// assert.Equal([]rune(".wwhello.world"), q.Insert([]rune("."), 1))
	// assert.Equal([]rune(".wwh.ello.world"), q.Insert([]rune("."), 4))
	// assert.Equal([]rune(".wwh.ello.worldg"), q.Insert([]rune("g"), 15))
	// assert.Equal([]rune(".wwh.ello.worldg"), q.Insert([]rune("a"), 20))
}
func TestQueryStringInsert(t *testing.T) {
	var assert = assert.New(t)
	q := NewQueryWithString("helloworld")

	assert.Equal("whelloworld", q.StringInsert("w", 0))
	assert.Equal("wwhelloworld", q.StringInsert("w", 1))
	// assert.Equal("wwhelloworld", q.StringInsert("w", 1))
	// assert.Equal("wwhelloworld", q.StringInsert(".", 1))
	// assert.Equal("wwhelloworld", q.StringInsert(".", 4))
	// assert.Equal("wwhelloworlda", q.StringInsert("a", 15))
	// assert.Equal("wwhelloworlda", q.StringInsert("a", 20))
}

func TestQueryClear(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("test")
	q := NewQuery(v)

	assert.Equal(q.Clear(), []rune(""))
}

func TestQueryDelete(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("helloworld")
	q := NewQuery(v)

	assert.Equal([]rune("d"), q.Delete(-1))
	assert.Equal([]rune("helloworl"), q.Get())
	assert.Equal([]rune("l"), q.Delete(-1))
	assert.Equal([]rune("hellowor"), q.Get())
	assert.Equal([]rune("or"), q.Delete(-2))
	assert.Equal([]rune("hellow"), q.Get())
	assert.Equal([]rune("hellow"), q.Delete(-6))
	assert.Equal([]rune(""), q.Get())

	q = NewQuery([]rune("helloworld"))
	assert.Equal([]rune("h"), q.Delete(0))
	assert.Equal([]rune("elloworld"), q.Get())
	assert.Equal([]rune("l"), q.Delete(1))
	assert.Equal([]rune("eloworld"), q.Get())
	assert.Equal([]rune{}, q.Delete(9))
	assert.Equal([]rune("eloworld"), q.Get())
	assert.Equal([]rune("o"), q.Delete(2))
	assert.Equal([]rune("elworld"), q.Get())
}

func TestQueryStringGet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("test")
	q := NewQuery(v)

	assert.Equal(q.StringGet(), "test")
}

func TestQueryStringSet(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("hello")
	q := NewQuery(v)

	assert.Equal(q.StringSet("world"), "world")
}

func TestQueryStringAdd(t *testing.T) {
	var assert = assert.New(t)

	v := []rune("hello")
	q := NewQuery(v)

	assert.Equal(q.StringAdd("world"), "helloworld")
}
