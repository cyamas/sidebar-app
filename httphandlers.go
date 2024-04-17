package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := "OK"
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func signinUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form from post request", err)
	}
	username := r.Form["username"][0]
	log.Println(username)
	var newUser *User
	allUsernames := hub.getAllUsernames()
	_, ok := allUsernames[username]
	if !ok {
		newUser = hub.createUser(username)
	}

	tmpl, err := template.ParseFiles("templates/create-room-btn.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("error parsing templates/create-room-btn.html")
		return
	}
	if err := tmpl.Execute(w, newUser.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
