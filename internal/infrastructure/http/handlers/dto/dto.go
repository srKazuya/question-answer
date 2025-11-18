package handlerdto

import (
	"time"
	resp "question-answer/pkg/validator"
)

type AddQuestionRequest struct {
	Text string `json:"text" validate:"required,min=3,max=500"`
}

type AddQuestionResponse struct {
	resp.ValidationResponse
	Text      string    `json:"text" validate:"required,min=3,max=500"`
	CreatedAt time.Time `json:"created_at"`
}
