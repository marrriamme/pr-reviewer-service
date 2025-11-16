package team

import (
	"encoding/json"
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/transport/dto"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/utils/response"
	"github.com/marrria_mme/pr-reviewer-service/internal/usecase"
)

type TeamHandler struct {
	usecase usecase.ITeamUsecase
}

func NewTeamHandler(usecase usecase.ITeamUsecase) *TeamHandler {
	return &TeamHandler{usecase: usecase}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dto.TeamRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "invalid request body")
		return
	}

	team := dto.ToTeamModel(req)
	if err := h.usecase.CreateTeam(r.Context(), team); err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.ToTeamResponseDTO(team)
	response.SendJSONResponse(r.Context(), w, http.StatusCreated, responseDTO)
}

func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "BAD_REQUEST", "team_name parameter is required")
		return
	}

	team, err := h.usecase.GetTeam(r.Context(), teamName)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.ToTeamResponseDTO(team)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
