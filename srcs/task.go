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
	mid, _ := s.ChannelMessageSend(r.ChannelID, "채점을 등록하였습니다.")
	gitUrl := submitURLs[r.UserID]
	subjectID := submitSIDs[r.UserID]
	code := c.Submit(subjectID, r.UserID, gitUrl)
	response := schema.Success{}
	_ = json.Unmarshal(code, &response)
	result := schema.EvalResult{
		Course: subjectID,
		Pass:   response.Success,
	}
	ctx := context.Background()
	curUser := schema.Person{}
	c.MDB.Collection("people").FindOne(ctx, bson.D{{Key: "name", Value: r.UserID}}).Decode(&curUser)
	curUser.Score = append(curUser.Score, result)
	c.MDB.Collection("people").UpdateOne(ctx, bson.D{{Key: "name", Value: r.UserID}}, curUser)
	s.ChannelMessageEdit(r.ChannelID, mid.ID, "채점이 완료되었습니다..!")
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
		scoreEmbed.AddField(result.Course, "[ OK ]")
	} else {
		scoreEmbed.AddField(result.Course, "[ KO ]")
	}
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
			"***주의*** 평가가 매칭된 후, 제출을 취소할 수 없음!",
			"평가받을 Git Repo: "+command[1]+"\n"+
				"평가받을 Subject : "+subjectName+"\n"+
				"아직 매칭되지 않은 평가를 취소하고 싶다면 "+prefix+"제출취소 명령어를 사용하세요",
		),
	)
	submitMIDs[m.Author.ID] = submitMsg.ID
	submitURLs[m.Author.ID] = command[1]
	submitSIDs[m.Author.ID] = subjectName
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "⭕")
	s.MessageReactionAdd(dmChan.ID, submitMsg.ID, "❌")
}
