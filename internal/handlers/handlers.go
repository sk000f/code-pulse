package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"code-pulse/internal/models"
)

type Handlers struct {
	db *sql.DB
}

func New(db *sql.DB) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handlers) GetGithubMetrics(w http.ResponseWriter, r *http.Request) {
	repository := r.URL.Query().Get("repository")
	daysStr := r.URL.Query().Get("days")
	
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	query := `
		SELECT repository, workflow_name, status, duration, created_at, completed_at
		FROM github_workflows
		WHERE ($1 = '' OR repository = $1)
		AND created_at >= $2
		ORDER BY created_at DESC`

	since := time.Now().AddDate(0, 0, -days)
	rows, err := h.db.Query(query, repository, since)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workflows []models.GithubWorkflow
	for rows.Next() {
		var workflow models.GithubWorkflow
		err := rows.Scan(&workflow.Repository, &workflow.WorkflowName, &workflow.Status,
			&workflow.Duration, &workflow.CreatedAt, &workflow.CompletedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Scan error: %v", err), http.StatusInternalServerError)
			return
		}
		workflows = append(workflows, workflow)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

func (h *Handlers) GetSonarqubeMetrics(w http.ResponseWriter, r *http.Request) {
	projectKey := r.URL.Query().Get("project_key")
	metricKey := r.URL.Query().Get("metric_key")
	daysStr := r.URL.Query().Get("days")
	
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	query := `
		SELECT project_key, metric_key, value, component, collected_at
		FROM sonarqube_metrics
		WHERE ($1 = '' OR project_key = $1)
		AND ($2 = '' OR metric_key = $2)
		AND collected_at >= $3
		ORDER BY collected_at DESC`

	since := time.Now().AddDate(0, 0, -days)
	rows, err := h.db.Query(query, projectKey, metricKey, since)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var metrics []models.SonarqubeMetric
	for rows.Next() {
		var metric models.SonarqubeMetric
		err := rows.Scan(&metric.ProjectKey, &metric.MetricKey, &metric.Value,
			&metric.Component, &metric.CollectedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Scan error: %v", err), http.StatusInternalServerError)
			return
		}
		metrics = append(metrics, metric)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func (h *Handlers) GetJiraMetrics(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	assignee := r.URL.Query().Get("assignee")
	daysStr := r.URL.Query().Get("days")
	
	days := 30
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil {
			days = d
		}
	}

	query := `
		SELECT ticket_key, summary, status, priority, assignee, created_at, updated_at, resolved_at
		FROM jira_tickets
		WHERE ($1 = '' OR status = $1)
		AND ($2 = '' OR assignee = $2)
		AND created_at >= $3
		ORDER BY created_at DESC`

	since := time.Now().AddDate(0, 0, -days)
	rows, err := h.db.Query(query, status, assignee, since)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tickets []models.JiraTicket
	for rows.Next() {
		var ticket models.JiraTicket
		err := rows.Scan(&ticket.TicketKey, &ticket.Summary, &ticket.Status,
			&ticket.Priority, &ticket.Assignee, &ticket.CreatedAt, &ticket.UpdatedAt, &ticket.ResolvedAt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Scan error: %v", err), http.StatusInternalServerError)
			return
		}
		tickets = append(tickets, ticket)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tickets)
}