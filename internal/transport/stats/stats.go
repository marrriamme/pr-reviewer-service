package stats

import (
	"net/http"

	"github.com/marrria_mme/pr-reviewer-service/internal/transport/dto"
	"github.com/marrria_mme/pr-reviewer-service/internal/transport/utils/response"
	"github.com/marrria_mme/pr-reviewer-service/internal/usecase"
)

type StatsHandler struct {
	usecase usecase.IStatsUsecase
}

func NewStatsHandler(usecase usecase.IStatsUsecase) *StatsHandler {
	return &StatsHandler{usecase: usecase}
}

func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	statsModel, err := h.usecase.GetStats(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err)
		return
	}

	responseDTO := dto.ToStatsResponseDTO(statsModel)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, responseDTO)
}
