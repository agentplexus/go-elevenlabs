package elevenlabs

import (
	"context"
	"errors"
	"io"

	ht "github.com/ogen-go/ogen/http"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// MusicService handles music composition and generation.
type MusicService struct {
	client *Client
}

// ErrMusicCompositionUnavailable is returned when music composition endpoints
// are called but the underlying API client doesn't support them.
// This happens when the OpenAPI spec uses complex anyOf schemas that ogen cannot generate.
var ErrMusicCompositionUnavailable = errors.New("music composition endpoints are temporarily unavailable due to API schema complexity; use SeparateStems or VideoToMusic instead")

// MusicRequest contains options for music generation.
type MusicRequest struct {
	// Prompt is a simple text description of the music to generate.
	// Cannot be used with CompositionPlan.
	Prompt string

	// DurationMs is the length of the song in milliseconds (3000-600000).
	// If not provided, the model will choose based on the prompt.
	DurationMs int

	// ForceInstrumental ensures the song has no vocals.
	ForceInstrumental bool

	// Seed for deterministic generation (optional).
	Seed int
}

// MusicResponse contains the music generation result.
type MusicResponse struct {
	// Audio is the generated music.
	Audio io.Reader

	// SongID is the unique identifier for this song.
	SongID string
}

// Generate creates music from a text prompt.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
// The ElevenLabs API uses a union type (anyOf) that the code generator cannot handle.
// Use SeparateStems or VideoToMusic instead.
func (s *MusicService) Generate(ctx context.Context, req *MusicRequest) (*MusicResponse, error) {
	return nil, ErrMusicCompositionUnavailable
}

// GenerateStream creates music with streaming output.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *MusicService) GenerateStream(ctx context.Context, req *MusicRequest) (*MusicResponse, error) {
	return nil, ErrMusicCompositionUnavailable
}

// Simple generates music from a prompt with default settings.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *MusicService) Simple(ctx context.Context, prompt string) (io.Reader, error) {
	return nil, ErrMusicCompositionUnavailable
}

// GenerateInstrumental generates instrumental music from a prompt.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *MusicService) GenerateInstrumental(ctx context.Context, prompt string, durationMs int) (io.Reader, error) {
	return nil, ErrMusicCompositionUnavailable
}

// CompositionPlan represents a detailed music composition plan.
// This can be used with GenerateDetailed for fine-grained control over music generation.
type CompositionPlan struct {
	// PositiveGlobalStyles are styles that should be present throughout the song.
	PositiveGlobalStyles []string

	// NegativeGlobalStyles are styles that should NOT be present in the song.
	NegativeGlobalStyles []string

	// Sections defines the structure of the song with individual sections.
	Sections []SongSection
}

// SongSection represents a section of a song in a composition plan.
type SongSection struct {
	// SectionName is the name of the section (e.g., "intro", "verse", "chorus").
	SectionName string

	// DurationMs is the duration in milliseconds (3000-120000).
	DurationMs int

	// Lines are the lyrics for this section (max 200 chars per line).
	Lines []string

	// PositiveLocalStyles are styles for this specific section.
	PositiveLocalStyles []string

	// NegativeLocalStyles are styles to avoid in this section.
	NegativeLocalStyles []string
}

// CompositionPlanRequest contains options for generating a composition plan.
type CompositionPlanRequest struct {
	// Prompt is the text description of the music to plan.
	Prompt string

	// DurationMs is the target duration in milliseconds (3000-600000).
	DurationMs int

	// SourcePlan is an optional existing plan to use as a starting point.
	SourcePlan *CompositionPlan
}

// GeneratePlan creates a composition plan from a text prompt.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *MusicService) GeneratePlan(ctx context.Context, req *CompositionPlanRequest) (*CompositionPlan, error) {
	return nil, ErrMusicCompositionUnavailable
}

// MusicDetailedRequest contains options for detailed music generation.
type MusicDetailedRequest struct {
	// Prompt is a simple text description (cannot be used with CompositionPlan).
	Prompt string

	// CompositionPlan is a detailed plan (cannot be used with Prompt).
	CompositionPlan *CompositionPlan

	// DurationMs is the length in milliseconds (only used with Prompt).
	DurationMs int

	// ForceInstrumental ensures no vocals (only used with Prompt).
	ForceInstrumental bool

	// Seed for deterministic generation.
	Seed int

	// WithTimestamps returns word timestamps in the response.
	WithTimestamps bool
}

