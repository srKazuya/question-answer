// Package handlers provides HTTP handler functions for the service.
package handlers

import (
	"question-answer/internal/domain/qa"
	dto "question-answer/internal/infrastructure/http/handlers/dto"
	"question-answer/internal/infrastructure/http/middleware"
	"question-answer/internal/infrastructure/http/transport"
	"question-answer/pkg/sl_logger/sl"
	validateResp "question-answer/pkg/validator"
	"strconv"
	"time"

	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
)

// POST
func NewAddQuestionHandler(log *slog.Logger, svc qa.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		const op = "handlers.question.add"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		var req dto.AddQuestionRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			log.Error("bad request",
				slog.String("type", transport.ErrEmptyReqBody.Error()),
				sl.Err(err),
			)
			addQuestionResponseErr(w, transport.ErrEmptyReqBody.Error())
			return
		}
		if err != nil {
			log.Error("bad request",
				slog.String("type", transport.ErrFailedToDecodeReqBody.Error()),
				sl.Err(err),
			)
			addQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			_ = transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(validateErr))
			return
		}

		respQuestion := qa.Question{
			Text: req.Text,
		}

		reqQuestion, err := svc.CreateQuestion(respQuestion)
		if err != nil {
			log.Error("failed to add quest",
				sl.Err(err),
			)
			addQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("quest added", slog.Any("title", reqQuestion.Text))

		addQuestionResponseOK(w, reqQuestion.Text, reqQuestion.CreatedAt)
	}
}

// Get
func NewGetQuestionHandler(log *slog.Logger, svc qa.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const op = "hanlers.question.get"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		reqQuestions, err := svc.GetAllQuestions()
		if err != nil {
			log.Error("failed to add quest",
				sl.Err(err),
			)
			getQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("quest added", slog.Any("quest arr geeted", len(reqQuestions)))

		getQuestionResponseOK(w, reqQuestions)
	}

}

// get
func NewGetAllQuestionHandler(log *slog.Logger, svc qa.Service, idStr string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		const op = "hanlers.question.getWiothAnswer"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to convert string",
				sl.Err(err),
			)
			http.Error(w, "intreranl server error", http.StatusMethodNotAllowed)
			return
		}

		question, answers, err := svc.GetQuestionWithAnswers(uint64(id))
		if err != nil {
			log.Error("failed to get quest",
				sl.Err(err),
			)
			getQAResponseErr(w, "failed to get quest")
			return
		}

		log.Info("quest added", slog.Any("question-answer", question.Text))

		getQAResponseOK(w, *question, answers)
	})
}

// delete
func NewDeleteQuestionHandler(log *slog.Logger, svc qa.Service, idStr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		const op = "handlers.delete.question"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("failed to convert string",
				sl.Err(err),
			)
			http.Error(w, "intreranl server error", http.StatusMethodNotAllowed)
			return
		}
		
		if err = svc.DeleteQuestion(uint64(id)); err != nil {
			log.Error("failed to delete quest",
				sl.Err(err),
			)
			deleteQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("quest deleted", slog.Any("id", id))

		deleteQuestionResponseOK(w)
	}

}

// Post Quest
func addQuestionResponseOK(w http.ResponseWriter, text string, created time.Time) {
	r := dto.AddQuestionResponse{
		ValidationResponse: validateResp.OK(),
		Text:               text,
		CreatedAt:          created,
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func addQuestionResponseErr(w http.ResponseWriter, q string) {
	r := dto.AddQuestionResponse{
		ValidationResponse: validateResp.Error(q),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

// Get Q
func getQuestionResponseOK(w http.ResponseWriter, q []qa.Question) {
	data := make([]dto.AddQuestionResponse, 0)
	for _, v := range q {
		data = append(data, dto.AddQuestionResponse{
			Text:      v.Text,
			CreatedAt: v.CreatedAt,
		})
	}
	r := dto.GetQuestionResponse{
		ValidationResponse: validateResp.OK(),
		Data:               data,
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func getQuestionResponseErr(w http.ResponseWriter, q string) {
	r := dto.GetQuestionResponse{
		ValidationResponse: validateResp.Error(q),
		Data:               nil,
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

// Get QA
func getQAResponseOK(w http.ResponseWriter, q qa.Question, a []qa.Answer) {
	answers := make([]dto.AnswerResponse, 0, len(a))
	for _, v := range a {
		answers = append(answers, dto.AnswerResponse{
			Text:      v.Text,
			CreatedAt: v.CreatedAt,
		})
	}
	r := dto.QAResponse{
		ValidationResponse: validateResp.OK(),
		Data: dto.QAData{
			Question: dto.QuestionResponse{
				Text:      q.Text,
				CreatedAt: q.CreatedAt,
			},
			Answers: answers,
		},
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func getQAResponseErr(w http.ResponseWriter, q string) {
	r := dto.GetQuestionResponse{
		ValidationResponse: validateResp.Error(q),
		Data:               nil,
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

func deleteQuestionResponseErr(w http.ResponseWriter, q string) {
	r := dto.DeleteQuestionResponse{
		ValidationResponse: validateResp.Error(q),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

func deleteQuestionResponseOK(w http.ResponseWriter) {
	r := dto.DeleteQuestionResponse{
		ValidationResponse: validateResp.OK(),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}
