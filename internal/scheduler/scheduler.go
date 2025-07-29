package scheduler

import (
	"log"
	"time"

	"code-pulse/internal/config"
	"code-pulse/internal/services"
	
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron           *cron.Cron
	metricsService *services.MetricsService
	config         *config.Config
}

func New(metricsService *services.MetricsService, config *config.Config) *Scheduler {
	return &Scheduler{
		cron:           cron.New(),
		metricsService: metricsService,
		config:         config,
	}
}

func (s *Scheduler) Start() error {
	if s.config.CollectionSchedule == "" {
		log.Println("No collection schedule configured, skipping scheduled data collection")
		return nil
	}

	_, err := s.cron.AddFunc(s.config.CollectionSchedule, s.collectMetrics)
	if err != nil {
		return err
	}

	s.cron.Start()
	log.Printf("Scheduler started with schedule: %s", s.config.CollectionSchedule)
	
	log.Println("Running initial metrics collection...")
	go s.collectMetrics()
	
	return nil
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

func (s *Scheduler) collectMetrics() {
	log.Println("Starting scheduled metrics collection...")
	startTime := time.Now()

	s.collectGithubMetrics()
	s.collectSonarqubeMetrics()
	s.collectJiraMetrics()

	duration := time.Since(startTime)
	log.Printf("Metrics collection completed in %v", duration)
}

func (s *Scheduler) collectGithubMetrics() {
	if s.config.GithubToken == "" || s.config.GithubOrg == "" {
		log.Println("GitHub configuration incomplete, skipping GitHub metrics collection")
		return
	}

	log.Println("Collecting GitHub metrics...")
	
	for _, repo := range s.config.GithubRepos {
		for _, workflow := range repo.Workflows {
			log.Printf("Collecting metrics for %s/%s workflow: %s", s.config.GithubOrg, repo.Name, workflow)
			
			if err := s.metricsService.CollectGithubMetricsByWorkflow(s.config.GithubOrg, repo.Name, workflow); err != nil {
				log.Printf("Error collecting GitHub metrics for %s/%s workflow %s: %v", 
					s.config.GithubOrg, repo.Name, workflow, err)
				continue
			}
			
			time.Sleep(1 * time.Second)
		}
	}
	
	log.Println("GitHub metrics collection completed")
}

func (s *Scheduler) collectSonarqubeMetrics() {
	if s.config.SonarqubeURL == "" || s.config.SonarqubeToken == "" {
		log.Println("SonarQube configuration incomplete, skipping SonarQube metrics collection")
		return
	}

	log.Println("SonarQube metrics collection not yet implemented in scheduler")
}

func (s *Scheduler) collectJiraMetrics() {
	if s.config.JiraURL == "" || s.config.JiraEmail == "" || s.config.JiraToken == "" {
		log.Println("Jira configuration incomplete, skipping Jira metrics collection")
		return
	}

	log.Println("Jira metrics collection not yet implemented in scheduler")
}