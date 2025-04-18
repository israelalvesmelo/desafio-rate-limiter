package handler

// Este arquivo contém trechos de código licenciados sob a Licença MIT,
// originalmente criados por devfullcycle. Veja o arquivo LICENSE para detalhes.

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
)

type RateLimitConfigHandler struct {
	useCase usecase.CreateRateLimitConfigUseCase
}

func NewRateLimitConfigHandler(useCase usecase.CreateRateLimitConfigUseCase) *RateLimitConfigHandler {
	return &RateLimitConfigHandler{
		useCase: useCase,
	}
}

func (h *RateLimitConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	input := dto.RateLimitConfigInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("error decoding input data:", err.Error())
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateInput(input); err != nil {
		log.Println("erro de validação:", err)
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.useCase.Execute(r.Context(), usecase.CreateRateLimitConfigInput{
		Key:             infra.GetClientIP(r),
		IsToken:         input.IsToken,
		MaxRequests:     input.MaxRequests,
		TimeWindow:      input.TimeWindow,
		BlockedDuration: input.BlockDuration,
	})

	if err != nil {
		log.Println("error creating rate limit config:", err.Error())
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *RateLimitConfigHandler) validateInput(input dto.RateLimitConfigInput) error {
	validate := validator.New()
	return validate.Struct(input)
}
