package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"question-answer/internal/domain/qa"
	dto "question-answer/internal/infrastructure/http/handlers/dto"
	"question-answer/internal/infrastructure/http/middleware"
	"question-answer/internal/infrastructure/http/transport"
	"question-answer/pkg/sl_logger/sl"
	validators "question-answer/pkg/validator"

	"log/slog"
	"net/http"
	"strconv"
)

func NewAddAnswerHandler(log *slog.Logger, svc qa.Service, strID string) http.HandlerFunc {
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

		questID, err := strconv.Atoi(strID)
		if err != nil {
			log.Error("failed to convert string",
				sl.Err(err),
			)
			http.Error(w, "intreranl server error", http.StatusInternalServerError)
			return
		}

		var req dto.AnswerRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			log.Error("bad request",
				slog.String("type", transport.ErrEmptyReqBody.Error()),
				sl.Err(err),
			)
			addAnswerResponseErr(w, transport.ErrEmptyReqBody.Error())
			return
		}
		if err != nil {
			log.Error("bad request",
				slog.String("type", transport.ErrFailedToDecodeReqBody.Error()),
				sl.Err(err),
			)
			addAnswerResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		reqAnswer := qa.Answer{
			UserID:     1,
			Text:       req.Text,
			QuestionID: uint64(questID),
		}

		answerID, err := svc.CreateAnswer(reqAnswer)
		if err != nil {
			log.Error("failed to add Answer",
				sl.Err(err),
			)
			addQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("answer added", slog.Any("title", reqAnswer.Text))

		addAnswerResponseOK(w, answerID)
	}
}
func NewGetAnswerHandler(log *slog.Logger, svc qa.Service, strID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.question.add"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Error("mehthod not allowed")
			return
		}

		answerID, err := strconv.Atoi(strID)
		if err != nil {
			log.Error("failde to convert string",
				sl.Err(err),
			)
			http.Error(w, "internal server erro", http.StatusInternalServerError)
			return
		}

		answer, err := svc.GetAnswer(uint64(answerID))
		if err != nil {
			log.Error("failed to get Answer",
				sl.Err(err),
			)
			getQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("answer getted", slog.Any("title", answer.Text))

		getAnswerResponseOK(w, *answer)

	}
}
func NewDeleteAnswerHandler(log *slog.Logger, svc qa.Service, strID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.question.delete"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetRequestID(r)),
		)

		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			log.Error("mehthod not allowed")
			return
		}

		answerID, err := strconv.Atoi(strID)
		if err != nil {
			log.Error("failde to convert string",
				sl.Err(err),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = svc.DeleteAnswer(uint64(answerID))
		if err != nil {
			log.Error("failed to get Answer",
				sl.Err(err),
			)
			deleteQuestionResponseErr(w, transport.ErrFailedToDecodeReqBody.Error())
			return
		}

		log.Info("answer deleted")

		deleteAnswerResponseOK(w)

	}
}

func addAnswerResponseErr(w http.ResponseWriter, e string) {
	r := dto.AddAnswerResponse{
		ValidationResponse: validators.Error(e),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

func addAnswerResponseOK(w http.ResponseWriter, id uint64) {
	r := dto.AddAnswerResponse{
		ValidationResponse: validators.OK(),
		ID:                 id,
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func getAnswerResponseErr(w http.ResponseWriter, e string) {
	r := dto.AddAnswerResponse{
		ValidationResponse: validators.Error(e),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

func getAnswerResponseOK(w http.ResponseWriter, ans qa.Answer) {
	r := dto.GetAnswerResponse{
		ValidationResponse: validators.OK(),
		AnswerResponse: dto.AnswerResponse{
			Text:      ans.Text,
			CreatedAt: ans.CreatedAt,
		},
	}
	transport.WriteJSON(w, http.StatusOK, r)
}

func deleteAnswerResponseErr(w http.ResponseWriter, e string) {
	r := dto.AddAnswerResponse{
		ValidationResponse: validators.Error(e),
	}
	transport.WriteJSON(w, http.StatusBadRequest, r)
}

func deleteAnswerResponseOK(w http.ResponseWriter) {
	r := dto.GetAnswerResponse{
		ValidationResponse: validators.OK(),
	}
	transport.WriteJSON(w, http.StatusOK, r)
}
