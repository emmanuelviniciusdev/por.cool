package rabbitmq

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/porcool/ingestion/internal/config"
)

func TestIngestionMessageUnmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected IngestionMessage
		wantErr  bool
	}{
		{
			name:  "valid message",
			input: `{"successfullyIngestedFirestoreDocsID": "doc123"}`,
			expected: IngestionMessage{
				SuccessfullyIngestedFirestoreDocsID: "doc123",
			},
			wantErr: false,
		},
		{
			name:  "empty document ID",
			input: `{"successfullyIngestedFirestoreDocsID": ""}`,
			expected: IngestionMessage{
				SuccessfullyIngestedFirestoreDocsID: "",
			},
			wantErr: false,
		},
		{
			name:     "invalid JSON",
			input:    `{invalid json}`,
			expected: IngestionMessage{},
			wantErr:  true,
		},
		{
			name:  "missing field",
			input: `{}`,
			expected: IngestionMessage{
				SuccessfullyIngestedFirestoreDocsID: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg IngestionMessage
			err := json.Unmarshal([]byte(tt.input), &msg)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if msg.SuccessfullyIngestedFirestoreDocsID != tt.expected.SuccessfullyIngestedFirestoreDocsID {
				t.Errorf("expected SuccessfullyIngestedFirestoreDocsID=%q, got %q",
					tt.expected.SuccessfullyIngestedFirestoreDocsID,
					msg.SuccessfullyIngestedFirestoreDocsID)
			}
		})
	}
}

func TestIngestionMessageMarshal(t *testing.T) {
	msg := IngestionMessage{
		SuccessfullyIngestedFirestoreDocsID: "test-doc-id",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"successfullyIngestedFirestoreDocsID":"test-doc-id"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNewConsumer_InvalidURI(t *testing.T) {
	cfg := config.RabbitMQConfig{
		URI:       "amqp://invalid:invalid@localhost:9999/",
		QueueName: "test-queue",
	}

	handler := func(ctx context.Context, msg IngestionMessage) error {
		return nil
	}

	_, err := NewConsumer(cfg, handler)
	if err == nil {
		t.Error("expected error for invalid URI, got none")
	}
}

func TestConsumer_IsRunning(t *testing.T) {
	c := &Consumer{
		running: false,
	}

	if c.IsRunning() {
		t.Error("expected IsRunning to return false")
	}

	c.running = true
	if !c.IsRunning() {
		t.Error("expected IsRunning to return true")
	}
}
