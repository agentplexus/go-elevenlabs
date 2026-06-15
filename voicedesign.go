package elevenlabs

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// VoiceDesignService handles AI voice generation and design.
type VoiceDesignService struct {
	client *Client
}

// ErrVoiceDesignLegacyUnavailable is returned when legacy voice design methods are called.
// The ElevenLabs API now uses text-based voice descriptions instead of gender/age/accent parameters.
var ErrVoiceDesignLegacyUnavailable = errors.New("legacy voice design API (gender/age/accent) is no longer available; use GenerateFromDescription instead")

// VoiceGender represents the gender options for voice generation (legacy).
// Deprecated: Use VoiceDesignRequest.VoiceDescription instead.
type VoiceGender string

const (
	VoiceGenderFemale VoiceGender = "female"
	VoiceGenderMale   VoiceGender = "male"
)

// VoiceAge represents the age options for voice generation (legacy).
// Deprecated: Use VoiceDesignRequest.VoiceDescription instead.
type VoiceAge string

const (
	VoiceAgeYoung      VoiceAge = "young"
	VoiceAgeMiddleAged VoiceAge = "middle_aged"
	VoiceAgeOld        VoiceAge = "old"
)

// VoiceAccent represents accent options for voice generation (legacy).
// Deprecated: Use VoiceDesignRequest.VoiceDescription instead.
type VoiceAccent string

const (
	VoiceAccentBritish    VoiceAccent = "british"
	VoiceAccentAmerican   VoiceAccent = "american"
	VoiceAccentAfrican    VoiceAccent = "african"
	VoiceAccentAustralian VoiceAccent = "australian"
	VoiceAccentIndian     VoiceAccent = "indian"
)

// VoiceDesignRequest contains options for generating a voice (legacy).
// Deprecated: Use VoiceDescriptionRequest instead.
type VoiceDesignRequest struct {
	Gender         VoiceGender
	Age            VoiceAge
	Accent         VoiceAccent
	AccentStrength float64
	Text           string
}

// VoiceDescriptionRequest contains options for the new text-based voice design.
type VoiceDescriptionRequest struct {
	// VoiceDescription describes the voice characteristics you want.
	// Example: "A warm, friendly female voice with a slight British accent"
	VoiceDescription string

	// Text is the sample text for voice preview (100-1000 characters).
	// If empty and AutoGenerateText is true, text will be generated automatically.
	Text string

	// AutoGenerateText generates suitable preview text automatically.
	AutoGenerateText bool

	// GuidanceScale controls how closely the AI follows the prompt (0-100).
	// Lower = more creative, Higher = more literal.
	GuidanceScale float64

	// Seed for reproducible generation (optional).
	Seed int

	// ModelID specifies the model to use (optional).
	// Options: "eleven_multilingual_ttv_v2", "eleven_ttv_v3"
	ModelID string

	// ShouldEnhance uses AI to expand the voice description.
	ShouldEnhance bool
}

// VoicePreview represents a single generated voice preview.
type VoicePreview struct {
	// GeneratedVoiceID can be used to save this voice permanently.
	GeneratedVoiceID string

	// AudioBase64 is the base64-encoded audio data.
	AudioBase64 string

	// DurationSecs is the duration of the preview in seconds.
	DurationSecs float64

	// Language is the detected language of the preview.
	Language string

	// MediaType is the audio format (e.g., "audio/mpeg").
	MediaType string
}

// VoiceDesignResponse contains the generated voice previews.
type VoiceDesignResponse struct {
	// Previews contains the generated voice previews.
	// The API may return multiple previews for the same description.
	Previews []VoicePreview

	// Text is the text that was used to generate the previews.
	Text string

	// GeneratedVoiceID is the ID of the first preview (convenience field).
	// Use this to save the voice via SaveVoice.
	GeneratedVoiceID string
}

// SaveVoiceRequest contains options for saving a generated voice.
type SaveVoiceRequest struct {
	// GeneratedVoiceID from the design response.
	GeneratedVoiceID string

	// VoiceName is the name for the saved voice.
	VoiceName string

	// VoiceDescription describes the voice.
	VoiceDescription string

	// Labels are optional metadata tags.
	Labels map[string]string
}

// GeneratePreview creates a voice preview based on design parameters.
// Deprecated: Use GenerateFromDescription instead. The legacy API (gender/age/accent)
// is no longer available.
func (s *VoiceDesignService) GeneratePreview(ctx context.Context, req *VoiceDesignRequest) (*VoiceDesignResponse, error) {
	return nil, ErrVoiceDesignLegacyUnavailable
}

