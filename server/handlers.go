package server

import (
	"context"
	"encoding/json"
	"main/entities"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *APIServer) CreateProviderHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	newProvider := &entities.ProviderCreationRequest{}
	if err := json.NewDecoder(r.Body).Decode(newProvider); err != nil {
		return err
	}
	if err := newProvider.Validate(); err != nil {
		return err
	}
	if err := s.Service.CreateProvider(newProvider.ToSchema()); err != nil {
		return err
	}
	return writeJSON(w, http.StatusCreated, map[string]string{"message": "created"})
}

func (s *APIServer) GetProvidersHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	providers, err := s.Storage.GetAllProviders()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, providers)
}

func (s APIServer) getProviderById(ctx context.Context, r *http.Request) (*entities.Provider, error) {
	providerId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return nil, err
	}
	return s.Storage.GetProviderById(providerId)
}

func (s *APIServer) GetProviderByIdHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	provider, err := s.getProviderById(ctx, r)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, provider)
}

func (s *APIServer) SendSMSHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	providerId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	reqBody := &entities.SendSMSRequest{}
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		return err
	}
	if err := reqBody.Validate(); err != nil {
		return err
	}

	return s.Service.SendSMSWith(providerId, reqBody.Message, reqBody.PhoneNumber)
}
