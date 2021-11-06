package client

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"piscine-golang-interact/schema"
	"runtime"
	"strconv"

	"sync"
)

// SubjectInfoMap 은 sid 를 기반으로 해당 서브젝트의 정보 구조체를 반환하는 맵이다.
var SubjectInfoMap map[string]SubjectInfo

// IntervieweeList 는 피평가자의 uid 를 이용하는 Queue 이다.
var IntervieweeList *[]string

// InterviewerList 는 평가자의 uid 를 이용하는 Queue 이다.
var InterviewerList *[]string

// QueueMutex 는 대기열의 동기화를 위한 Mutex 이다.
var QueueMutex sync.Mutex

// MatchInfo 구조체는 평가 매칭이 성공했을 때 전달하는 평가 정보 구조체이다.
type MatchInfo struct {
	// Code 는 매칭 성공시 true, 매칭 취소시 false 이다.
	// InterviewerID 는 평가자의 uid 이다.
	// IntervieweeID 는 피평가자의 uid 이다.
	Code          bool
	InterviewerID string
	IntervieweeID string
	Subject       SubjectInfo
}

// Client 구조체는 Piscine Golang 서브젝트의 평가 매칭을 관리하는 오브젝트이다.
type Client struct {
	// MatchMap 은 uid 를 key 로 하여,
	// 해당 유저가 매칭 성공시에 상대의 uid 를 받기 위한 채널을 value 로 한다.
	MatchMap            map[string]chan MatchInfo
	SubmittedSubjectMap map[string]SubjectInfo
	MDB                 *mongo.Database
}

func queueLock() {
	log.Println("try queue lock acquire")
	log.Println(runtime.Caller(2))
	QueueMutex.Lock()
	log.Warn("queue lock is acquire")
}

func queueUnlock() {
	log.Println("try queue lock release")
	log.Println(runtime.Caller(2))
	QueueMutex.Unlock()
	log.Warn("queue lock is released")
}

func init() {
	SubjectInfoMap = make(map[string]SubjectInfo)
	InitSubject(SubjectInfoMap)
	tempInterviewee := make([]string, 0, 100)
	IntervieweeList = &tempInterviewee
	tempInterviewer := make([]string, 0, 100)
	InterviewerList = &tempInterviewer
	QueueMutex = sync.Mutex{}
}

func removeClient(list *[]string, i int) {
	if len(*list) == i {
		*list = (*list)[:i]
		return
	}
	*list = append((*list)[:i], (*list)[i+1:]...)
}

// NewClient 함수는 Client 구조체의 생성자이다.
func NewClient(mDB *mongo.Database) (ret *Client) {
	ret = &Client{}
	ret.MDB = mDB
	ret.MatchMap = make(map[string]chan MatchInfo)
	ret.SubmittedSubjectMap = make(map[string]SubjectInfo)
	return ret
}

// SignUp 함수는 uid(userID) intraID를 받아 DB 에 추가하는 함수이다.
// DB 에 추가하기 전에 기존에 가입된 intraID 라면 가입이 되지 않는다.
func (c *Client) SignUp(uid, name string, ctx context.Context) (msg string) {
	searchPerson := schema.Person{}
	err := c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{"password", uid}},
		).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Password != "" {
		return "이미 등록된 디스코드 계정입니다."
	}
	err = c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{"name", name}},
		).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Name != "" {
		return "이미 등록된 IntraID 입니다."
	}
	curUser := schema.Person{
		Name: name,
		Password: uid,
		Course: 0,
		Point: 5,
	}
	_, err = c.MDB.Collection("people").InsertOne(ctx, curUser)
	if err != nil {
		log.Error(err)
	}
	return "인트라 등록이 완료되었습니다."
}

// IsUserInQ 함수는 uid 를 바탕으로 사용자가 큐에 등록했는지를 확인하는 함수이다.
func (c *Client) IsUserInQ(uid string) bool {
	queueLock()
	defer queueUnlock()
	for _, item := range *InterviewerList {
		if item == uid {
			return true
		}
	}
	for _, item := range *IntervieweeList {
		if item == uid {
			return true
		}
	}
	return false
}

