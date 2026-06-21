# Conversational AI Agents

Manage ElevenLabs Conversational AI agents programmatically.

## Overview

The Agents service enables:

- **Agent Management**: Create, read, update, and delete conversational AI agents
- **Branch Management**: Version control with branches for agent configurations
- **Deployment**: Deploy agent branches with traffic splitting
- **Analytics**: Access conversation topics and agent insights
- **Testing**: Organize and manage agent test cases

## Basic Usage

### List Agents

```go
// List all agents
resp, err := client.Agents().List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, agent := range resp.Agents {
    fmt.Printf("Agent: %s (%s)\n", agent.Name, agent.AgentID)
}
```

### With Pagination and Filtering

```go
import "github.com/plexusone/elevenlabs-go/agents"

resp, err := client.Agents().List(ctx, &agents.ListAgentsOptions{
    PageSize:        10,
    Search:          "support",
    SortBy:          "created_at",
    SortDirection:   "desc",
    CreatedByUserID: "@me", // Only your agents
})

if resp.HasMore {
    // Fetch next page
    nextPage, _ := client.Agents().List(ctx, &agents.ListAgentsOptions{
        Cursor: resp.NextCursor,
    })
}
```

## Creating Agents

```go
import "github.com/plexusone/elevenlabs-go/agents"

agent, err := client.Agents().Create(ctx, &agents.CreateAgentRequest{
    Name: "Customer Support Agent",
    Tags: []string{"support", "production"},
    ConversationConfig: map[string]any{
        "agent": map[string]any{
            "prompt": map[string]any{
                "prompt": "You are a helpful customer support agent for Acme Corp.",
            },
            "first_message": "Hello! How can I help you today?",
            "language":      "en",
        },
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created agent: %s\n", agent.AgentID)
```

## Reading and Updating Agents

### Get Agent

```go
agent, err := client.Agents().Get(ctx, agentID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Agent: %s\n", agent.Name)
fmt.Printf("Created: %d\n", agent.CreatedAtUnixSecs)
```

### Update Agent

```go
updated, err := client.Agents().Update(ctx, agentID, &agents.UpdateAgentRequest{
    Name: "Customer Support Agent v2",
    Tags: []string{"support", "production", "v2"},
})
```

### Delete Agent

```go
err := client.Agents().Delete(ctx, agentID)
```

### Duplicate Agent

```go
copy, err := client.Agents().Duplicate(ctx, agentID, "Agent Copy")
fmt.Printf("Duplicated to: %s\n", copy.AgentID)
```

## Branch Management

Branches enable version control for agent configurations.

### List Branches

```go
branches, err := client.Agents().ListBranches(ctx, agentID)
if err != nil {
    log.Fatal(err)
}

for _, b := range branches {
    fmt.Printf("Branch: %s (%s) - %.0f%% live\n",
        b.Name, b.ID, b.CurrentLivePercentage)
}
```

### Get Branch Details

```go
branch, err := client.Agents().GetBranch(ctx, agentID, branchID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Branch: %s\n", branch.Name)
fmt.Printf("Description: %s\n", branch.Description)
fmt.Printf("Last committed: %d\n", branch.LastCommittedAt)
```

### Create Branch

```go
branch, err := client.Agents().CreateBranch(ctx, agentID, &agents.CreateBranchRequest{
    Name:         "experiment-v2",
    Description:  "Testing new prompt structure",
    BaseBranchID: mainBranchID, // Fork from existing branch
})
```

### Update Branch

```go
updated, err := client.Agents().UpdateBranch(ctx, agentID, branchID,
    &agents.UpdateBranchRequest{
        Name: "experiment-v2-final",
    })
```

### Merge Branches

```go
// Merge source branch into target branch
err := client.Agents().MergeBranch(ctx, agentID, sourceBranchID, targetBranchID)
```

## Deployment

Deploy branches with traffic splitting for A/B testing.

```go
// Deploy main branch at 80%, experiment at 20%
err := client.Agents().Deploy(ctx, agentID, []agents.DeploymentRequest{
    {BranchID: mainBranchID, Percentage: 80.0},
    {BranchID: experimentBranchID, Percentage: 20.0},
})

// Deploy single branch at 100%
err := client.Agents().Deploy(ctx, agentID, []agents.DeploymentRequest{
    {BranchID: mainBranchID, Percentage: 100.0},
})
```

## Agent Analytics

### Get Shareable Link

```go
link, err := client.Agents().GetLink(ctx, agentID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Token: %s\n", link.SignedURL)
fmt.Printf("Expires: %d\n", link.ExpiresAtUnixMillis)
```

