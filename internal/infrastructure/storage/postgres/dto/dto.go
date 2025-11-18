// Package pgdto provides DTO models for GORM and mapping to domain models.
package pgdto

import (
    "question-answer/internal/domain/auth"
    "question-answer/internal/domain/qa"
    "time"
)

// GORM DTO MODELS
type QuestionDTO struct {
    ID        uint64    `gorm:"primaryKey;autoIncrement"`
    Text      string    `gorm:"type:varchar(500);not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type AnswerDTO struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement"`
    QuestionID uint64    `gorm:"index;not null"`
    UserID     uint64    `gorm:"not null"`
    Text       string    `gorm:"type:varchar(1000);not null"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
}

type UserDTO struct {
    ID           uint64 `gorm:"primaryKey;autoIncrement"`
    Username     string `gorm:"type:varchar(32);unique;not null"`
    PasswordHash string `gorm:"type:varchar(255);not null"`
}

// MAPPERS — Question
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

// MAPPERS — Answer
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

// MAPPERS — User (domain/auth.User ↔ UserDTO)
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
