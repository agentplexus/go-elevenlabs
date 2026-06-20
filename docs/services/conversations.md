# Conversations

Access and manage conversation history for Conversational AI agents.

## Overview

The Conversations service enables:

- **History Access**: Retrieve conversation records with transcripts
- **Audio Playback**: Get conversation audio recordings
- **Analysis**: Run and retrieve conversation analysis
- **Search**: Search across conversation transcripts
- **Feedback**: Submit user feedback for conversations

## Basic Usage

### List Conversations

```go
// List recent conversations
resp, err := client.Conversations().List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, conv := range resp.Conversations {
    fmt.Printf("Conversation: %s (Agent: %s)\n", conv.ConversationID, conv.AgentName)
    fmt.Printf("  Duration: %ds, Messages: %d\n", conv.CallDurationSecs, conv.MessageCount)
}
```

### With Filters

```go
// Filter by agent and success status
callStartAfter := int(time.Now().Add(-24 * time.Hour).Unix())
resp, err := client.Conversations().List(ctx, &elevenlabs.ListConversationsOptions{
    AgentID:            agentID,
    CallSuccessful:     "success",
    CallStartAfterUnix: &callStartAfter,
    PageSize:           20,
    IncludeSummary:     true,
})

if resp.HasMore {
    // Fetch next page
    nextPage, _ := client.Conversations().List(ctx, &elevenlabs.ListConversationsOptions{
        Cursor: resp.NextCursor,
    })
}
```

## Get Conversation Details

```go
conv, err := client.Conversations().Get(ctx, conversationID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Conversation: %s\n", conv.ConversationID)
fmt.Printf("Agent: %s\n", conv.AgentName)
fmt.Printf("Status: %s\n", conv.Status)
fmt.Printf("Has Audio: %v\n", conv.HasAudio)

// Print transcript
for _, msg := range conv.Transcript {
    fmt.Printf("[%s] %s\n", msg.Role, msg.Message)

    // Check for tool calls
    for _, tc := range msg.ToolCalls {
        fmt.Printf("  -> Tool: %s\n", tc.ToolName)
    }
}

// Print analysis if available
if conv.Analysis != nil {
    fmt.Printf("Summary: %s\n", conv.Analysis.TranscriptSummary)
    fmt.Printf("Result: %s\n", conv.Analysis.CallSuccessful)
}
```

## Get Audio Recording

```go
audio, err := client.Conversations().GetAudio(ctx, conversationID)
if err != nil {
    log.Fatal(err)
}
defer audio.Close()

// Save to file
f, _ := os.Create("conversation.mp3")
defer f.Close()
io.Copy(f, audio)
```

## Search Transcripts

```go
// Search for conversations mentioning specific terms
results, err := client.Conversations().Search(ctx, "refund policy", &elevenlabs.ListConversationsOptions{
    AgentID:  agentID,
    PageSize: 10,
})
if err != nil {
    log.Fatal(err)
}

for _, conv := range results.Conversations {
    fmt.Printf("Found: %s - %s\n", conv.ConversationID, conv.CallSummaryTitle)
}
```

## Submit Feedback

```go
// Submit positive feedback
err := client.Conversations().SubmitFeedback(ctx, conversationID, "like")

// Submit negative feedback
err = client.Conversations().SubmitFeedback(ctx, conversationID, "dislike")
```

## Run Analysis

```go
// Trigger analysis on a conversation
err := client.Conversations().RunAnalysis(ctx, conversationID)
if err != nil {
    log.Fatal(err)
}

// Fetch updated conversation to get analysis results
conv, _ := client.Conversations().Get(ctx, conversationID)
fmt.Printf("Analysis: %s\n", conv.Analysis.TranscriptSummary)
```

## Delete Conversation

```go
err := client.Conversations().Delete(ctx, conversationID)
if err != nil {
    log.Fatal(err)
}
```

## Request Types

### ListConversationsOptions

| Field | Type | Description |
|-------|------|-------------|
| `Cursor` | string | Pagination cursor |
| `AgentID` | string | Filter by agent ID |
| `CallSuccessful` | string | Filter by result: "success", "failure", "unknown" |
| `CallStartBeforeUnix` | *int | Filter before this Unix timestamp |
| `CallStartAfterUnix` | *int | Filter after this Unix timestamp |
| `CallDurationMinSecs` | *int | Minimum call duration in seconds |
| `CallDurationMaxSecs` | *int | Maximum call duration in seconds |
| `RatingMin` | *int | Minimum rating (1-5) |
| `RatingMax` | *int | Maximum rating (1-5) |
| `HasFeedbackComment` | *bool | Filter by presence of feedback |
| `UserID` | string | Filter by user ID |
| `BranchID` | string | Filter by branch ID |
| `TopicIDs` | []string | Filter by topic IDs |
| `ToolNames` | []string | Filter by tool names used |
| `MainLanguages` | []string | Filter by languages |
| `ExcludeStatuses` | []string | Exclude these statuses |
| `PageSize` | int | Max results (max 100, default 30) |
| `IncludeSummary` | bool | Include transcript summaries |

### Conversation

| Field | Type | Description |
|-------|------|-------------|
| `ConversationID` | string | Unique identifier |
| `AgentID` | string | Agent that handled the conversation |
| `AgentName` | string | Agent display name |
| `BranchID` | string | Branch used for the conversation |
| `Status` | string | Conversation status |
| `StartTimeUnixSecs` | int | Start time (Unix timestamp) |
| `CallDurationSecs` | int | Duration in seconds |
| `MessageCount` | int | Number of messages |
| `CallSuccessful` | string | Success evaluation |
| `CallSummaryTitle` | string | Short summary title |
| `TranscriptSummary` | string | Full transcript summary |
| `Rating` | *float64 | User rating (1-5) |
| `MainLanguage` | string | Primary language |
| `TerminationReason` | string | How the call ended |
| `Direction` | string | "inbound" or "outbound" |
| `ToolNames` | []string | Tools used in conversation |

### TranscriptMessage

| Field | Type | Description |
|-------|------|-------------|
| `Role` | string | "user" or "agent" |
| `Message` | string | Message content |
| `TimeInCallSecs` | int | Time offset in the call |
| `Interrupted` | bool | Whether message was interrupted |
| `ToolCalls` | []ToolCall | Tool invocations in this message |

## Example: Export Conversations

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "time"

    elevenlabs "github.com/plexusone/elevenlabs-go"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Export last 7 days of conversations
    weekAgo := int(time.Now().Add(-7 * 24 * time.Hour).Unix())

    var allConversations []*elevenlabs.ConversationDetail
    var cursor string

    for {
        resp, err := client.Conversations().List(ctx, &elevenlabs.ListConversationsOptions{
            CallStartAfterUnix: &weekAgo,
            Cursor:             cursor,
            PageSize:           100,
        })
        if err != nil {
            log.Fatal(err)
        }

        // Fetch full details for each conversation
        for _, conv := range resp.Conversations {
            detail, err := client.Conversations().Get(ctx, conv.ConversationID)
            if err != nil {
                log.Printf("Warning: failed to get %s: %v", conv.ConversationID, err)
                continue
            }
            allConversations = append(allConversations, detail)
        }

        if !resp.HasMore {
            break
        }
        cursor = resp.NextCursor
    }

    // Export to JSON
    f, _ := os.Create("conversations-export.json")
    defer f.Close()

    enc := json.NewEncoder(f)
    enc.SetIndent("", "  ")
    enc.Encode(allConversations)

    log.Printf("Exported %d conversations", len(allConversations))
}
```
