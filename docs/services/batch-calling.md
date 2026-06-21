# Batch Calling

Schedule and manage bulk outbound calls with Conversational AI agents.

!!! note "v0.13.0 API Change"
    As of v0.13.0, Batch Calling methods are accessed via `client.Agents()` instead of `client.BatchCalling()`.

## Overview

The Batch Calling service enables:

- **Bulk Campaigns**: Schedule calls to multiple recipients
- **Job Management**: Create, monitor, and control batch jobs
- **Scheduling**: Set specific times and timezones for calls
- **Concurrency Control**: Limit simultaneous calls
- **Retry Handling**: Retry failed calls automatically

## Creating Batch Calls

### Basic Batch Call

```go
import "github.com/plexusone/elevenlabs-go/agents"

batch, err := client.Agents().CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
    Name:               "Customer Outreach",
    AgentID:            agentID,
    AgentPhoneNumberID: phoneNumberID,
    Recipients: []agents.BatchCallRecipient{
        {PhoneNumber: "+14155551234"},
        {PhoneNumber: "+14155555678"},
        {PhoneNumber: "+14155559012"},
    },
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created batch: %s (ID: %s)\n", batch.Name, batch.ID)
fmt.Printf("Scheduled: %d calls\n", batch.TotalCallsScheduled)
```

### Scheduled Batch Call

```go
// Schedule for tomorrow at 9 AM EST
scheduledTime := int(time.Now().Add(24 * time.Hour).Unix())

batch, err := client.Agents().CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
    Name:               "Appointment Reminders",
    AgentID:            agentID,
    AgentPhoneNumberID: phoneNumberID,
    ScheduledTimeUnix:  &scheduledTime,
    Timezone:           "America/New_York",
    Recipients: []agents.BatchCallRecipient{
        {PhoneNumber: "+14155551234"},
        {PhoneNumber: "+14155555678"},
    },
})
```

### With Concurrency Limit

```go
concurrencyLimit := 5 // Max 5 simultaneous calls

batch, err := client.Agents().CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
    Name:                   "High-Volume Campaign",
    AgentID:                agentID,
    AgentPhoneNumberID:     phoneNumberID,
    TargetConcurrencyLimit: &concurrencyLimit,
    Recipients: []agents.BatchCallRecipient{
        {PhoneNumber: "+14155551234"},
        // ... many more recipients
    },
})
```

### With Custom Data per Recipient

```go
batch, err := client.Agents().CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
    Name:               "Personalized Outreach",
    AgentID:            agentID,
    AgentPhoneNumberID: phoneNumberID,
    Recipients: []agents.BatchCallRecipient{
        {
            PhoneNumber: "+14155551234",
            ID:          "customer-001",
            DynamicVariables: map[string]any{
                "name":        "John Smith",
                "account_id":  "ACC-12345",
                "appointment": "2024-01-15 10:00",
            },
        },
        {
            PhoneNumber: "+14155555678",
            ID:          "customer-002",
            DynamicVariables: map[string]any{
                "name":        "Jane Doe",
                "account_id":  "ACC-67890",
                "appointment": "2024-01-15 14:00",
            },
        },
    },
})
```

## Listing Batch Calls

### List All Batches

```go
resp, err := client.Agents().ListBatchCalls(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, batch := range resp.BatchCalls {
    fmt.Printf("%s: %s (Status: %s)\n", batch.ID, batch.Name, batch.Status)
    fmt.Printf("  Progress: %d/%d calls completed\n",
        batch.TotalCallsFinished, batch.TotalCallsScheduled)
}
```

### Filter by Agent

```go
resp, err := client.Agents().ListBatchCalls(ctx, &agents.ListBatchCallsOptions{
    AgentID:  agentID,
    PageSize: 20,
})
```

### Pagination

```go
resp, err := client.Agents().ListBatchCalls(ctx, &agents.ListBatchCallsOptions{
    PageSize: 10,
})

// Fetch next page
if resp.NextCursor != "" {
    nextPage, _ := client.Agents().ListBatchCalls(ctx, &agents.ListBatchCallsOptions{
        Cursor:   resp.NextCursor,
        PageSize: 10,
    })
}
```

## Get Batch Details

