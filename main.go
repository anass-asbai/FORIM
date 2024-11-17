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
		Addr:    ":8081",
		Handler: routes(),
	}
	log.Println("Listening on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/", handlers.GetHome)
	mux.HandleFunc("/comment", handlers.GetComment)
	mux.HandleFunc("/create", handlers.CreatePost)
	mux.HandleFunc("/like_post", handlers.Like_post)

	mux.HandleFunc("/newcomment",handlers.NewComment)
	return mux
}
