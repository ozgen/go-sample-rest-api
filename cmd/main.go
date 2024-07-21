package main

import (
	"go-sample-rest-api/cmd/api"
	"go-sample-rest-api/config"
	"go-sample-rest-api/db"
	"go-sample-rest-api/logging"
	"go-sample-rest-api/storage"
)

func main() {
	log := logging.GetLogger()
	cfg := config.Envs
	db, err := db.NewPostgresStorageConn(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	azStorage := storage.NewAzureStorage(cfg.AzureStorageAccountName, cfg.AzureContainerAccessKey, cfg.AzureContainerName)
	serverAddress := ":" + cfg.ServerPort
	server := api.NewAPIServer(serverAddress, db, azStorage)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
