package match_client

// MatchClient 구조체는 각 go-piscine 서브젝트의 평가 매칭을 관리하는 오브젝트이다.
type MatchClient struct {
	// MatchMap 은 uid 를 key 로 하여,
	// 해당 유저가 매칭 성공시에 상대의 uid 를 받기 위한 채널을 value 로 한다.
	MatchMap map[string]chan string
}

// NewMatchClient 함수는 MatchClient 구조체의 생성자이다.
func NewMatchClient() (ret *MatchClient) {
	ret = &MatchClient{}
	ret.MatchMap = make(map[string]chan string)
	return ret
}

// Submit 함수는 sid(subject id) uid(userID) url(github repo link)와
// 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 서브젝트 제출을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
func (mc *MatchClient) Submit(sid, uid, url string, matchedUserId chan string) (msg string) {
	return ""
}

// SubmitCancel 함수는 uid 를 인자로 받아 해당 유저의 제출을 취소하는 함수이다.
// 제출 취소의 성공/실패 여부를 msg 로 리턴한다.
func (mc *MatchClient) SubmitCancel(uid string) (msg string) {
	return ""
}

// RegisterEval 함수는 uid 와 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 평가 등록을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
func (mc *MatchClient) RegisterEval(uid string, matchedUid chan string) (msg string) {
	return ""
}

// EvalCancel 함수는 uid 를 인자로 받아 해당 유저의 평가 등록을 취소하는 함수이다.
// 평가 등록 취소의 성공/실패 여부를 msg 로 리턴한다.
func (mc *MatchClient) EvalCancel(sid, uid string) (msg string) {
	return ""
}

// MyGrade 함수는 uid 를 인자로 받아 해당 유저의 점수 정보를 리턴하는 함수이다.
func (mc *MatchClient) MyGrade(uid string) (grades EmbedInfo) {
	return
}

// MatchState 함수는 uid 를 인자로 받아 해당 유저의 매칭 상태와 현재 대기중인 평가자/피평가자 수를 리턴하는 함수이다.
func (mc *MatchClient) MatchState() (matchState EmbedInfo) {
	return
}

// EmbedRow 구조체는 name 과 lines 를 가진다.
// name, lines 를 반환하는 게터 함수들 역시 가진다.
type EmbedRow struct {
	name string
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
	title string
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
