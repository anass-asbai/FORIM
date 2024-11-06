package main

import (
	"log"
	"net/http"

	"forim/database"
	"forim/handlers"
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
	mux.HandleFunc("/comment", handlers.GetComment)
	mux.HandleFunc("/post/create", handlers.CreatePost)
//	mux.HandleFunc("POST /like_post", handlers.Like_post)
	return mux
}
