package logging

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/porcool/ingestion/internal/config"
)

// OpenSearchClient handles communication with OpenSearch
type OpenSearchClient struct {
	cfg        config.OpenSearchConfig
	httpClient *http.Client
	available  bool
	mu         sync.RWMutex
}

// LogEntry represents a single log entry to be sent to OpenSearch
type LogEntry struct {
	Timestamp   time.Time              `json:"@timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	Host        string                 `json:"host,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// NewOpenSearchClient creates a new OpenSearch client
func NewOpenSearchClient(cfg config.OpenSearchConfig) *OpenSearchClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		MaxIdleConns:        10,
		IdleConnTimeout:     30 * time.Second,
		DisableCompression:  false,
	}

	client := &OpenSearchClient{
		cfg: cfg,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
		available: false,
	}

	return client
}

// Connect tests the connection to OpenSearch and sets up the index template
func (c *OpenSearchClient) Connect(ctx context.Context) error {
	if !c.cfg.Enabled || c.cfg.URL == "" {
		return nil
	}

	// Test connection with a simple request
	if err := c.ping(ctx); err != nil {
		c.setAvailable(false)
		return fmt.Errorf("failed to connect to OpenSearch: %w", err)
	}

	// Create index template
	if err := c.createIndexTemplate(ctx); err != nil {
		c.setAvailable(false)
		return fmt.Errorf("failed to create index template: %w", err)
	}

	// Create ISM policy for retention (non-blocking, uses separate timeout)
	go func() {
		ismCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := c.createISMPolicy(ismCtx); err != nil {
			fmt.Printf("Warning: failed to create ISM policy (retention may need manual setup): %v\n", err)
		}
	}()

	c.setAvailable(true)
	return nil
}

// ping checks if OpenSearch is reachable
func (c *OpenSearchClient) ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.cfg.URL, nil)
	if err != nil {
		return err
	}

	if c.cfg.Username != "" && c.cfg.Password != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenSearch returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// createIndexTemplate creates an index template for log indices
func (c *OpenSearchClient) createIndexTemplate(ctx context.Context) error {
	templateName := c.cfg.IndexPrefix + "-template"
	templateURL := fmt.Sprintf("%s/_index_template/%s", c.cfg.URL, templateName)

	// OpenSearch index template
	// Retention is handled by ISM (Index State Management) policy created separately
	template := map[string]interface{}{
		"index_patterns": []string{c.cfg.IndexPrefix + "-*"},
		"template": map[string]interface{}{
			"settings": map[string]interface{}{
				"number_of_shards":   1,
				"number_of_replicas": 0,
			},
			"mappings": map[string]interface{}{
				"properties": map[string]interface{}{
					"@timestamp": map[string]interface{}{
						"type": "date",
					},
					"level": map[string]interface{}{
						"type": "keyword",
					},
					"message": map[string]interface{}{
						"type": "text",
					},
					"service": map[string]interface{}{
						"type": "keyword",
					},
					"host": map[string]interface{}{
						"type": "keyword",
					},
					"fields": map[string]interface{}{
						"type":    "object",
						"enabled": true,
					},
				},
			},
		},
		"priority": 100,
	}

	body, err := json.Marshal(template)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", templateURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.Username != "" && c.cfg.Password != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create index template: status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// createISMPolicy creates an Index State Management policy for log retention
func (c *OpenSearchClient) createISMPolicy(ctx context.Context) error {
	policyName := c.cfg.IndexPrefix + "-ism-policy"
	policyURL := fmt.Sprintf("%s/_plugins/_ism/policies/%s", c.cfg.URL, policyName)

	// OpenSearch ISM policy format
	policy := map[string]interface{}{
		"policy": map[string]interface{}{
			"description":   fmt.Sprintf("ISM policy for %s with %d days retention", c.cfg.IndexPrefix, c.cfg.RetentionDays),
			"default_state": "hot",
			"states": []map[string]interface{}{
				{
					"name":    "hot",
					"actions": []map[string]interface{}{},
					"transitions": []map[string]interface{}{
						{
							"state_name": "delete",
							"conditions": map[string]interface{}{
								"min_index_age": fmt.Sprintf("%dd", c.cfg.RetentionDays),
							},
						},
					},
				},
				{
					"name": "delete",
					"actions": []map[string]interface{}{
						{
							"delete": map[string]interface{}{},
						},
					},
					"transitions": []map[string]interface{}{},
				},
			},
			"ism_template": []map[string]interface{}{
				{
					"index_patterns": []string{c.cfg.IndexPrefix + "-*"},
					"priority":       100,
				},
			},
		},
	}

	body, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", policyURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.Username != "" && c.cfg.Password != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle cases where policy already exists:
	// - 409 Conflict: standard "already exists" response
	// - 400 with "matching existing policy templates": policy with same pattern already exists
	if resp.StatusCode >= 400 && resp.StatusCode != 409 {
		respBody, _ := io.ReadAll(resp.Body)
		bodyStr := string(respBody)
		// If the error is about matching existing policy templates, the policy already exists
		if resp.StatusCode == 400 && strings.Contains(bodyStr, "matching existing policy templates") {
			return nil
		}
		return fmt.Errorf("failed to create ISM policy: status %d: %s", resp.StatusCode, bodyStr)
	}

	return nil
}

// IndexLog sends a log entry to OpenSearch
func (c *OpenSearchClient) IndexLog(ctx context.Context, entry LogEntry) error {
	if !c.IsAvailable() {
		return fmt.Errorf("OpenSearch is not available")
	}

	// Use date-based index naming for easier management
	indexName := fmt.Sprintf("%s-%s", c.cfg.IndexPrefix, time.Now().Format("2006.01.02"))
	indexURL := fmt.Sprintf("%s/%s/_doc", c.cfg.URL, indexName)

	body, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", indexURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.Username != "" && c.cfg.Password != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.setAvailable(false)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		// Don't mark as unavailable for document-level errors
		return fmt.Errorf("failed to index log: status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// IndexLogAsync sends a log entry to OpenSearch asynchronously
func (c *OpenSearchClient) IndexLogAsync(entry LogEntry) {
	if !c.IsAvailable() {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// Ignore errors in async mode - they will be logged to fallback
		_ = c.IndexLog(ctx, entry)
	}()
}

// IsAvailable returns whether OpenSearch is currently available
func (c *OpenSearchClient) IsAvailable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cfg.Enabled && c.available
}

// setAvailable sets the availability status
func (c *OpenSearchClient) setAvailable(available bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.available = available
}

// Close closes the OpenSearch client
func (c *OpenSearchClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}

// Reconnect attempts to reconnect to OpenSearch
func (c *OpenSearchClient) Reconnect(ctx context.Context) error {
	return c.Connect(ctx)
}
