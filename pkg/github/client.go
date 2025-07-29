package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

type Workflow struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type WorkflowsResponse struct {
	TotalCount int        `json:"total_count"`
	Workflows  []Workflow `json:"workflows"`
}

type WorkflowRun struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	DisplayTitle string    `json:"display_title"`
	Status       string    `json:"status"`
	Conclusion   string    `json:"conclusion"`
	WorkflowID   int       `json:"workflow_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	RunStartedAt time.Time `json:"run_started_at"`
}

type WorkflowRunsResponse struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []WorkflowRun `json:"workflow_runs"`
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    "https://api.github.com",
	}
}

func (c *Client) GetWorkflows(owner, repo string) ([]Workflow, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows", c.baseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response WorkflowsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Workflows, nil
}

func (c *Client) GetWorkflowRunsByName(owner, repo, workflowName string) ([]WorkflowRun, error) {
	workflows, err := c.GetWorkflows(owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows: %w", err)
	}

	var workflowID int
	for _, workflow := range workflows {
		if workflow.Name == workflowName {
			workflowID = workflow.ID
			break
		}
	}

	if workflowID == 0 {
		return nil, fmt.Errorf("workflow '%s' not found in repository %s/%s", workflowName, owner, repo)
	}

	return c.GetWorkflowRuns(owner, repo, workflowID)
}

func (c *Client) GetWorkflowRuns(owner, repo string, workflowID int) ([]WorkflowRun, error) {
	params := url.Values{}
	params.Set("per_page", "100")

	url := fmt.Sprintf("%s/repos/%s/%s/actions/workflows/%d/runs?%s",
		c.baseURL, owner, repo, workflowID, params.Encode())

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response WorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.WorkflowRuns, nil
}
