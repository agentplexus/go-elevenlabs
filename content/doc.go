// Package content provides content generation services including music, dubbing, and projects.
//
// This package combines music composition, video dubbing, and Studio Projects for long-form content.
//
// Example:
//
//	import (
//	    "github.com/plexusone/elevenlabs-go"
//	    "github.com/plexusone/elevenlabs-go/content"
//	)
//
//	client, _ := elevenlabs.NewClient()
//	stems, _ := client.Content().SeparateStems(ctx, &content.StemSeparationRequest{...})
//	projects, _ := client.Content().ListProjects(ctx)
package content
