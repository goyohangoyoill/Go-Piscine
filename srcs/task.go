package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"piscine-golang-interact/client"
	"piscine-golang-interact/schema"

	"github.com/bwmarrin/discordgo"
	embed "github.com/clinet/discordgo-embed"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func submissionResponse(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	s.ChannelMessageSend(r.ChannelID, "채점을 등록하였습니다. 잠시 기다려주세요...")
	gitUrl := submitURLs[r.UserID]
	subjectID := submitSIDs[r.UserID]
	code := c.Submit(subjectID, r.UserID, gitUrl)
	response := schema.Success{}
	_ = json.Unmarshal(code, &response)
	result := schema.EvalResult{
		Course: subjectID,
		Pass:   response.Success,
	}
	if response.Error {
		s.ChannelMessageSend(r.ChannelID, response.Content)
		return
	}
	ctx := context.Background()
	curUser := schema.Person{}
	c.MDB.Collection("people").FindOne(ctx, bson.D{{Key: "password", Value: r.UserID}}).Decode(&curUser)
	if curUser.Score == nil {
		curUser.Score = make([]schema.EvalResult, 0)
	}
	curUser.Score = append(curUser.Score, result)
	_, err := c.MDB.Collection("people").UpdateOne(ctx, bson.D{{Key: "password", Value: r.UserID}},
	bson.D{{Key: "$set", Value: bson.D{{Key: "score", Value: curUser.Score}}}})
	if err != nil {
		log.Warn(err)
	}
	s.ChannelMessageSend(r.ChannelID, "채점이 완료되었습니다..!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(r.ChannelID, "...\n")
	time.Sleep(time.Second)
	s.ChannelMessageSend(r.ChannelID, "...\n")
	time.Sleep(time.Second)
	s.ChannelMessageSend(r.ChannelID, "...\n")
	time.Sleep(time.Second)

	scoreEmbed := embed.NewEmbed()
	scoreEmbed.SetTitle(curUser.Name + "의 채점 결과")
	if result.Pass {
		scoreEmbed.AddField(result.Course, "[ ✅ ]")
	} else {
		scoreEmbed.AddField(result.Course, "[ 💥 ]")
	}
	scoreEmbed.AddField(result.Course + "의 평가지", client.SubjectInfoMap[result.Course].EvalGuideURL + "\n" +
		"📔 최소 2명 이상의 동료에게 동료평가를 받으십시오.")
	s.ChannelMessageSendEmbed(r.ChannelID, scoreEmbed.MessageEmbed)
}

func submissionTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	if c.FindIntraByUID(m.Author.ID) == "가입하지 않은 사용자" {
		if mode {
			s.ChannelMessageSendEmbed(m.ChannelID,
				embed.NewGenericEmbed("피신 등록을 진행하지 않은 사용자입니다",
					prefix+"인트라등록 명령어를 이용해\n"+
						"인트라 등록을 진행해 주시기 바랍니다."))
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID,
				embed.NewGenericEmbed("피신 등록을 진행하지 않은 사용자입니다.",
					"인트라등록 기간이 지났습니다.\n"+
						"관리자에게 문의 바랍니다."))
		}
		return
	}
	command := strings.Split(m.Content, " ")
	if len(command) != 3 {
		log.Println("uid:", m.Author.ID, ", 포맷과 다른 제출")
		submitEmbed := embed.NewEmbed()
		submitEmbed.SetTitle("제출 명령어는 다음과 같이 입력해야 합니다.")
		submitEmbed.AddField(
			"<명령어 예시>",
			prefix+"제출 <github repo url> <subject name>\n"+
				prefix+"제출 https://github.com/example123/ExampleRepo Day01")
		s.ChannelMessageSendEmbed(m.ChannelID, submitEmbed.MessageEmbed)
		return
	}
	subjectName := client.ConvSubjectName(command[2])
	if subjectName == "존재하지 않는 서브젝트" {
		s.ChannelMessageSend(m.ChannelID, "서브젝트명 "+command[2]+" 는 존재하지 않습니다.")
		return
	}
	dmChan, _ := s.UserChannelCreate(m.Author.ID)
	submitMsg, _ := s.ChannelMessageSendEmbed(dmChan.ID,
		embed.NewGenericEmbed(
			"**주의** 제출된 깃허브 레포는 Public 이어야 합니다.",
			"평가받을 Subject : "+subjectName+"\n"+
				"평가받을 Git Repo: "+command[1]+"\n"+
				"평가가 끝난 후 레포를 Private 으로 바꾸는 것 잊지 마십시오!",
		),
	)
	submitMIDs[m.Author.ID] = submitMsg.ID
	submitURLs[m.Author.ID] = command[1]
	submitSIDs[m.Author.ID] = subjectName
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "⭕")
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "❌")
}
