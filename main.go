package main

import (
	"forim/database"
	"forim/handlers"
	"log"
	"net/http"
)

func main() {
	err := database.InitializeDB("./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	srv := http.Server{
		Addr:    ":8080",
		Handler: routes(),
	}
	log.Println("Listening on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Login)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/post", handlers.GetHome)
	mux.HandleFunc("/post/create", handlers.CreatePost)
	return mux
}
