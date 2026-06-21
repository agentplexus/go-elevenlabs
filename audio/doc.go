// Package audio provides audio processing services including isolation, alignment, and sound effects.
//
// This package combines audio isolation (vocal extraction), forced alignment, and sound effect generation.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/audio"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	isolated, _ := client.Audio().Isolate(ctx, &audio.IsolationRequest{...})
//	sfx, _ := client.Audio().GenerateSoundEffect(ctx, &audio.SoundEffectRequest{...})
package audio
