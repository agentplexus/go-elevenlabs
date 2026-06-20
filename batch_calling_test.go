package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestBatchCallingService_Validation tests input validation without API calls.
func TestBatchCallingService_Validation(t *testing.T) {
	client, err := elevenlabs.NewClient(elevenlabs.WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"Create nil request", func() error { _, err := client.BatchCalling().Create(ctx, nil); return err }},
		{"Create empty name", func() error {
			_, err := client.BatchCalling().Create(ctx, &elevenlabs.CreateBatchCallRequest{
				AgentID:    "agent123",
				Recipients: []elevenlabs.BatchCallRecipient{{PhoneNumber: "+1234567890"}},
			})
			return err
		}},
		{"Create empty agent ID", func() error {
			_, err := client.BatchCalling().Create(ctx, &elevenlabs.CreateBatchCallRequest{
				Name:       "Test Batch",
				Recipients: []elevenlabs.BatchCallRecipient{{PhoneNumber: "+1234567890"}},
			})
			return err
		}},
		{"Create no recipients", func() error {
			_, err := client.BatchCalling().Create(ctx, &elevenlabs.CreateBatchCallRequest{
				Name:    "Test Batch",
				AgentID: "agent123",
			})
			return err
		}},
		{"Get empty ID", func() error { _, err := client.BatchCalling().Get(ctx, ""); return err }},
		{"Cancel empty ID", func() error { return client.BatchCalling().Cancel(ctx, "") }},
		{"Retry empty ID", func() error { return client.BatchCalling().Retry(ctx, "") }},
		{"Delete empty ID", func() error { return client.BatchCalling().Delete(ctx, "") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestBatchCallingService_Integration tests batch calling operations with live API.
func TestBatchCallingService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test listing batch calls
	t.Run("List", func(t *testing.T) {
		resp, err := client.BatchCalling().List(ctx, &elevenlabs.ListBatchCallsOptions{
			Limit: 5,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		t.Logf("Found %d batch calls", len(resp.BatchCalls))
	})
}
