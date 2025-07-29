-- GitHub Actions workflow runs
CREATE TABLE IF NOT EXISTS github_workflows (
    id SERIAL PRIMARY KEY,
    repository VARCHAR(255) NOT NULL,
    workflow_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    duration INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NOT NULL,
    UNIQUE(repository, workflow_name, created_at)
);

CREATE INDEX IF NOT EXISTS idx_github_workflows_repository ON github_workflows(repository);
CREATE INDEX IF NOT EXISTS idx_github_workflows_status ON github_workflows(status);
CREATE INDEX IF NOT EXISTS idx_github_workflows_created_at ON github_workflows(created_at);

-- SonarQube metrics
CREATE TABLE IF NOT EXISTS sonarqube_metrics (
    id SERIAL PRIMARY KEY,
    project_key VARCHAR(255) NOT NULL,
    metric_key VARCHAR(255) NOT NULL,
    value TEXT NOT NULL,
    component VARCHAR(255) NOT NULL,
    collected_at TIMESTAMP NOT NULL,
    UNIQUE(project_key, metric_key, component, collected_at)
);

CREATE INDEX IF NOT EXISTS idx_sonarqube_metrics_project_key ON sonarqube_metrics(project_key);
CREATE INDEX IF NOT EXISTS idx_sonarqube_metrics_metric_key ON sonarqube_metrics(metric_key);
CREATE INDEX IF NOT EXISTS idx_sonarqube_metrics_collected_at ON sonarqube_metrics(collected_at);

-- Jira tickets
CREATE TABLE IF NOT EXISTS jira_tickets (
    id SERIAL PRIMARY KEY,
    ticket_key VARCHAR(50) NOT NULL UNIQUE,
    summary TEXT NOT NULL,
    status VARCHAR(100) NOT NULL,
    priority VARCHAR(50) NOT NULL,
    assignee VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    resolved_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_jira_tickets_status ON jira_tickets(status);
CREATE INDEX IF NOT EXISTS idx_jira_tickets_assignee ON jira_tickets(assignee);
CREATE INDEX IF NOT EXISTS idx_jira_tickets_created_at ON jira_tickets(created_at);