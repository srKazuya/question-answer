package handlerdto

import (
	resp "question-answer/pkg/validator"
	"time"
)

type QuestionResponse struct {
	Text      string    `json:"text" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at"`
}
type AddQuestionRequest struct {
	Text string `json:"text" validate:"required,min=3,max=500"`
}

type AddQuestionResponse struct {
	resp.ValidationResponse
	Text      string    `json:"text" validate:"required,min=3,max=500"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteQuestionResponse struct {
	resp.ValidationResponse
}

type DeleteQuestionRequest struct {
	resp.ValidationResponse
}

type GetQuestionResponse struct {
	resp.ValidationResponse
	Data []AddQuestionResponse `json:"data"`
}
