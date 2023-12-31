package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/squashd/blog-aggregator/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	config := &apiConfig{
		DB: database.New(db),
	}

	port := os.Getenv("PORT")

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1 := chi.NewRouter()
	v1.Get("/readiness", config.ReadinessHandler)
	v1.Get("/err", config.ErrorHandler)

	v1.Post("/users", config.handleCreateUser)
	v1.Get("/users", config.middlewareAuth(config.handleGetUser))

	v1.Get("/feeds", config.handleGetFeeds)
	v1.Post("/feeds", config.middlewareAuth(config.handleCreateFeed))

	v1.Get("/posts", config.middlewareAuth(config.handleGetPosts))

	v1.Get("/feed-follows", config.middlewareAuth(config.handleGetAllFeedsForUser))
	v1.Post("/feed-follows", config.middlewareAuth(config.handleCreateFeedFollow))
	v1.Delete("/feed-follows/{feedFollowId}", config.handleDeleteFeedFollow)
	router.Mount("/v1", v1)

	go startScraping(config.DB, 10, 10*time.Minute)

	server := &http.Server{
		Handler: router,
		Addr:    "localhost:" + port,
	}

	fmt.Printf("Server listening on port %s\n", port)
	_ = server.ListenAndServe()
}
