package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"vault/database"
	"vault/server"
)

var (
	envPort       = flag.Int("port", 8080, "Port for the HTTP server")
	envMongodbUri = flag.String("mongodb", "", "MongoDB URI")
)

// @title Vault API
// @version 1.0
// @description Vault is a project to have a light-weight, secure and multi-tenant solution for encrypted password storage. It provides a simple Rest API where you can manage your environments and tokens. It uses a MongoDB database as the storage backend.

// @BasePath /
// @securityDefinitions.apikey X-Auth-Token
// @in X-Auth-Token
// @name Authorization
// @description Secure token

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	flag.Parse()

	listenPort := 8080
	if os.Getenv("PORT") != "" {
		listenPort, _ = strconv.Atoi(os.Getenv("PORT"))
		if listenPort == 0 {
			listenPort = *envPort
		}
	} else {
		listenPort = *envPort
	}

	mongodbConnectUri := ""
	if os.Getenv("MONGODB_URI") != "" {
		mongodbConnectUri = os.Getenv("MONGODB_URI")
	} else {
		mongodbConnectUri = *envMongodbUri
	}
	db := database.NewMongoDB(mongodbConnectUri)
	server.InitDatabase(&db)

	mux := server.NewCustomMux(&db)
	server.Load(&mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      &mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error while listening on port: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}

}
