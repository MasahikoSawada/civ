package civ

type QueryLine struct {
	query     []rune
	curCursor int // starts from 0
	mode      int
}

// Modes
const (
	MODE_VIEW    = 1
	MODE_COMMAND = 2
	MODE_SEARCH  = 3
	MODE_FILTER  = 4
)

func NewQueryLine() *QueryLine {
	q := &QueryLine{
		query:     nil,
		curCursor: 0,
		mode:      MODE_VIEW,
	}
	return q
}

func (q *QueryLine) ClearQuery() {
	q.query = nil
	q.curCursor = 0
}

func (q *QueryLine) ClearAll() {
	q.query = nil
	q.curCursor = 0
	q.mode = MODE_VIEW
}

func (q *QueryLine) QueryLen() int {
	return len(q.query)
}

func (q *QueryLine) InputChar(r rune) {
	if len(q.query) <= q.curCursor {
		q.query = append(q.query, r)
	} else {
		_q := make([]rune, q.curCursor)
		copy(_q, q.query[:q.curCursor])
		q.query = append(append(_q, r), q.query[q.curCursor:]...)
	}

	q.curCursor++
}

func (q *QueryLine) BackwardChar() {
	if q.curCursor == 0 {
		return
	}

	_q := make([]rune, q.curCursor-1)
	copy(_q, q.query[:q.curCursor])
	q.query = append(_q, q.query[q.curCursor:]...)

	q.curCursor--
}

func (q *QueryLine) DeleteChar() {
	if q.curCursor >= len(q.query) {
		return
	}

	_q := make([]rune, q.curCursor)
	copy(_q, q.query[:q.curCursor])
	q.query = append(_q, q.query[(q.curCursor+1):]...)
}

func (q *QueryLine) TruncateChars() {
	if q.curCursor >= len(q.query) {
		return
	}

	_q := make([]rune, q.curCursor)
	copy(_q, q.query[:q.curCursor])
	q.query = _q
}

func (q *QueryLine) MoveForward() {
	if len(q.query) > q.curCursor {
		q.curCursor++
	}
}

func (q *QueryLine) MoveBackward() {
	if len(q.query) >= q.curCursor && q.curCursor > 0 {
		q.curCursor--
	}
}

func (q *QueryLine) MoveToTop() {
	q.curCursor = 0
}

func (q *QueryLine) MoveToEnd() {
	q.curCursor = len(q.query)
}
