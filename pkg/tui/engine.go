package stardict

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/rongyi/stardict/pkg/dump"
	"github.com/rongyi/stardict/pkg/parser"
)

const (
	DefaultY     int    = 1
	FilterPrompt string = "[Word]> "
)

type Engine struct {
	queryCursorIdx int
	query          *Query
	term           *Terminal
	contentOffset  int
	input          []string
	dict           *parser.Dictionary
	prevSave       string
	saveFile       *os.File
}

func NewEngine(ifo, idx, d io.Reader) (*Engine, error) {
	var fflow []string
	dict, err := parser.NewDictionary(ifo, idx, d)
	if err != nil {
		return nil, err
	}
	saveFD, err := dump.OpenFile()
	if err != nil {
		return nil, err
	}

	e := &Engine{
		queryCursorIdx: 0,
		query:          NewQuery([]rune("")),
		term:           NewTerminal(FilterPrompt, DefaultY),
		contentOffset:  0,
		input:          fflow,
		dict:           dict,
		saveFile:       saveFD,
	}
	e.queryCursorIdx = e.query.Length()

	return e, nil
}

func (e *Engine) inputChar(ch rune) {
	_ = e.query.Insert([]rune{ch}, e.queryCursorIdx)
	e.queryCursorIdx++
}

func (e *Engine) deleteChar() {
	if i := e.queryCursorIdx - 1; i >= 0 {
		_ = e.query.Delete(i)
		e.queryCursorIdx--
	}
}

func (e *Engine) clearChar() {
	for i := e.queryCursorIdx - 1; i >= 0; i-- {
		_ = e.query.Delete(i)
		e.queryCursorIdx--
	}
}

func (e *Engine) moveCursorBackward() {
	if i := e.queryCursorIdx - 1; i >= 0 {
		e.queryCursorIdx--
	}
}

func (e *Engine) moveCursorForward() {
	if e.query.Length() > e.queryCursorIdx {
		e.queryCursorIdx++
	}
}
func (e *Engine) moveCursorToTop() {
	e.queryCursorIdx = 0
}
func (e *Engine) moveCursorToEnd() {
	e.queryCursorIdx = e.query.Length()
}

func (e *Engine) scrollToBelow() {
	e.contentOffset++
}

func (e *Engine) scrollToAbove() {
	if o := e.contentOffset - 1; o >= 0 {
		e.contentOffset = o
	}
}

func (e *Engine) scrollToBottom(rownum int) {
	e.contentOffset = rownum - 1
}

func (e *Engine) scrollToTop() {
	e.contentOffset = 0
}

func (e *Engine) getContents() []string {
	word := e.query.StringGet()
	if word == "" {
		return e.input
	}
	return e.dict.GetFormatedMeaning(word)
}

func (e *Engine) Run() []string {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	defer e.saveFile.Close()

	var contents []string
mainloop:
	for {
		bl := len(contents)
		contents = e.getContents()
		if bl != len(contents) {
			e.contentOffset = 0
		}

		ta := &TerminalAttributes{
			Query:           e.query.StringGet(),
			CursorOffset:    e.query.IndexOffset(e.queryCursorIdx),
			Contents:        contents,
			ContentsOffsetY: e.contentOffset,
		}
		err = e.term.Draw(ta)
		if err != nil {
			panic(err)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 0, termbox.KeySpace:
				e.inputChar(ev.Ch)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				e.deleteChar()
			case termbox.KeyCtrlU:
				e.clearChar()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				e.moveCursorBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				e.moveCursorForward()
			case termbox.KeyHome, termbox.KeyCtrlA:
				e.moveCursorToTop()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				e.moveCursorToEnd()
			case termbox.KeyCtrlK:
				e.scrollToAbove()
			case termbox.KeyCtrlJ:
				e.scrollToBelow()
			case termbox.KeyCtrlG:
				e.scrollToBottom(len(contents))
			case termbox.KeyCtrlT:
				e.scrollToTop()
			case termbox.KeyCtrlS:
				e.save()
			case termbox.KeyCtrlC:
				break mainloop
			}
		case termbox.EventError:
			break mainloop
		}
	}

	return contents
}

func (e *Engine) RunWithOutput() int {
	filterOutput := e.Run()
	if len(filterOutput) > 0 {
		fmt.Println(strings.Join(filterOutput, "\n"))
	}

	return 0
}

func (e *Engine) save() {
	word := e.query.StringGet()
	// duplicate save
	if word == e.prevSave {
		return
	}
	if e.prevSave == "" && word != "" {
		e.prevSave = word
	}
	meaning := e.dict.GetFormatedMeaning(word)
	if word != "" && len(meaning) != 0 {
		content := word + "\n" + strings.Join(meaning, "\n")
		e.dump(content)
	}
}

func (e *Engine) dump(s string) error {
	io := bufio.NewWriter(e.saveFile)
	io.Write(dump.WordSeperator)
	io.WriteByte('\n')
	io.WriteString(s)
	io.WriteByte('\n')
	io.Flush()
	return nil
}
