package pgdto

import (
	"question-answer/internal/domain/users"
	"question-answer/internal/domain/qa"
	"time"
)

type QuestionDTO struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	Text      string    `gorm:"type:varchar(500);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (QuestionDTO) TableName() string {
	return "questions"
}

type AnswerDTO struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	QuestionID uint64    `gorm:"index;not null"`
	UserID     uint64    `gorm:"not null"`
	Text       string    `gorm:"type:varchar(1000);not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

func (AnswerDTO) TableName() string {
	return "answers"
}

type UserDTO struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"type:varchar(32);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
}

func (UserDTO) TableName() string {
	return "users"
}

//Question

func ToDomainQuestion(q QuestionDTO) qa.Question {
	return qa.Question{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
	}
}

func ToDTOQuestion(q qa.Question) QuestionDTO {
	return QuestionDTO{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
	}
}

// Answer

func ToDomainAnswer(a AnswerDTO) qa.Answer {
	return qa.Answer{
		ID:         a.ID,
		QuestionID: a.QuestionID,
		UserID:     a.UserID,
		Text:       a.Text,
		CreatedAt:  a.CreatedAt,
	}
}

func ToDTOAnswer(a qa.Answer) AnswerDTO {
	return AnswerDTO{
		ID:         a.ID,
		QuestionID: a.QuestionID,
		UserID:     a.UserID,
		Text:       a.Text,
		CreatedAt:  a.CreatedAt,
	}
}

// User

func ToDomainUser(u UserDTO) auth.User {
	return auth.User{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}

func ToDTOUser(u auth.User) UserDTO {
	return UserDTO{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
	}
}
