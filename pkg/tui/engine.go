package tui

import (
	// "log"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/rongyi/stardict/pkg/sql"
)

const (
	DefaultY     int    = 1
	FilterPrompt string = "[Search]>"
)

type EngineInterface interface {
	Run() EngineResultInterface
	GetQuery() QueryInterface
}

type EngineResultInterface interface {
	GetQueryString() string
	GetContent() string
	GetError() error
}

type Engine struct {
	// manager        *JsonManager
	query          QueryInterface
	queryCursorIdx int
	term           *Terminal
	candidates     []string
	candidatemode  bool
	candidateidx   int
	contentOffset  int
	queryConfirm   bool
	// prettyResult   bool
	db *sql.Database
}

type EngineAttribute struct {
	DefaultQuery string
	Monochrome   bool
	// PrettyResult bool
}

func NewEngine(liteFile string, ea *EngineAttribute) (EngineInterface, error) {
	db, err := sql.NewDatabase(liteFile)
	if err != nil {
		return nil, err
	}
	e := &Engine{
		term:          NewTerminal(FilterPrompt, DefaultY, ea.Monochrome),
		query:         NewQuery([]rune(ea.DefaultQuery)),
		candidates:    []string{},
		candidatemode: false,
		candidateidx:  0,
		contentOffset: 0,
		queryConfirm:  false,
		// prettyResult:  ea.PrettyResult,
		db: db,
	}
	e.queryCursorIdx = e.query.Length()
	return e, nil
}

type EngineResult struct {
	content string
	qs      string
	err     error
}

func (er *EngineResult) GetQueryString() string {
	return er.qs
}

func (er *EngineResult) GetContent() string {
	return er.content
}
func (er *EngineResult) GetError() error {
	return er.err
}

func (e *Engine) GetQuery() QueryInterface {
	return e.query
}

func (e *Engine) Run() EngineResultInterface {
	defer e.db.Close()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	var contents []string

	for {

		if e.query.StringGet() == "" {
			e.queryCursorIdx = e.query.Length()
		}

		bl := len(contents)
		contents = e.getContents()
		e.setCandidateData()
		e.queryConfirm = false
		if bl != len(contents) {
			e.contentOffset = 0
		}

		ta := &TerminalDrawAttributes{
			Query:           e.query.StringGet(),
			Contents:        contents,
			CandidateIndex:  e.candidateidx,
			ContentsOffsetY: e.contentOffset,
			Candidates:      e.candidates,
			CursorOffset:    e.query.IndexOffset(e.queryCursorIdx),
		}
		err = e.term.Draw(ta)
		if err != nil {
			panic(err)
		}

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case 0:
				e.inputChar(ev.Ch)
			case termbox.KeySpace:
				e.inputChar(' ')
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				e.deleteChar()
			case termbox.KeyTab:
				e.tabAction()
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
			case termbox.KeyCtrlU:
				e.deleteLineQuery()
			case termbox.KeyCtrlW:
				e.clearInput()
			case termbox.KeyEsc:
				e.escapeCandidateMode()
			case termbox.KeyEnter:
				if !e.candidatemode {
					return &EngineResult{
						content: strings.Join(contents, "\n"),
						qs:      ta.Query,
						err:     nil,
					}
				}
				e.confirmCandidate()
			case termbox.KeyCtrlC:
				return &EngineResult{}
			default:
			}
		case termbox.EventError:
			panic(ev.Err)
			break
		default:
		}
	}
}

func (e *Engine) tabAction() {
	if !e.candidatemode {
		e.candidatemode = true
	} else {
		e.candidateidx = e.candidateidx + 1
	}
	e.queryCursorIdx = e.query.Length()
}

func (e *Engine) setCandidateData() {
	if l := len(e.candidates); l >= 1 {
		if e.candidateidx >= l {
			e.candidateidx = 0
		}
	} else {
		e.candidatemode = false
	}
	if !e.candidatemode {
		e.candidateidx = 0
		e.candidates = []string{}
	}
}

// getContents has side effect!
// it set candidates
func (e *Engine) getContents() []string {
	var contents []string
	input := e.query.StringGet()

	if e.queryConfirm || len(input) < 5 {
		e.candidates = []string{}
		if explain, err := e.db.Exact(input); err == nil {
			contents = strings.Split(explain, "\n")
		}
	} else {
		// too much candidates, we do nothing
		pres, err := e.db.Prefix(input)
		if err != nil {
			e.candidates = []string{}
		} else {
			e.candidates = pres
		}
	}

	return contents
}

func (e *Engine) confirmCandidate() {
	// delete all and put the candidate on
	e.clearInput()

	_ = e.query.StringAdd(e.candidates[e.candidateidx])
	e.queryCursorIdx = e.query.Length()
	e.queryConfirm = true
}

func (e *Engine) deleteChar() {
	if i := e.queryCursorIdx - 1; i >= 0 {
		_ = e.query.Delete(i)
		e.queryCursorIdx--
	}

}

func (e *Engine) deleteLineQuery() {
	_ = e.query.StringSet("")
	e.queryCursorIdx = 0
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

// func (e *Engine) toggleKeyFinished() {
// 	e.keyFinished = !e.keyFinished
// }

func (e *Engine) clearInput() {
	// just reset to empty
	e.query.Set([]rune(""))
	e.queryCursorIdx = e.query.Length()
}

func (e *Engine) escapeCandidateMode() {
	e.candidatemode = false
}
func (e *Engine) inputChar(ch rune) {
	_ = e.query.Insert([]rune{ch}, e.queryCursorIdx)
	e.queryCursorIdx++
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

func (e *Engine) moveCursorWordBackwark() {
}
func (e *Engine) moveCursorWordForward() {
}
func (e *Engine) moveCursorToTop() {
	e.queryCursorIdx = 0
}
func (e *Engine) moveCursorToEnd() {
	e.queryCursorIdx = e.query.Length()
}