// MusicDetailedResponse contains the detailed music generation result.
type MusicDetailedResponse struct {
	// Audio is the generated music.
	Audio io.Reader

	// SongID is the unique identifier for this song.
	SongID string
}

// GenerateDetailed creates music with detailed options and metadata.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *MusicService) GenerateDetailed(ctx context.Context, req *MusicDetailedRequest) (*MusicDetailedResponse, error) {
	return nil, ErrMusicCompositionUnavailable
}

// StemSeparationRequest contains options for stem separation.
type StemSeparationRequest struct {
	// File is the audio file to separate.
	File io.Reader

	// Filename is the name of the file.
	Filename string

	// StemVariation specifies which stem variation to use.
	// Options: "two_stems_v1" (vocals + music), "six_stems_v1" (vocals, drums, bass, other - default)
	StemVariation string
}

// SeparateStems separates a song into individual stems (vocals, instruments, etc.).
//
// Example:
//
//	f, _ := os.Open("song.mp3")
//	stems, err := client.Music().SeparateStems(ctx, &StemSeparationRequest{
//	    File:     f,
//	    Filename: "song.mp3",
//	})
//	// Save the separated stems (returned as a zip file)
//	output, _ := os.Create("stems.zip")
//	io.Copy(output, stems)
func (s *MusicService) SeparateStems(ctx context.Context, req *StemSeparationRequest) (io.Reader, error) {
	if req.File == nil {
		return nil, &ValidationError{Field: "file", Message: "cannot be nil"}
	}
	if req.Filename == "" {
		return nil, &ValidationError{Field: "filename", Message: "cannot be empty"}
	}

	body := &api.BodyStemSeparationV1MusicStemSeparationPostMultipart{
		File: ht.MultipartFile{
			Name: req.Filename,
			File: req.File,
		},
	}

	if req.StemVariation != "" {
		body.StemVariationID = api.NewOptBodyStemSeparationV1MusicStemSeparationPostMultipartStemVariationID(
			api.BodyStemSeparationV1MusicStemSeparationPostMultipartStemVariationID(req.StemVariation))
	}

	resp, err := s.client.apiClient.SeparateSongStems(ctx, body, api.SeparateSongStemsParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.SeparateSongStemsOKHeaders:
		return r.Response.Data, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// VideoToMusicRequest contains options for converting a video to music.
type VideoToMusicRequest struct {
	// Videos is one or more video files to convert.
	Videos []VideoFile

	// Description is an optional text description of the music.
	Description string

	// Tags are optional style tags (e.g., ["upbeat", "cinematic"]).
	Tags []string

	// ModelID specifies the model to use (optional).
	ModelID string
}

// VideoFile represents a video file for VideoToMusic.
type VideoFile struct {
	// Name is the filename.
	Name string
	// File is the video content.
	File io.Reader
}

// VideoToMusic converts a video to music.
//
// Example:
//
//	f, _ := os.Open("video.mp4")
//	music, err := client.Music().VideoToMusic(ctx, &VideoToMusicRequest{
//	    Videos: []VideoFile{{Name: "video.mp4", File: f}},
//	})
func (s *MusicService) VideoToMusic(ctx context.Context, req *VideoToMusicRequest) (io.Reader, error) {
	if len(req.Videos) == 0 {
		return nil, &ValidationError{Field: "videos", Message: "at least one video is required"}
	}

	var videos []ht.MultipartFile
	for _, v := range req.Videos {
		if v.File == nil {
			return nil, &ValidationError{Field: "videos", Message: "video file cannot be nil"}
		}
		if v.Name == "" {
			return nil, &ValidationError{Field: "videos", Message: "video filename cannot be empty"}
		}
		videos = append(videos, ht.MultipartFile{
			Name: v.Name,
			File: v.File,
		})
	}

	body := &api.BodyVideoToMusicV1MusicVideoToMusicPostMultipart{
		Videos: videos,
		Tags:   req.Tags,
	}

	if req.Description != "" {
		body.Description = api.NewOptNilString(req.Description)
	}

	if req.ModelID != "" {
		body.ModelID = api.NewOptBodyVideoToMusicV1MusicVideoToMusicPostMultipartModelID(
			api.BodyVideoToMusicV1MusicVideoToMusicPostMultipartModelID(req.ModelID))
	}

	resp, err := s.client.apiClient.VideoToMusic(ctx, body, api.VideoToMusicParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.VideoToMusicOKHeaders:
		return r.Response.Data, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}
