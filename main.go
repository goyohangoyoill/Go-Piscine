package main

// Go-Piscine-EvalBot 은 42서울 해커톤 Go? Ahead! 팀의 평가 매칭 시스템을 서포트하기 위한
// Discord Bot 서버를 구동하는 프로젝트입니다.

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	mc "github.com/goyohangoyoill/Go-Piscine-EvalBot/match_client"
)

var (
	matchClient *mc.MatchClient
	// Token 은 해당 디스코드 봇의 토큰 값입니다.
	Token string
)

func init() {
	matchClient = mc.NewMatchClient()
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "!명령어" {
		sendCommandDetail(s, m)
		return
	}
	if m.Content == "!제출" {
		submitEmbed := embed.NewEmbed()
		submitEmbed.SetTitle("제출 명령어는 다음과 같이 입력해야 합니다.")
		submitEmbed.AddField(
			"<명령어 예시>",
			"!제출 <github repo url> <subject name>\n" +
				"!제출 https://github.com/example123/ExampleRepo Day01")
		s.ChannelMessageSendEmbed(m.ChannelID, submitEmbed.MessageEmbed)
		return
	}
	if strings.HasPrefix(m.Content, "!제출 ") {	// !제출 <github repo url> <subject name>
		submissionTask(s, m)
		return
	}
	if m.Content == "!제출취소" {
		submissionCancelTask(s, m)
		return
	}
	if m.Content == "!평가등록" {
		registerEvalTask(s, m)
		return
	}
	if m.Content == "!평가취소" {
		evalCancelTask(s, m)
		return
	}
	if m.Content == "!매칭상태" {
		state := matchClient.MatchState()
		sendEmbedPretty(s, m.ChannelID, state)
		return
	}
	if m.Content == "!내점수" {
		grade := matchClient.MyGrade(m.Author.ID)
		sendEmbedPretty(s, m.ChannelID, grade)
	}
}

func sendEmbedPretty(s *discordgo.Session, cid string, info mc.EmbedInfo) {
	answer := embed.NewEmbed()
	answer.SetTitle(info.Title())
	fields := info.EmbedRows()
	for _, row := range fields {
		name := row.Name()
		value := ""
		for _, line := range row.Lines() {
			value += line + "\n"
		}
		answer.AddField(name, value)
	}
	s.ChannelMessageSendEmbed(cid, answer.MessageEmbed)
}

func registerEvalTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	matchedUserID := make(chan mc.MatchInfo, 2)
	msg := matchClient.RegisterEval(m.Author.ID, matchedUserID)
	s.ChannelMessageSend(m.ChannelID, msg)
	switch evalInfo := <- matchedUserID; evalInfo.Code {
	case false:
		close(matchedUserID)
	case true:
		dmChan, _ := s.UserChannelCreate(m.Author.ID)
		matchSuccessEmbed := embed.NewEmbed()
		matchSuccessEmbed.SetTitle("평가 매칭 성공!")
		matchSuccessEmbed.AddField(
			"피평가자 intra ID:",
			 matchClient.FindIntraByUID(evalInfo.IntervieweeID),
			)
		matchSuccessEmbed.AddField(
			"평가할 서브젝트:",
			evalInfo.SubjectName + "\n" +
			evalInfo.SubjectURL,
			)
		matchSuccessEmbed.AddField(
			"평가표 링크:",
			evalInfo.EvalGuideURL,
			)
		s.ChannelMessageSendEmbed(dmChan.ID, matchSuccessEmbed.MessageEmbed)
	}
}

func evalCancelTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	userChannel := matchClient.MatchMap[m.Author.ID]
	if userChannel == nil {
		s.ChannelMessageSend(m.ChannelID, "현재 평가 등록을 하지 않은 사용자입니다.")
		return
	}
	userChannel <- mc.MatchInfo{Code: false}
	s.ChannelMessageSend(m.ChannelID, "정상적으로 평가 등록이 취소되었습니다.")
}

func submissionCancelTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	userChannel := matchClient.MatchMap[m.Author.ID]
	if userChannel == nil {
		s.ChannelMessageSend(m.ChannelID, "현재 제출을 하지 않은 사용자입니다.")
		return
	}
	userChannel <- mc.MatchInfo{Code: false}
	s.ChannelMessageSend(m.ChannelID, "정상적으로 서브젝트 제출이 취소되었습니다.")
}

func submissionTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	matchedUserID := make(chan mc.MatchInfo, 2)
	if len(command) != 3 {
		submitEmbed := embed.NewEmbed()
		submitEmbed.SetTitle("제출 명령어는 다음과 같이 입력해야 합니다.")
		submitEmbed.AddField(
			"<명령어 예시>",
			"!제출 <github repo url> <subject name>\n" +
			"!제출 https://github.com/example123/ExampleRepo Day01")
		s.ChannelMessageSendEmbed(m.ChannelID, submitEmbed.MessageEmbed)
		return
	}
	msg := matchClient.Submit(command[2], m.Author.ID, command[1], matchedUserID)
	s.ChannelMessageSend(m.ChannelID, msg)
	switch evalInfo := <- matchedUserID; evalInfo.Code {
	case false:
		close(matchedUserID)
	case true:
		dmChan, _ := s.UserChannelCreate(m.Author.ID)
		matchSuccessEmbed := embed.NewEmbed()
		matchSuccessEmbed.SetTitle("평가 매칭 성공!")
		matchSuccessEmbed.AddField(
			"평가자 intra ID:",
			matchClient.FindIntraByUID(evalInfo.InterviewerID),
		)
		matchSuccessEmbed.AddField(
			"평가할 서브젝트:",
			evalInfo.SubjectName + "\n" +
				evalInfo.SubjectURL,
		)
		s.ChannelMessageSendEmbed(dmChan.ID, matchSuccessEmbed.MessageEmbed)
	}
}

// sendCommandDetail 함수는 명령어 정보를 전부 전송하는 함수이다.
func sendCommandDetail(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandDetailEmbed := embed.NewEmbed()
	commandDetailEmbed.SetTitle("명령어 목록")
	commandDetailEmbed.AddField(
		"도움말 명령어",
		"!도움말" + "\n" +
			"!명령어",
			)
	commandDetailEmbed.AddField(
		"제출 명령어",
		"!제출 <GitRepoURL> <SubjectID>" + "\n" +
			"!제출취소",
		)
	commandDetailEmbed.AddField(
		"평가자 등록 명령어",
		"!평가등록" + "\n" +
			"!평가취소",
		)
	commandDetailEmbed.AddField(
		"정보 확인 명령어",
		"!내점수" + "\n" +
			"!매칭상태",
		)
	s.ChannelMessageSendEmbed(m.ChannelID, commandDetailEmbed.MessageEmbed)
}