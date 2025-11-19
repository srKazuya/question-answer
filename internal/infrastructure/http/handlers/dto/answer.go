package handlerdto

import (
	"time"
)

type AnswerResponse struct {
	Text      string    `json:"text" validate:"required,min=3,max=500"`
	CreatedAt time.Time `json:"created_at"`
}
