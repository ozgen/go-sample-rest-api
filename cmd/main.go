package main

import (
	"go-sample-rest-api/cmd/api"
	"go-sample-rest-api/config"
	db2 "go-sample-rest-api/db"
	"go-sample-rest-api/logging"
	"go-sample-rest-api/storage"
)

func SetupServer() (*api.APIServer, error) {
	log := logging.GetLogger()
	cfg := config.Envs
	db, err := db2.NewPostgresStorageConn(cfg)
	if err != nil {
		log.Error("Failed to connect to database:", err)
		return nil, err
	}
	sqldb := db2.NewSQLDB(db)
	azStorage := storage.NewAzureStorage(cfg.AzureStorageAccountName, cfg.AzureContainerAccessKey, cfg.AzureContainerName)
	serverAddress := ":" + cfg.ServerPort
	server := api.NewAPIServer(serverAddress, sqldb, azStorage)
	return server, nil
}

func main() {
	server, err := SetupServer()
	if err != nil {
		logging.GetLogger().Fatal(err)
		return
	}

	if err := server.Run(); err != nil {
		logging.GetLogger().Fatal(err)
	}
}
