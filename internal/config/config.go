package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
	
	GithubToken    string
	GithubOrg      string
	GithubRepos    []RepoConfig
	
	SonarqubeURL   string
	SonarqubeToken string
	JiraURL        string
	JiraEmail      string
	JiraToken      string
	
	CollectionSchedule string
}

type RepoConfig struct {
	Name      string   `json:"name"`
	Workflows []string `json:"workflows"`
}

func Load() *Config {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://localhost/codepulse?sslmode=disable"),
		
		GithubToken:        getEnv("GITHUB_TOKEN", ""),
		GithubOrg:          getEnv("GITHUB_ORG", ""),
		
		SonarqubeURL:       getEnv("SONARQUBE_URL", ""),
		SonarqubeToken:     getEnv("SONARQUBE_TOKEN", ""),
		JiraURL:            getEnv("JIRA_URL", ""),
		JiraEmail:          getEnv("JIRA_EMAIL", ""),
		JiraToken:          getEnv("JIRA_TOKEN", ""),
		
		CollectionSchedule: getEnv("COLLECTION_SCHEDULE", "0 */6 * * *"), // Every 6 hours by default
	}
	
	if reposJSON := getEnv("GITHUB_REPOS", ""); reposJSON != "" {
		if err := json.Unmarshal([]byte(reposJSON), &cfg.GithubRepos); err == nil {
		} else {
			cfg.GithubRepos = []RepoConfig{}
		}
	}
	
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}