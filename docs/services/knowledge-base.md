# Knowledge Base

Manage knowledge base documents for RAG (Retrieval-Augmented Generation) in Conversational AI agents.

!!! note "v0.13.0 API Change"
    As of v0.13.0, Knowledge Base methods are accessed via `client.Agents()` instead of `client.KnowledgeBase()`.

## Overview

The Knowledge Base service enables:

- **Document Management**: Upload files, create text documents, or index URLs
- **Folder Organization**: Organize documents into folders
- **RAG Indexing**: Manage RAG indexes for document retrieval
- **Chunk Access**: Retrieve document chunks for debugging
- **Agent Integration**: Track which agents use each document

## Document Types

| Type | Description | Use Case |
|------|-------------|----------|
| `file` | Uploaded files (PDF, DOCX, TXT, etc.) | Product manuals, policies |
| `text` | Plain text content | FAQs, custom content |
| `url` | Indexed web pages | Documentation sites, help centers |
| `folder` | Container for documents | Organization |

## Creating Documents

### Upload a File

```go
import "github.com/plexusone/elevenlabs-go/agents"

f, err := os.Open("product-manual.pdf")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

doc, err := client.Agents().CreateFileDocument(ctx, &agents.CreateFileDocumentRequest{
    File:     f,
    Filename: "product-manual.pdf",
    Name:     "Product Manual v2.0",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created document: %s (ID: %s)\n", doc.Name, doc.ID)
```

### Create Text Document

```go
doc, err := client.Agents().CreateTextDocument(ctx, &agents.CreateTextDocumentRequest{
    Content: `
# Frequently Asked Questions

## What are your business hours?
We're open Monday through Friday, 9 AM to 5 PM EST.

## How do I reset my password?
Click "Forgot Password" on the login page and follow the instructions.
    `,
    Name: "Company FAQ",
})
```

### Index a URL

```go
doc, err := client.Agents().CreateURLDocument(ctx, &agents.CreateURLDocumentRequest{
    URL:  "https://docs.example.com/api-reference",
    Name: "API Documentation",
})
```

## Listing Documents

### List All Documents

```go
resp, err := client.Agents().ListDocuments(ctx, nil)
if err != nil {
    log.Fatal(err)
}

for _, doc := range resp.Documents {
    fmt.Printf("%s: %s (%s)\n", doc.Type, doc.Name, doc.ID)
}
```

### With Pagination and Filters

```go
resp, err := client.Agents().ListDocuments(ctx, &agents.ListDocumentsOptions{
    PageSize:       20,
    Search:         "manual",
    ParentFolderID: folderID,
})

if resp.HasMore {
    nextPage, _ := client.Agents().ListDocuments(ctx, &agents.ListDocumentsOptions{
        Cursor: resp.NextCursor,
    })
}
```

## Reading Documents

### Get Document Metadata

```go
doc, err := client.Agents().GetDocument(ctx, documentID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", doc.Name)
fmt.Printf("Type: %s\n", doc.Type)
fmt.Printf("Size: %d bytes\n", doc.SizeBytes)
fmt.Printf("Created: %d\n", doc.CreatedAtUnix)
```

### Get Document Content

```go
content, err := client.Agents().GetDocumentContent(ctx, documentID)
if err != nil {
    log.Fatal(err)
}
defer content.Close()

// Read content
data, _ := io.ReadAll(content)
fmt.Printf("Content: %s\n", string(data))
```

### Get Document Chunks

```go
chunks, err := client.Agents().GetDocumentChunks(ctx, documentID)
if err != nil {
    log.Fatal(err)
}

for i, chunk := range chunks {
    fmt.Printf("Chunk %d: %s\n", i+1, chunk.Name)
    fmt.Printf("  Content: %s...\n", chunk.Content[:100])
}
```

### Get Specific Chunk

```go
chunk, err := client.Agents().GetDocumentChunk(ctx, documentID, chunkID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Chunk: %s\n", chunk.Name)
fmt.Printf("Content: %s\n", chunk.Content)
```

