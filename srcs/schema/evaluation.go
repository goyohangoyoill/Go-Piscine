package schema

import (
	"time"
)

type Match struct {
	ID            int       `json:"id"`
	IntervieweeID string    `json:"intervieweeID"`
	InterviewerID string    `json:"interviewerID"`
	Course        string    `json:"course"`
	Score         int       `json:"score"`
	Pass          bool      `json:"pass"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}
