package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
)

type request struct {
	FileName string `json:"fileName"`
	Message  string `json:"message"`
}

const (
	//S3Region - define the region where bucket exists
	S3Region = "us-east-1"
	//S3Bucket - the bucket where it holds the documents
	S3Bucket = "example"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile).Methods(http.MethodPost)
	fmt.Println("Server should be availabe at http :9090")
	fmt.Println(http.ListenAndServe(":9090", r))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3Region)})
	if err != nil {
		fmt.Printf("error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	// Open the file for use
	file, err := os.Create(req.FileName)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(req.FileName),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		fmt.Printf("error:%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Success")))

}
