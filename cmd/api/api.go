package api

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"go-sample-rest-api/logging"
	auth2 "go-sample-rest-api/service/auth"
	"go-sample-rest-api/service/camerametadata"
	"go-sample-rest-api/service/user"
	"go-sample-rest-api/storage"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	address      string
	db           *sql.DB
	azureStorage *storage.AzureStorage
}

func NewAPIServer(addr string, db *sql.DB, azureStorage *storage.AzureStorage) *APIServer {
	return &APIServer{
		address:      addr,
		db:           db,
		azureStorage: azureStorage,
	}
}

func (s *APIServer) Run() error {
	log := logging.GetLogger()

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//user
	userStore := user.NewStore(s.db)
	auth := auth2.NewAuthenticator()
	userService := user.NewHandler(userStore, auth)
	userService.RegisterRoutes(subrouter)

	// cameraMetadata
	cameraMetadataStore := camerametadata.NewStore(s.db)
	cameraMetadataService := camerametadata.NewHandler(cameraMetadataStore, s.azureStorage)
	cameraMetadataService.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.WithFields(logrus.Fields{
		"address": s.address,
	}).Info("Listening on")

	return http.ListenAndServe(s.address, router)
}
