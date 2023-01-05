package main

import (
	"log"
	"main/database"
	"main/entities"
	"main/server"
	"main/service"
	"main/service/sender"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	entities.NewValidator()
	authenticator := server.NewNoneAuthenticator()
	storage, err := database.NewSqliteStorage(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	smsSender := sender.NewMelipayamakSender()
	senderMap := sender.SenderMap{entities.Melipayamak: smsSender}
	smsService := service.NewSMSService(storage, senderMap)
	apiServer := server.APIServer{
		Addr:    config.Hostname,
		Auth:    authenticator,
		Storage: storage,
		Service: smsService,
	}
	log.Fatal(apiServer.Run())
}