```go
batch, err := client.Agents().GetBatchCall(ctx, batchID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Batch: %s\n", batch.Name)
fmt.Printf("Agent: %s (%s)\n", batch.AgentName, batch.AgentID)
fmt.Printf("Status: %s\n", batch.Status)
fmt.Printf("Scheduled: %d\n", batch.TotalCallsScheduled)
fmt.Printf("Dispatched: %d\n", batch.TotalCallsDispatched)
fmt.Printf("Finished: %d\n", batch.TotalCallsFinished)
fmt.Printf("Retry Count: %d\n", batch.RetryCount)
```

## Managing Batch Jobs

### Cancel a Batch

```go
err := client.Agents().CancelBatchCall(ctx, batchID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Batch cancelled")
```

### Retry Failed Calls

```go
_, err := client.Agents().RetryBatchCall(ctx, batchID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Retrying failed calls")
```

### Delete a Batch

```go
err := client.Agents().DeleteBatchCall(ctx, batchID)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Batch deleted")
```

## Request Types

### CreateBatchCallRequest

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Batch job name |
| `AgentID` | string | Agent to use for calls |
| `BranchID` | string | Optional branch ID |
| `PhoneNumberID` | string | Phone number for outbound calls |
| `Recipients` | []BatchCallRecipient | List of recipients |
| `ScheduledTimeUnix` | *int | Optional scheduled time (Unix) |
| `Timezone` | string | Timezone (e.g., "America/New_York") |
| `Environment` | string | Environment (e.g., "production") |
| `TargetConcurrencyLimit` | *int | Max concurrent calls |

### BatchCallRecipient

| Field | Type | Description |
|-------|------|-------------|
| `PhoneNumber` | string | Phone number (E.164 format) |
| `ID` | string | Optional tracking ID |
| `WhatsAppUserID` | string | Optional WhatsApp user ID |
| `CustomData` | map[string]any | Custom data for conversation |

### ListBatchCallsOptions

| Field | Type | Description |
|-------|------|-------------|
| `Limit` | int | Maximum results |
| `LastDoc` | string | Pagination cursor |
| `AgentID` | string | Filter by agent |

### BatchCall

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Batch ID |
| `Name` | string | Batch name |
| `AgentID` | string | Agent ID |
| `AgentName` | string | Agent name |
| `BranchID` | string | Branch ID |
| `BranchName` | string | Branch name |
| `Status` | string | Job status |
| `PhoneNumberID` | string | Phone number ID |
| `PhoneProvider` | string | Phone provider |
| `Environment` | string | Environment |
| `Timezone` | string | Timezone |
| `ScheduledTimeUnix` | int | Scheduled time |
| `CreatedAtUnix` | int | Creation time |
| `LastUpdatedAtUnix` | int | Last update time |
| `TotalCallsScheduled` | int | Total calls planned |
| `TotalCallsDispatched` | int | Calls started |
| `TotalCallsFinished` | int | Calls completed |
| `RetryCount` | int | Number of retries |
| `TargetConcurrencyLimit` | *int | Concurrency limit |

## Example: Campaign with Monitoring

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    elevenlabs "github.com/plexusone/elevenlabs-go"
    "github.com/plexusone/elevenlabs-go/agents"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create batch call
    batch, err := client.Agents().CreateBatchCall(ctx, &agents.CreateBatchCallRequest{
        Name:               "Customer Survey",
        AgentID:            "your-agent-id",
        AgentPhoneNumberID: "your-phone-number-id",
        Recipients: []agents.BatchCallRecipient{
            {PhoneNumber: "+14155551234", ID: "survey-001"},
            {PhoneNumber: "+14155555678", ID: "survey-002"},
            {PhoneNumber: "+14155559012", ID: "survey-003"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Started batch: %s\n", batch.ID)

    // Monitor progress
    for {
        batchDetail, err := client.Agents().GetBatchCall(ctx, batch.ID)
        if err != nil {
            log.Fatal(err)
        }

        progress := float64(batchDetail.TotalCallsFinished) / float64(batchDetail.TotalCallsScheduled) * 100
        fmt.Printf("Progress: %.1f%% (%d/%d calls completed)\n",
            progress, batchDetail.TotalCallsFinished, batchDetail.TotalCallsScheduled)

        if batchDetail.Status == "completed" || batchDetail.Status == "cancelled" {
            fmt.Printf("Batch %s\n", batchDetail.Status)
            break
        }

        time.Sleep(10 * time.Second)
    }
}
```
