package client

// SubjectInfo 구조체는 서브젝트 관련 정보들을 담고 있는 구조체이다.
type SubjectInfo struct {
	// SubjectName 는 Subject 의 이름이다.
	// SubjectID 는 Subject 의 ID 이다.
	// SubjectURL 는 해당 서브젝트의 공식 문서 url 이다.
	// EvalGuideURL 은 해당 서브젝트 평가표의 url 이다.
	SubjectName  string
	SubjectID    int
	SubjectURL   string
	EvalGuideURL string
}

func InitSubject(sInfos map[string]SubjectInfo) {
	sInfos["Day00"] = SubjectInfo{"Day00", 0, "대충 Day00 URL", "대충 Day00 가이드 URL"}
	sInfos["Day01"] = SubjectInfo{"Day01", 1, "대충 Day01 URL", "대충 Day01 가이드 URL"}
	sInfos["Day02"] = SubjectInfo{"Day02", 2, "대충 Day02 URL", "대충 Day02 가이드 URL"}
	sInfos["Day03"] = SubjectInfo{"Day03", 3, "대충 Day03 URL", "대충 Day03 가이드 URL"}
	sInfos["Day04"] = SubjectInfo{"Day04", 4, "대충 Day04 URL", "대충 Day04 가이드 URL"}
	sInfos["Day05"] = SubjectInfo{"Day05", 5, "대충 Day05 URL", "대충 Day05 가이드 URL"}
	sInfos["Rush00"] = SubjectInfo{"Rush00", 100, "대충 Rush00 URL", "대충 Rush00 가이드 URL"}
}
