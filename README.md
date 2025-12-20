# ElevenLabs Go SDK

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

Go SDK for the [ElevenLabs API](https://elevenlabs.io/).

## Features

- 🗣️ **Text-to-Speech**: Convert text to realistic speech with multiple voices and models
- 📝 **Speech-to-Text**: Transcribe audio with speaker diarization support
- 🔊 **Sound Effects**: Generate sound effects from text descriptions
- 🎨 **Voice Design**: Create custom AI voices with specific characteristics
- 🎵 **Music Composition**: Generate music from text prompts
- 🎙️ **Audio Isolation**: Extract vocals/speech from audio
- ⏱️ **Forced Alignment**: Get word-level timestamps for audio
- 💬 **Text-to-Dialogue**: Generate multi-speaker conversations
- 🌍 **Dubbing**: Translate and dub video/audio content
- 📚 **Projects**: Manage long-form audio content (audiobooks, podcasts)
- 📖 **Pronunciation Dictionaries**: Control pronunciation of specific terms

## Installation

```bash
go get github.com/grokify/go-elevenlabs
```

## Quick Start

### Basic Text-to-Speech

```go
package main

import (
    "context"
    "io"
    "log"
    "os"

    elevenlabs "github.com/grokify/go-elevenlabs"
)

func main() {
    // Create client (uses ELEVENLABS_API_KEY env var)
    client, err := elevenlabs.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List available voices
    voices, err := client.Voices().List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Found %d voices", len(voices))

    // Generate speech
    if len(voices) > 0 {
        audio, err := client.TextToSpeech().Simple(ctx,
            voices[0].VoiceID,
            "Hello from the ElevenLabs Go SDK!")
        if err != nil {
            log.Fatal(err)
        }

        // Save to file
        f, _ := os.Create("hello.mp3")
        defer f.Close()
        io.Copy(f, audio)
    }
}
```

### With Custom Options

```go
client, err := elevenlabs.NewClient(
    elevenlabs.WithAPIKey("your-api-key"),
    elevenlabs.WithTimeout(5 * time.Minute),
)
```

## Services

### Text-to-Speech

```go
// Simple generation
audio, err := client.TextToSpeech().Simple(ctx, voiceID, "Hello world")

// With full options
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
```

### Speech-to-Text

```go
// Transcribe from URL
result, err := client.SpeechToText().TranscribeURL(ctx, "https://example.com/audio.mp3")
fmt.Printf("Text: %s\n", result.Text)
fmt.Printf("Language: %s\n", result.LanguageCode)

// With speaker diarization
result, err := client.SpeechToText().TranscribeWithDiarization(ctx, audioURL)
for _, word := range result.Words {
    fmt.Printf("[%s] %s (%.2fs - %.2fs)\n", word.Speaker, word.Text, word.Start, word.End)
}
```

### Sound Effects

```go
// Simple sound effect
audio, err := client.SoundEffects().Simple(ctx, "thunder and rain storm")

// With options
sfx, err := client.SoundEffects().Generate(ctx, &elevenlabs.SoundEffectRequest{
    Text:            "spaceship engine humming",
    DurationSeconds: 10,
    PromptInfluence: 0.5,
})
```

### Music Composition

```go
// Generate music from prompt
resp, err := client.Music().Generate(ctx, &elevenlabs.MusicRequest{
    Prompt:     "upbeat electronic music for a tech video",
    DurationMs: 30000,
})

// Instrumental only
audio, err := client.Music().GenerateInstrumental(ctx, "calm piano melody", 60000)

// Generate with composition plan for fine-grained control
plan, _ := client.Music().GeneratePlan(ctx, &elevenlabs.CompositionPlanRequest{
    Prompt:     "pop song about summer",
    DurationMs: 180000,
})
resp, err := client.Music().GenerateDetailed(ctx, &elevenlabs.MusicDetailedRequest{
    CompositionPlan: plan,
})

// Separate stems (vocals, drums, bass, etc.)
f, _ := os.Open("song.mp3")
stems, err := client.Music().SeparateStems(ctx, &elevenlabs.StemSeparationRequest{
    File:     f,
    Filename: "song.mp3",
})
```

### Audio Isolation

```go
// Extract vocals from audio file
f, _ := os.Open("mixed_audio.mp3")
isolated, err := client.AudioIsolation().IsolateFile(ctx, f, "mixed_audio.mp3")
```

### Forced Alignment

```go
// Get word-level timestamps
f, _ := os.Open("speech.mp3")
result, err := client.ForcedAlignment().AlignFile(ctx, f, "speech.mp3",
    "The text that was spoken in the audio")

for _, word := range result.Words {
    fmt.Printf("%s: %.2fs - %.2fs\n", word.Text, word.Start, word.End)
}
```

### Text-to-Dialogue

```go
// Generate multi-speaker dialogue
audio, err := client.TextToDialogue().Simple(ctx, []elevenlabs.DialogueInput{
    {Text: "Hello, how are you?", VoiceID: "voice1"},
    {Text: "I'm doing great, thanks!", VoiceID: "voice2"},
})
```

### Voice Design

```go
// Generate a custom voice
resp, err := client.VoiceDesign().GeneratePreview(ctx, &elevenlabs.VoiceDesignRequest{
    Gender:         elevenlabs.VoiceGenderFemale,
    Age:            elevenlabs.VoiceAgeYoung,
    Accent:         elevenlabs.VoiceAccentAmerican,
    AccentStrength: 1.0,
    Text:           "This is a preview of the generated voice. It should be at least one hundred characters long for best results.",
})
```

### Pronunciation Dictionaries

```go
// Create from a map
dict, err := client.Pronunciation().CreateFromMap(ctx, "Tech Terms", map[string]string{
    "API":     "A P I",
    "kubectl": "kube control",
    "nginx":   "engine X",
})

// Create from JSON file
dict, err := client.Pronunciation().CreateFromJSON(ctx, "Terms", "pronunciation.json")
```

### Dubbing

```go
// Create dubbing job
dub, err := client.Dubbing().Create(ctx, &elevenlabs.DubbingRequest{
    SourceURL:      "https://example.com/video.mp4",
    TargetLanguage: "es",
    Name:           "Video - Spanish",
})

// Check status
status, err := client.Dubbing().GetStatus(ctx, dub.DubbingID)
```

### Projects (Studio)

```go
// Create a project for long-form content
project, err := client.Projects().Create(ctx, &elevenlabs.CreateProjectRequest{
    Name:                    "My Audiobook",
    DefaultModelID:          "eleven_multilingual_v2",
    DefaultParagraphVoiceID: voiceID,
})

// Convert to audio
err = client.Projects().Convert(ctx, project.ProjectID)
```

## Error Handling

```go
audio, err := client.TextToSpeech().Simple(ctx, voiceID, text)
if err != nil {
    if elevenlabs.IsRateLimitError(err) {
        log.Println("Rate limited, waiting...")
        time.Sleep(time.Minute)
    } else if elevenlabs.IsUnauthorizedError(err) {
        log.Fatal("Invalid API key")
    } else if elevenlabs.IsNotFoundError(err) {
        log.Fatal("Voice not found")
    } else {
        log.Fatalf("Error: %v", err)
    }
}
```

## Environment Variables

- `ELEVENLABS_API_KEY`: Your ElevenLabs API key (used automatically if not provided via `WithAPIKey`)

## Documentation

- [API Reference](https://grokify.github.io/go-elevenlabs/)
- [ElevenLabs API Docs](https://elevenlabs.io/docs)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License

 [build-status-svg]: https://github.com/grokify/go-elevenlabs/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/go-elevenlabs/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/go-elevenlabs/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/go-elevenlabs/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/go-elevenlabs
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/go-elevenlabs
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/go-elevenlabs
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/go-elevenlabs
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/go-elevenlabs/blob/master/LICENSE
