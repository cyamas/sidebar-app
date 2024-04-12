package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	router := chi.NewRouter()
	fileServer := http.FileServer(http.Dir("static"))
	router.Use(middleware.Logger)
	router.Get("/", home)
	router.Post("/signin", signinUser)
	router.Get("/ws", handleWS)
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	http.ListenAndServe(":6699", router)
}
