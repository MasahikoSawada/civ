package civ

import (
	"strings"
)

func (c *Civ) ExecuteCommand() {
	// Quick exit if there is no actual query
	if len(c.ql.query) <= 0 {
		return
	}

	a := strings.Fields(string(c.ql.query))

	command := a[0]
	vargs := a[1:]

	if strings.HasPrefix("show", command) {
		c.executeShowCommand(vargs)
	} else if strings.HasPrefix("show_only", command) {
		c.executeShowOnlyCommand(vargs)
	} else if strings.HasPrefix("hide", command) {
		c.executeHideCommand(vargs)
	} else if strings.HasPrefix("reset", command) {
		c.executeResetCommand(vargs)
	} else if strings.HasPrefix("exit", command) {
		c.executeExitCommand(vargs)
	}
}

func (c *Civ) executeExitCommand(vargs []string) {
	c.table.outputStdout = true
	c.terminate = true
}

// Hide the given columns
func (c *Civ) executeHideCommand(vargs []string) {
	for _, a := range vargs {
		if i, ok := c.table.FindColName(a); ok {
			c.table.AddDisabledCol(i)
		}
	}
}

// Show only the given colums by hidding other columns
func (c *Civ) executeShowOnlyCommand(vargs []string) {
	var hideCols []string

	for _, a := range c.table.header.cols {
		show := false
		for _, s := range vargs {
			if a.data == s {
				show = true
			}
		}

		if !show {
			hideCols = append(hideCols, a.data)
		}
	}

	c.executeHideCommand(hideCols)
}

func (c *Civ) executeResetCommand(vargs []string) {
	// Make all columsn visible
	c.table.ResetDisabledCol()

	// Make all rows visible
	c.table.ResetVisibility()
}

func (c *Civ) executeShowCommand(vargs []string) {
	for _, a := range vargs {
		if i, ok := c.table.FindColName(a); ok {
			c.table.RemoveDisabledCol(i)
		}
	}
}

func (c *Civ) executeSearchCommand() {
	searchWord := string(c.ql.query)

	for _r, row := range c.table.contents {
		matched := false
		for _c, cell := range row.cols {
			if idx := strings.Index(cell.data, searchWord); idx != -1 {
				c.table.SetMatched(_r, _c, idx, idx+len(searchWord))
				matched = true
			} else {
				c.table.SetMatched(_r, _c, -1, -1)
			}
		}
		row.hasMatched = matched
	}
}

func (c *Civ) executeFilterCommand() {
	searchWord := string(c.ql.query)

	for _r, row := range c.table.contents {
		found := false

		for _, cell := range row.cols {
			if idx := strings.Index(cell.data, searchWord); idx != -1 {
				found = true
				break
			}
		}

		c.table.SetRowVisibility(_r, found)
	}
}
