//go:build integration

package agents_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	elevenlabs "github.com/plexusone/elevenlabs-go"
	"github.com/plexusone/elevenlabs-go/agents"
)

func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("ELEVENLABS_API_KEY") == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping integration test")
	}
}

// skipOn401 skips the test if the error is a 401 unauthorized error.
func skipOn401(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "401") {
		t.Skipf("API key does not have access to this endpoint: %v", err)
	}
}

// skipOn403 skips the test if the error is a 403 forbidden error.
func skipOn403(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), "403") {
		t.Skipf("API key does not have permission for this endpoint: %v", err)
	}
}

func TestAgentsList_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Agents().List(ctx, &agents.ListAgentsOptions{PageSize: 10})
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().List: %v", err)
	}

	t.Logf("Successfully listed %d agents", len(resp.Agents))

	for i, agent := range resp.Agents {
		if i >= 3 {
			t.Logf("... and %d more agents", len(resp.Agents)-3)
			break
		}
		t.Logf("Agent %d: %s (ID: %s)", i+1, agent.Name, agent.AgentID)
	}
}

func TestConversationsList_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Agents().ListConversations(ctx, &agents.ListConversationsOptions{PageSize: 10})
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().ListConversations: %v", err)
	}

	t.Logf("Successfully listed %d conversations", len(resp.Conversations))
}

func TestDocumentsList_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Agents().ListDocuments(ctx, &agents.ListDocumentsOptions{PageSize: 10})
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().ListDocuments: %v", err)
	}

	t.Logf("Successfully listed %d documents", len(resp.Documents))
}

func TestBatchCallsList_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Agents().ListBatchCalls(ctx, &agents.ListBatchCallsOptions{PageSize: 10})
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().ListBatchCalls: %v", err)
	}

	t.Logf("Successfully listed %d batch calls", len(resp.BatchCalls))
}

func TestTestsList_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Agents().ListTests(ctx, &agents.ListTestsOptions{PageSize: 10})
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().ListTests: %v", err)
	}

	t.Logf("Successfully listed %d tests", len(resp.Tests))
}

func TestLiveCount_Integration(t *testing.T) {
	skipIfNoAPIKey(t)

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := client.Agents().GetLiveCount(ctx)
	skipOn401(t, err)
	skipOn403(t, err)
	if err != nil {
		t.Fatalf("Agents().GetLiveCount: %v", err)
	}

	t.Logf("Successfully got live count: %d", count)
}
