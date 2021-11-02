package match_client

type MatchClient struct {

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
func (mc *MatchClient) MyGrade(uid string) (msg string) {
	return ""
}

// MatchState 함수는 uid 를 인자로 받아 해당 유저의 매칭 상태와 현재 대기중인 평가자/피평가자 수를 리턴하는 함수이다.
func (mc *MatchClient) MatchState(uid string) (msg string) {
	return ""
}
