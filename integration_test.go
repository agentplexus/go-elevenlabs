// Integration tests that require ELEVENLABS_API_KEY to be set.
// These tests verify that ogen-generated code correctly handles nullable fields
// returned by the ElevenLabs API.
//
// Run with:
//   ELEVENLABS_API_KEY=your_key go test -v -tags=integration ./...
//
// These tests are placed in the root directory (not internal/api/) so they
// won't be overwritten when ogen regenerates code with --clean.

//go:build integration

package elevenlabs

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("ELEVENLABS_API_KEY") == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping integration test")
	}
}

// skipOn401 skips the test if the error is a 401 unauthorized error.
// Some API keys may not have access to all endpoints.
func skipOn401(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "401") {
		t.Skipf("API key does not have access to this endpoint: %v", err)
	}
}

// TestVoicesListNullHandling tests that the voices list endpoint correctly
// handles nullable fields like manual_verification, settings, and sharing.
// This catches ogen null handling issues like https://github.com/ogen-go/ogen/issues/1358
func TestVoicesListNullHandling(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	voices, err := client.Voices().List(ctx)
	if err != nil {
		t.Fatalf("Voices().List: %v", err)
	}

	t.Logf("Successfully listed %d voices", len(voices))

	// Verify we got some voices back
	if len(voices) == 0 {
		t.Log("Warning: no voices returned, but API call succeeded")
	}

	// Log some details to help debug if issues arise
	for i, v := range voices {
		if i >= 3 {
			t.Logf("... and %d more voices", len(voices)-3)
			break
		}
		t.Logf("Voice %d: %s (ID: %s)", i+1, v.Name, v.VoiceID)
	}
}

// TestModelsListNullHandling tests that the models list endpoint correctly
// handles nullable fields.
func TestModelsListNullHandling(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models, err := client.Models().List(ctx)
	skipOn401(t, err)
	if err != nil {
		t.Fatalf("Models().List: %v", err)
	}

	t.Logf("Successfully listed %d models", len(models))

	if len(models) == 0 {
		t.Log("Warning: no models returned, but API call succeeded")
	}
}

// TestUserGetNullHandling tests that the user info endpoint correctly
// handles nullable fields in the subscription response.
func TestUserGetNullHandling(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := client.User().GetInfo(ctx)
	skipOn401(t, err)
	if err != nil {
		t.Fatalf("User().GetInfo: %v", err)
	}

	t.Logf("Successfully got user info for: %s", user.FirstName)
}

// TestHistoryListNullHandling tests that the history list endpoint correctly
// handles nullable fields.
func TestHistoryListNullHandling(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Just fetch first page with a small limit
	history, err := client.History().List(ctx, &HistoryListOptions{PageSize: 10})
	skipOn401(t, err)
	if err != nil {
		t.Fatalf("History().List: %v", err)
	}

	t.Logf("Successfully listed %d history items", len(history.Items))
}