### Get Conversation Topics

```go
topics, err := client.Agents().GetTopics(ctx, agentID)
if err != nil {
    log.Fatal(err)
}

for _, topic := range topics {
    fmt.Printf("Topic: %s (%d conversations)\n",
        topic.Label, topic.ConversationCount)
}
```

### Get Widget Configuration

```go
widget, err := client.Agents().GetWidget(ctx, agentID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Widget config: %+v\n", widget.WidgetConfig)
```

## Avatar Upload

```go
f, err := os.Open("avatar.png")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

err = client.Agents().UploadAvatar(ctx, agentID, f)
```

## Test Folder Management

Test folder management is part of the unified `agents` package.

```go
import "github.com/plexusone/elevenlabs-go/agents"

// Create test folder
folder, err := client.Agents().CreateTestFolder(ctx, &agents.CreateTestFolderRequest{
    Name: "Regression Tests",
})

// List tests
tests, err := client.Agents().ListTests(ctx, nil)

// Bulk move tests
err := client.Agents().BulkMoveTests(ctx, &agents.BulkMoveTestsRequest{
    TestIDs:        []string{"test1", "test2", "test3"},
    TargetFolderID: folderID,
})
```

See [Agent Testing](agent-testing.md) for full documentation.

## Conversation Simulation

Simulate conversations for testing agents.

### Non-Streaming Simulation

```go
import "github.com/plexusone/elevenlabs-go/agents"

messages, err := client.Agents().SimulateConversation(ctx, agentID,
    &agents.SimulateConversationOptions{
        SimulationPersona: "A frustrated customer asking about refunds",
        InitialMessage:    "I want my money back!",
        MaxTurns:          5,
    })
if err != nil {
    log.Fatal(err)
}

for _, msg := range messages {
    fmt.Printf("[%s] %s\n", msg.Role, msg.Content)
}
```

### Streaming Simulation

```go
conn, err := client.Agents().SimulateConversationStream(ctx, agentID,
    &agents.SimulateConversationOptions{
        SimulationPersona: "A curious customer",
    })
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Receive messages
go func() {
    for msg := range conn.Messages() {
        fmt.Printf("[%s] %s\n", msg.Role, msg.Content)
    }
}()

// Handle errors
for err := range conn.Errors() {
    log.Printf("Error: %v", err)
}
```

## Request Types

### ListAgentsOptions

| Field | Type | Description |
|-------|------|-------------|
| `PageSize` | int | Max agents to return (max 100, default 30) |
| `Search` | string | Filter by agent name |
| `Archived` | *bool | Filter by archived status |
| `SortBy` | string | Sort field (name, created_at, etc.) |
| `SortDirection` | string | Sort direction (asc, desc) |
| `Cursor` | string | Pagination cursor |
| `CreatedByUserID` | string | Filter by creator (@me for self) |

### CreateAgentRequest

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Agent display name |
| `Tags` | []string | Categorization tags |
| `ConversationConfig` | map[string]any | Agent configuration |
| `PlatformSettings` | map[string]any | Platform-specific settings |

### DeploymentRequest

| Field | Type | Description |
|-------|------|-------------|
| `BranchID` | string | Branch to deploy |
| `Percentage` | float64 | Traffic percentage (0-100) |

## Example: Full Agent Lifecycle

```go
package main

import (
    "context"
    "fmt"
    "log"

    elevenlabs "github.com/plexusone/elevenlabs-go"
    "github.com/plexusone/elevenlabs-go/agents"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create agent
    agent, err := client.Agents().Create(ctx, &agents.CreateAgentRequest{
        Name: "Demo Agent",
        Tags: []string{"demo"},
        ConversationConfig: map[string]any{
            "agent": map[string]any{
                "prompt": map[string]any{
                    "prompt": "You are a helpful demo assistant.",
                },
                "first_message": "Hello! I'm a demo agent.",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created: %s\n", agent.AgentID)

    // List branches
    branches, _ := client.Agents().ListBranches(ctx, agent.AgentID)
    fmt.Printf("Branches: %d\n", len(branches))

    // Get shareable link
    link, _ := client.Agents().GetLink(ctx, agent.AgentID)
    fmt.Printf("Link token available: %v\n", link.SignedURL != "")

    // Duplicate for testing
    copy, _ := client.Agents().Duplicate(ctx, agent.AgentID, "Demo Agent Copy")
    fmt.Printf("Duplicated: %s\n", copy.AgentID)

    // Clean up
    client.Agents().Delete(ctx, copy.AgentID)
    client.Agents().Delete(ctx, agent.AgentID)
    fmt.Println("Cleaned up")
}
```
