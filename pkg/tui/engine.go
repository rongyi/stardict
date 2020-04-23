package tui

import (
	// "log"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/rongyi/stardict/pkg/sql"
)

const (
	FilterPrompt    string = ">>"
	DefaultSaveFile        = ".startdict.txt"
)

type Engine struct {
	db            *sql.Database
	tapp          *tview.Application
	f             io.WriteCloser
	lastWriteWord string
}

func NewEngine(dbfile string) *Engine {
	db, err := sql.NewDatabase(dbfile)
	if err != nil {
		panic(err)
	}
	f, err := openf()
	if err != nil {
		panic(err)
	}
	ret := &Engine{
		db:   db,
		tapp: tview.NewApplication(),
		f:    f,
	}

	ret.createUI()

	return ret
}

func (e *Engine) createUI() {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).SetFullScreen(true)

	inputField := tview.NewInputField().
		SetLabel(FilterPrompt).
		SetFieldWidth(80)
	inputField.SetAutocompleteFunc(func(currentText string) []string {
		if len(currentText) < 4 {
			return nil
		}
		ws, err := e.db.Prefix(currentText)
		if err != nil {
			return nil
		}
		return ws
	})

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetTextAlign(tview.AlignLeft).
		SetChangedFunc(func() {
			e.tapp.Draw()
		})
	textView.SetBorder(false)

	flex.AddItem(inputField, 1, 1, true)
	flex.AddItem(textView, 100, 1, false) // ration 40: 1 to inputField

	inputField.SetChangedFunc(func(w string) {
		if len(w) < 4 {
			return
		}
		meaning, err := e.db.Exact(w)
		if err != nil {
			return
		}
		textView.SetText(meaning)
	})
	// overide some keys
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			w := inputField.GetText()
			if e.needWrite(w) {
				meaning := textView.GetText(true)
				e.f.Write([]byte("========\n"))
				e.f.Write([]byte(w))
				e.f.Write([]byte("\n"))
				e.f.Write([]byte(meaning))
				e.f.Write([]byte("\n"))
			}
		case tcell.KeyEnter:
			inputField.Autocomplete()
		}
		return event
	})

	e.tapp.SetRoot(flex, true)
}

func openf() (io.WriteCloser, error) {
	cur, _ := user.Current()
	// cur.HomeDir
	fname := filepath.Join(cur.HomeDir, DefaultSaveFile)
	return os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0600)
}

func (e *Engine) Stop() {
	e.db.Close()
	e.f.Close()
}

func (e *Engine) Run() error {
	return e.tapp.Run()
}

func (e *Engine) needWrite(current string) bool {
	defer func() {
		e.lastWriteWord = current
	}()

	if e.lastWriteWord == "" || e.lastWriteWord != current {
		return true
	}
	return false
}
