package user

import (
	"encoding/json"
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/models/errs"
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
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedUser, err := h.usecase.SetUserActivity(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if err == errs.ErrNotFound {
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_FOUND"
			errorResp.Error.Message = "user not found"
			response.SendJSONResponse(r.Context(), w, http.StatusNotFound, errorResp)
			return
		}
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.ToUserResponseDTO(updatedUser)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}

func (h *UserHandler) GetUserReviewPRs(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	prs, err := h.usecase.GetUserReviewPRs(r.Context(), userID)
	if err != nil {
		if err == errs.ErrNotFound {
			errorResp := dto.ErrorResponseDTO{}
			errorResp.Error.Code = "NOT_FOUND"
			errorResp.Error.Message = "user not found"
			response.SendJSONResponse(r.Context(), w, http.StatusNotFound, errorResp)
			return
		}
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.UserReviewResponseDTO{
		UserID:       userID,
		PullRequests: prs,
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
