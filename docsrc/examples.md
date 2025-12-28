# Examples

Complete working examples for common use cases.

## Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"

    elevenlabs "github.com/agentplexus/go-elevenlabs"
)

func main() {
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List voices
    voices, _ := client.Voices().List(ctx)
    fmt.Printf("Found %d voices\n", len(voices))

    // Generate speech
    if len(voices) > 0 {
        audio, _ := client.TextToSpeech().Simple(ctx,
            voices[0].VoiceID,
            "Hello from go-elevenlabs!")

        f, _ := os.Create("hello.mp3")
        defer f.Close()
        io.Copy(f, audio)
    }
}
```

## Text-to-Speech with Options

```go
resp, err := client.TextToSpeech().Generate(ctx, &elevenlabs.TTSRequest{
    VoiceID: "21m00Tcm4TlvDq8ikWAM",
    Text:    "Hello with custom settings!",
    ModelID: "eleven_multilingual_v2",
    VoiceSettings: &elevenlabs.VoiceSettings{
        Stability:       0.6,
        SimilarityBoost: 0.8,
        Style:           0.1,
        SpeakerBoost:    true,
    },
    OutputFormat: "mp3_44100_192",
})
if err != nil {
    log.Fatal(err)
}

f, _ := os.Create("custom.mp3")
defer f.Close()
io.Copy(f, resp.Audio)
```

## Sound Effects

```go
// Simple sound effect
thunder, _ := client.SoundEffects().Simple(ctx, "thunder and rain storm")

// With options
sfx, _ := client.SoundEffects().Generate(ctx, &elevenlabs.SoundEffectRequest{
    Text:            "spaceship engine humming",
    DurationSeconds: 10,
    PromptInfluence: 0.5,
    Loop:            true,
})

// Looping background
ambience, _ := client.SoundEffects().GenerateLoop(ctx,
    "peaceful forest with birds", 30)
```

## Pronunciation Dictionary

```go
// From a map (simplest)
dict, _ := client.Pronunciation().CreateFromMap(ctx, "Tech Terms", map[string]string{
    "API":     "A P I",
    "kubectl": "kube control",
    "nginx":   "engine X",
})

// From JSON file
dict, _ := client.Pronunciation().CreateFromJSON(ctx, "Terms", "terms.json")

// With full options
rules := elevenlabs.PronunciationRules{
    {Grapheme: "API", Alias: "A P I"},
    {Grapheme: "nginx", Phoneme: "ˈɛndʒɪnˈɛks"},
}

dict, _ := client.Pronunciation().Create(ctx, &elevenlabs.CreatePronunciationDictionaryRequest{
    Name:        "Custom Terms",
    Description: "Technical vocabulary",
    Rules:       rules,
    Language:    "en-US",
})
```

## Projects (Long-form Content)

```go
// Create project
project, _ := client.Projects().Create(ctx, &elevenlabs.CreateProjectRequest{
    Name:                    "My Audiobook",
    DefaultModelID:          "eleven_multilingual_v2",
    DefaultParagraphVoiceID: "21m00Tcm4TlvDq8ikWAM",
    DefaultTitleVoiceID:     "21m00Tcm4TlvDq8ikWAM",
})

// List chapters
chapters, _ := client.Projects().ListChapters(ctx, project.ProjectID)

// Convert project to audio
client.Projects().Convert(ctx, project.ProjectID)

// Download completed audio
snapshots, _ := client.Projects().ListSnapshots(ctx, project.ProjectID)
if len(snapshots) > 0 {
    reader, _ := client.Projects().DownloadSnapshotArchive(ctx,
        project.ProjectID, snapshots[0].ProjectSnapshotID)

    f, _ := os.Create("audiobook.zip")
    io.Copy(f, reader)
    f.Close()
}
```

## Dubbing

```go
// Create dubbing job
dub, _ := client.Dubbing().Create(ctx, &elevenlabs.DubbingRequest{
    SourceURL:      "https://example.com/video.mp4",
    TargetLanguage: "es",
    Name:           "Video - Spanish",
})

// Poll for completion
for {
    status, _ := client.Dubbing().GetStatus(ctx, dub.DubbingID)
    if status.Status == "dubbed" {
        break
    }
    time.Sleep(30 * time.Second)
}

// Download dubbed file
audio, _ := client.Dubbing().GetDubbedFile(ctx, dub.DubbingID, "es")
f, _ := os.Create("video_spanish.mp4")
io.Copy(f, audio)
f.Close()
```

## Usage Monitoring

```go
// Check subscription
sub, _ := client.User().GetSubscription(ctx)

fmt.Printf("Tier: %s\n", sub.Tier)
fmt.Printf("Used: %d / %d\n", sub.CharacterCount, sub.CharacterLimit)
fmt.Printf("Remaining: %d\n", sub.CharactersRemaining())

// Pre-generation check
func checkAndGenerate(client *elevenlabs.Client, text string) error {
    sub, _ := client.User().GetSubscription(context.Background())
    if sub.CharactersRemaining() < len(text) {
        return errors.New("insufficient characters")
    }
    // Proceed with generation...
    return nil
}
```

## Error Handling

```go
audio, err := client.TextToSpeech().Simple(ctx, voiceID, text)
if err != nil {
    if elevenlabs.IsRateLimitError(err) {
        log.Println("Rate limited, waiting...")
        time.Sleep(time.Minute)
        // Retry...
    } else if elevenlabs.IsUnauthorizedError(err) {
        log.Fatal("Invalid API key")
    } else if elevenlabs.IsNotFoundError(err) {
        log.Fatal("Voice not found")
    } else {
        log.Fatalf("Error: %v", err)
    }
}
```

## Course Generation Workflow

```go
func generateCourse() {
    client, _ := elevenlabs.NewClient()
    ctx := context.Background()

    // 1. Set up pronunciation
    client.Pronunciation().CreateFromMap(ctx, "Course Terms", map[string]string{
        "API": "A P I",
        "SDK": "S D K",
    })

    // 2. Generate intro
    intro, _ := client.SoundEffects().Simple(ctx, "professional course intro")
    saveAudio(intro, "intro.mp3")

    // 3. Generate chapters
    chapters := []string{
        "Welcome to this course...",
        "In this chapter...",
    }

    for i, text := range chapters {
        audio, _ := client.TextToSpeech().Simple(ctx, voiceID, text)
        saveAudio(audio, fmt.Sprintf("chapter%d.mp3", i+1))
    }
}

func saveAudio(r io.Reader, filename string) {
    f, _ := os.Create(filename)
    defer f.Close()
    io.Copy(f, r)
}
```

## More Examples

See the [examples directory](https://github.com/agentplexus/go-elevenlabs/tree/main/examples) in the repository for complete working examples.
