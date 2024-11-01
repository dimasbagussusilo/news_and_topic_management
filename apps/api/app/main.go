package main

import (
	"database/sql"
	"fmt"
	"github.com/bxcodec/go-clean-arch/topic"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/bxcodec/go-clean-arch/app/docs"
	postgresRepo "github.com/bxcodec/go-clean-arch/internal/repository/postgres" // Update the repository package
	"github.com/bxcodec/go-clean-arch/internal/rest"
	"github.com/bxcodec/go-clean-arch/internal/rest/middleware"
	"github.com/bxcodec/go-clean-arch/news"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	defaultTimeout = 30 * time.Second
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
func main() {
	// Prepare database connection
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	connection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	dbConn, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("Failed to open connection to database:", err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal("Error closing the DB connection:", err)
		}
	}()

	// Prepare context timeout
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("Failed to parse timeout, using default timeout")
		timeout = int(defaultTimeout.Seconds())
	}
	timeoutContext := time.Duration(timeout) * time.Second

	// Initialize repositories and services
	newsRepo := postgresRepo.NewNewsRepository(dbConn)
	authorRepo := postgresRepo.NewAuthorRepository(dbConn)
	topicRepo := postgresRepo.NewTopicRepository(dbConn)
	newsTopicRepo := postgresRepo.NewNewsTopicRepository(dbConn)

	ns := news.NewService(newsRepo, authorRepo, topicRepo, newsTopicRepo)
	ts := topic.NewService(topicRepo)

	// Initialize handlers with standard http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprint(w, `{"status": "ok"}`)
		if err != nil {
			return
		}
	})

	// Swagger endpoint
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	rest.NewNewsHandler(mux, ns)
	rest.NewTopicHandler(mux, ts)

	// Middleware setup
	handlerWithMiddleware := middleware.CORS(mux)
	timeoutMiddleware := middleware.SetRequestContextWithTimeout(timeoutContext)
	handlerWithTimeout := timeoutMiddleware(handlerWithMiddleware)

	// Start server
	server := &http.Server{
		Addr:         address,
		Handler:      handlerWithTimeout,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", address)
	log.Fatal(server.ListenAndServe())
}
