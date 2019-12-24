package civ

import (
	"fmt"
	termbox "github.com/nsf/termbox-go"
	"os"
)

const (
	SCROLL_SIZE = 1
)

type Civ struct {
	ql        *QueryLine
	terminal  *Terminal
	table     *Table
	mode      int
	terminate bool
}

const promptDefault string = "Search> "

func NewCiv(data [][]string, dummyHeader bool) *Civ {
	c := &Civ{
		ql:       NewQueryLine(),
		terminal: NewTerminal(),
		table:    NewTable(data, dummyHeader),
	}

	return c
}

// Key event used by view mode
func (c *Civ) keyEventView(event termbox.Event) {

	switch event.Key {
	case 0:
		if !c.maybeChangeMode(event.Ch) {
			_, maxY := GetMaxXY()
			windowMoveSize := maxY / 2

			switch event.Ch {
			case 'b':
				c.table.MoveUp(windowMoveSize)
			case 'F':
				c.table.MoveDown(windowMoveSize)
			case 'f':
				c.table.MoveDown(windowMoveSize)
			case 'e':
				c.table.MoveDown(SCROLL_SIZE)
			case 'y':
				c.table.MoveUp(SCROLL_SIZE)
			case 'd':
				c.table.MoveDown(windowMoveSize / 2)
			case 'u':
				c.table.MoveUp(windowMoveSize / 2)
			case 'g':
				c.table.SetOffsetRow(0)
			case 'G':
				h := c.table.computeHeight()
				moveTo := 0
				if h > maxY {
					moveTo = h - maxY
				}
				c.table.SetOffsetRow(moveTo)
			}
		}
	case termbox.KeyArrowRight:
		c.table.MoveRight(SCROLL_SIZE)
	case termbox.KeyArrowLeft:
		c.table.MoveLeft(SCROLL_SIZE)
	case termbox.KeyArrowDown:
		c.table.MoveDown(SCROLL_SIZE)
	case termbox.KeyArrowUp:
		c.table.MoveUp(SCROLL_SIZE)
	case termbox.KeySpace:
		_, maxY := GetMaxXY()
		windowMoveSize := maxY / 2
		c.table.MoveDown(windowMoveSize)
	case termbox.KeyEnter:
		c.table.MoveDown(SCROLL_SIZE)
	}
}

// Key event used by command mode and search mode
func (c *Civ) keyEventInput(event termbox.Event) {

	// When press enter key, execute the command and clear query
	// line string while keeping mode and state of each cell (matches,
	// visibility). Also we return immediately after cleared query line
	// so that we can leave the current state.
	if event.Key == termbox.KeyEnter {
		c.ExecuteCommand()
		c.ql.ClearQuery()
		return
	}

	switch event.Key {
	case 0:
		if !c.maybeChangeMode(event.Ch) {
			// mode is not changed, input char
			c.ql.InputChar(event.Ch)

		}
	case termbox.KeySpace:
		c.ql.InputChar(' ')
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		c.ql.BackwardChar()
	case termbox.KeyCtrlD:
		c.ql.DeleteChar()
	case termbox.KeyCtrlK:
		c.ql.TruncateChars()
	case termbox.KeyArrowRight, termbox.KeyCtrlF:
		c.ql.MoveForward()
	case termbox.KeyArrowLeft, termbox.KeyCtrlB:
		c.ql.MoveBackward()
	case termbox.KeyHome, termbox.KeyCtrlA:
		c.ql.MoveToTop()
	case termbox.KeyEnd, termbox.KeyCtrlE:
		c.ql.MoveToEnd()
	default:
		// Don't execute increment search/filter when get invalid
		// key.
		return
	}

	// Do incremental search and filter
	if c.ql.mode == MODE_SEARCH {
		c.executeSearchCommand()
	} else if c.ql.mode == MODE_FILTER {
		c.executeFilterCommand()
	}
}

func (c *Civ) handleKeyEvent(event termbox.Event) {
	if c.ql.mode == MODE_VIEW {
		c.keyEventView(event)
	} else {
		// MODE_COMMAND, MODE_SEARCH
		c.keyEventInput(event)
	}
}

func (c *Civ) maybeChangeMode(key rune) bool {
	// If this is not first character, it's not mode change
	if ql := c.ql.QueryLen(); ql != 0 {
		return false
	}

	if key == '/' && c.ql.mode != MODE_SEARCH {
		c.ql.mode = MODE_SEARCH
		return true
	} else if key == '@' && c.ql.mode != MODE_COMMAND {
		c.ql.mode = MODE_COMMAND
		return true
	} else if key == '^' && c.ql.mode != MODE_FILTER {
		c.ql.mode = MODE_FILTER
		return true
	} else if key == ':' && c.ql.mode != MODE_VIEW {
		c.ql.mode = MODE_VIEW
		return true
	}

	return false
}

func (c *Civ) Draw() {
	c.terminal.Draw(c.ql, c.table)
}

func (c *Civ) Run() bool {
	for {
		// Terminate if commanded
		if c.terminate {
			return true
		}

		c.terminal.Draw(c.ql, c.table)

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:

			if ev.Key == termbox.KeyCtrlC {
				// Ctrl-c is always used to terminate
				return false
			} else if ev.Key == termbox.KeyCtrlG {
				c.ql.ClearAll()
			}

			// Dispatch key input even to the current mode
			c.handleKeyEvent(ev)

		case termbox.EventError:
			fmt.Printf("detected an error event from termbox\n")
			os.Exit(1)
			break
		default:
		}
	}

	return false
}
