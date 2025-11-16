package user

import (
	"encoding/json"
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/models"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/dto"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/utils/response"
	"github.com/marrria_mme/pr-reviewer-service/internal/usecase"
)

type UserHandler struct {
	usecase usecase.IUserUsecase
}

func NewUserHandler(usecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) SetUserActivity(w http.ResponseWriter, r *http.Request) {
	var req dto.UserActivityRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	updatedUser, err := h.usecase.SetUserActivity(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.ToUserResponseDTO(updatedUser)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}

func (h *UserHandler) GetUserReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "user_id parameter is required")
		return
	}

	prs, err := h.usecase.GetUserReviewPRs(r.Context(), userID)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	if prs == nil {
		prs = []models.PullRequestShort{}
	}

	responseDTO := dto.NewUserReviewResponseDTO(userID, prs)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
