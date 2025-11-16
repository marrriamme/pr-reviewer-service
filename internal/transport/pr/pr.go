package pr

import (
	"encoding/json"
	"net/http"

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
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	prModel := dto.ToPRModel(req)
	createdPR, err := h.usecase.CreatePR(r.Context(), prModel)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
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
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	mergedPR, err := h.usecase.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
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
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	updatedPR, newUserID, err := h.usecase.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	if updatedPR.AssignedReviewers == nil {
		updatedPR.AssignedReviewers = []string{}
	}

	responseDTO := dto.ToReassignResponseDTO(updatedPR, newUserID)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
