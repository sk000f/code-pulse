package models

import (
	"time"
)

type GithubWorkflow struct {
	ID           int       `json:"id" db:"id"`
	Repository   string    `json:"repository" db:"repository"`
	WorkflowName string    `json:"workflow_name" db:"workflow_name"`
	Status       string    `json:"status" db:"status"`
	Duration     int       `json:"duration" db:"duration"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	CompletedAt  time.Time `json:"completed_at" db:"completed_at"`
}

type SonarqubeMetric struct {
	ID           int       `json:"id" db:"id"`
	ProjectKey   string    `json:"project_key" db:"project_key"`
	MetricKey    string    `json:"metric_key" db:"metric_key"`
	Value        string    `json:"value" db:"value"`
	Component    string    `json:"component" db:"component"`
	CollectedAt  time.Time `json:"collected_at" db:"collected_at"`
}

type JiraTicket struct {
	ID          int       `json:"id" db:"id"`
	TicketKey   string    `json:"ticket_key" db:"ticket_key"`
	Summary     string    `json:"summary" db:"summary"`
	Status      string    `json:"status" db:"status"`
	Priority    string    `json:"priority" db:"priority"`
	Assignee    string    `json:"assignee" db:"assignee"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at" db:"resolved_at"`
}