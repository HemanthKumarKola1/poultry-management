package handlers

import (
	"encoding/json"
	"net/http"

	"poultry-management.com/place-holder/archive/tenants"
)

func CreateTenantHandler(s *tenants.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		tenant, err := s.CreateTenant(req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // More specific error handling in production
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tenant)
	}
}