// Simple generates a voice preview with common defaults.
// Deprecated: Use GenerateFromDescription instead.
func (s *VoiceDesignService) Simple(ctx context.Context, gender VoiceGender, age VoiceAge, accent VoiceAccent, previewText string) (*VoiceDesignResponse, error) {
	return nil, ErrVoiceDesignLegacyUnavailable
}

// GenerateFromDescription creates a voice preview based on a text description.
// This is the new API that replaces the legacy gender/age/accent approach.
//
// Example:
//
//	resp, err := client.VoiceDesign().GenerateFromDescription(ctx, &VoiceDescriptionRequest{
//	    VoiceDescription: "A warm, friendly female voice with a slight British accent",
//	    Text:             "Hello, this is a sample of my voice...", // 100-1000 chars
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Save the voice
//	voice, err := client.VoiceDesign().SaveVoice(ctx, &SaveVoiceRequest{
//	    GeneratedVoiceID: resp.GeneratedVoiceID,
//	    VoiceName:        "My Custom Voice",
//	})
func (s *VoiceDesignService) GenerateFromDescription(ctx context.Context, req *VoiceDescriptionRequest) (*VoiceDesignResponse, error) {
	if req.VoiceDescription == "" {
		return nil, &ValidationError{Field: "voice_description", Message: "cannot be empty"}
	}

	body := &api.VoiceDesignRequestModel{
		VoiceDescription: req.VoiceDescription,
	}

	if req.Text != "" {
		body.Text = api.NewOptNilString(req.Text)
	}
	if req.AutoGenerateText {
		body.AutoGenerateText = api.NewOptBool(true)
	}
	if req.GuidanceScale > 0 {
		body.GuidanceScale = api.NewOptFloat64(req.GuidanceScale)
	}
	if req.Seed > 0 {
		body.Seed = api.NewOptNilInt(req.Seed)
	}
	if req.ModelID != "" {
		body.ModelID = api.NewOptVoiceDesignRequestModelModelID(
			api.VoiceDesignRequestModelModelID(req.ModelID))
	}
	if req.ShouldEnhance {
		body.ShouldEnhance = api.NewOptBool(true)
	}

	resp, err := s.client.apiClient.TextToVoiceDesign(ctx, body, api.TextToVoiceDesignParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.VoicePreviewsResponseModel:
		result := &VoiceDesignResponse{
			Text:     r.Text,
			Previews: make([]VoicePreview, 0, len(r.Previews)),
		}
		for _, p := range r.Previews {
			preview := VoicePreview{
				GeneratedVoiceID: p.GeneratedVoiceID,
				AudioBase64:      p.AudioBase64,
				DurationSecs:     p.DurationSecs,
				MediaType:        p.MediaType,
			}
			if !p.Language.Null {
				preview.Language = p.Language.Value
			}
			result.Previews = append(result.Previews, preview)
		}
		// Set convenience field to first preview's ID
		if len(result.Previews) > 0 {
			result.GeneratedVoiceID = result.Previews[0].GeneratedVoiceID
		}
		return result, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// SaveVoice saves a previously generated voice to your voice library.
func (s *VoiceDesignService) SaveVoice(ctx context.Context, req *SaveVoiceRequest) (*Voice, error) {
	if req.GeneratedVoiceID == "" {
		return nil, &ValidationError{Field: "generated_voice_id", Message: "cannot be empty"}
	}
	if req.VoiceName == "" {
		return nil, &ValidationError{Field: "voice_name", Message: "cannot be empty"}
	}

	body := &api.BodyCreateANewVoiceFromVoicePreviewV1TextToVoicePost{
		GeneratedVoiceID: req.GeneratedVoiceID,
		VoiceName:        req.VoiceName,
		VoiceDescription: req.VoiceDescription,
	}

	if len(req.Labels) > 0 {
		labels := api.BodyCreateANewVoiceFromVoicePreviewV1TextToVoicePostLabels(req.Labels)
		body.Labels = api.NewOptNilBodyCreateANewVoiceFromVoicePreviewV1TextToVoicePostLabels(labels)
	}

	resp, err := s.client.apiClient.CreateVoice(ctx, body, api.CreateVoiceParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.VoiceResponseModel:
		return &Voice{
			VoiceID:     r.VoiceID,
			Name:        r.Name,
			Description: r.Description.Value,
			Category:    string(r.Category),
		}, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}
