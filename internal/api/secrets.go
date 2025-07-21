package api

import (
	"net/http"
	"time"

	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateSecretRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SecretResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ListSecretsResponse struct {
	Secrets []SecretResponse `json:"secrets"`
}

func (hs *HTTPServer) CreateSecret(w http.ResponseWriter, r *http.Request) {
	var req CreateSecretRequest
	if err := parseJSON(r, &req); err != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Name == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Secret name is required",
		})
		return
	}

	if req.Value == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Secret value is required",
		})
		return
	}

	// For now, we'll use a dummy user ID. In a real implementation
	userID := pgtype.UUID{}
	userID.Scan("00000000-0000-0000-0000-000000000001")

	encrypted, err := hs.secretManager.Encrypt(req.Value)
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to encrypt secret",
		})
		return
	}

	secret, err := hs.store.CreateSecret(r.Context(), db.CreateSecretParams{
		UserID:         userID,
		Name:           req.Name,
		EncryptedValue: encrypted,
	})
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to store secret",
		})
		return
	}

	respondJSON(w, 201, SecretResponse{
		ID:        secret.ID.String(),
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt.Time.Format(time.RFC3339),
	})
}

func (hs *HTTPServer) ListSecrets(w http.ResponseWriter, r *http.Request) {
	// For now, we'll use a dummy user ID
	userID := pgtype.UUID{}
	userID.Scan("00000000-0000-0000-0000-000000000001")

	secrets, err := hs.store.ListSecrets(r.Context(), userID)
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to retrieve secrets",
		})
		return
	}

	var response []SecretResponse
	for _, secret := range secrets {
		response = append(response, SecretResponse{
			ID:        secret.ID.String(),
			Name:      secret.Name,
			CreatedAt: secret.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	respondJSON(w, 200, ListSecretsResponse{
		Secrets: response,
	})
}

func (hs *HTTPServer) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	secretID := r.URL.Query().Get("id")
	if secretID == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Secret ID is required",
		})
		return
	}

	parsedSecretID, err := utils.ParseUUID(secretID)
	if err != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Invalid secret ID",
		})
		return
	}

	// For now, we'll use a dummy user ID
	userID := pgtype.UUID{}
	userID.Scan("00000000-0000-0000-0000-000000000001")

	err = hs.store.DeleteSecret(r.Context(), db.DeleteSecretParams{
		ID:     parsedSecretID,
		UserID: userID,
	})
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to delete secret",
		})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"message": "Secret deleted successfully",
	})
}

func (hs *HTTPServer) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretID := r.URL.Query().Get("id")
	if secretID == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Secret ID is required",
		})
		return
	}

	parsedSecretID, err := utils.ParseUUID(secretID)
	if err != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Invalid secret ID",
		})
		return
	}

	var req struct {
		Value string `json:"value"`
	}
	if err := parseJSON(r, &req); err != nil {
		respondJSON(w, 400, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Value == "" {
		respondJSON(w, 400, map[string]interface{}{
			"error": "Secret value is required",
		})
		return
	}

	// For now, we'll use a dummy user ID
	userID := pgtype.UUID{}
	userID.Scan("00000000-0000-0000-0000-000000000001")

	encrypted, err := hs.secretManager.Encrypt(req.Value)
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to encrypt secret",
		})
		return
	}

	err = hs.store.UpdateSecret(r.Context(), db.UpdateSecretParams{
		ID:             parsedSecretID,
		EncryptedValue: encrypted,
		UserID:         userID,
	})
	if err != nil {
		respondJSON(w, 500, map[string]interface{}{
			"error": "Failed to update secret",
		})
		return
	}

	respondJSON(w, 200, map[string]interface{}{
		"message": "Secret updated successfully",
	})
}
