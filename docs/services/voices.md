# Voices

The Voices service manages voice selection and settings.

!!! note "v0.13.0 API Change"
    As of v0.13.0, Voice methods are accessed via `client.Voice()` instead of `client.Voices()`.

## List All Voices

```go
voices, err := client.Voice().List(ctx)
if err != nil {
    log.Fatal(err)
}

for _, v := range voices {
    fmt.Printf("%s: %s (%s)\n", v.VoiceID, v.Name, v.Category)
}
```

## Get a Specific Voice

```go
voice, err := client.Voice().Get(ctx, "21m00Tcm4TlvDq8ikWAM")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", voice.Name)
fmt.Printf("Category: %s\n", voice.Category)
fmt.Printf("Description: %s\n", voice.Description)
fmt.Printf("Labels: %v\n", voice.Labels)
```

## Voice Object

| Field | Type | Description |
|-------|------|-------------|
| `VoiceID` | string | Unique identifier |
| `Name` | string | Display name |
| `Category` | string | `premade`, `cloned`, `generated` |
| `Description` | string | Voice description |
| `Labels` | map | Metadata (accent, age, gender, etc.) |
| `PreviewURL` | string | URL to preview audio |

## Get Voice Settings

```go
settings, err := client.Voice().GetSettings(ctx, voiceID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Stability: %f\n", settings.Stability)
fmt.Printf("Similarity Boost: %f\n", settings.SimilarityBoost)
```

## Get Default Settings

```go
defaults, err := client.Voice().GetDefaultSettings(ctx)
```

## Popular Pre-made Voices

| Voice ID | Name | Description |
|----------|------|-------------|
| `21m00Tcm4TlvDq8ikWAM` | Rachel | Calm, narration |
| `AZnzlk1XvdvUeBnXmlld` | Domi | Strong, confident |
| `EXAVITQu4vr4xnSDxMaL` | Bella | Soft, gentle |
| `ErXwobaYiN019PkySvjV` | Antoni | Well-rounded |
| `MF3mGyEYCl7XYWbV9V6O` | Elli | Emotional range |

## Finding the Right Voice

```go
voices, _ := client.Voice().List(ctx)

// Filter by category
for _, v := range voices {
    if v.Category == "premade" {
        // Pre-made voices
    }
}

// Filter by labels
for _, v := range voices {
    if accent, ok := v.Labels["accent"]; ok && accent == "american" {
        // American accent voices
    }
}
```

## Voice Selection Tips

1. **For narration**: Use calm, neutral voices (Rachel, Antoni)
2. **For characters**: Match voice personality to character
3. **For multilingual**: Check voice language support
4. **Test first**: Use preview URLs before committing
