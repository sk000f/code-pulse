package config

import (
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	
	GithubToken    string
	SonarqubeURL   string
	SonarqubeToken string
	JiraURL        string
	JiraEmail      string
	JiraToken      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://localhost/codepulse?sslmode=disable"),
		
		GithubToken:    getEnv("GITHUB_TOKEN", ""),
		SonarqubeURL:   getEnv("SONARQUBE_URL", ""),
		SonarqubeToken: getEnv("SONARQUBE_TOKEN", ""),
		JiraURL:        getEnv("JIRA_URL", ""),
		JiraEmail:      getEnv("JIRA_EMAIL", ""),
		JiraToken:      getEnv("JIRA_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}