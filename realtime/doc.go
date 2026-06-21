// Package realtime provides WebSocket-based real-time TTS and STT services.
//
// This package provides low-latency streaming capabilities for text-to-speech
// and speech-to-text using WebSocket connections.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/realtime"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	conn, _ := client.Realtime().ConnectTTS(ctx, "voice_id", &realtime.TTSOptions{...})
//	conn.SendText("Hello, world!")
package realtime
