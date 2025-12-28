# go-elevenlabs

A Go SDK for the [ElevenLabs](https://elevenlabs.io/) Text-to-Speech API.

## Features

- **Text-to-Speech** - Generate natural-sounding speech from text
- **Voice Selection** - Access pre-made and cloned voices
- **Sound Effects** - Generate sound effects from text descriptions
- **Projects (Studio)** - Create long-form content with chapters
- **Pronunciation Dictionaries** - Ensure correct pronunciation of technical terms
- **Dubbing** - Translate audio/video to other languages
- **History** - Access and manage generated audio history

## Installation

```bash
go get github.com/agentplexus/go-elevenlabs
```

## Quick Example

```go
package main

import (
    "context"
    "io"
    "os"

    elevenlabs "github.com/agentplexus/go-elevenlabs"
)

func main() {
    // Create client (uses ELEVENLABS_API_KEY env var)
    client, _ := elevenlabs.NewClient()
    ctx := context.Background()

    // Generate speech
    audio, _ := client.TextToSpeech().Simple(ctx,
        "21m00Tcm4TlvDq8ikWAM",  // Voice ID
        "Hello, welcome to ElevenLabs!")

    // Save to file
    f, _ := os.Create("output.mp3")
    defer f.Close()
    io.Copy(f, audio)
}
```

## Use Cases

This SDK is particularly well-suited for:

- **Online Courses** - Generate professional narration for Udemy, LMS platforms
- **Audiobooks** - Create chapter-organized audio content
- **Podcasts** - Produce consistent, high-quality audio
- **Video Production** - Add voiceovers and sound effects
- **Accessibility** - Convert text content to audio format

## Documentation

- [Getting Started](getting-started/installation.md) - Installation and setup
- [Services](services/text-to-speech.md) - API service documentation
- [Guides](guides/lms-courses.md) - Use case guides
- [Examples](examples.md) - Code examples

## License

MIT License - see [LICENSE](https://github.com/agentplexus/go-elevenlabs/blob/main/LICENSE) for details.
