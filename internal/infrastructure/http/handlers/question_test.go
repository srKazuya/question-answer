package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"question-answer/internal/domain/qa"

	"question-answer/internal/infrastructure/http/handlers"
	dto "question-answer/internal/infrastructure/http/handlers/dto"
	"question-answer/internal/infrastructure/http/handlers/mocks"
	slogdiscard "question-answer/pkg/sl_logger/slog_discard"
	validateresp "question-answer/pkg/validator"

	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddQuestionHandler(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name           string
		reqBody        string
		mockReturnQ    *qa.Question
		mockReturnErr  error
		expectedStatus int
		expectedResp   dto.AddQuestionResponse
	}{
		{
			name:    "Success",
			reqBody: `{"text": "Почему небо голубое?"}`,
			mockReturnQ: &qa.Question{
				ID:        123,
				Text:      "Почему небо голубое?",
				CreatedAt: fixedTime,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
			expectedResp: dto.AddQuestionResponse{
				ValidationResponse: validateresp.OK(), 
				Text:               "Почему небо голубое?",
				CreatedAt:          fixedTime,
			},
		},
		{
			name:           "Invalid JSON",
			reqBody:        `{"text": "valid"`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: dto.AddQuestionResponse{
				ValidationResponse: validateresp.Error("failed to decode request body"),
			},
		},
		{
			name:           "Service error",
			reqBody:        `{"text": "Этот вопрос упадёт"}`,
			mockReturnQ:    nil,
			mockReturnErr:  errors.New("db down"),
			expectedStatus: http.StatusBadRequest,
			expectedResp: dto.AddQuestionResponse{
				ValidationResponse: validateresp.Error("failed to decode request body"),
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svcMock := mocks.NewService(t)

			if tc.name == "Success" || tc.name == "Service error" {
				svcMock.On("CreateQuestion", mock.AnythingOfType("qa.Question")).
					Return(tc.mockReturnQ, tc.mockReturnErr).
					Once()
			}

			handler := handlers.NewAddQuestionHandler(slogdiscard.NewDiscardLogger(), svcMock)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tc.reqBody)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			var resp dto.AddQuestionResponse
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tc.expectedResp.Status, resp.Status)
			require.Equal(t, tc.expectedResp.Errors, resp.Errors)
			require.Equal(t, tc.expectedResp.Text, resp.Text)

			if tc.expectedStatus == http.StatusOK {
				require.WithinDuration(t, tc.expectedResp.CreatedAt, resp.CreatedAt, time.Second)
			}

			svcMock.AssertExpectations(t)
		})
	}
}
