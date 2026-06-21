// Package stt provides speech-to-text transcription services.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/stt"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	resp, _ := client.STT().Transcribe(ctx, &stt.Request{
//	    FileURL: "https://example.com/audio.mp3",
//	})
package stt
