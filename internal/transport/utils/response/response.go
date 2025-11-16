package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/dto"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/middleware/logctx"
)

func SendJSONError(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResp := dto.ErrorResponseDTO{}
	errorResp.Error.Code = getErrorCode(statusCode)
	errorResp.Error.Message = message

	resp, err := json.Marshal(errorResp)
	if err != nil {
		logctx.GetLogger(ctx).Error("failed to marshal response: ", err.Error())
		return
	}

	if _, err = w.Write(resp); err != nil {
		logctx.GetLogger(ctx).Error("failed to write response: ", err.Error())
	}
}

func SendJSONResponse(ctx context.Context, w http.ResponseWriter, statusCode int, body any) {
	if body == nil {
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logctx.GetLogger(ctx).Error("failed to marshal response", err.Error())
		return
	}

	w.WriteHeader(statusCode)
	if _, err = w.Write(resp); err != nil {
		logctx.GetLogger(ctx).Error("failed to write response", err.Error())
	}
}

func HandleDomainError(ctx context.Context, w http.ResponseWriter, err error) {
	log := logctx.GetLogger(ctx)

	switch {
	case errors.Is(err, errs.ErrTeamExists):
		SendJSONError(ctx, w, http.StatusConflict, "team_name already exists")
		log.Debug("team_name already exists: ", err.Error())

	case errors.Is(err, errs.ErrPRExists):
		SendJSONError(ctx, w, http.StatusConflict, "PR already exists")
		log.Debug("PR already exists: ", err.Error())

	case errors.Is(err, errs.ErrPRMerged):
		SendJSONError(ctx, w, http.StatusConflict, "PR is already merged")
		log.Debug("PR is merged: ", err.Error())

	case errors.Is(err, errs.ErrNotAssigned):
		SendJSONError(ctx, w, http.StatusConflict, "reviewer is not assigned to this PR")
		log.Debug("reviewer not assigned: ", err.Error())

	case errors.Is(err, errs.ErrNoCandidate):
		SendJSONError(ctx, w, http.StatusConflict, "no active replacement candidate available")
		log.Debug("no candidate: ", err.Error())

	case errors.Is(err, errs.ErrUserNotActive):
		SendJSONError(ctx, w, http.StatusConflict, "user is not active")
		log.Debug("user not active: ", err.Error())

	case errors.Is(err, errs.ErrUserNoTeam):
		SendJSONError(ctx, w, http.StatusConflict, "user does not belong to any team")
		log.Debug("user has no team: ", err.Error())

	case errors.Is(err, errs.ErrInvalidID):
		SendJSONError(ctx, w, http.StatusBadRequest, "invalid id format")
		log.Debug("invalid id format: ", err.Error())

	case errors.Is(err, errs.ErrNotFound):
		SendJSONError(ctx, w, http.StatusNotFound, "resource not found")
		log.Debug("resource not found: ", err.Error())

	default:
		SendJSONError(ctx, w, http.StatusInternalServerError, "internal server error")
		log.Error("unexpected error: ", err.Error())
	}
}

func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}
