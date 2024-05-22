package main

import (
	"log"
	"net/http"
	"os"

	// "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/geraldhoxha/resume-backend/config"
	"github.com/geraldhoxha/resume-backend/directives"
	"github.com/geraldhoxha/resume-backend/graph"

	"github.com/geraldhoxha/resume-backend/middlewares"
	"github.com/geraldhoxha/resume-backend/migration"
	"github.com/geraldhoxha/resume-backend/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var ALLOWED_ORIGINS = []string{
	"http://localhost:3000",
	"http://localhost:8080",
	"http://localhosr:3000/refresh",
	"http://localhost:8080/refresh",
	"http://192.168.1.167:3000",
}

const defaultPort = "8080"

func main() {
	// AutoMigration
	// AutoMigration
	migration.MigrateTable()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// defer database
	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Create a new router
	router := mux.NewRouter()

	// Apply Auth Middleware to all routes
	router.Use(middlewares.AuthMiddleware)

	// GraphQL server setup
	c := graph.Config{Resolvers: &graph.Resolver{}}
	c.Directives.Auth = directives.Auth
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	// Define routes
	router.Handle("/", playground.Handler("GraphQL playground", "/query")).Methods("GET", "POST", "DELETE")
	router.Handle("/query", srv).Methods("POST")
	router.HandleFunc("/refresh", service.RefreshToken).Methods("POST")

	// CORS configuration
	CORS := handlers.CORS(
		handlers.AllowedOrigins(ALLOWED_ORIGINS),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(3600),
	)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	// Start the server
	log.Fatal(http.ListenAndServe(":"+port, CORS(router)))
	
}
