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

func UploadImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeBucket(w, bucketCh)

	imageId, path, imageFile := InitializeImageDataFromPostRequest(w, r)

	bucket := <-bucketCh
	StoreImageInBucket(w, bucket, path, imageId, imageFile)

	Data := CreateResponseDataAfterStoreImage(imageId, path)
	jsonData := MakeSuccessResponse(&Data)
	ReturnResponse(w, jsonData)
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeBucket(w, bucketCh)

	imageId, path, imageFile := InitializeImageDataFromPutRequest(w, r)

	bucket := <-bucketCh
	StoreImageInBucket(w, bucket, path, imageId, imageFile)
	Data := CreateResponseDataAfterStoreImage(imageId, path)
	jsonData := MakeSuccessResponse(&Data)
	ReturnResponse(w, jsonData)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	bucketCh := make(chan *storage.BucketHandle)
	go initializeBucket(w, bucketCh)

	id, path := InitializeImageDataFromDeleteRequest(w, r)

	bucket := <-bucketCh
	DeleteImageFromBucket(w, bucket, path, id)
	Data := CreateResponseDataAfterDeleteImage("Image deleted")
	json := MakeSuccessResponse(&Data)
	ReturnResponse(w, json)
}

func InitializeImageDataFromPostRequest(w http.ResponseWriter, r *http.Request) (string, string, io.Reader) {
	imageId := uuid.New().String()
	path := os.Getenv("IMG_PATH")
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		go HandleError(w, err, "Error reading image file from request")
	}
	defer imageFile.Close()
	return imageId, path, imageFile
}

func InitializeImageDataFromPutRequest(w http.ResponseWriter, r *http.Request) (string, string, io.Reader) {
	imageId := r.FormValue("id")
	path := os.Getenv("IMG_PATH")
	imageFile, _, err := r.FormFile("image")
	if err != nil {
		go HandleError(w, err, "Error reading image file from request")
	}
	defer imageFile.Close()
	return imageId, path, imageFile
}

func InitializeImageDataFromDeleteRequest(w http.ResponseWriter, r *http.Request) (string, string) {
	imageId := r.URL.Query().Get("id")
	path := os.Getenv("IMG_PATH")
	return imageId, path
}

func CreateResponseDataAfterStoreImage(imageId string, path string) interface{} {
	imgPath := url.PathEscape(path + imageId)
	url := os.Getenv("FIREBASE_DOMAIN") + os.Getenv("BUCKET_NAME") + "/o/" + imgPath + "?alt=media&token=" + imageId
	Data := struct {
		Url string
		Id  string
	}{
		Url: url,
		Id:  imageId,
	}
	return Data
}

func CreateResponseDataAfterDeleteImage(msg string) interface{} {
	Data := struct {
		Msg string
	}{
		Msg: "Image deleted",
	}
	return Data
}

func initializeBucket(w http.ResponseWriter, buCh chan *storage.BucketHandle) {
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

func StoreImageInBucket(w http.ResponseWriter, bucket *storage.BucketHandle, path string, id string, imageFile io.Reader) {
	objectHandle := bucket.Object(path + id)
	writer := objectHandle.NewWriter(context.Background())
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id}
	defer writer.Close()
	if _, err := io.Copy(writer, imageFile); err != nil {
		HandleError(w, err, "Error uploading image")
	}
}

func DeleteImageFromBucket(w http.ResponseWriter, bucket *storage.BucketHandle, path string, id string) {
	objectHandle := bucket.Object(path + id)
	err := objectHandle.Delete(context.Background())
	if err != nil {
		go HandleError(w, err, "Error deleting image")
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
