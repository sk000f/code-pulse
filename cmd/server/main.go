package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"code-pulse/internal/config"
	"code-pulse/internal/database"
	"code-pulse/internal/handlers"
	"code-pulse/internal/scheduler"
	"code-pulse/internal/services"
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

	metricsService := services.NewMetricsService(
		db,
		cfg.GithubToken,
		cfg.SonarqubeURL,
		cfg.SonarqubeToken,
		cfg.JiraURL,
		cfg.JiraEmail,
		cfg.JiraToken,
	)

	schedulerService := scheduler.New(metricsService, cfg)
	if err := schedulerService.Start(); err != nil {
		log.Fatal("Failed to start scheduler:", err)
	}
	defer schedulerService.Stop()

	h := handlers.New(db)
	
	http.HandleFunc("/api/metrics/github", h.GetGithubMetrics)
	http.HandleFunc("/api/metrics/sonarqube", h.GetSonarqubeMetrics)
	http.HandleFunc("/api/metrics/jira", h.GetJiraMetrics)
	http.HandleFunc("/api/health", h.Health)

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
			log.Fatal("Server failed:", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}