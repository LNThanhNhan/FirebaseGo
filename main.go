package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	router := mux.NewRouter()
	LoadFileEnv()
	imageRouter := router.PathPrefix("/img").Subrouter()
	imageRouter.HandleFunc("/create", UploadImage).Methods("POST")   // Request need a image file from form-data
	imageRouter.HandleFunc("/update", UpdateImage).Methods("PUT")    // Request need a image file from form-data and id of image
	imageRouter.HandleFunc("/delete", DeleteImage).Methods("DELETE") // Request need id of image
	http.ListenAndServe("127.0.0.1:8080", router)
}

func LoadFileEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Print("Error loading .env file")
	}
}
