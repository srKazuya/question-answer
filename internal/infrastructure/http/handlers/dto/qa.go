package handlerdto

import (
	resp "question-answer/pkg/validator"
)

type QAResponse struct {
	resp.ValidationResponse
	Data QAData
}
type QAData struct {
	Question QuestionResponse `json:"question"`
	Answers  []AnswerResponse `json:"answers"`
}
