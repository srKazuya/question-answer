// Package handlers provides HTTP handler functions for the service.
package handlers

import (
	"question-answer/internal/domain/qa"
	dto "question-answer/internal/infrastructure/http/handlers/dto"
	"question-answer/internal/infrastructure/http/middleware"
	"question-answer/internal/infrastructure/http/transport"
	"question-answer/pkg/sl_logger/sl"
	validateResp "question-answer/pkg/validator"
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
			transport.WriteJSON(w, http.StatusBadRequest, validateResp.ValidationError(validateErr))
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

		addEventResponseOK(w, reqQuestion.Text, reqQuestion.CreatedAt)
	}
}

// Get
func NewGetQuestionHandler(log *slog.Logger, svc qa.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

// get
func NewGetAllQuestionHandler(log *slog.Logger, svc qa.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

// delete
func NewDeleteQuestionHandler(log *slog.Logger, svc qa.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func addEventResponseOK(w http.ResponseWriter, text string, created time.Time) {
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
