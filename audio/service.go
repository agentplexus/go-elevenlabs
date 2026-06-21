package audio

import (
	"context"
	"fmt"
	"io"

	ht "github.com/ogen-go/ogen/http"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Service handles audio processing including isolation, alignment, and sound effects.
type Service struct {
	apiClient *api.Client
}

// New creates a new audio service.
func New(apiClient *api.Client) *Service {
	return &Service{apiClient: apiClient}
}

// --- Audio Isolation ---

// IsolationRequest contains options for audio isolation.
type IsolationRequest struct {
	// Audio is the audio file to process (required).
	Audio io.Reader

	// Filename is the name of the file (required).
	Filename string
}

// Isolate extracts vocals/speech from audio, removing background noise.
// Returns an io.Reader containing the isolated audio.
func (s *Service) Isolate(ctx context.Context, req *IsolationRequest) (io.Reader, error) {
	if req.Audio == nil {
		return nil, fmt.Errorf("audio: audio cannot be nil")
	}

	body := &api.BodyAudioIsolationV1AudioIsolationPostMultipart{
		Audio: ht.MultipartFile{
			Name: req.Filename,
			File: req.Audio,
		},
	}

	resp, err := s.apiClient.AudioIsolation(ctx, body, api.AudioIsolationParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AudioIsolationOK:
		return r.Data, nil
	default:
		return nil, fmt.Errorf("audio: unexpected response type")
	}
}

// IsolateFile is a convenience method to isolate vocals from an audio file.
func (s *Service) IsolateFile(ctx context.Context, audio io.Reader, filename string) (io.Reader, error) {
	return s.Isolate(ctx, &IsolationRequest{
		Audio:    audio,
		Filename: filename,
	})
}

// IsolateStream extracts vocals/speech from audio with streaming output.
// Returns an io.Reader for streaming the isolated audio.
func (s *Service) IsolateStream(ctx context.Context, req *IsolationRequest) (io.Reader, error) {
	if req.Audio == nil {
		return nil, fmt.Errorf("audio: audio cannot be nil")
	}

	body := &api.BodyAudioIsolationStreamV1AudioIsolationStreamPostMultipart{
		Audio: ht.MultipartFile{
			Name: req.Filename,
			File: req.Audio,
		},
	}

	resp, err := s.apiClient.AudioIsolationStream(ctx, body, api.AudioIsolationStreamParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AudioIsolationStreamOK:
		return r.Data, nil
	default:
		return nil, fmt.Errorf("audio: unexpected response type")
	}
}

// --- Forced Alignment ---

// AlignmentRequest contains options for forced alignment.
type AlignmentRequest struct {
	// File is the audio file to align (required).
	File io.Reader

	// Filename is the name of the file (required).
	Filename string

	// Text is the text to align with the audio (required).
	Text string
}

// AlignmentResponse contains the alignment result.
type AlignmentResponse struct {
	// Words contains word-level timing information.
	Words []AlignmentWord

	// Characters contains character-level timing information.
	Characters []AlignmentCharacter

	// Loss is the average alignment confidence score.
	Loss float64
}

// AlignmentWord represents a word with timing information.
type AlignmentWord struct {
	// Text is the word text.
	Text string

	// Start is the start time in seconds.
	Start float64

	// End is the end time in seconds.
	End float64

	// Loss is the confidence score for this word.
	Loss float64
}

// AlignmentCharacter represents a character with timing information.
type AlignmentCharacter struct {
	// Text is the character text.
	Text string

	// Start is the start time in seconds.
	Start float64

	// End is the end time in seconds.
	End float64
}

// Align performs forced alignment between audio and text.
// This is useful for generating word-level timestamps for captions and subtitles.
func (s *Service) Align(ctx context.Context, req *AlignmentRequest) (*AlignmentResponse, error) {
	if req.File == nil {
		return nil, fmt.Errorf("audio: file cannot be nil")
	}
	if req.Text == "" {
		return nil, fmt.Errorf("audio: text cannot be empty")
	}

	body := &api.BodyCreateForcedAlignmentV1ForcedAlignmentPostMultipart{
		File: ht.MultipartFile{
			Name: req.Filename,
			File: req.File,
		},
		Text: req.Text,
	}

	resp, err := s.apiClient.ForcedAlignment(ctx, body, api.ForcedAlignmentParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.ForcedAlignmentResponseModel:
		result := &AlignmentResponse{
			Loss: r.Loss,
		}

		// Convert words
		for _, w := range r.Words {
			result.Words = append(result.Words, AlignmentWord{
				Text:  w.Text,
				Start: w.Start,
				End:   w.End,
				Loss:  w.Loss,
			})
		}

		// Convert characters
		for _, c := range r.Characters {
			result.Characters = append(result.Characters, AlignmentCharacter{
				Text:  c.Text,
				Start: c.Start,
				End:   c.End,
			})
		}

		return result, nil
	default:
		return nil, fmt.Errorf("audio: unexpected response type")
	}
}

// AlignFile is a convenience method to align audio from a file reader with text.
func (s *Service) AlignFile(ctx context.Context, file io.Reader, filename, text string) (*AlignmentResponse, error) {
	return s.Align(ctx, &AlignmentRequest{
		File:     file,
		Filename: filename,
		Text:     text,
	})
}

// --- Sound Effects ---

// SoundEffectRequest contains options for generating a sound effect.
type SoundEffectRequest struct {
	// Text is the description of the sound effect to generate.
	// Examples: "car engine starting", "thunder and rain", "crowd cheering"
	Text string

	// DurationSeconds is the target duration (0.5 to 30 seconds).
	// If not set, the optimal duration will be guessed from the prompt.
	DurationSeconds float64

	// PromptInfluence controls how closely the generation follows the prompt (0.0 to 1.0).
	// Higher values = more faithful to prompt but less variation.
	// Default is 0.3.
	PromptInfluence float64

	// Loop creates a sound effect that loops smoothly.
	Loop bool

	// OutputFormat specifies the audio format (e.g., "mp3_44100_128").
	OutputFormat string
}

// Validate validates the sound effect request.
func (r *SoundEffectRequest) Validate() error {
	if r.Text == "" {
		return fmt.Errorf("audio: text cannot be empty")
	}
	if r.DurationSeconds != 0 && (r.DurationSeconds < 0.5 || r.DurationSeconds > 30) {
		return fmt.Errorf("audio: duration_seconds must be between 0.5 and 30")
	}
	if r.PromptInfluence != 0 && (r.PromptInfluence < 0 || r.PromptInfluence > 1) {
		return fmt.Errorf("audio: prompt_influence must be between 0 and 1")
	}
	return nil
}

// SoundEffectResponse contains the generated sound effect.
type SoundEffectResponse struct {
	// Audio is the generated sound effect data.
	Audio io.Reader
}

// GenerateSoundEffect creates a sound effect from a text description.
func (s *Service) GenerateSoundEffect(ctx context.Context, req *SoundEffectRequest) (*SoundEffectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	body := &api.BodySoundGenerationV1SoundGenerationPost{
		Text: req.Text,
	}

	if req.DurationSeconds > 0 {
		body.DurationSeconds = api.NewOptNilFloat64(req.DurationSeconds)
	}
	if req.PromptInfluence > 0 {
		body.PromptInfluence = api.NewOptNilFloat64(req.PromptInfluence)
	}
	if req.Loop {
		body.Loop = api.NewOptBool(true)
	}

	params := api.SoundGenerationParams{}
	if req.OutputFormat != "" {
		params.OutputFormat = api.NewOptAllowedOutputFormats(
			api.AllowedOutputFormats(req.OutputFormat),
		)
	}

	resp, err := s.apiClient.SoundGeneration(ctx, body, params)
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.SoundGenerationOKHeaders:
		return &SoundEffectResponse{Audio: r.Response.Data}, nil
	default:
		return nil, fmt.Errorf("audio: unexpected response type")
	}
}

// SimpleSoundEffect generates a sound effect with minimal configuration.
func (s *Service) SimpleSoundEffect(ctx context.Context, description string) (io.Reader, error) {
	resp, err := s.GenerateSoundEffect(ctx, &SoundEffectRequest{
		Text: description,
	})
	if err != nil {
		return nil, err
	}
	return resp.Audio, nil
}

// GenerateLoopingSoundEffect generates a looping sound effect.
func (s *Service) GenerateLoopingSoundEffect(ctx context.Context, description string, durationSeconds float64) (io.Reader, error) {
	resp, err := s.GenerateSoundEffect(ctx, &SoundEffectRequest{
		Text:            description,
		DurationSeconds: durationSeconds,
		Loop:            true,
	})
	if err != nil {
		return nil, err
	}
	return resp.Audio, nil
}
