package logging

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/porcool/ingestion/internal/config"
)

func TestNewOpenSearchClient(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           "http://localhost:9200",
		Username:      "admin",
		Password:      "admin",
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	client := NewOpenSearchClient(cfg)
	if client == nil {
		t.Fatal("NewOpenSearchClient returned nil")
	}
	if client.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
	if client.cfg.URL != cfg.URL {
		t.Errorf("URL = %s, want %s", client.cfg.URL, cfg.URL)
	}
	if client.cfg.RetentionDays != 90 {
		t.Errorf("RetentionDays = %d, want 90", client.cfg.RetentionDays)
	}
}

func TestOpenSearchClient_IsAvailable_Disabled(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: false,
	}

	client := NewOpenSearchClient(cfg)
	if client.IsAvailable() {
		t.Error("IsAvailable should return false when disabled")
	}
}

func TestOpenSearchClient_IsAvailable_NotConnected(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: true,
		URL:     "http://localhost:9200",
	}

	client := NewOpenSearchClient(cfg)
	// Not connected yet
	if client.IsAvailable() {
		t.Error("IsAvailable should return false when not connected")
	}
}

func TestOpenSearchClient_Connect_Disabled(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: false,
	}

	client := NewOpenSearchClient(cfg)
	err := client.Connect(context.Background())
	if err != nil {
		t.Errorf("Connect should not return error when disabled: %v", err)
	}
}

func TestOpenSearchClient_Connect_EmptyURL(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: true,
		URL:     "",
	}

	client := NewOpenSearchClient(cfg)
	err := client.Connect(context.Background())
	if err != nil {
		t.Errorf("Connect should not return error when URL is empty: %v", err)
	}
}

func TestOpenSearchClient_Connect_MockServer(t *testing.T) {
	// Create a mock OpenSearch server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/":
			// Ping endpoint
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name":         "mock-opensearch",
				"cluster_name": "test",
				"version": map[string]interface{}{
					"number": "2.0.0",
				},
			})
		case r.Method == "PUT" && strings.Contains(r.URL.Path, "_index_template"):
			// Index template creation
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"acknowledged": true})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           server.URL,
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	client := NewOpenSearchClient(cfg)
	err := client.Connect(context.Background())
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	if !client.IsAvailable() {
		t.Error("Client should be available after successful connection")
	}
}

func TestOpenSearchClient_IndexLog_MockServer(t *testing.T) {
	var receivedBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"name": "mock"})
		case r.Method == "PUT" && strings.Contains(r.URL.Path, "_index_template"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"acknowledged": true})
		case r.Method == "POST" && strings.Contains(r.URL.Path, "_doc"):
			// Document indexing
			json.NewDecoder(r.Body).Decode(&receivedBody)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"_id":    "test-id",
				"result": "created",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           server.URL,
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	client := NewOpenSearchClient(cfg)
	err := client.Connect(context.Background())
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Service:   "test-service",
		Host:      "localhost",
	}

	err = client.IndexLog(context.Background(), entry)
	if err != nil {
		t.Errorf("IndexLog failed: %v", err)
	}

	// Verify the received body
	if receivedBody["level"] != "INFO" {
		t.Errorf("level = %v, want INFO", receivedBody["level"])
	}
	if receivedBody["message"] != "test message" {
		t.Errorf("message = %v, want test message", receivedBody["message"])
	}
}

func TestOpenSearchClient_IndexLog_NotAvailable(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: true,
		URL:     "http://localhost:9200",
	}

	client := NewOpenSearchClient(cfg)
	// Not connected, so not available

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test",
	}

	err := client.IndexLog(context.Background(), entry)
	if err == nil {
		t.Error("IndexLog should return error when not available")
	}
}

func TestOpenSearchClient_Close(t *testing.T) {
	cfg := config.OpenSearchConfig{
		Enabled: true,
		URL:     "http://localhost:9200",
	}

	client := NewOpenSearchClient(cfg)
	err := client.Close()
	if err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestOpenSearchClient_BasicAuth(t *testing.T) {
	var authHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"name": "mock"})
	}))
	defer server.Close()

	cfg := config.OpenSearchConfig{
		Enabled:  true,
		URL:      server.URL,
		Username: "testuser",
		Password: "testpass",
	}

	client := NewOpenSearchClient(cfg)
	_ = client.ping(context.Background())

	if authHeader == "" {
		t.Error("Authorization header should be set")
	}
	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Error("Authorization should use Basic auth")
	}
}

func TestOpenSearchClient_Reconnect(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch {
		case r.Method == "GET" && r.URL.Path == "/":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"name": "mock"})
		case r.Method == "PUT":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"acknowledged": true})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           server.URL,
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	client := NewOpenSearchClient(cfg)

	// First connect
	err := client.Connect(context.Background())
	if err != nil {
		t.Fatalf("First Connect failed: %v", err)
	}

	initialCallCount := callCount

	// Reconnect
	err = client.Reconnect(context.Background())
	if err != nil {
		t.Errorf("Reconnect failed: %v", err)
	}

	// Should have made additional calls
	if callCount <= initialCallCount {
		t.Error("Reconnect should have made additional HTTP calls")
	}
}

func TestOpenSearchClient_ISMPolicyCreation(t *testing.T) {
	ismPolicyCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"name": "mock"})
		case r.Method == "PUT" && strings.Contains(r.URL.Path, "_index_template"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"acknowledged": true})
		case r.Method == "PUT" && strings.Contains(r.URL.Path, "_plugins/_ism/policies"):
			ismPolicyCalled = true
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"_id": "test-policy"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := config.OpenSearchConfig{
		Enabled:       true,
		URL:           server.URL,
		IndexPrefix:   "test-logs",
		RetentionDays: 90,
	}

	client := NewOpenSearchClient(cfg)
	err := client.Connect(context.Background())

	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	// Wait a bit for the async ISM policy creation
	time.Sleep(100 * time.Millisecond)

	if !ismPolicyCalled {
		t.Error("ISM policy creation should have been called")
	}
}
