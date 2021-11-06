package main

// Pisicne Golang Interact 는 42서울 해커톤 Go? Ahead! 팀의 평가 매칭 시스템을 서포트하기 위한
// Discord Bot 서버를 구동하는 프로젝트입니다.

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"piscine-golang-interact/client"
	"piscine-golang-interact/mongodb"
)

const (
	prefix = "="
)

var (
	mDB *mongo.Database

	c            *client.Client
	registerMIDs map[string]string

	submitMIDs map[string]string
	submitURLs map[string]string
	submitSIDs map[string]string

	signupMIDs map[string]string
	modifyMIDs map[string]string
	IntraIDs   map[string]string
	mode       = false
)

func init() {
	registerMIDs = make(map[string]string)

	submitMIDs = make(map[string]string)
	submitURLs = make(map[string]string)
	submitSIDs = make(map[string]string)

	signupMIDs = make(map[string]string)
	modifyMIDs = make(map[string]string)
	IntraIDs = make(map[string]string)
	mClient, _ := mongodb.MongoConn()
	if mClient == nil {
		log.Errorf("mongoDB connection Failed")
	}
	c = client.NewClient(mClient.Database("k8sdb"))
}

func main() {
	//if err := record.Connection(); err != nil {
	//	log.Println("error creating DB connection", err)
	//	return
	//}

	dg, err := discordgo.New("Bot " + (viper.Get("BOT_TOKEN")).(string))
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactionAdd)
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	log.Println("uid:", r.UserID, "add reaction")
	switch r.MessageID {
	case registerMIDs[r.UserID]:
		switch r.Emoji.Name {
		case "⭕":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			registerEvalResponse(s, r)
		case "❌":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			s.ChannelMessageSend(r.ChannelID, "평가자 등록을 취소하셨습니다.")
		}
		return
	case submitMIDs[r.UserID]:
		switch r.Emoji.Name {
		case "⭕":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			submissionResponse(s, r)
		case "❌":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			s.ChannelMessageSend(r.ChannelID, "제출을 취소하셨습니다.")
		}
		return
	case signupMIDs[r.UserID]:
		switch r.Emoji.Name {
		case "⭕":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			msg := c.SignUp(r.UserID, IntraIDs[r.UserID], context.Background())
			s.ChannelMessageSend(r.ChannelID, msg)
		case "❌":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			s.ChannelMessageSend(r.ChannelID, "등록을 취소하셨습니다.")
		}
		return
	case modifyMIDs[r.UserID]:
		switch r.Emoji.Name {
		case "⭕":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			msg := c.ModifyId(r.UserID, IntraIDs[r.UserID])
			s.ChannelMessageSend(r.ChannelID, msg)
		case "❌":
			s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			s.ChannelMessageSend(r.ChannelID, "수정을 취소하셨습니다.")
		}
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	log.Println(m.Content, "is sent by uid: ", m.Author.ID)
	if m.Content == prefix+"명령어" {
		sendCommandDetail(s, m)
		return
	}
	if m.Content == prefix+"제출" {
		if c.FindIntraByUID(m.Author.ID) == "가입되지 않은 사용자" {
			s.ChannelMessageSend(m.ChannelID, "피신에 참여중인 참가자가 아닙니다.")
			return
		}
		submitEmbed := embed.NewEmbed()
		submitEmbed.SetTitle("제출 명령어는 다음과 같이 입력해야 합니다.")
		submitEmbed.AddField(
			"<명령어 예시>",
			prefix+"제출 <github repo url> <subject name>\n"+
				prefix+"제출 https://github.com/example123/ExampleRepo Day01")
		s.ChannelMessageSendEmbed(m.ChannelID, submitEmbed.MessageEmbed)
		return
	}
	if strings.HasPrefix(m.Content, prefix+"제출 ") { // !제출 <github repo url> <subject name>
		if c.FindIntraByUID(m.Author.ID) == "가입되지 않은 사용자" {
			s.ChannelMessageSend(m.ChannelID, "피신에 참여중인 참가자가 아닙니다.")
			return
		}
		if c.IsUserInQ(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "이미 큐에 등록된 유저입니다.")
			return
		}
		submissionTask(s, m)
		return
	}
	if m.Content == prefix+"제출취소" {
		submissionCancelTask(s, m)
		return
	}
	if m.Content == prefix+"평가등록" {
		if c.FindIntraByUID(m.Author.ID) == "가입되지 않은 사용자" {
			s.ChannelMessageSend(m.ChannelID, "피신에 참여중인 참가자가 아닙니다.")
			return
		}
		if c.IsUserInQ(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "이미 큐에 등록된 유저입니다.")
			return
		}
		registerEvalTask(s, m)
		return
	}
	if m.Content == prefix+"평가취소" {
		RegisterCancelTask(s, m)
		return
	}
	if m.Content == prefix+"매칭상태" {
		state := c.MatchState()
		sendEmbedPretty(s, m.ChannelID, state)
		return
	}
	if m.Content == prefix+"내점수" {
		grade := c.MyGrade(m.Author.ID)
		sendEmbedPretty(s, m.ChannelID, grade)
		return
	}
	if m.Content == prefix+"GOPISCINEREGISTERMODE" && (m.Author.ID == "318743234601811969" ||
		m.Author.ID == "905699384581312542" || m.Author.ID == "382356905990815744" ||
		m.Author.ID == "383847223504666626") {
		mode = !mode
		if mode {
			s.ChannelMessageSend(m.ChannelID, "사용자 등록모드 시작")
		} else {
			s.ChannelMessageSend(m.ChannelID, "사용자 등록모드 종료")
		}
	}
	if mode && strings.HasPrefix(m.Content, prefix+"인트라등록") {
		command := strings.Split(m.Content, " ")
		if len(command) != 2 {
			s.ChannelMessageSend(m.ChannelID, "사용방법: "+prefix+"인트라등록 <intraID>")
			return
		}
		dmChan, _ := s.UserChannelCreate(m.Author.ID)
		regMsg, _ := s.ChannelMessageSend(dmChan.ID, "**주의** 등록된 정보는 등록 기간이 지나면 바꿀 수 없음!\n"+
			"등록된 인트라 ID 를 바꾸고 싶다면 "+prefix+"인트라수정 명령어를 사용하세요\n"+
			"당신의 인트라 ID 가 "+command[1]+" 이(가) 맞습니까?")
		signupMIDs[m.Author.ID] = regMsg.ID
		IntraIDs[m.Author.ID] = command[1]
		s.MessageReactionAdd(dmChan.ID, regMsg.ID, "⭕")
		s.MessageReactionAdd(dmChan.ID, regMsg.ID, "❌")
	}
	if mode && strings.HasPrefix(m.Content, prefix+"인트라수정") {
		command := strings.Split(m.Content, " ")
		if len(command) != 2 {
			s.ChannelMessageSend(m.ChannelID, "사용방법: "+prefix+"인트라수정 <intraID>")
			return
		}
		dmChan, _ := s.UserChannelCreate(m.Author.ID)
		regMsg, _ := s.ChannelMessageSend(dmChan.ID, "**주의** 등록된 정보는 등록 기간이 지나면 바꿀 수 없음!\n"+
			"당신의 인트라 ID 가 "+command[1]+" 이(가) 맞습니까?")
		modifyMIDs[m.Author.ID] = regMsg.ID
		IntraIDs[m.Author.ID] = command[1]
		s.MessageReactionAdd(dmChan.ID, regMsg.ID, "⭕")
		s.MessageReactionAdd(dmChan.ID, regMsg.ID, "❌")
	}
}

