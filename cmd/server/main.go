package main

import (
	"log"
	"net/http"

	"code-pulse/internal/config"
	"code-pulse/internal/database"
	"code-pulse/internal/handlers"
)

func main() {
	cfg := config.Load()
	
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	h := handlers.New(db)
	
	http.HandleFunc("/api/metrics/github", h.GetGithubMetrics)
	http.HandleFunc("/api/metrics/sonarqube", h.GetSonarqubeMetrics)
	http.HandleFunc("/api/metrics/jira", h.GetJiraMetrics)
	http.HandleFunc("/api/health", h.Health)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}