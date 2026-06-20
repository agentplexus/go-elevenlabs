# Agent Testing

Manage test cases and test folders for Conversational AI agent validation.

## Overview

The Agent Testing service enables:

- **Test Organization**: Create folders to organize test cases
- **Test Management**: View, delete, and move tests between folders
- **Test Summaries**: Get summaries for multiple tests at once
- **Bulk Operations**: Move multiple tests in a single operation

## Test Folder Management

### Create Test Folder

```go
// Create top-level folder
folder, err := client.AgentTesting().CreateFolder(ctx, &elevenlabs.CreateTestFolderRequest{
    Name: "Regression Tests",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created folder: %s (ID: %s)\n", folder.Name, folder.ID)
```

### Create Nested Folder

```go
// Create folder inside another folder
subFolder, err := client.AgentTesting().CreateFolder(ctx, &elevenlabs.CreateTestFolderRequest{
    Name:           "Edge Cases",
    ParentFolderID: parentFolderID,
})
```

### Get Folder Details

```go
folder, err := client.AgentTesting().GetFolder(ctx, folderID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Folder: %s\n", folder.Name)
fmt.Printf("Children: %d\n", folder.ChildrenCount)
fmt.Printf("Path: ")
for _, seg := range folder.Path {
    fmt.Printf("/%s", seg.Name)
}
fmt.Println()
```

### Update Folder

```go
folder, err := client.AgentTesting().UpdateFolder(ctx, folderID, &elevenlabs.UpdateTestFolderRequest{
    Name: "Regression Tests (Updated)",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Updated folder: %s\n", folder.Name)
```

### Delete Folder

```go
err := client.AgentTesting().DeleteFolder(ctx, folderID)
if err != nil {
    log.Fatal(err)
}
```

## Listing Tests

### List All Tests

```go
resp, err := client.AgentTesting().ListTests(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, test := range resp.Tests {
    fmt.Printf("Test: %s (%s)\n", test.Name, test.ID)
    fmt.Printf("  Type: %s\n", test.EntityType)
    fmt.Printf("  Folder: %s\n", test.FolderID)
}

if resp.HasMore {
    // Fetch next page
    nextPage, _ := client.AgentTesting().ListTests(ctx, &elevenlabs.ListTestsOptions{
        Cursor: resp.NextCursor,
    })
}
```

### List Tests in Folder

```go
resp, err := client.AgentTesting().ListTests(ctx, &elevenlabs.ListTestsOptions{
    ParentFolderID: folderID,
    PageSize:       20,
})
```

### Search Tests

```go
resp, err := client.AgentTesting().ListTests(ctx, &elevenlabs.ListTestsOptions{
    Search:   "refund",
    PageSize: 10,
})
```

## Working with Tests

### Get Response Test

```go
test, err := client.AgentTesting().GetResponseTest(ctx, testID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Test: %s\n", test.Name)
fmt.Printf("Type: %s\n", test.Type)
fmt.Printf("Folder: %s\n", test.FolderID)
```

### Get Multiple Test Summaries

```go
summaries, err := client.AgentTesting().GetResponseTestSummaries(ctx, []string{
    "test-id-1",
    "test-id-2",
    "test-id-3",
})
if err != nil {
    log.Fatal(err)
}

for id, summary := range summaries {
    fmt.Printf("%s: %s (Folder: %s)\n", id, summary.Name, summary.FolderID)
}
```

### Delete Test

```go
err := client.AgentTesting().DeleteTest(ctx, testID)
if err != nil {
    log.Fatal(err)
}
```

## Bulk Operations

### Move Tests to Folder

```go
err := client.AgentTesting().BulkMoveTests(ctx, &elevenlabs.BulkMoveTestsRequest{
    TestIDs:        []string{"test-1", "test-2", "test-3"},
    TargetFolderID: targetFolderID,
})
if err != nil {
    log.Fatal(err)
}

fmt.Println("Tests moved successfully")
```

### Move Tests to Root

```go
// Leave TargetFolderID empty to move to root
err := client.AgentTesting().BulkMoveTests(ctx, &elevenlabs.BulkMoveTestsRequest{
    TestIDs: []string{"test-1", "test-2"},
    // TargetFolderID: "" - moves to root
})
```

## Request Types

### CreateTestFolderRequest

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Folder name |
| `ParentFolderID` | string | Optional parent folder ID |

### UpdateTestFolderRequest

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | New folder name |

### ListTestsOptions

| Field | Type | Description |
|-------|------|-------------|
| `PageSize` | int | Maximum results |
| `Cursor` | string | Pagination cursor |
| `ParentFolderID` | string | Filter by folder |
| `Search` | string | Search query |

### BulkMoveTestsRequest

| Field | Type | Description |
|-------|------|-------------|
| `TestIDs` | []string | Test IDs to move |
| `TargetFolderID` | string | Destination folder (empty for root) |

### AgentTestFolder

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Folder ID |
| `Name` | string | Folder name |
| `ChildrenCount` | int | Number of children |
| `Path` | []AgentTestFolderSegment | Path from root |

### TestSummary

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Test ID |
| `Name` | string | Test name |
| `FolderID` | string | Parent folder ID |
| `CreatedAtUnix` | int | Creation timestamp |
| `EntityType` | string | Type of test entity |
| `Path` | []AgentTestFolderSegment | Full path |

### AgentResponseTest

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Test ID |
| `Name` | string | Test name |
| `Type` | string | "llm", "tool", or "simulation" |
| `FolderID` | string | Parent folder ID |

## Example: Organize Tests by Category

```go
package main

import (
    "context"
    "log"

    elevenlabs "github.com/plexusone/elevenlabs-go"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create category folders
    categories := []string{"Unit Tests", "Integration Tests", "Regression Tests"}
    folderIDs := make(map[string]string)

    for _, category := range categories {
        folder, err := client.AgentTesting().CreateFolder(ctx, &elevenlabs.CreateTestFolderRequest{
            Name: category,
        })
        if err != nil {
            log.Printf("Warning: failed to create %s folder: %v", category, err)
            continue
        }
        folderIDs[category] = folder.ID
        log.Printf("Created folder: %s (ID: %s)", folder.Name, folder.ID)
    }

    // List all tests
    resp, err := client.AgentTesting().ListTests(ctx, &elevenlabs.ListTestsOptions{
        PageSize: 100,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Categorize tests based on naming convention
    unitTests := []string{}
    integrationTests := []string{}

    for _, test := range resp.Tests {
        switch {
        case len(test.Name) > 5 && test.Name[:5] == "unit_":
            unitTests = append(unitTests, test.ID)
        case len(test.Name) > 5 && test.Name[:5] == "intg_":
            integrationTests = append(integrationTests, test.ID)
        }
    }

    // Move tests to appropriate folders
    if len(unitTests) > 0 {
        err := client.AgentTesting().BulkMoveTests(ctx, &elevenlabs.BulkMoveTestsRequest{
            TestIDs:        unitTests,
            TargetFolderID: folderIDs["Unit Tests"],
        })
        if err != nil {
            log.Printf("Warning: failed to move unit tests: %v", err)
        } else {
            log.Printf("Moved %d unit tests", len(unitTests))
        }
    }

    if len(integrationTests) > 0 {
        err := client.AgentTesting().BulkMoveTests(ctx, &elevenlabs.BulkMoveTestsRequest{
            TestIDs:        integrationTests,
            TargetFolderID: folderIDs["Integration Tests"],
        })
        if err != nil {
            log.Printf("Warning: failed to move integration tests: %v", err)
        } else {
            log.Printf("Moved %d integration tests", len(integrationTests))
        }
    }
}
```