func sendEmbedPretty(s *discordgo.Session, cid string, info client.EmbedInfo) {
	answer := embed.NewEmbed()
	answer.SetTitle(info.Title())
	fields := info.EmbedRows()
	for _, row := range fields {
		name := row.Name()
		lines := row.Lines()
		value := ""
		for i := 0; i < len(lines)-1; i++ {
			value += lines[i] + "\n"
		}
		value += lines[len(lines)-1]
		answer.AddField(name, value)
	}
	s.ChannelMessageSendEmbed(cid, answer.MessageEmbed)
}

// sendCommandDetail 함수는 명령어 정보를 전부 전송하는 함수이다.
func sendCommandDetail(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandDetailEmbed := embed.NewEmbed()
	commandDetailEmbed.SetTitle("명령어 목록")
	commandDetailEmbed.AddField(
		"도움말 명령어",
		prefix+"도움말"+"\n"+
			prefix+"명령어",
	)
	commandDetailEmbed.AddField(
		"제출 명령어",
		prefix+"제출 <GitRepoURL> <SubjectID>"+"\n"+
			prefix+"제출취소",
	)
	commandDetailEmbed.AddField(
		"평가자 등록 명령어",
		prefix+"평가등록"+"\n"+
			prefix+"평가취소",
	)
	commandDetailEmbed.AddField(
		"정보 확인 명령어",
		prefix+"내점수"+"\n"+
			prefix+"매칭상태",
	)
	s.ChannelMessageSendEmbed(m.ChannelID, commandDetailEmbed.MessageEmbed)
}
