# Dubbing

Translate and dub audio/video content into other languages while preserving the original speaker's voice characteristics.

!!! note "v0.13.0 API Change"
    As of v0.13.0, Dubbing methods are accessed via `client.Content()` instead of `client.Dubbing()`.

## Creating a Dub

### From URL

```go
import "github.com/plexusone/elevenlabs-go/content"

dub, err := client.Content().CreateDubbing(ctx, &content.DubbingRequest{
    SourceURL:      "https://example.com/video.mp4",
    TargetLanguage: "es",  // Spanish
    Name:           "My Video - Spanish",
})
```

### Request Options

| Option | Description |
|--------|-------------|
| `SourceURL` | URL to video/audio file |
| `TargetLanguage` | Target language code |
| `Name` | Name for the dubbing project |
| `SourceLanguage` | Source language (auto-detected if not set) |
| `NumSpeakers` | Number of speakers (auto-detected if not set) |
| `Watermark` | Add watermark to output |
| `StartTime` | Start time in seconds |
| `EndTime` | End time in seconds |

## Checking Status

```go
status, err := client.Content().GetDubbingStatus(ctx, dubbingID)

fmt.Printf("Status: %s\n", status.Status)
fmt.Printf("Target Languages: %v\n", status.TargetLanguages)
```

## Dubbing Status Values

| Status | Description |
|--------|-------------|
| `dubbing` | In progress |
| `dubbed` | Complete |
| `failed` | Failed |

## Downloading Dubbed Audio

```go
audio, err := client.Content().GetDubbedFile(ctx, dubbingID, "es")
if err != nil {
    log.Fatal(err)
}

f, _ := os.Create("dubbed_spanish.mp4")
defer f.Close()
io.Copy(f, audio)
```

## Deleting a Dub

```go
err := client.Content().DeleteDubbing(ctx, dubbingID)
```

## Supported Languages

Common language codes:

| Code | Language |
|------|----------|
| `en` | English |
| `es` | Spanish |
| `fr` | French |
| `de` | German |
| `it` | Italian |
| `pt` | Portuguese |
| `pl` | Polish |
| `hi` | Hindi |
| `ja` | Japanese |
| `ko` | Korean |
| `zh` | Chinese |

## Workflow Example

```go
// 1. Create dubbing job
dub, err := client.Content().CreateDubbing(ctx, &content.DubbingRequest{
    SourceURL:      "https://example.com/course-intro.mp4",
    TargetLanguage: "es",
    Name:           "Course Intro - Spanish",
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Dubbing ID: %s\n", dub.DubbingID)

// 2. Poll for completion
for {
    status, _ := client.Content().GetDubbingStatus(ctx, dub.DubbingID)

    if status.Status == "dubbed" {
        fmt.Println("Dubbing complete!")
        break
    } else if status.Status == "failed" {
        log.Fatal("Dubbing failed")
    }

    fmt.Printf("Status: %s, waiting...\n", status.Status)
    time.Sleep(30 * time.Second)
}

// 3. Download dubbed file
audio, _ := client.Content().GetDubbedFile(ctx, dub.DubbingID, "es")
f, _ := os.Create("intro_spanish.mp4")
io.Copy(f, audio)
f.Close()
```

## Best Practices

1. **Check source quality** - Better input = better output
2. **Specify speaker count** - Helps with voice separation
3. **Review output** - AI dubbing may need manual review
4. **Consider cultural context** - Some content may need localization beyond translation