// ModifyId 함수는 uid 를 기반으로 intraID 를 변경하는 함수이다.
func (c *Client) ModifyId(uid, name string) (msg string) {
	//tx, tErr := record.DB.Begin()
	//if tErr != nil {
	//	return "인트라 ID 수정오류: 트랜잭션 초기화"
	//}
	//defer tx.Rollback()
	//// ret nil, err nil    : 사용자 없음
	//// ret nil, err        : 에러
	//// ret, err nil        : 사용자 있음
	//// ret, err            : 에러
	//if ret, qErr := tx.Query(
	//	`SELECT id FROM people WHERE password = ? ;`,
	//	uid); qErr != nil {
	//	if ret != nil {
	//		return "인트라 ID 수정오류: 매칭되는 사용자가 없음"
	//	}
	//	if _, eErr := tx.Exec(
	//		`UPDATE people SET name = ? WHERE password = ? ;`,
	//		name, uid); eErr != nil {
	//		return "인트라 ID 수정오류: 수정 실패" + name + uid
	//	}
	//}
	//tErr = tx.Commit()
	//if tErr != nil {
	//	return "인트라 ID 수정오류: 트랜잭션 적용"
	//} else {
	//	return "인트라 ID 수정 완료"
	//}
	return
}

// Submit 함수는 sid(subject id) uid(userID) url(github repo link)와
// 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 서브젝트 제출을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
// Eval Queue 에 사용자가 있는지 Mutex 를 걸고 확인한 후에 있다면 매칭을 진행해야한다. ** MUTEX 활용 필수!!
func (c *Client) Submit(sName, uid, url string, matchedUserId chan MatchInfo) (msg string) {
	log.Println("Submit called")
	defer log.Println("Submit ended")
	// convertID := SubjectStrMap[sid]
	if c.MatchMap[uid] != nil {
		return "이미 큐에 등록된 사용자입니다."
	}
	queueLock()
	defer queueUnlock()
	if len(*InterviewerList) == 0 {
		*IntervieweeList = append(*IntervieweeList, uid)
		c.MatchMap[uid] = matchedUserId
		c.SubmittedSubjectMap[uid] = SubjectInfoMap[sName]
	} else {
		matchedInterviewerID := (*InterviewerList)[0]
		myMatchInfo := MatchInfo{
			Code:          true,
			IntervieweeID: uid,
			InterviewerID: matchedInterviewerID,
			Subject:       SubjectInfoMap[sName],
		}
		go func() {
			c.MatchMap[matchedInterviewerID] <- myMatchInfo
		}()
		go func() {
			matchedUserId <- myMatchInfo
		}()
		removeClient(InterviewerList, 0)
	}
	return "제출완료"
}

// SubmitCancel 함수는 uid 를 인자로 받아 해당 유저의 제출을 취소하는 함수이다.
// 제출 취소의 성공/실패 여부를 msg 로 리턴한다.
func (c *Client) SubmitCancel(uid string) (msg string) {
	queueLock()
	defer queueUnlock()
	for i, v := range *InterviewerList {
		if v == uid {
			c.MatchMap[uid] = nil
			c.SubmittedSubjectMap[uid] = SubjectInfo{}
			removeClient(IntervieweeList, i)
			return "취소완료"
		}
	}
	return "취소오류"
}

// Register 함수는 uid 와 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 평가 등록을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
// Submit Queue 에 사용자가 있는지 Mutex 를 걸고 확인한 후에 있다면 매칭을 진행해야한다. ** MUTEX 활용 필수!!
func (c *Client) Register(uid string, matchedUid chan MatchInfo) (msg string) {
	if c.MatchMap[uid] != nil {
		return "이미 큐에 등록된 사용자입니다."
	}
	log.Println("Register called")
	defer log.Println("Register ended")
	queueLock()
	defer queueUnlock()
	if len(*IntervieweeList) == 0 {
		*InterviewerList = append(*InterviewerList, uid)
		c.MatchMap[uid] = matchedUid
	} else {
		matchedIntervieweeID := (*IntervieweeList)[0]
		myMatchInfo := MatchInfo{
			Code:          true,
			IntervieweeID: matchedIntervieweeID,
			InterviewerID: uid,
			Subject:       c.SubmittedSubjectMap[matchedIntervieweeID],
		}
		go func() {
			c.MatchMap[matchedIntervieweeID] <- myMatchInfo
		}()
		go func() {
			matchedUid <- myMatchInfo
		}()
		removeClient(IntervieweeList, 0)
	}
	return "평가등록완료"
}

