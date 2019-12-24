package civ

import (
	"fmt"
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

type Terminal struct {
	prompt string
}

func NewTerminal() *Terminal {
	t := &Terminal{
		prompt: ":",
	}

	return t
}

func (t *Terminal) Draw(ql *QueryLine, tb *Table) error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	drawQueryLine(ql)

	drawTable(tb)

	termbox.Flush()

	return nil
}

// This function creates a list of termbox.Cell from rowIdx'th row of
// the table.
func formatRow(tb *Table, rowIdx int, isHeader bool) []termbox.Cell {
	var cells []termbox.Cell
	var row Row

	if isHeader {
		row = tb.header
	} else {
		row = tb.contents[rowIdx]
	}

	nWritten := 0
	for i, cell := range row.cols {
		field := formatCell(cell, tb.maxLen[i], isHeader)

		if isHeader {
			// Set length of formatted column including a delimiter
			tb.header.cols[i].width = len(field) + 1
		}

		// Skip if this col is disabled
		if !tb.IsColEnabled(i) {
			continue
		}

		// Skip until offset
		if i < tb.offsetCol {
			continue
		}

		cells = append(cells, field...)
		nWritten++

		// Don't need delimiter after the last column
		if nWritten != tb.NEnabledCols() {
			cells = append(cells, termbox.Cell{
				Ch: '|',
				Fg: termbox.ColorDefault,
				Bg: termbox.ColorDefault,
			})
		}
	}

	return cells
}

func formatCell(cell Cell, maxSize int, isHeader bool) []termbox.Cell {
	var tcells []termbox.Cell
	var beforeSpaces int
	var afterSpaces int

	if isHeader {
		beforeSpaces = ((maxSize + 2) - len(cell.data) + 1) / 2
		afterSpaces = maxSize + 2 - len(cell.data) - beforeSpaces
	} else {
		beforeSpaces = 1
		afterSpaces = ((maxSize + 1) - len(cell.data))
	}

	// pad before data by space
	tcells = padSpaces(tcells, beforeSpaces)

	// Prepare termbox cells for table's cells
	for i, r := range cell.data {
		// Determine cell color
		fg := termbox.ColorDefault
		bg := termbox.ColorDefault

		if cell.matchBegin != -1 &&
			cell.matchBegin <= i &&
			cell.matchEnd > i {
			fg = termbox.ColorBlack | termbox.AttrBold
			bg = termbox.ColorCyan | termbox.AttrBold
		}

		tcells = append(tcells, termbox.Cell{
			Ch: r,
			Fg: fg,
			Bg: bg,
		})
	}

	// pad after data by space
	tcells = padSpaces(tcells, afterSpaces)

	return tcells
}

func drawTable(tb *Table) {
	var lineCells []termbox.Cell

	// draw header cells
	headerCells := formatRow(tb, 0, true)

	if tb.outputStdout {
		outputCells(headerCells)
	} else {
		drawCells(0, 1, headerCells)
	}

	// draw line cells
	for _, c := range headerCells {
		r := '-'
		if c.Ch == '|' {
			r = '+'
		}
		lineCells = append(lineCells, termbox.Cell{
			Ch: r,
			Fg: termbox.ColorDefault,
			Bg: termbox.ColorDefault,
		})
	}
	if tb.outputStdout {
		outputCells(lineCells)
	} else {
		drawCells(0, 2, lineCells)
	}

	// draw table
	nRows := 0
	for i, r := range tb.contents {
		// Skip invisible row
		if !r.isVisible {
			continue
		}
		// Skip until offset
		if i < tb.offsetRow {
			continue
		}

		cells := formatRow(tb, i, false)

		// Remember the length of row
		if tb.outputStdout {
			outputCells(cells)
		} else {
			drawCells(0, 3+nRows, cells)
		}
		nRows++
	}
}

// Return the padded cell
func padSpaces(cells []termbox.Cell, nSpaces int) []termbox.Cell {
	// Add spaces before data
	for _i := 0; _i < nSpaces; _i++ {
		cells = append(cells, termbox.Cell{
			Ch: ' ',
			Fg: termbox.ColorDefault,
			Bg: termbox.ColorDefault,
		})
	}

	return cells
}

// Function to draw both prompt and query line
func drawQueryLine(ql *QueryLine) {
	var cells []termbox.Cell

	if ql == nil {
		return
	}

	// Print prefix
	modeRune := ':'
	if ql.mode == MODE_SEARCH {
		modeRune = '/'
	} else if ql.mode == MODE_COMMAND {
		modeRune = '@'
	} else if ql.mode == MODE_FILTER {
		modeRune = '^'
	} else if ql.mode == MODE_VIEW {
		modeRune = ':'
	}
	cells = append(cells, termbox.Cell{
		Ch: modeRune,
		Fg: termbox.ColorDefault,
		Bg: termbox.ColorDefault,
	})

	// Print query line
	for _, r := range ql.query {
		cells = append(cells, termbox.Cell{
			Ch: r,
			Fg: termbox.ColorDefault,
			Bg: termbox.ColorDefault,
		})
	}

	drawCells(0, 0, cells)
	termbox.SetCursor(1+ql.curCursor, 0)
}

// Actual drawing the given cells to the terminal
func drawCells(x int, y int, cells []termbox.Cell) {
	i := 0
	maxX, _ := termbox.Size()

	for _, c := range cells {
		if i >= maxX-2 {
			// Before reaching the right edge we show '.'
			// instead of data to indicate that the data
			// is continuing, like 'hel..'.
			termbox.SetCell(x+i, y, rune('.'), c.Fg, c.Bg)
		} else {
			termbox.SetCell(x+i, y, c.Ch, c.Fg, c.Bg)
		}

		w := runewidth.RuneWidth(c.Ch)
		if w == 0 || w == 2 && runewidth.IsAmbiguousWidth(c.Ch) {
			w = 1
		}

		i += w
	}
}

func GetMaxXY() (maxX int, maxY int) {
	return termbox.Size()

}

func outputCells(cells []termbox.Cell) {
	for _, c := range cells {
		fmt.Printf("%s", string(c.Ch))
	}
	fmt.Println("")
}
