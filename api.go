package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("Initializing API")

	// initialize and validate database connection
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("Unable to create database connection: ", err.Error())
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Unable to connect to database: ", err.Error())
	}

	router := mux.NewRouter()

	// Register endpoints
	router.HandleFunc("/pages", searchPagesHandler(client)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/pages/{page-id}/content", viewPageContentHandler(client)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/count-pages", countPagesHandler(client)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/pages-stream", pagesStreamHandler(client)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/crawled-urls", getCrawledUrls(client)).Methods(http.MethodGet, http.MethodOptions)

	router.Use(mux.CORSMethodMiddleware(router))

	log.Println("API will listen on: http://0.0.0.0:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", router); err != nil {
		log.Println("fatal error while trying to listen: " + err.Error())
	}
}