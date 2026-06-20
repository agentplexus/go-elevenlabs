# Analytics

Access real-time analytics for Conversational AI operations.

## Overview

The Analytics service provides:

- **Live Monitoring**: Track active ongoing conversations in real-time

## Getting Live Conversation Count

```go
count, err := client.Analytics().GetLiveCount(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Active conversations: %d\n", count)
```

## Use Cases

### Real-Time Dashboard

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    elevenlabs "github.com/plexusone/elevenlabs-go"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Monitor live conversation count
    for {
        count, err := client.Analytics().GetLiveCount(ctx)
        if err != nil {
            log.Printf("Error: %v", err)
        } else {
            fmt.Printf("[%s] Active conversations: %d\n",
                time.Now().Format("15:04:05"), count)
        }

        time.Sleep(5 * time.Second)
    }
}
```

### Capacity Alerting

```go
package main

import (
    "context"
    "log"
    "time"

    elevenlabs "github.com/plexusone/elevenlabs-go"
)

const MaxConcurrentCalls = 100
const AlertThreshold = 0.8 // 80%

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    for {
        count, err := client.Analytics().GetLiveCount(ctx)
        if err != nil {
            log.Printf("Error getting live count: %v", err)
            time.Sleep(10 * time.Second)
            continue
        }

        utilization := float64(count) / float64(MaxConcurrentCalls)

        if utilization >= AlertThreshold {
            log.Printf("WARNING: High utilization %.1f%% (%d/%d calls)",
                utilization*100, count, MaxConcurrentCalls)
            // Send alert to monitoring system
        }

        time.Sleep(10 * time.Second)
    }
}
```

### Combining with Conversation Data

```go
package main

import (
    "context"
    "fmt"
    "log"

    elevenlabs "github.com/plexusone/elevenlabs-go"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get current live count
    liveCount, err := client.Analytics().GetLiveCount(ctx)
    if err != nil {
        log.Fatal(err)
    }

    // Get recent conversations
    resp, err := client.Conversations().List(ctx, &elevenlabs.ListConversationsOptions{
        PageSize: 100,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Calculate statistics
    var totalDuration int
    successCount := 0
    for _, conv := range resp.Conversations {
        totalDuration += conv.CallDurationSecs
        if conv.CallSuccessful == "success" {
            successCount++
        }
    }

    avgDuration := 0
    successRate := 0.0
    if len(resp.Conversations) > 0 {
        avgDuration = totalDuration / len(resp.Conversations)
        successRate = float64(successCount) / float64(len(resp.Conversations)) * 100
    }

    fmt.Printf("=== Conversation Analytics ===\n")
    fmt.Printf("Live conversations: %d\n", liveCount)
    fmt.Printf("Recent conversations: %d\n", len(resp.Conversations))
    fmt.Printf("Average duration: %ds\n", avgDuration)
    fmt.Printf("Success rate: %.1f%%\n", successRate)
}
```

## Methods

### GetLiveCount

Returns the number of currently active conversations.

```go
func (s *AnalyticsService) GetLiveCount(ctx context.Context) (int, error)
```

**Returns:**

- `int`: Number of active ongoing conversations
- `error`: Error if the request fails

## Related Services

For more detailed analytics, see:

- [Conversations](conversations.md) - Access conversation history and transcripts
- [Agents](agents.md) - Get conversation topics per agent
