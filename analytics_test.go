package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestAnalyticsService_Integration tests analytics operations with live API.
func TestAnalyticsService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test getting live count
	t.Run("GetLiveCount", func(t *testing.T) {
		count, err := client.Analytics().GetLiveCount(ctx)
		if err != nil {
			t.Fatalf("GetLiveCount() error = %v", err)
		}
		t.Logf("Live conversations: %d", count)
	})
}
