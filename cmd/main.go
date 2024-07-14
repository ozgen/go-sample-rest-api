package main

import (
	"go-sample-rest-api/cmd/api"
	"go-sample-rest-api/config"
	"go-sample-rest-api/db"
	"log"
)

func main() {

	db, err := db.NewPostgresStorageConn(config.Envs)
	if err != nil {
		log.Fatal(err)
		return
	}

	server := api.NewAPIServer(":8080", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
