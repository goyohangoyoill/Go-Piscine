package client

// EmbedRow 구조체는 name 과 lines 를 가진다.
// name, lines 를 반환하는 게터 함수들 역시 가진다.
type EmbedRow struct {
	name  string
	lines []string
}

// Name 함수는 name 을 반환하는 게터이다.
func (si EmbedRow) Name() string {
	return si.name
}

// Lines 함수는 lines 을 반환하는 게터이다.
func (si EmbedRow) Lines() []string {
	return si.lines
}

// EmbedInfo 구조체는 title 과 subjectGrades 를 가진다.
// title, subjectGrades 를 반환하는 게터 함수들 역시 가진다.
type EmbedInfo struct {
	title     string
	embedRows []EmbedRow
}

// Title 함수는 title 을 반환하는 게터이다.
func (gi EmbedInfo) Title() string {
	return gi.title
}

// EmbedRows 함수는 embedRows 을 반환하는 게터이다.
func (gi EmbedInfo) EmbedRows() []EmbedRow {
	return gi.embedRows
}
