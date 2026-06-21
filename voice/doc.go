// Package voice provides voice management and design services.
//
// This package combines voice listing, settings, and AI voice design/generation.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	voices, _ := client.Voice().List(ctx)
//	settings, _ := client.Voice().GetSettings(ctx, "voice_id")
package voice
