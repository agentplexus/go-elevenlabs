# OmniVoice Integration

OmniVoice is a vendor-agnostic abstraction layer for voice AI services. The `omnivoice/` subpackage provides ElevenLabs implementations of the OmniVoice interfaces, allowing you to use ElevenLabs as a backend while keeping your application code portable.

## What is OmniVoice?

[OmniVoice](https://github.com/plexusone/omnivoice-core) defines standard interfaces for:

- **Text-to-Speech (TTS)** - Convert text to spoken audio
- **Speech-to-Text (STT)** - Transcribe audio to text
- **Voice Agents** - Interactive voice sessions with real-time audio

By coding against OmniVoice interfaces, you can swap voice providers (ElevenLabs, OpenAI, Google, etc.) without changing application logic.

## Installation

The OmniVoice providers are included in the go-elevenlabs SDK:

```bash
go get github.com/plexusone/elevenlabs-go
```

You'll also need the OmniVoice interfaces:

```bash
go get github.com/plexusone/omnivoice-core
```

## Quick Start

### TTS Provider

```go
import (
    "github.com/plexusone/omnivoice-core/tts"
    eleventts "github.com/plexusone/elevenlabs-go/omnivoice/tts"
)

// Create provider (uses ELEVENLABS_API_KEY env var)
provider, err := eleventts.New()
if err != nil {
    log.Fatal(err)
}

// Use with OmniVoice client
client := tts.NewClient(provider)
result, err := client.Synthesize(ctx, "Hello world", tts.SynthesisConfig{
    VoiceID: "21m00Tcm4TlvDq8ikWAM", // Rachel
})
```

### STT Provider

```go
import (
    "github.com/plexusone/omnivoice-core/stt"
    elevenstt "github.com/plexusone/elevenlabs-go/omnivoice/stt"
)

// Create provider
provider, err := elevenstt.New()
if err != nil {
    log.Fatal(err)
}

// Transcribe audio
result, err := provider.Transcribe(ctx, audioBytes, stt.TranscriptionConfig{
    Language: "en",
})
fmt.Println(result.Text)
```

### Agent Provider

```go
import (
    "github.com/plexusone/omnivoice-core/agent"
    elevenagent "github.com/plexusone/elevenlabs-go/omnivoice/agent"
)

// Create provider
provider, err := elevenagent.New()
if err != nil {
    log.Fatal(err)
}

// Create and start session
session, err := provider.CreateSession(ctx, agent.Config{
    VoiceID:  "21m00Tcm4TlvDq8ikWAM",
    Language: "en",
})
if err != nil {
    log.Fatal(err)
}

if err := session.Start(ctx); err != nil {
    log.Fatal(err)
}
defer session.Stop(ctx)

// Handle events
for event := range session.Events() {
    switch event.Type {
    case agent.EventUserTranscript:
        fmt.Printf("User: %s\n", event.Data)
    case agent.EventAgentTranscript:
        fmt.Printf("Agent: %s\n", event.Data)
    }
}
```

## Configuration

All providers accept optional configuration:

```go
// With explicit API key
provider, err := eleventts.New(
    eleventts.WithAPIKey("your-api-key"),
)

// With custom base URL (for proxies or testing)
provider, err := eleventts.New(
    eleventts.WithBaseURL("https://custom-api.example.com"),
)
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ELEVENLABS_API_KEY` | ElevenLabs API key (used by default) |

## When to Use OmniVoice

Use OmniVoice when you want:

- **Vendor portability** - Switch providers without code changes
- **Consistent API** - Same interface across all voice services
- **Testing flexibility** - Mock providers for unit tests

Use the SDK directly when you need:

- **ElevenLabs-specific features** - Voice cloning, pronunciation dictionaries, projects
- **Maximum control** - Direct access to all API parameters
- **Performance optimization** - Avoid abstraction overhead

## Provider Details

- [Capabilities](capabilities.md) - Full capability matrix and conformance test status
- [TTS Provider](tts.md) - Text-to-speech synthesis
- [STT Provider](stt.md) - Speech-to-text transcription
- [Agent Provider](agent.md) - Interactive voice sessions
