package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/geraldhoxha/resume-backend/config"
	"github.com/geraldhoxha/resume-backend/directives"
	"github.com/geraldhoxha/resume-backend/graph"
	"github.com/geraldhoxha/resume-backend/middlewares"
	"github.com/geraldhoxha/resume-backend/migration"
	"github.com/gorilla/handlers"
	// "github.com/gorilla/mux"
)

var ALLOWED_ORIGINS = []string{
	"http://localhost:3000",
	"http://192.168.1.167:3000/query",
}

const defaultPort = "8080"

func main() {
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

	// router := mux.NewRouter()
	// router.Use(middlewares.AuthMiddleware)

	c := graph.Config{Resolvers: &graph.Resolver{}}
	c.Directives.Auth = directives.Auth

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middlewares.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)

	CORS := handlers.CORS(
		handlers.AllowedOrigins(ALLOWED_ORIGINS),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(3600),
	)

	log.Fatal(http.ListenAndServe(":"+port, CORS(srv)))
}
