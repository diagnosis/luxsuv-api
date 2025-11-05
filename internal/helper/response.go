package helper

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/diagnosis/luxsuv-api-v2/internal/apperror"
	"github.com/diagnosis/luxsuv-api-v2/internal/logger"
)

type ErrorResponse struct {
	Error struct {
		Code          string    `json:"code"`
		Message       string    `json:"message"`
		CorrelationID string    `json:"correlation_id,omitempty"`
		Timestamp     time.Time `json:"timestamp"`
	} `json:"error"`
}

type SuccessResponse struct {
	Data          any       `json:"data,omitempty"`
	Message       string    `json:"message,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()
	correlationID := logger.GetCorrelationID(ctx)

	appErr := apperror.AsAppError(err)

	logger.Error(ctx, "handler error",
		"error_code", appErr.Code,
		"error_message", appErr.Message,
		"http_status", appErr.HTTPStatus,
		"underlying_error", appErr.Err,
	)

	response := ErrorResponse{}
	response.Error.Code = string(appErr.Code)
	response.Error.Message = appErr.Message
	response.Error.CorrelationID = correlationID
	response.Error.Timestamp = time.Now().UTC()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(response)
}

func RespondJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	ctx := r.Context()
	correlationID := logger.GetCorrelationID(ctx)

	response := SuccessResponse{
		Data:          data,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func RespondMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	ctx := r.Context()
	correlationID := logger.GetCorrelationID(ctx)

	response := SuccessResponse{
		Message:       message,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
