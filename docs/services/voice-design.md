# Voice Design

Generate custom AI voices with specific characteristics like gender, age, and accent.

!!! note "v0.13.0 API Change"
    As of v0.13.0, Voice Design methods are accessed via `client.Voice()` instead of `client.VoiceDesign()`.

## Basic Usage

```go
import "github.com/plexusone/elevenlabs-go/voice"

// Generate a voice preview
resp, err := client.Voice().DesignSimple(ctx,
    voice.GenderFemale,
    voice.AgeYoung,
    voice.AccentAmerican,
    "This is a preview of the generated voice. It should be at least one hundred characters long for the best quality results.",
)
if err != nil {
    log.Fatal(err)
}

// Listen to the preview
f, _ := os.Create("voice_preview.mp3")
io.Copy(f, resp.Audio)
```

## Full Options

```go
resp, err := client.Voice().Design(ctx, &voice.DesignRequest{
    Gender:         voice.GenderFemale,
    Age:            voice.AgeYoung,
    Accent:         voice.AccentBritish,
    AccentStrength: 1.5,  // 0.3 to 2.0
    Text:           "This is a preview text that must be between one hundred and one thousand characters long for optimal voice generation quality.",
})
```

## Save Generated Voice

Once you like a preview, save it to your voice library:

```go
// Generate preview
preview, _ := client.Voice().Design(ctx, &voice.DesignRequest{
    Gender: voice.GenderMale,
    Age:    voice.AgeMiddleAged,
    Accent: voice.AccentBritish,
    Text:   sampleText,
})

// Save to library
savedVoice, err := client.Voice().SaveDesign(ctx, &voice.SaveDesignRequest{
    GeneratedVoiceID: preview.GeneratedVoiceID,
    VoiceName:        "British Narrator",
    VoiceDescription: "Professional British male voice for narration",
    Labels: map[string]string{
        "use_case": "narration",
        "style":    "professional",
    },
})

fmt.Printf("Saved voice ID: %s\n", savedVoice.VoiceID)
```

## Voice Options

### Gender

```go
voice.GenderFemale
voice.GenderMale
```

### Age

```go
voice.AgeYoung       // Young adult
voice.AgeMiddleAged  // Middle-aged
voice.AgeOld         // Elderly
```

### Accent

```go
voice.AccentAmerican
voice.AccentBritish
voice.AccentAustralian
voice.AccentIndian
voice.AccentAfrican
```

### Accent Strength

| Value | Effect |
|-------|--------|
| 0.3 | Subtle accent |
| 1.0 | Normal (default) |
| 1.5 | Strong accent |
| 2.0 | Very strong accent |

## Request Structure

```go
type DesignRequest struct {
    Gender         Gender  // Required
    Age            Age     // Required
    Accent         Accent  // Required
    AccentStrength float64 // 0.3 to 2.0 (default: 1.0)
    Text           string  // 100-1000 characters
}

type DesignResponse struct {
    Audio            io.Reader // Preview audio
    GeneratedVoiceID string    // ID to save the voice
}

type SaveDesignRequest struct {
    GeneratedVoiceID string            // From preview response
    VoiceName        string            // Name for saved voice
    VoiceDescription string            // Optional description
    Labels           map[string]string // Optional metadata
}
```

## Use Cases

### Create Character Voices

```go
characters := []struct {
    Name   string
    Gender voice.Gender
    Age    voice.Age
    Accent voice.Accent
}{
    {"Hero", voice.GenderMale, voice.AgeYoung, voice.AccentAmerican},
    {"Mentor", voice.GenderMale, voice.AgeOld, voice.AccentBritish},
    {"Sidekick", voice.GenderFemale, voice.AgeYoung, voice.AccentAustralian},
}

sampleText := "This is a sample of how this character will sound in the story. The text needs to be long enough for quality generation."

for _, char := range characters {
    preview, _ := client.Voice().DesignSimple(ctx, char.Gender, char.Age, char.Accent, sampleText)

    // Save if satisfied
    savedVoice, _ := client.Voice().SaveDesign(ctx, &voice.SaveDesignRequest{
        GeneratedVoiceID: preview.GeneratedVoiceID,
        VoiceName:        char.Name,
        Labels:           map[string]string{"project": "audiobook"},
    })

    fmt.Printf("Created %s voice: %s\n", char.Name, savedVoice.VoiceID)
}
```

### A/B Test Voices

```go
// Generate multiple previews with same parameters
var previews []*voice.DesignResponse

for i := 0; i < 3; i++ {
    preview, _ := client.Voice().Design(ctx, &voice.DesignRequest{
        Gender: voice.GenderFemale,
        Age:    voice.AgeYoung,
        Accent: voice.AccentAmerican,
        Text:   sampleText,
    })
    previews = append(previews, preview)

    // Save preview audio for comparison
    f, _ := os.Create(fmt.Sprintf("preview_%d.mp3", i))
    io.Copy(f, preview.Audio)
    f.Close()
}

// Listen to all previews and save the best one
```

### Brand Voice Creation

```go
// Define brand voice characteristics
brandVoice := voice.DesignRequest{
    Gender:         voice.GenderFemale,
    Age:            voice.AgeMiddleAged,
    Accent:         voice.AccentAmerican,
    AccentStrength: 0.5,  // Subtle accent for professionalism
    Text:           "Welcome to our service. We're here to help you succeed. Our team is dedicated to providing the best experience possible for all our customers.",
}

preview, _ := client.Voice().Design(ctx, &brandVoice)

savedVoice, _ := client.Voice().SaveDesign(ctx, &voice.SaveDesignRequest{
    GeneratedVoiceID: preview.GeneratedVoiceID,
    VoiceName:        "Brand Voice - Main",
    VoiceDescription: "Official brand voice for customer communications",
    Labels: map[string]string{
        "brand":    "true",
        "approved": "true",
        "use_case": "customer_service",
    },
})
```

## Text Requirements

The preview text must be:

- **Minimum:** 100 characters
- **Maximum:** 1,000 characters
- **Content:** Representative of intended use

Good preview text:

```go
text := `This is a sample of how this voice will sound when reading longer content.
The text should be representative of the actual content you plan to generate,
including the style, tone, and type of vocabulary you'll be using.`
```

## Best Practices

1. **Use representative text** - Preview text should match your intended use case
2. **Generate multiple previews** - Each generation is unique; try several
3. **Test accent strength** - Adjust for natural-sounding results
4. **Add descriptive labels** - Makes organizing voices easier
5. **Save good voices immediately** - Generated voice IDs may expire