## Folder Management

### Create Folder

```go
// Create top-level folder
folder, err := client.Agents().CreateKnowledgeBaseFolder(ctx, &agents.CreateKnowledgeBaseFolderRequest{
    Name: "Product Documentation",
})

// Create nested folder
subFolder, err := client.Agents().CreateKnowledgeBaseFolder(ctx, &agents.CreateKnowledgeBaseFolderRequest{
    Name:           "API Docs",
    ParentFolderID: folder.ID,
})
```

### Upload to Folder

```go
f, _ := os.Open("api-guide.pdf")
defer f.Close()

doc, err := client.Agents().CreateFileDocument(ctx, &agents.CreateFileDocumentRequest{
    File:           f,
    Filename:       "api-guide.pdf",
    Name:           "API Guide",
    ParentFolderID: folderID,
})
```

## Agent Dependencies

### Find Agents Using a Document

```go
dependentAgents, err := client.Agents().GetDocumentDependentAgents(ctx, documentID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Document is used by %d agents:\n", len(dependentAgents))
for _, agent := range dependentAgents {
    fmt.Printf("  - %s (%s)\n", agent.Name, agent.ID)
}
```

## Delete Document

```go
err := client.Agents().DeleteDocument(ctx, documentID)
if err != nil {
    log.Fatal(err)
}
```

## Request Types

### CreateFileDocumentRequest

| Field | Type | Description |
|-------|------|-------------|
| `File` | io.Reader | File content to upload |
| `FileName` | string | Name of the file |
| `Name` | string | Optional display name |
| `ParentFolderID` | string | Optional parent folder |

### CreateTextDocumentRequest

| Field | Type | Description |
|-------|------|-------------|
| `Content` | string | Text content |
| `Name` | string | Document name |
| `ParentFolderID` | string | Optional parent folder |

### CreateURLDocumentRequest

| Field | Type | Description |
|-------|------|-------------|
| `URL` | string | URL to index |
| `Name` | string | Optional display name |
| `ParentFolderID` | string | Optional parent folder |

### ListDocumentsOptions

| Field | Type | Description |
|-------|------|-------------|
| `Cursor` | string | Pagination cursor |
| `PageSize` | int | Maximum results |
| `ParentFolderID` | string | Filter by folder |
| `Search` | string | Search query |

### KnowledgeBaseDocument

| Field | Type | Description |
|-------|------|-------------|
| `ID` | string | Document ID |
| `Name` | string | Display name |
| `Type` | string | "file", "text", "url", "folder" |
| `SourceURL` | string | URL for URL documents |
| `Content` | string | Content for text documents |
| `ParentFolderID` | string | Parent folder ID |
| `CreatedAtUnix` | int | Creation timestamp |
| `LastUpdatedUnix` | int | Last update timestamp |
| `SizeBytes` | int | Content size in bytes |

## Example: Bulk Import Documents

```go
package main

import (
    "context"
    "log"
    "os"
    "path/filepath"

    elevenlabs "github.com/plexusone/elevenlabs-go"
    "github.com/plexusone/elevenlabs-go/agents"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Create folder for imports
    folder, err := client.Agents().CreateKnowledgeBaseFolder(ctx, &agents.CreateKnowledgeBaseFolderRequest{
        Name: "Imported Docs",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Import all PDFs from a directory
    files, _ := filepath.Glob("./docs/*.pdf")
    for _, path := range files {
        f, err := os.Open(path)
        if err != nil {
            log.Printf("Warning: failed to open %s: %v", path, err)
            continue
        }

        doc, err := client.Agents().CreateFileDocument(ctx, &agents.CreateFileDocumentRequest{
            File:           f,
            Filename:       filepath.Base(path),
            ParentFolderID: folder.ID,
        })
        f.Close()

        if err != nil {
            log.Printf("Warning: failed to upload %s: %v", path, err)
            continue
        }

        log.Printf("Uploaded: %s (ID: %s)", doc.Name, doc.ID)
    }
}
```
