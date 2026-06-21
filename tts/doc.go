// Package tts provides text-to-speech, speech-to-speech, and dialogue generation services.
//
// This package combines the TTS, STS, and dialogue services for audio generation.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/tts"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	resp, _ := client.TTS().Generate(ctx, &tts.Request{
//	    VoiceID: "EXAVITQu4vr4xnSDxMaL",
//	    Text:    "Hello, world!",
//	})
package tts
