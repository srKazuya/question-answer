package handlerdto

import (
	resp "question-answer/pkg/validator"
	"time"
)

type AnswerResponse struct {
	Text      string    `json:"text" validate:"required,min=3,max=500"`
	CreatedAt time.Time `json:"created_at"`
}

type AnswerRequest struct {
	Text string `json:"text" validate:"required,min=3,max=500"`
}

type AddAnswerResponse struct {
	resp.ValidationResponse
	ID uint64
}
