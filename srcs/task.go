package main

import (
	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	"piscine-golang-interact/client"
	"strings"
)

func registerEvalResponse(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	matchedUserID := make(chan client.MatchInfo, 2)
	msg := c.Register(r.UserID, matchedUserID)
	s.ChannelMessageSend(r.ChannelID, msg)
	switch evalInfo := <-matchedUserID; evalInfo.Code {
	case false:
		// pass
	case true:
		dmChan, _ := s.UserChannelCreate(r.UserID)
		matchSuccessEmbed := embed.NewEmbed()
		matchSuccessEmbed.SetTitle("평가 매칭 성공!")
		matchSuccessEmbed.AddField(
			"피평가자 intra ID:",
			c.FindIntraByUID(evalInfo.IntervieweeID),
		)
		matchSuccessEmbed.AddField(
			"평가할 서브젝트:",
			evalInfo.Subject.SubjectName+"\n"+
				evalInfo.Subject.SubjectURL,
		)
		matchSuccessEmbed.AddField(
			"평가표 링크:",
			evalInfo.Subject.EvalGuideURL,
		)
		s.ChannelMessageSendEmbed(dmChan.ID, matchSuccessEmbed.MessageEmbed)
	}
	delete(registerMIDs, r.UserID)
}

func registerEvalTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	dmChan, _ := s.UserChannelCreate(m.Author.ID)
	regMsg, _ := s.ChannelMessageSend(dmChan.ID, "**주의** 평가가 매칭된 후, 평가는 취소할 수 없음!\n" +
		"아직 매칭되지 않은 평가를 취소하고 싶다면 " + prefix + "평가취소 명령어를 사용하세요")
	registerMIDs[m.Author.ID] = regMsg.ID
	s.MessageReactionAdd(dmChan.ID, regMsg.ID, "⭕")
	s.MessageReactionAdd(dmChan.ID, regMsg.ID, "❌")
}

func RegisterCancelTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	userChannel := c.MatchMap[m.Author.ID]
	if userChannel == nil {
		s.ChannelMessageSend(m.ChannelID, "현재 평가 등록을 하지 않은 사용자입니다.")
		return
	}
	userChannel <- client.MatchInfo{Code: false}
	s.ChannelMessageSend(m.ChannelID, "정상적으로 평가 등록이 취소되었습니다.")
}

func submissionResponse(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	gitUrl := submitURLs[r.UserID]
	subjectID := submitSIDs[r.UserID]
	matchedUserID := make(chan client.MatchInfo, 2)
	msg := c.Submit(subjectID, r.UserID, gitUrl, matchedUserID)
	s.ChannelMessageSend(r.ChannelID, msg)
	switch evalInfo := <-matchedUserID; evalInfo.Code {
	case false:
		// pass
	case true:
		dmChan, _ := s.UserChannelCreate(r.UserID)
		matchSuccessEmbed := embed.NewEmbed()
		matchSuccessEmbed.SetTitle("평가 매칭 성공!")
		matchSuccessEmbed.AddField(
			"평가자 intra ID:",
			c.FindIntraByUID(evalInfo.InterviewerID),
		)
		matchSuccessEmbed.AddField(
			"평가할 서브젝트:",
			evalInfo.Subject.SubjectName+"\n"+
				evalInfo.Subject.SubjectURL,
		)
		s.ChannelMessageSendEmbed(dmChan.ID, matchSuccessEmbed.MessageEmbed)
	}
	delete(submitMIDs, r.UserID)
	delete(submitURLs, r.UserID)
	delete(submitSIDs, r.UserID)
}

func submissionTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := strings.Split(m.Content, " ")
	if len(command) != 3 {
		submitEmbed := embed.NewEmbed()
		submitEmbed.SetTitle("제출 명령어는 다음과 같이 입력해야 합니다.")
		submitEmbed.AddField(
			"<명령어 예시>",
			prefix + "제출 <github repo url> <subject name>\n"+
				prefix + "제출 https://github.com/example123/ExampleRepo Day01")
		s.ChannelMessageSendEmbed(m.ChannelID, submitEmbed.MessageEmbed)
		return
	}
	dmChan, _ := s.UserChannelCreate(m.Author.ID)
	submitMsg, _ := s.ChannelMessageSend(dmChan.ID, "**주의** 평가가 매칭된 후, 제출을 취소할 수 없음!\n" +
		"아직 매칭되지 않은 평가를 취소하고 싶다면 " + prefix + "제출취소 명령어를 사용하세요")
	submitMIDs[m.Author.ID] = submitMsg.ID
	submitURLs[m.Author.ID] = command[1]
	submitSIDs[m.Author.ID] = command[2]
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "⭕")
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "❌")
}

func submissionCancelTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	userChannel := c.MatchMap[m.Author.ID]
	if userChannel == nil {
		s.ChannelMessageSend(m.ChannelID, "현재 제출을 하지 않은 사용자입니다.")
		return
	}
	userChannel <- client.MatchInfo{Code: false}
	s.ChannelMessageSend(m.ChannelID, "정상적으로 서브젝트 제출이 취소되었습니다.")
}
