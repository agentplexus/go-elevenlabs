package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestConversationsService_Validation tests input validation without API calls.
func TestConversationsService_Validation(t *testing.T) {
	client, err := elevenlabs.NewClient(elevenlabs.WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"Get empty ID", func() error { _, err := client.Conversations().Get(ctx, ""); return err }},
		{"Delete empty ID", func() error { return client.Conversations().Delete(ctx, "") }},
		{"GetAudio empty ID", func() error { _, err := client.Conversations().GetAudio(ctx, ""); return err }},
		{"SubmitFeedback empty ID", func() error { return client.Conversations().SubmitFeedback(ctx, "", "like") }},
		{"SubmitFeedback invalid feedback", func() error { return client.Conversations().SubmitFeedback(ctx, "conv123", "invalid") }},
		{"RunAnalysis empty ID", func() error { return client.Conversations().RunAnalysis(ctx, "") }},
		{"Search empty query", func() error { _, err := client.Conversations().Search(ctx, "", nil); return err }},
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

// TestConversationsService_Integration tests conversation operations with live API.
func TestConversationsService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test listing conversations
	t.Run("List", func(t *testing.T) {
		resp, err := client.Conversations().List(ctx, &elevenlabs.ListConversationsOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		t.Logf("Found %d conversations, HasMore=%v", len(resp.Conversations), resp.HasMore)
	})
}
