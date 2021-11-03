package client

// MatchInfo 구조체는 평가 매칭이 성공했을 때 전달하는 평가 정보 구조체이다.
type MatchInfo struct {
	// Code 는 매칭 성공시 true, 매칭 취소시 false 이다.
	// InterviewerID 는 평가자의 uid 이다.
	// IntervieweeID 는 피평가자의 uid 이다.
	// SubjectName 는 Subject 의 이름이다.
	// SubjectURL 는 해당 서브젝트의 공식 문서 url 이다.
	// EvalGuideURL 은 해당 서브젝트 평가표의 url 이다.
	Code          bool
	InterviewerID string
	IntervieweeID string
}

// Client 구조체는 각 go-piscine 서브젝트의 평가 매칭을 관리하는 오브젝트이다.
type Client struct {
	// MatchMap 은 uid 를 key 로 하여,
	// 해당 유저가 매칭 성공시에 상대의 uid 를 받기 위한 채널을 value 로 한다.
	MatchMap map[string]chan MatchInfo
}

// NewClient 함수는 Client 구조체의 생성자이다.
func NewClient() (ret *Client) {
	ret = &Client{}
	ret.MatchMap = make(map[string]chan MatchInfo)
	return ret
}

// Submit 함수는 sid(subject id) uid(userID) url(github repo link)와
// 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 서브젝트 제출을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
// Eval Queue 에 사용자가 있는지 Mutex 를 걸고 확인한 후에 있다면 매칭을 진행해야한다. ** MUTEX 활용 필수!!
func (c *Client) Submit(sid, uid, url string, matchedUserId chan MatchInfo) (msg string) {
	return "제출완료"
}

// SubmitCancel 함수는 uid 를 인자로 받아 해당 유저의 제출을 취소하는 함수이다.
// 제출 취소의 성공/실패 여부를 msg 로 리턴한다.
func (c *Client) SubmitCancel(uid string) (msg string) {
	return "취소완료"
}

// Register 함수는 uid 와 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 평가 등록을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
// Submit Queue 에 사용자가 있는지 Mutex 를 걸고 확인한 후에 있다면 매칭을 진행해야한다. ** MUTEX 활용 필수!!
func (c *Client) Register(uid string, matchedUid chan MatchInfo) (msg string) {
	return "평가등록완료"
}

// RegisterCancel 함수는 uid 를 인자로 받아 해당 유저의 평가 등록을 취소하는 함수이다.
// 평가 등록 취소의 성공/실패 여부를 msg 로 리턴한다.
func (c *Client) RegisterCancel(sid, uid string) (msg string) {
	return "평가취소완료"
}

// MatchState 함수는 uid 를 인자로 받아 해당 유저의 매칭 상태와 현재 대기중인 평가자/피평가자 수를 리턴하는 함수이다.
func (c *Client) MatchState() (matchState EmbedInfo) {
	return
}

// MyGrade 함수는 uid 를 인자로 받아 해당 유저의 점수 정보를 리턴하는 함수이다.
func (c *Client) MyGrade(uid string) (grades EmbedInfo) {
	return
}

// FindIntraByUID 함수는 uid 를 인자로 받아 intraID 를 반환하는 함수이다.
func (c *Client) FindIntraByUID(uid string) (intraID string) {
	return
}

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