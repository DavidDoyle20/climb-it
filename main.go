package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"climb_it/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	DB        *database.Queries
	secretKey string
}

//go:embed static/*
var staticFiles embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configurations. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable is not set")
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Println("SECRET_KEY environment variable is not set")
	}

	apiCfg := apiConfig{}
	apiCfg.secretKey = secretKey

	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		log.Println("DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("libsql", dbURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbURL, err)
			os.Exit(1)
		}
		dbQueries := database.New(db)
		apiCfg.DB = dbQueries
		log.Println("Connected to database!")
		defer db.Close()
	}

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(CORSMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFiles.Open("static/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	v1Router := chi.NewRouter()

	if apiCfg.DB != nil {

		v1Router.Post("/login", apiCfg.handlerUsersLogin)
		v1Router.Post("/logout", apiCfg.handlerUsersLogout)
		v1Router.Post("/refresh-token", apiCfg.handlerRefresh)

		v1Router.Post("/users", apiCfg.handlerUsersCreate)

		v1Router.Post("/habits", apiCfg.handlerHabitsCreate)
		v1Router.Get("/habits", apiCfg.handlerHabitsGet)
		v1Router.Delete("/habits/{habitID}", apiCfg.handlerHabitsDelete)

	}
	v1Router.Get("/healthzv", handlerReadiness)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 30 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
