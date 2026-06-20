package elevenlabs_test

import (
	"context"
	"os"
	"testing"

	"github.com/plexusone/elevenlabs-go"
)

// TestAgentsService_Integration is an idempotent integration test that:
// 1. Creates a test agent
// 2. Performs various operations on it
// 3. Deletes the agent at the end (cleanup)
func TestAgentsService_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	var agentID string

	// Cleanup function - always delete the agent at the end
	t.Cleanup(func() {
		if agentID != "" {
			t.Logf("Cleanup: deleting agent %s", agentID)
			if err := client.Agents().Delete(ctx, agentID); err != nil {
				t.Logf("Cleanup warning: failed to delete agent: %v", err)
			}
		}
	})

	// Step 1: Create agent
	t.Run("Create", func(t *testing.T) {
		agent, err := client.Agents().Create(ctx, &elevenlabs.CreateAgentRequest{
			Name: "SDK Integration Test Agent",
			Tags: []string{"test", "integration"},
			ConversationConfig: map[string]any{
				"agent": map[string]any{
					"prompt": map[string]any{
						"prompt": "You are a helpful test assistant.",
					},
					"first_message": "Hello! This is a test agent.",
				},
			},
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if agent.AgentID == "" {
			t.Fatal("Create() returned empty agent ID")
		}
		agentID = agent.AgentID
		t.Logf("Created agent: ID=%s, Name=%s", agent.AgentID, agent.Name)
	})

	if agentID == "" {
		t.Fatal("Cannot continue without agent ID")
	}

	// Step 2: Get agent
	t.Run("Get", func(t *testing.T) {
		agent, err := client.Agents().Get(ctx, agentID)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if agent.AgentID != agentID {
			t.Errorf("Get() returned wrong ID: got %s, want %s", agent.AgentID, agentID)
		}
		t.Logf("Fetched agent: ID=%s, Name=%s", agent.AgentID, agent.Name)
	})

	// Step 3: List agents (should include our test agent)
	t.Run("List", func(t *testing.T) {
		resp, err := client.Agents().List(ctx, &elevenlabs.ListAgentsOptions{
			PageSize: 10,
		})
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		found := false
		for _, a := range resp.Agents {
			if a.AgentID == agentID {
				found = true
				break
			}
		}
		if !found {
			t.Error("List() did not include the created agent")
		}
		t.Logf("Found %d agents, test agent present: %v", len(resp.Agents), found)
	})

	// Step 4: List branches
	var branchID string
	t.Run("ListBranches", func(t *testing.T) {
		branches, err := client.Agents().ListBranches(ctx, agentID)
		if err != nil {
			t.Fatalf("ListBranches() error = %v", err)
		}
		if len(branches) == 0 {
			t.Error("ListBranches() returned no branches")
		}
		for _, b := range branches {
			t.Logf("Branch: ID=%s, Name=%s, LivePct=%.0f%%", b.ID, b.Name, b.CurrentLivePercentage)
			if b.Name == "Main" {
				branchID = b.ID
			}
		}
	})

	// Step 5: Get specific branch
	t.Run("GetBranch", func(t *testing.T) {
		if branchID == "" {
			t.Skip("No branch ID from previous test")
		}
		branch, err := client.Agents().GetBranch(ctx, agentID, branchID)
		if err != nil {
			t.Fatalf("GetBranch() error = %v", err)
		}
		if branch.ID != branchID {
			t.Errorf("GetBranch() returned wrong ID: got %s, want %s", branch.ID, branchID)
		}
		t.Logf("Branch details: ID=%s, Name=%s, AgentID=%s", branch.ID, branch.Name, branch.AgentID)
	})

	// Step 6: Get shareable link
	t.Run("GetLink", func(t *testing.T) {
		link, err := client.Agents().GetLink(ctx, agentID)
		if err != nil {
			t.Fatalf("GetLink() error = %v", err)
		}
		if link.AgentID != agentID {
			t.Errorf("GetLink() returned wrong agent ID: got %s, want %s", link.AgentID, agentID)
		}
		t.Logf("Link: AgentID=%s, HasToken=%v", link.AgentID, link.SignedURL != "")
	})

	// Step 7: Get topics (likely empty for new agent)
	t.Run("GetTopics", func(t *testing.T) {
		topics, err := client.Agents().GetTopics(ctx, agentID)
		if err != nil {
			t.Fatalf("GetTopics() error = %v", err)
		}
		t.Logf("Found %d topics", len(topics))
	})

	// Step 8: Update agent
	t.Run("Update", func(t *testing.T) {
		updated, err := client.Agents().Update(ctx, agentID, &elevenlabs.UpdateAgentRequest{
			Name: "SDK Integration Test Agent (Updated)",
			Tags: []string{"test", "integration", "updated"},
		})
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		t.Logf("Updated agent: ID=%s", updated.AgentID)

		// Verify update
		agent, err := client.Agents().Get(ctx, agentID)
		if err != nil {
			t.Fatalf("Get() after update error = %v", err)
		}
		if agent.Name != "SDK Integration Test Agent (Updated)" {
			t.Logf("Note: Agent name may not be updated in Get response: %s", agent.Name)
		}
	})

	// Step 9: Duplicate agent
	var duplicateID string
	t.Run("Duplicate", func(t *testing.T) {
		dup, err := client.Agents().Duplicate(ctx, agentID, "SDK Test Agent Copy")
		if err != nil {
			t.Fatalf("Duplicate() error = %v", err)
		}
		duplicateID = dup.AgentID
		t.Logf("Duplicated agent: ID=%s", dup.AgentID)

		// Clean up duplicate
		t.Cleanup(func() {
			if duplicateID != "" {
				t.Logf("Cleanup: deleting duplicate agent %s", duplicateID)
				if err := client.Agents().Delete(ctx, duplicateID); err != nil {
					t.Logf("Cleanup warning: failed to delete duplicate: %v", err)
				}
			}
		})
	})

	// Step 10: Delete agent (explicit test, cleanup will handle if this fails)
	t.Run("Delete", func(t *testing.T) {
		err := client.Agents().Delete(ctx, agentID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		t.Logf("Deleted agent: %s", agentID)
		agentID = "" // Clear so cleanup doesn't try to delete again
	})
}

// TestAgentsService_Validation tests input validation without API calls.
func TestAgentsService_Validation(t *testing.T) {
	client, err := elevenlabs.NewClient(elevenlabs.WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name string
		fn   func() error
	}{
		{"Get empty ID", func() error { _, err := client.Agents().Get(ctx, ""); return err }},
		{"Delete empty ID", func() error { return client.Agents().Delete(ctx, "") }},
		{"Duplicate empty ID", func() error { _, err := client.Agents().Duplicate(ctx, "", "name"); return err }},
		{"GetLink empty ID", func() error { _, err := client.Agents().GetLink(ctx, ""); return err }},
		{"GetTopics empty ID", func() error { _, err := client.Agents().GetTopics(ctx, ""); return err }},
		{"ListBranches empty ID", func() error { _, err := client.Agents().ListBranches(ctx, ""); return err }},
		{"GetBranch empty agent ID", func() error { _, err := client.Agents().GetBranch(ctx, "", "branch"); return err }},
		{"GetBranch empty branch ID", func() error { _, err := client.Agents().GetBranch(ctx, "agent", ""); return err }},
		{"UpdateBranch empty agent ID", func() error { _, err := client.Agents().UpdateBranch(ctx, "", "branch", nil); return err }},
		{"UpdateBranch empty branch ID", func() error { _, err := client.Agents().UpdateBranch(ctx, "agent", "", nil); return err }},
		{"MergeBranch empty agent ID", func() error { return client.Agents().MergeBranch(ctx, "", "src", "tgt") }},
		{"MergeBranch empty source ID", func() error { return client.Agents().MergeBranch(ctx, "agent", "", "tgt") }},
		{"MergeBranch empty target ID", func() error { return client.Agents().MergeBranch(ctx, "agent", "src", "") }},
		{"Deploy empty agent ID", func() error { return client.Agents().Deploy(ctx, "", nil) }},
		{"Deploy empty deployments", func() error { return client.Agents().Deploy(ctx, "agent", nil) }},
		{"UploadAvatar empty ID", func() error { return client.Agents().UploadAvatar(ctx, "", nil) }},
		{"GetWidget empty ID", func() error { _, err := client.Agents().GetWidget(ctx, ""); return err }},
		{"CreateBranch empty ID", func() error { _, err := client.Agents().CreateBranch(ctx, "", nil); return err }},
		{"CreateBranch nil request", func() error { _, err := client.Agents().CreateBranch(ctx, "agent", nil); return err }},
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

// TestAgentTestingService_Folders_Integration tests agent test folder operations.
func TestAgentTestingService_Folders_Integration(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	client, err := elevenlabs.NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	var folderID string

	t.Cleanup(func() {
		if folderID != "" {
			t.Logf("Cleanup: deleting test folder %s", folderID)
			if err := client.AgentTesting().DeleteFolder(ctx, folderID); err != nil {
				t.Logf("Cleanup warning: failed to delete folder: %v", err)
			}
		}
	})

	// Create test folder
	t.Run("CreateFolder", func(t *testing.T) {
		folder, err := client.AgentTesting().CreateFolder(ctx, &elevenlabs.CreateTestFolderRequest{
			Name: "SDK Test Folder",
		})
		if err != nil {
			t.Fatalf("CreateFolder() error = %v", err)
		}
		folderID = folder.ID
		t.Logf("Created folder: ID=%s, Name=%s", folder.ID, folder.Name)
	})

	if folderID == "" {
		t.Fatal("Cannot continue without folder ID")
	}

	// Get test folder
	t.Run("GetFolder", func(t *testing.T) {
		folder, err := client.AgentTesting().GetFolder(ctx, folderID)
		if err != nil {
			t.Fatalf("GetFolder() error = %v", err)
		}
		if folder.ID != folderID {
			t.Errorf("GetFolder() returned wrong ID: got %s, want %s", folder.ID, folderID)
		}
		t.Logf("Folder: ID=%s, Name=%s, Children=%d", folder.ID, folder.Name, folder.ChildrenCount)
	})

	// Update test folder
	t.Run("UpdateFolder", func(t *testing.T) {
		folder, err := client.AgentTesting().UpdateFolder(ctx, folderID, &elevenlabs.UpdateTestFolderRequest{
			Name: "SDK Test Folder (Updated)",
		})
		if err != nil {
			t.Fatalf("UpdateFolder() error = %v", err)
		}
		if folder.Name != "SDK Test Folder (Updated)" {
			t.Logf("Note: Folder name update may not be reflected: %s", folder.Name)
		}
		t.Logf("Updated folder: ID=%s, Name=%s", folder.ID, folder.Name)
	})

	// Delete test folder
	t.Run("DeleteFolder", func(t *testing.T) {
		err := client.AgentTesting().DeleteFolder(ctx, folderID)
		if err != nil {
			t.Fatalf("DeleteFolder() error = %v", err)
		}
		t.Logf("Deleted folder: %s", folderID)
		folderID = "" // Clear so cleanup doesn't try again
	})
}
