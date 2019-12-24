package civ

import (
	"fmt"
	"strconv"
)

type Table struct {
	header Row

	// Maximum length of each columns data including both header
	// and contents
	maxLen []int

	// A list of column number that is disabled to display
	disabledCols []int

	// Table contents excluding header
	contents []Row

	// View setting, starting from 0
	offsetCol int
	offsetRow int

	outputStdout bool
}

type Row struct {
	cols       []Cell
	hasMatched bool
	isVisible  bool
}

type Cell struct {
	data       string
	matchBegin int
	matchEnd   int
	width      int // formatted size, available on only header
}

func NewTable(indata [][]string, dummyHeader bool) *Table {
	t := &Table{
		offsetCol: 0,
		offsetRow: 0,
	}
	nCols := len(indata[0])

	// Load header data
	t.maxLen = make([]int, nCols)
	for i, cell := range indata[0] {
		c := Cell{
			matchBegin: -1,
			matchEnd:   -1,
		}

		if dummyHeader {
			// make dummy header
			c.data = "col_" + strconv.Itoa(i)
			t.header.cols = append(t.header.cols, c)
		} else {
			c.data = cell
			t.header.cols = append(t.header.cols, c)
		}

		// initialize the maximum lenght of column by
		// the header data
		t.maxLen[i] = len(cell)
	}

	// Load table data
	csvdata := indata[1:]
	if dummyHeader {
		// If using dummy header, csv data starts from indata[0]
		csvdata = indata
	}
	for _, row := range csvdata {
		var r Row
		r.isVisible = true

		for i, cell := range row {
			c := &Cell{
				data:       cell,
				matchBegin: -1,
				matchEnd:   -1,
			}
			r.cols = append(r.cols, *c)

			if t.maxLen[i] < len(cell) {
				t.maxLen[i] = len(cell)
			}
		}
		t.contents = append(t.contents, r)
	}

	return t
}

// Return true and position if the given name is in col names
func (t *Table) FindColName(name string) (int, bool) {
	for i, c := range t.header.cols {
		if c.data == name {
			return i, true
		}
	}

	return -1, false
}

func (t *Table) NEnabledCols() int {
	return len(t.header.cols) - len(t.disabledCols)
}

func (t *Table) IsColEnabled(colNum int) bool {
	for _, n := range t.disabledCols {
		if n == colNum {
			return false
		}
	}

	return true
}

// Add the idx'th column to disabled column list.
// That is, make it invisible
func (t *Table) AddDisabledCol(idx int) {
	t.disabledCols = append(t.disabledCols, idx)
}

// Remove the idx'th column from disabled column list.
// That is, make it visible
func (t *Table) RemoveDisabledCol(idx int) {
	res := []int{}

	for _, c := range t.disabledCols {
		if c != idx {
			res = append(res, c)
		}
	}

	t.disabledCols = res
}

func (t *Table) SetRowVisibility(rowIdx int, visible bool) {
	t.contents[rowIdx].isVisible = visible
}

func (t *Table) SetMatched(r int, c int, b int, e int) {
	t.contents[r].cols[c].matchBegin = b
	t.contents[r].cols[c].matchEnd = e
}

func (t *Table) ResetDisabledCol() {
	t.disabledCols = []int{}
}

func (t *Table) ResetVisibility() {
	for i := 0; i < len(t.contents); i++ {
		t.contents[i].isVisible = true
	}
}

// Return the total size of visible cols. Note that c.width is the size
// of formatted cell, i.e. includes padding.
func (t *Table) computeWidth() (width int) {
	width = 0
	for i, c := range t.header.cols {
		// Skip until offset
		if i < t.offsetCol {
			continue
		}

		// Skip if this col is disabled
		if !t.IsColEnabled(i) {
			continue
		}

		width += c.width
	}

	return width
}

// Return the total size of visible rows.
func (t *Table) computeHeight() (height int) {
	height = 0
	for i, r := range t.contents {
		// Skip until offset
		if i < t.offsetRow {
			continue
		}

		if !r.isVisible {
			continue
		}

		height++
	}

	// always include both header and line
	return height + 2
}

// 0: up, 1: right, 2:down, 3: left
func (t *Table) isMovable(direction int) bool {
	maxX, maxY := GetMaxXY()
	x := t.computeWidth()
	y := t.computeHeight()

	if x >= maxX && direction == 1 {
		// Movable to right and want to move right
		return true
	} else if t.offsetCol > 0 && direction == 3 {
		// Viewing right part and want to move left
		return true
	} else if y > maxY && direction == 2 {
		// Movable to down ana want ot move down
		return true
	} else if t.offsetRow > 0 && direction == 0 {
		// Viewing bottom part and want to move up
		return true
	}
	return false
}

func (t *Table) SetOffsetRow(r int) {
	t.offsetRow = r
}

func (t *Table) MoveRight(move int) {
	if !t.isMovable(1) {
		return
	}

	t.offsetCol += move

	if t.offsetCol > len(t.header.cols) {
		t.offsetRow = len(t.header.cols)
	}
}

func (t *Table) MoveLeft(move int) {
	if !t.isMovable(3) {
		return
	}

	t.offsetCol -= move

	if t.offsetCol < 0 {
		t.offsetCol = 0
	}
}

func (t *Table) MoveUp(move int) {
	if !t.isMovable(0) {
		return
	}

	// Skip invisible rows
	for i := t.offsetRow; !t.contents[i].isVisible && i > 0; i-- {
		if !t.contents[i].isVisible {
			move++
		}
	}
	t.offsetRow -= move

	if t.offsetRow < 0 {
		t.offsetRow = 0
	}
}

func (t *Table) MoveDown(move int) {
	if !t.isMovable(2) {
		return
	}

	// Skip invisible rows
	for i := t.offsetRow; !t.contents[i].isVisible && i < len(t.contents); i++ {
		if !t.contents[i].isVisible {
			move++
		}
	}
	t.offsetRow += move

	if t.offsetRow > len(t.contents) {
		t.offsetRow = len(t.contents)
	}
}

func (t *Table) Debugdump() {
	fmt.Println("---- Header ----")
	for _, cell := range t.header.cols {
		fmt.Print(cell.data, " ")
	}
	fmt.Println("")
	fmt.Println("---- Contents -----------")
	for _, row := range t.contents {
		for _, cell := range row.cols {
			fmt.Print(cell.data, " ")
		}
		fmt.Println("")
	}

	fmt.Println("maxLen: ", t.maxLen)
}
