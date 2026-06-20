package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestKnowledgeBaseService_Validation tests input validation without API calls.
func TestKnowledgeBaseService_Validation(t *testing.T) {
	client, err := elevenlabs.NewClient(elevenlabs.WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"CreateFileDocument nil request", func() error { _, err := client.KnowledgeBase().CreateFileDocument(ctx, nil); return err }},
		{"CreateTextDocument nil request", func() error { _, err := client.KnowledgeBase().CreateTextDocument(ctx, nil); return err }},
		{"CreateURLDocument nil request", func() error { _, err := client.KnowledgeBase().CreateURLDocument(ctx, nil); return err }},
		{"Get empty ID", func() error { _, err := client.KnowledgeBase().Get(ctx, ""); return err }},
		{"Delete empty ID", func() error { return client.KnowledgeBase().Delete(ctx, "") }},
		{"GetContent empty ID", func() error { _, err := client.KnowledgeBase().GetContent(ctx, ""); return err }},
		{"GetChunks empty ID", func() error { _, err := client.KnowledgeBase().GetChunks(ctx, ""); return err }},
		{"GetChunk empty doc ID", func() error { _, err := client.KnowledgeBase().GetChunk(ctx, "", "chunk123"); return err }},
		{"GetChunk empty chunk ID", func() error { _, err := client.KnowledgeBase().GetChunk(ctx, "doc123", ""); return err }},
		{"GetDependentAgents empty ID", func() error { _, err := client.KnowledgeBase().GetDependentAgents(ctx, ""); return err }},
		{"GetRAGIndexes empty ID", func() error { _, err := client.KnowledgeBase().GetRAGIndexes(ctx, ""); return err }},
		{"DeleteRAGIndex empty doc ID", func() error { return client.KnowledgeBase().DeleteRAGIndex(ctx, "", "idx123") }},
		{"DeleteRAGIndex empty idx ID", func() error { return client.KnowledgeBase().DeleteRAGIndex(ctx, "doc123", "") }},
		{"CreateFolder empty name", func() error { _, err := client.KnowledgeBase().CreateFolder(ctx, "", ""); return err }},
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

// TestKnowledgeBaseService_Integration tests knowledge base operations with live API.
func TestKnowledgeBaseService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test listing documents
	t.Run("List", func(t *testing.T) {
		resp, err := client.KnowledgeBase().List(ctx, &elevenlabs.ListDocumentsOptions{
			PageSize: 5,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		t.Logf("Found %d documents, HasMore=%v", len(resp.Documents), resp.HasMore)
	})
}
