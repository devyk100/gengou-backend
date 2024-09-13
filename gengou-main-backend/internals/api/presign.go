package api

import (
	"context"
	"encoding/json"
	"fmt"
	"gengou-main-backend/internals/redis"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

var S3Client *s3.Client
var bucketName = "gengou-bucket"

func InitPresigner() {

	var accountId = "817c5692c2feefe1b588ca939ec4a599"
	var accessKeyId = "1cc21aa5be9175c7a0ba4c122ec73ac1"
	var accessKeySecret = "48ba1aa07ec55f2998bdf0d7e40e2e1c22bad2553bb2ec1f16afaf4829490b8a"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal(err)
	}

	S3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	//listObjectsOutput, err := S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
	//	Bucket: &bucketName,
	//})
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, object := range listObjectsOutput.Contents {
	//	obj, _ := json.MarshalIndent(object, "", "\t")
	//	fmt.Println(string(obj))
	//}
	//
	//listBucketsOutput, err := S3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, object := range listBucketsOutput.Buckets {
	//	obj, _ := json.MarshalIndent(object, "", "\t")
	//	fmt.Println(string(obj))
	//}

}

type PresignerResponse struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
}

type PresignerRequest struct {
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
	Type        string `json:"type"`
}

type PresignerPreviewRequest struct {
	FileName string `json:"fileName"`
}

type PresignerPreviewResponse struct {
	FileUrl string `json:"fileUrl"`
}

func PresignerApiRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/put-image", func(w http.ResponseWriter, r *http.Request) {
		fileName := uuid.New()
		var body PresignerRequest
		presignClient := s3.NewPresignClient(S3Client)
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			panic(err.Error())
			return
		}
		if body.Size > 1024*1024*15 {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("Fucked up"))
		}
		presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:        aws.String(bucketName),
			Key:           aws.String(fileName.String()),
			ContentType:   aws.String(body.ContentType),
			ContentLength: aws.Int64(body.Size),
		}, s3.WithPresignExpires(15*time.Minute))
		if err != nil {
			panic("Couldn't get presigned URL for PutObject")
		}
		//presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		//	Bucket: nil,
		//	Key:    nil,
		//})
		json, err := json.Marshal(PresignerResponse{
			Url:      presignResult.URL,
			FileName: fileName.String(),
		})
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(json)
		if err != nil {
			panic(err.Error())
			return
		}
		fmt.Printf("Presigned URL For object: %s\n", presignResult.URL, presignResult.Method)
		return
	})
	router.Post("/get-image", func(w http.ResponseWriter, r *http.Request) {
		var body PresignerPreviewRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		val, err := redis.Instance.Get(body.FileName)
		if err != nil {
			return
		}
		if val != "" {
			marshal, err := json.Marshal(PresignerPreviewResponse{FileUrl: val})
			if err != nil {
				panic(err.Error())
			}
			w.WriteHeader(http.StatusOK)
			_, err = w.Write(marshal)
		}
		presignClient := s3.NewPresignClient(S3Client)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(body.FileName, "Is the filename")
		presignResult, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(body.FileName),
		}, s3.WithPresignExpires(2*time.Hour*24))
		/**
		WE MUST STORE JUST THE KEYS OF THE FILES INSIDE OF THE DB.
		THIS IS THE SAME AMOUNT OF TIME WE WILL CACHE THIS IMAGE URL IN REDIS.
		*/
		if err != nil {
			panic(err.Error())
		}
		err = redis.Instance.Set(body.FileName, presignResult.URL, time.Hour*2*24)
		if err != nil {
			panic(err.Error())
			return
		}
		marshal, err := json.Marshal(PresignerPreviewResponse{FileUrl: presignResult.URL})
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(marshal)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		//presignResult.URL
	})
	return router
}
