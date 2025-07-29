package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	email      string
	token      string
	httpClient *http.Client
	baseURL    string
}

type Issue struct {
	Key    string `json:"key"`
	Fields Fields `json:"fields"`
}

type Fields struct {
	Summary   string    `json:"summary"`
	Status    Status    `json:"status"`
	Priority  Priority  `json:"priority"`
	Assignee  *User     `json:"assignee"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Resolved  *time.Time `json:"resolutiondate"`
}

type Status struct {
	Name string `json:"name"`
}

type Priority struct {
	Name string `json:"name"`
}

type User struct {
	DisplayName string `json:"displayName"`
}

type SearchResponse struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

func NewClient(baseURL, email, token string) *Client {
	return &Client{
		email:      email,
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) SearchIssues(jql string) ([]Issue, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=%s", c.baseURL, jql)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Issues, nil
}