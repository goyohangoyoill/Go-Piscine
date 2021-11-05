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

const (
	sDay00  = "https://drive.google.com/file/d/1gaInmQbHqp75T-9XasAK1qzBPvTJa0EN/view?usp=sharing"
	sDay01  = "https://drive.google.com/file/d/19qO4jprHdnJS4AP9HF7WN11csaeJR_0x/view?usp=sharing"
	sDay02  = "https://drive.google.com/file/d/1c2TLc_g6lGBwYEQ7qVjfFKgmeXAz1LCm/view?usp=sharing"
	sDay03  = "https://drive.google.com/file/d/1dZsEY158-A9MOjRxxCr1Py4a5YB_PnmT/view?usp=sharing"
	sDay04  = "https://drive.google.com/file/d/1rnPYdYdT9cd5H_ohYKsfpMvXn3PRECn0/view?usp=sharing"
	sDay05  = "https://drive.google.com/file/d/1ASESn4hHM-mLy2Q0MIgz03kbCMGUKqDI/view?usp=sharing"
	sRush00 = "https://drive.google.com/file/d/1rJ6eaxiJJj9OZET373eTqRKnh7U8-_p7/view?usp=sharing"
)

const (
	vDay00  = "https://drive.google.com/file/d/14Emsu_11_1YsE_kYX2iComJFXYFqSxEu/view?usp=sharing"
	vDay01  = "https://drive.google.com/file/d/1Ns9QvPkTgrrNq2Lo9xk6CV6rcC_-aaOT/view?usp=sharing"
	vDay02  = "https://drive.google.com/file/d/1IDk6_cmtfJwwZs6YArfHJi35yGh2KJQ7/view?usp=sharing"
	vDay03  = "https://drive.google.com/file/d/1wsibyfrIB5-6e_K7GeZn96iA7DphUbRk/view?usp=sharing"
	vDay04  = "https://drive.google.com/file/d/1SHaCDxoxihAs5GIzDEwvj5_cqFvOJGPU/view?usp=sharing"
	vDay05  = "https://drive.google.com/file/d/1qYvWJCaTv5yscpsGEJFgzuAIjSLDG16N/view?usp=sharing"
	vRush00 = "https://drive.google.com/file/d/1jn_yENWtvf4XbUi4uxozNWyWkEDj8Kez/view?usp=sharing"
)

func InitSubject(sInfos map[string]SubjectInfo) {
	sInfos["Day00"] = SubjectInfo{"Day00", 0, sDay00, vDay00}
	sInfos["Day01"] = SubjectInfo{"Day01", 1, sDay01, vDay01}
	sInfos["Day02"] = SubjectInfo{"Day02", 2, sDay02, vDay02}
	sInfos["Day03"] = SubjectInfo{"Day03", 3, sDay03, vDay03}
	sInfos["Day04"] = SubjectInfo{"Day04", 4, sDay04, vDay04}
	sInfos["Day05"] = SubjectInfo{"Day05", 5, sDay05, vDay05}
	sInfos["Rush00"] = SubjectInfo{"Rush00", 100, sRush00, vRush00}
}
