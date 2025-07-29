package sonarqube

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

type Metric struct {
	Metric    string `json:"metric"`
	Value     string `json:"value"`
	Component string `json:"component"`
}

type Component struct {
	Key     string   `json:"key"`
	Name    string   `json:"name"`
	Metrics []Metric `json:"measures"`
}

type MeasuresResponse struct {
	Component Component `json:"component"`
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) GetProjectMetrics(projectKey string, metricKeys []string) ([]Metric, error) {
	params := url.Values{}
	params.Set("component", projectKey)
	for _, key := range metricKeys {
		params.Add("metricKeys", key)
	}

	url := fmt.Sprintf("%s/api/measures/component?%s", c.baseURL, params.Encode())
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var response MeasuresResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for i := range response.Component.Metrics {
		response.Component.Metrics[i].Component = projectKey
	}

	return response.Component.Metrics, nil
}