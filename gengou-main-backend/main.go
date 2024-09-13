package main

import (
	"fmt"
	"gengou-main-backend/internals/api"
	"gengou-main-backend/internals/auth"
	"gengou-main-backend/internals/database"
	"gengou-main-backend/internals/redis"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	database.DbInit()
	redis.RedisInit()
	api.InitPresigner()
	defer func() {
		database.DbClose()
		redis.RedisClose()
	}()
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.DefaultLogger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(auth.AuthenticateUser)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world the request came")
		_, err := w.Write([]byte("hello world the request came"))
		if err != nil {
			return
		}
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world the request came")
		_, err := w.Write([]byte("hello world the request came"))
		if err != nil {
			return
		}
	})
	r.Mount("/flashcard", api.FlashcardApiRouter())
	r.Mount("/presign", api.PresignerApiRouter())
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
