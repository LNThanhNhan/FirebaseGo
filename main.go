package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const JSON_FILE = "./serviceAccountKey.json"

func main() {
	fmt.Println("Server is running...")
	router := mux.NewRouter()
	LoadFileEnv()
	imageRouter := router.PathPrefix("/img").Subrouter()
	imageRouter.HandleFunc("/create", UploadImage).Methods("POST")   // Request need a image file from form-data
	imageRouter.HandleFunc("/update", UpdateImage).Methods("PUT")    // Request need a image file from form-data and id of image
	imageRouter.HandleFunc("/delete", DeleteImage).Methods("DELETE") // Request need id of image
	http.ListenAndServe("0.0.0.0:8090", router)
}

func LoadFileEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println(err.Error())
	}
}

func HandleError(w http.ResponseWriter, err error, msg string) {
	errMsg := msg + ": " + err.Error()
	jsonData := MakeFailResponse(errMsg)
	ReturnResponse(w, jsonData)
}
