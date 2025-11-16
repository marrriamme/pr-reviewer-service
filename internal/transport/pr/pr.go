package pr

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/dto"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/utils/response"
	"github.com/marrria_mme/pr-reviewer-service/internal/usecase"
)

type PRHandler struct {
	usecase usecase.IPRUsecase
}

func NewPRHandler(usecase usecase.IPRUsecase) *PRHandler {
	return &PRHandler{usecase: usecase}
}

func (h *PRHandler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePRRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	prModel := dto.ToPRModel(req)

	createdPR, err := h.usecase.CreatePR(r.Context(), prModel)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrPRExists):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "PR_EXISTS"
			errorResp.Error.Message = "PR id already exists"
			response.SendJSONResponse(r.Context(), w, http.StatusConflict, errorResp)
			return
		case errors.Is(err, errs.ErrNotFound):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_FOUND"
			errorResp.Error.Message = "author or team not found"
			response.SendJSONResponse(r.Context(), w, http.StatusNotFound, errorResp)
			return
		case errors.Is(err, errs.ErrUserNotActive), errors.Is(err, errs.ErrUserNoTeam):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "INVALID_AUTHOR"
			errorResp.Error.Message = err.Error()
			response.SendJSONResponse(r.Context(), w, http.StatusConflict, errorResp)
			return
		default:
			response.HandleDomainError(r.Context(), w, err)
			return
		}
	}

	if createdPR.AssignedReviewers == nil {
		createdPR.AssignedReviewers = []string{}
	}

	responseDTO := dto.ToPRResponseDTO(createdPR)
	response.SendJSONResponse(r.Context(), w, http.StatusCreated, responseDTO)
}

func (h *PRHandler) MergePR(w http.ResponseWriter, r *http.Request) {
	var req dto.MergePRRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	mergedPR, err := h.usecase.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_FOUND"
			errorResp.Error.Message = "PR not found"
			response.SendJSONResponse(r.Context(), w, http.StatusNotFound, errorResp)
			return
		}
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	if mergedPR.AssignedReviewers == nil {
		mergedPR.AssignedReviewers = []string{}
	}

	responseDTO := dto.ToPRResponseDTO(mergedPR)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}

func (h *PRHandler) ReassignReviewer(w http.ResponseWriter, r *http.Request) {
	var req dto.ReassignPRRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedPR, newUserID, err := h.usecase.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrPRMerged):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "PR_MERGED"
			errorResp.Error.Message = "cannot reassign on merged PR"
			response.SendJSONResponse(r.Context(), w, http.StatusConflict, errorResp)
			return
		case errors.Is(err, errs.ErrNotAssigned):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_ASSIGNED"
			errorResp.Error.Message = "reviewer is not assigned to this PR"
			response.SendJSONResponse(r.Context(), w, http.StatusConflict, errorResp)
			return
		case errors.Is(err, errs.ErrNoCandidate):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NO_CANDIDATE"
			errorResp.Error.Message = "no active replacement candidate in team"
			response.SendJSONResponse(r.Context(), w, http.StatusConflict, errorResp)
			return
		case errors.Is(err, errs.ErrNotFound):
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_FOUND"
			errorResp.Error.Message = "PR or user not found"
			response.SendJSONResponse(r.Context(), w, http.StatusNotFound, errorResp)
			return
		default:
			response.HandleDomainError(r.Context(), w, err)
			return
		}
	}

	if updatedPR.AssignedReviewers == nil {
		updatedPR.AssignedReviewers = []string{}
	}

	responseDTO := dto.ToReassignResponseDTO(updatedPR, newUserID)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
