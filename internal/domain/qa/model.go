// Package qa provides domain models for questions and answers.
package qa

import "time"

type Question struct {
	ID        uint64    `json:"id"`
	Text      string    `json:"text" validate:"required,min=3,max=500"`
	CreatedAt time.Time `json:"created_at"`
}

type Answer struct {
	ID         uint64    `json:"id"`
	QuestionID uint64    `json:"question_id" validate:"required"`
	UserID     uint64    `json:"user_id" validate:"required"`
	Text       string    `json:"text" validate:"required,min=1,max=1000"`
	CreatedAt  time.Time `json:"created_at"`
}
