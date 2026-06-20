package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestAgentTestingService_Validation tests input validation without API calls.
func TestAgentTestingService_Validation(t *testing.T) {
	client, err := elevenlabs.NewClient(elevenlabs.WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"CreateFolder nil request", func() error { _, err := client.AgentTesting().CreateFolder(ctx, nil); return err }},
		{"CreateFolder empty name", func() error {
			_, err := client.AgentTesting().CreateFolder(ctx, &elevenlabs.CreateTestFolderRequest{})
			return err
		}},
		{"GetFolder empty ID", func() error { _, err := client.AgentTesting().GetFolder(ctx, ""); return err }},
		{"UpdateFolder empty ID", func() error {
			_, err := client.AgentTesting().UpdateFolder(ctx, "", &elevenlabs.UpdateTestFolderRequest{Name: "test"})
			return err
		}},
		{"UpdateFolder nil request", func() error { _, err := client.AgentTesting().UpdateFolder(ctx, "folder123", nil); return err }},
		{"UpdateFolder empty name", func() error {
			_, err := client.AgentTesting().UpdateFolder(ctx, "folder123", &elevenlabs.UpdateTestFolderRequest{})
			return err
		}},
		{"DeleteFolder empty ID", func() error { return client.AgentTesting().DeleteFolder(ctx, "") }},
		{"GetResponseTest empty ID", func() error { _, err := client.AgentTesting().GetResponseTest(ctx, ""); return err }},
		{"GetResponseTestSummaries empty IDs", func() error { _, err := client.AgentTesting().GetResponseTestSummaries(ctx, nil); return err }},
		{"DeleteTest empty ID", func() error { return client.AgentTesting().DeleteTest(ctx, "") }},
		{"BulkMoveTests nil request", func() error { return client.AgentTesting().BulkMoveTests(ctx, nil) }},
		{"BulkMoveTests empty IDs", func() error { return client.AgentTesting().BulkMoveTests(ctx, &elevenlabs.BulkMoveTestsRequest{}) }},
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

// TestAgentTestingService_Integration tests agent testing operations with live API.
func TestAgentTestingService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test listing tests
	t.Run("ListTests", func(t *testing.T) {
		resp, err := client.AgentTesting().ListTests(ctx, &elevenlabs.ListTestsOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Fatalf("ListTests() error = %v", err)
		}
		t.Logf("Found %d tests, HasMore=%v", len(resp.Tests), resp.HasMore)
	})
}
