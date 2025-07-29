package services

import (
	"database/sql"
	"fmt"
	"time"

	"code-pulse/internal/models"
	"code-pulse/pkg/github"
	"code-pulse/pkg/jira"
	"code-pulse/pkg/sonarqube"
)

type MetricsService struct {
	db             *sql.DB
	githubClient   *github.Client
	sonarClient    *sonarqube.Client
	jiraClient     *jira.Client
}

func NewMetricsService(db *sql.DB, githubToken, sonarURL, sonarToken, jiraURL, jiraEmail, jiraToken string) *MetricsService {
	return &MetricsService{
		db:             db,
		githubClient:   github.NewClient(githubToken),
		sonarClient:    sonarqube.NewClient(sonarURL, sonarToken),
		jiraClient:     jira.NewClient(jiraURL, jiraEmail, jiraToken),
	}
}

func (s *MetricsService) CollectGithubMetrics(owner, repo string) error {
	runs, err := s.githubClient.GetWorkflowRuns(owner, repo, 0)
	if err != nil {
		return fmt.Errorf("failed to get workflow runs: %w", err)
	}

	for _, run := range runs {
		if err := s.saveGithubWorkflowRun(owner, repo, run); err != nil {
			return fmt.Errorf("failed to save workflow run: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) CollectGithubMetricsByWorkflow(owner, repo, workflowName string) error {
	runs, err := s.githubClient.GetWorkflowRunsByName(owner, repo, workflowName)
	if err != nil {
		return fmt.Errorf("failed to get workflow runs for %s: %w", workflowName, err)
	}

	for _, run := range runs {
		if err := s.saveGithubWorkflowRun(owner, repo, run); err != nil {
			return fmt.Errorf("failed to save workflow run: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) saveGithubWorkflowRun(owner, repo string, run github.WorkflowRun) error {
	var duration int
	if !run.RunStartedAt.IsZero() && !run.UpdatedAt.IsZero() {
		duration = int(run.UpdatedAt.Sub(run.RunStartedAt).Seconds())
	} else {
		duration = int(run.UpdatedAt.Sub(run.CreatedAt).Seconds())
	}

	status := run.Status
	if run.Conclusion != "" {
		status = run.Conclusion
	}

	workflow := &models.GithubWorkflow{
		Repository:   fmt.Sprintf("%s/%s", owner, repo),
		WorkflowName: run.Name,
		Status:       status,
		Duration:     duration,
		CreatedAt:    run.CreatedAt,
		CompletedAt:  run.UpdatedAt,
	}

	return s.saveGithubWorkflow(workflow)
}

func (s *MetricsService) CollectSonarqubeMetrics(projectKey string) error {
	metricKeys := []string{
		"ncloc", "coverage", "duplicated_lines_density",
		"bugs", "vulnerabilities", "code_smells",
		"reliability_rating", "security_rating", "sqale_rating",
	}

	metrics, err := s.sonarClient.GetProjectMetrics(projectKey, metricKeys)
	if err != nil {
		return fmt.Errorf("failed to get sonarqube metrics: %w", err)
	}

	for _, metric := range metrics {
		sonarMetric := &models.SonarqubeMetric{
			ProjectKey:  projectKey,
			MetricKey:   metric.Metric,
			Value:       metric.Value,
			Component:   metric.Component,
			CollectedAt: time.Now(),
		}

		if err := s.saveSonarqubeMetric(sonarMetric); err != nil {
			return fmt.Errorf("failed to save sonarqube metric: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) CollectJiraMetrics(jql string) error {
	issues, err := s.jiraClient.SearchIssues(jql)
	if err != nil {
		return fmt.Errorf("failed to search jira issues: %w", err)
	}

	for _, issue := range issues {
		assignee := ""
		if issue.Fields.Assignee != nil {
			assignee = issue.Fields.Assignee.DisplayName
		}

		ticket := &models.JiraTicket{
			TicketKey:  issue.Key,
			Summary:    issue.Fields.Summary,
			Status:     issue.Fields.Status.Name,
			Priority:   issue.Fields.Priority.Name,
			Assignee:   assignee,
			CreatedAt:  issue.Fields.Created,
			UpdatedAt:  issue.Fields.Updated,
			ResolvedAt: issue.Fields.Resolved,
		}

		if err := s.saveJiraTicket(ticket); err != nil {
			return fmt.Errorf("failed to save jira ticket: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) saveGithubWorkflow(workflow *models.GithubWorkflow) error {
	query := `
		INSERT INTO github_workflows (repository, workflow_name, status, duration, created_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (repository, workflow_name, created_at) DO NOTHING`
	
	_, err := s.db.Exec(query, workflow.Repository, workflow.WorkflowName, workflow.Status,
		workflow.Duration, workflow.CreatedAt, workflow.CompletedAt)
	return err
}

func (s *MetricsService) saveSonarqubeMetric(metric *models.SonarqubeMetric) error {
	query := `
		INSERT INTO sonarqube_metrics (project_key, metric_key, value, component, collected_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (project_key, metric_key, component, collected_at) DO NOTHING`
	
	_, err := s.db.Exec(query, metric.ProjectKey, metric.MetricKey, metric.Value,
		metric.Component, metric.CollectedAt)
	return err
}

func (s *MetricsService) saveJiraTicket(ticket *models.JiraTicket) error {
	query := `
		INSERT INTO jira_tickets (ticket_key, summary, status, priority, assignee, created_at, updated_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (ticket_key) DO UPDATE SET
			summary = $2, status = $3, priority = $4, assignee = $5, updated_at = $7, resolved_at = $8`
	
	_, err := s.db.Exec(query, ticket.TicketKey, ticket.Summary, ticket.Status,
		ticket.Priority, ticket.Assignee, ticket.CreatedAt, ticket.UpdatedAt, ticket.ResolvedAt)
	return err
}