package client

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	"piscine-golang-interact/schema"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SubjectInfoMap 은 sid 를 기반으로 해당 서브젝트의 정보 구조체를 반환하는 맵이다.
var SubjectInfoMap map[string]SubjectInfo

// Client 구조체는 Piscine Golang 서브젝트의 평가 매칭을 관리하는 오브젝트이다.
type Client struct {
	SubmittedSubjectMap map[string]SubjectInfo
	MDB                 *mongo.Database
}

func init() {
	SubjectInfoMap = make(map[string]SubjectInfo)
	InitSubject(SubjectInfoMap)
}

// NewClient 함수는 Client 구조체의 생성자이다.
func NewClient(mDB *mongo.Database) (ret *Client) {
	ret = &Client{}
	ret.MDB = mDB
	ret.SubmittedSubjectMap = make(map[string]SubjectInfo)
	return ret
}

// SignUp 함수는 uid(userID) intraID를 받아 DB 에 추가하는 함수이다.
// DB 에 추가하기 전에 기존에 가입된 intraID 라면 가입이 되지 않는다.
func (c *Client) SignUp(uid, name string, ctx context.Context) (msg string) {
	searchPerson := schema.Person{}
	err := c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{Key: "password", Value: uid}},
	).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Password != "" {
		return "이미 등록된 디스코드 계정입니다."
	}
	err = c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{Key: "name", Value: name}},
	).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Name != "" {
		return "이미 등록된 IntraID 입니다."
	}
	curUser := schema.Person{
		Name:     name,
		Password: uid,
		Course:   0,
		Point:    5,
	}
	_, err = c.MDB.Collection("people").InsertOne(ctx, curUser)
	if err != nil {
		log.Error(err)
	}
	return "인트라 등록이 완료되었습니다."
}

// ModifyId 함수는 uid 를 기반으로 intraID 를 변경하는 함수이다.
func (c *Client) ModifyId(uid, name string, ctx context.Context) (msg string) {
	searchPerson := schema.Person{}
	err := c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{Key: "password", Value: uid}},
	).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Password == "" {
		return "해당 계정으로 가입한 사용자가 없습니다."
	}
	err = c.MDB.Collection("people").FindOne(
		ctx,
		bson.D{{Key: "name", Value: name}},
	).Decode(&searchPerson)
	if err != nil {
		log.Error(err)
	}
	if searchPerson.Name != "" {
		return "같은 이름의 인트라ID가 존재합니다."
	}
	filter := bson.M{"password": uid}
	update := bson.M{
		"$set": bson.M{
			"name": name,
		},
	}
	_, err = c.MDB.Collection("people").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error(err)
	}
	return "인트라 ID 수정이 완료되었습니다."
}

// Submit 함수는 sid(subject id) uid(userID) url(github repo link)와
// 매칭된 상대방의 UID 를 공유할 matchedUserId channel 을 인자로 받아
// 서브젝트 제출을 수행하고 작업이 성공적으로 이루어졌는지 여부를 알리는 msg 를 반환하는 함수이다.
// Eval Queue 에 사용자가 있는지 Mutex 를 걸고 확인한 후에 있다면 매칭을 진행해야한다. ** MUTEX 활용 필수!!
func (c *Client) Submit(sName, uid, url string) (msg []byte) {
	log.Println("Submit called")
	defer log.Println("Submit ended")
	baseURL := os.Getenv("GRADE_SERVICE")
	log.Println("http Get try")
	resp, err := http.Get("http://" + baseURL + ":4242/grade/" + SubjectInfoMap[sName].SubjectName + "?URL=" + url)
	log.Println("http Get ended")
	if err != nil {
		log.Println(c.FindIntraByUID(uid) + " got error")
		log.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body
}

// MyGrade 함수는 uid 를 인자로 받아 해당 유저의 점수 정보를 리턴하는 함수이다.
func (c *Client) MyGrade(uid string) (grades EmbedInfo) {
	grades.title = "서브젝트 채점 현황"
	ctx := context.Background()
	curPerson := schema.Person{}
	err := c.MDB.Collection("people").FindOne(ctx, bson.D{{Key: "password", Value: uid}}).Decode(&curPerson)
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
			name:  SubjectInfoMap[item.Course].SubjectName,
			lines: []string{},
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

// FindIntraByUID 함수는 uid 를 인자로 받아 intraID 를 반환하는 함수이다.
func (c *Client) FindIntraByUID(uid string) (intraID string) {
	ctx := context.Background()
	var curPerson schema.Person
	err := c.MDB.Collection("people").FindOne(ctx, bson.D{
		{Key: "password", Value: uid},
	}).Decode(&curPerson)
	if err != nil {
		return "가입하지 않은 사용자"
	}
	return curPerson.Name
}
