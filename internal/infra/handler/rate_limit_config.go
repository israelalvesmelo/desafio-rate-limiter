package handler

// Este arquivo contém trechos de código licenciados sob a Licença MIT,
// originalmente criados por devfullcycle. Veja o arquivo LICENSE para detalhes.

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.validateInput(input); err != nil {
		log.Println("Erro de validação:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.useCase.Execute(r.Context(), usecase.CreateRateLimitConfigInput{
		Ip:              h.getClientIP(r),
		IsToken:         input.IsToken,
		MaxRequests:     input.MaxRequests,
		TimeWindow:      input.TimeWindow,
		BlockedDuration: input.BlockDuration,
	})

	if err != nil {
		log.Println("error creating rate limit config:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// getClientIP extracts the client IP from the request
func (h *RateLimitConfigHandler) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := net.ParseIP(xff)
		if ips != nil {
			return ips.String()
		}
	}

	// Extract from RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (h *RateLimitConfigHandler) validateInput(input dto.RateLimitConfigInput) error {
	validate := validator.New()
	return validate.Struct(input)
}
