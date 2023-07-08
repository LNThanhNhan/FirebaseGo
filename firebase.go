package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

func initializeApp(w http.ResponseWriter, buCh chan *storage.BucketHandle) {
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		go HandleError(w, err, "Error initializing app")
		buCh <- nil
	}
	client, err := app.Storage(context.TODO())
	if err != nil {
		go HandleError(w, err, "Error initializing client")
		buCh <- nil
	}
	bucket, err := client.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		go HandleError(w, err, "Error initializing bucket")
		buCh <- nil
	}
	buCh <- bucket
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeApp(w, bucketCh)
	//Get image from request
	path := os.Getenv("IMG_PATH")
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		go HandleError(w, err, "Error reading image file from request")
		return
	}
	defer imageFile.Close()
	id := uuid.New().String()
	bucket := <-bucketCh
	if bucket == nil {
		return
	}
	objectHandle := bucket.Object(path + id)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}
	defer writer.Close()
	//Make url from image
	img_path := url.PathEscape(path + id)
	url := os.Getenv("FIREBASE_DOMAIN") + os.Getenv("BUCKET_NAME") + "/o/" + img_path + "?alt=media&token=" + id
	Data := struct {
		Url string
		Id  string
	}{
		Url: url,
		Id:  id,
	}
	jsonData := MakeSuccessResponse(&Data)
	ReturnResponse(w, jsonData)
	if _, err = io.Copy(writer, imageFile); err != nil {
		HandleError(w, err, "Error uploading image")
	}
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeApp(w, bucketCh)
	//Get image from request
	id := r.FormValue("id")
	path := os.Getenv("IMG_PATH")
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		go HandleError(w, err, "Error reading image file from request")
		return
	}
	defer imageFile.Close()
	bucket := <-bucketCh
	if bucket == nil {
		return
	}
	objectHandle := bucket.Object(path + id)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}
	defer writer.Close()
	//Make url from image
	img_path := url.PathEscape(path + id)
	url := os.Getenv("FIREBASE_DOMAIN") + os.Getenv("BUCKET_NAME") + "/o/" + img_path + "?alt=media&token=" + id
	Data := struct {
		Url string
		Id  string
	}{
		Url: url,
		Id:  id,
	}
	jsonData := MakeSuccessResponse(&Data)
	ReturnResponse(w, jsonData)
	if _, err = io.Copy(writer, imageFile); err != nil {
		HandleError(w, err, "Error uploading image")
	}
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeApp(w, bucketCh)
	id := r.URL.Query().Get("id")
	path := os.Getenv("IMG_PATH")
	bucket := <-bucketCh
	objectHandle := bucket.Object(path + id)
	err := objectHandle.Delete(context.Background())
	if err != nil {
		go HandleError(w, err, "Error deleting image")
	} else {
		Data := struct {
			Msg string
		}{
			Msg: "Image deleted",
		}
		json := MakeSuccessResponse(&Data)
		ReturnResponse(w, json)
	}
}

func HandleError(w http.ResponseWriter, err error, msg string) {
	errMsg := msg + ": " + err.Error()
	jsonData := MakeFailResponse(errMsg)
	ReturnResponse(w, jsonData)
}

func ReturnResponse(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
