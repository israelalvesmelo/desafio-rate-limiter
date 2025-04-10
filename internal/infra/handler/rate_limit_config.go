package handler

// Este arquivo contém trechos de código licenciados sob a Licença MIT,
// originalmente criados por devfullcycle. Veja o arquivo LICENSE para detalhes.

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/israelalvesmelo/desafio-rate-limiter/internal/domain/usecase"
	"github.com/israelalvesmelo/desafio-rate-limiter/internal/infra/dto"
)

type APIKeyHandler struct {
	useCase usecase.CreateRateLimitConfigUseCase
}

func NewAPIKeyHandler(useCase usecase.CreateRateLimitConfigUseCase) *APIKeyHandler {
	return &APIKeyHandler{
		useCase: useCase,
	}
}

func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	input := dto.RateLimitConfigInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("error decoding input data:", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.useCase.Execute(r.Context(), usecase.CreateRateLimitConfigInput{
		Ip:            h.getClientIP(r),
		IsToken:       input.IsToken,
		Limit:         input.Limit,
		BlockDuration: input.BlockDuration,
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
func (h *APIKeyHandler) getClientIP(r *http.Request) string {
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