// RegisterCancel 함수는 uid 를 인자로 받아 해당 유저의 평가 등록을 취소하는 함수이다.
// 평가 등록 취소의 성공/실패 여부를 msg 로 리턴한다.
func (c *Client) RegisterCancel(uid string) (msg string) {
	queueLock()
	defer queueUnlock()
	for _, val := range *InterviewerList {
		log.Println("before remove:", val)
	}
	for i, v := range *InterviewerList {
		if v == uid {
			c.MatchMap[uid] = nil
			removeClient(InterviewerList, i)
			return "평가취소완료"
		}
	}
	for _, val := range *InterviewerList {
		log.Println("after remove:", val)
	}
	return "평가취소오류"
}

// MyGrade 함수는 uid 를 인자로 받아 해당 유저의 점수 정보를 리턴하는 함수이다.
func (c *Client) MyGrade(uid string) (grades EmbedInfo) {
	grades.title = "서브젝트 채점 현황"
	ctx := context.Background()
	curPerson := schema.Person{}
	err := c.MDB.Collection("people").FindOne(ctx, bson.D{{"password", uid}}).Decode(&curPerson)
	if err != nil {
		log.Error(err)
	}
	curScores := curPerson.Score
	if len(curScores) == 0 {
		grades.embedRows = []EmbedRow{{
			"평가받은 과제가 없어요...",
			[]string{"평가받은 과제가 없습니다.", "Go? Ahead!"},
		}}
		return grades
	}
	for _, item := range curScores {
		itemRow := EmbedRow{
			name: SubjectInfoMap[item.Course].SubjectName,
			lines: []string{
				"평가자:\t" + c.FindIntraByUID(item.InterviewerID),
				"점수:\t" + strconv.Itoa(item.Score),
			},
		}
		if item.Pass {
			itemRow.lines = append(itemRow.lines, "[ OK ]")
		} else {
			itemRow.lines = append(itemRow.lines, "[ KO ]")
		}
		grades.embedRows = append(grades.embedRows, itemRow)
	}
	return
}

func matchEmbedRow(s string, p *map[string]string, l *[]string) EmbedRow {
	tempEmbedRow := EmbedRow{name: s}
	tempLines := make([]string, 0, 20)
	queueLock()
	for _, v := range *l {
		if i, ok := (*p)[v]; ok {
			tempLines = append(tempLines, i)
		}
	}
	queueUnlock()
	if len(tempLines) == 0 {
		tempLines = append(tempLines, "대기열 없음")
	}
	tempEmbedRow.lines = tempLines
	return tempEmbedRow
}

// MatchState 함수는 uid 를 인자로 받아 현재 큐 정보를 리턴하는 함수이다.
func (c *Client) MatchState() (grades EmbedInfo) {
	grades.title = "평가 및 피평가 매칭 현황"
	people := make(map[string]string)
	grades.embedRows = append(grades.embedRows, matchEmbedRow("평가자", &people, InterviewerList))
	grades.embedRows = append(grades.embedRows, matchEmbedRow("피평가자", &people, IntervieweeList))
	return
}

// FindIntraByUID 함수는 uid 를 인자로 받아 intraID 를 반환하는 함수이다.
func (c *Client) FindIntraByUID(uid string) (intraID string) {
	ctx := context.Background()
	var curPerson schema.Person
	err := c.MDB.Collection("people").FindOne(ctx, bson.D{
		{"password", uid},
	}).Decode(&curPerson)
	if err != nil {
		return "가입하지 않은 사용자"
	}
	return curPerson.Name
}
