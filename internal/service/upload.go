package service

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/harveytvt/movie-reservation-system/gen/go/api/movie_reservation/v1"
	"github.com/harveytvt/movie-reservation-system/internal/auth"
	"github.com/harveytvt/movie-reservation-system/internal/config"
)

func HandleHealthz(mux *runtime.ServeMux) {
	mux.HandlePath("GET", "/healthz", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.WriteHeader(http.StatusOK)
	})
}

func HandleUpload(mux *runtime.ServeMux) {
	mux.HandlePath("POST", "/upload/{dir}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// auth
		jwtPayload, err := auth.ParseJwtPayload(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if jwtPayload.Username == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// upload dir control based on user role.
		dir := pathParams["dir"]
		switch jwtPayload.Role {
		case movie_reservation.User_ROLE_USER:
			if dir != "avatar" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		case movie_reservation.User_ROLE_ADMIN, movie_reservation.User_ROLE_SUPER_ADMIN:
			if dir != "poster" && dir != "trailer" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		default:
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// multipart form file upload.
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		defer file.Close()

		newFileName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(header.Filename))

		var (
			bucket    = config.Get().Cloudflare.BucketName
			objectKey = fmt.Sprintf("movie_reservation/%s/%s", dir, newFileName)
		)

		err = r2Client.Put(r.Context(), bucket, objectKey, file, func(i *s3.PutObjectInput) {
			i.ContentType = aws.String(header.Header.Get("Content-Type"))
			i.ContentLength = aws.Int64(header.Size)
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("https://%s/%s", config.Get().Cloudflare.BucketDomain, objectKey)))
	})
}
