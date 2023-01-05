package server

import (
	"context"
	"encoding/json"
	"main/database"
	"main/service"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	Addr    string
	Auth    Authenticator
	Storage database.Storage
	Service service.SMSService
}

type apiFunc func(context.Context, http.ResponseWriter, *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		if err := f(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(s.AuthenticationMiddleware)

	router.HandleFunc("/providers", makeHTTPHandler(s.GetProvidersHandler)).Methods("GET")
	router.HandleFunc("/providers", makeHTTPHandler(s.CreateProviderHandler)).Methods("POST")
	router.HandleFunc("/providers/{id:[0-9]+}", makeHTTPHandler(s.GetProviderByIdHandler)).Methods("GET")
	router.HandleFunc("/providers/{id:[0-9]+}/send", makeHTTPHandler(s.SendSMSHandler)).Methods("POST")
	return http.ListenAndServe(s.Addr, router)
}
