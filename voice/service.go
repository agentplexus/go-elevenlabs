package voice

import (
	"context"
	"errors"
	"fmt"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Service handles voice management and design operations.
type Service struct {
	apiClient *api.Client
}

// New creates a new voice service.
func New(apiClient *api.Client) *Service {
	return &Service{apiClient: apiClient}
}

// VoiceSettings contains the voice configuration.
type VoiceSettings struct {
	// Stability determines how stable the voice is (0.0 to 1.0).
	Stability float64

	// SimilarityBoost determines how closely the AI should adhere to
	// the original voice (0.0 to 1.0).
	SimilarityBoost float64

	// Style determines the style exaggeration (0.0 to 1.0).
	Style float64

	// Speed adjusts the speed of the voice (0.25 to 4.0).
	Speed float64

	// UseSpeakerBoost boosts similarity to the original speaker.
	UseSpeakerBoost bool
}

// Voice represents an ElevenLabs voice.
type Voice struct {
	// VoiceID is the unique identifier for the voice.
	VoiceID string

	// Name is the display name of the voice.
	Name string

	// Category is the category of the voice (e.g., "premade", "cloned").
	Category string

	// Description is the description of the voice.
	Description string

	// PreviewURL is the URL to preview the voice.
	PreviewURL string

	// Labels contains additional metadata about the voice.
	Labels map[string]string
}

// List returns all available voices.
func (s *Service) List(ctx context.Context) ([]*Voice, error) {
	resp, err := s.apiClient.GetVoices(ctx, api.GetVoicesParams{})
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.GetVoicesResponseModel:
		voices := make([]*Voice, 0, len(r.Voices))
		for _, v := range r.Voices {
			voice := &Voice{
				VoiceID:  v.VoiceID,
				Name:     v.Name,
				Category: string(v.Category),
				Labels:   make(map[string]string),
			}
			if v.Description.Set && !v.Description.Null {
				voice.Description = v.Description.Value
			}
			if v.PreviewURL.Set && !v.PreviewURL.Null {
				voice.PreviewURL = v.PreviewURL.Value
			}
			// Convert labels
			for k, val := range v.Labels {
				voice.Labels[k] = val
			}
			voices = append(voices, voice)
		}
		return voices, nil
	default:
		return nil, fmt.Errorf("voice: unexpected response type")
	}
}

// Get returns a voice by ID.
func (s *Service) Get(ctx context.Context, voiceID string) (*Voice, error) {
	if voiceID == "" {
		return nil, fmt.Errorf("voice: voice ID is required")
	}

	resp, err := s.apiClient.GetVoiceByID(ctx, api.GetVoiceByIDParams{
		VoiceID: voiceID,
	})
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.VoiceResponseModel:
		voice := &Voice{
			VoiceID:  r.VoiceID,
			Name:     r.Name,
			Category: string(r.Category),
			Labels:   make(map[string]string),
		}
		if r.Description.Set && !r.Description.Null {
			voice.Description = r.Description.Value
		}
		if r.PreviewURL.Set && !r.PreviewURL.Null {
			voice.PreviewURL = r.PreviewURL.Value
		}
		// Convert labels
		for k, val := range r.Labels {
			voice.Labels[k] = val
		}
		return voice, nil
	default:
		return nil, fmt.Errorf("voice: unexpected response type")
	}
}

// GetSettings returns the settings for a voice.
func (s *Service) GetSettings(ctx context.Context, voiceID string) (*VoiceSettings, error) {
	if voiceID == "" {
		return nil, fmt.Errorf("voice: voice ID is required")
	}

	resp, err := s.apiClient.GetVoiceSettings(ctx, api.GetVoiceSettingsParams{
		VoiceID: voiceID,
	})
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.VoiceSettingsResponseModel:
		settings := &VoiceSettings{}
		if r.Stability.Set && !r.Stability.Null {
			settings.Stability = r.Stability.Value
		}
		if r.SimilarityBoost.Set && !r.SimilarityBoost.Null {
			settings.SimilarityBoost = r.SimilarityBoost.Value
		}
		if r.Style.Set && !r.Style.Null {
			settings.Style = r.Style.Value
		}
		if r.Speed.Set && !r.Speed.Null {
			settings.Speed = r.Speed.Value
		}
		return settings, nil
	default:
		return nil, fmt.Errorf("voice: unexpected response type")
	}
}

// GetDefaultSettings returns the default voice settings.
func (s *Service) GetDefaultSettings(ctx context.Context) (*VoiceSettings, error) {
	resp, err := s.apiClient.GetVoiceSettingsDefault(ctx)
	if err != nil {
		return nil, err
	}

	settings := &VoiceSettings{}
	if resp.Stability.Set && !resp.Stability.Null {
		settings.Stability = resp.Stability.Value
	}
	if resp.SimilarityBoost.Set && !resp.SimilarityBoost.Null {
		settings.SimilarityBoost = resp.SimilarityBoost.Value
	}
	if resp.Style.Set && !resp.Style.Null {
		settings.Style = resp.Style.Value
	}
	if resp.Speed.Set && !resp.Speed.Null {
		settings.Speed = resp.Speed.Value
	}
	return settings, nil
}

// Delete deletes a voice by ID.
func (s *Service) Delete(ctx context.Context, voiceID string) error {
	if voiceID == "" {
		return fmt.Errorf("voice: voice ID is required")
	}

	_, err := s.apiClient.DeleteVoice(ctx, api.DeleteVoiceParams{
		VoiceID: voiceID,
	})
	return err
}

// --- Voice Design ---

// ErrVoiceDesignLegacyUnavailable is returned when legacy voice design methods are called.
// The ElevenLabs API now uses text-based voice descriptions instead of gender/age/accent parameters.
var ErrVoiceDesignLegacyUnavailable = errors.New("voice: legacy voice design API (gender/age/accent) is no longer available; use DesignFromDescription instead")

// VoiceGender represents the gender options for voice generation (legacy).
// Deprecated: Use DesignRequest.VoiceDescription instead.
type VoiceGender string

const (
	VoiceGenderFemale VoiceGender = "female"
	VoiceGenderMale   VoiceGender = "male"
)

// VoiceAge represents the age options for voice generation (legacy).
// Deprecated: Use DesignRequest.VoiceDescription instead.
type VoiceAge string

const (
	VoiceAgeYoung      VoiceAge = "young"
	VoiceAgeMiddleAged VoiceAge = "middle_aged"
	VoiceAgeOld        VoiceAge = "old"
)

// VoiceAccent represents accent options for voice generation (legacy).
// Deprecated: Use DesignRequest.VoiceDescription instead.
type VoiceAccent string

const (
	VoiceAccentBritish    VoiceAccent = "british"
	VoiceAccentAmerican   VoiceAccent = "american"
	VoiceAccentAfrican    VoiceAccent = "african"
	VoiceAccentAustralian VoiceAccent = "australian"
	VoiceAccentIndian     VoiceAccent = "indian"
)

// DesignRequest contains options for the text-based voice design.
type DesignRequest struct {
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

// DesignResponse contains the generated voice previews.
type DesignResponse struct {
	// Previews contains the generated voice previews.
	// The API may return multiple previews for the same description.
	Previews []VoicePreview

	// Text is the text that was used to generate the previews.
	Text string

	// GeneratedVoiceID is the ID of the first preview (convenience field).
	// Use this to save the voice via SaveDesign.
	GeneratedVoiceID string
}

// SaveDesignRequest contains options for saving a generated voice.
type SaveDesignRequest struct {
	// GeneratedVoiceID from the design response.
	GeneratedVoiceID string

	// VoiceName is the name for the saved voice.
	VoiceName string

	// VoiceDescription describes the voice.
	VoiceDescription string

	// Labels are optional metadata tags.
	Labels map[string]string
}

// DesignFromDescription creates a voice preview based on a text description.
// This is the new API that replaces the legacy gender/age/accent approach.
//
// Example:
//
//	resp, err := client.Voice().DesignFromDescription(ctx, &voice.DesignRequest{
//	    VoiceDescription: "A warm, friendly female voice with a slight British accent",
//	    Text:             "Hello, this is a sample of my voice...", // 100-1000 chars
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Save the voice
//	voice, err := client.Voice().SaveDesign(ctx, &voice.SaveDesignRequest{
//	    GeneratedVoiceID: resp.GeneratedVoiceID,
//	    VoiceName:        "My Custom Voice",
//	})
func (s *Service) DesignFromDescription(ctx context.Context, req *DesignRequest) (*DesignResponse, error) {
	if req.VoiceDescription == "" {
		return nil, fmt.Errorf("voice: voice_description cannot be empty")
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

	resp, err := s.apiClient.TextToVoiceDesign(ctx, body, api.TextToVoiceDesignParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.VoicePreviewsResponseModel:
		result := &DesignResponse{
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
		return nil, fmt.Errorf("voice: unexpected response type")
	}
}

// SaveDesign saves a previously generated voice to your voice library.
func (s *Service) SaveDesign(ctx context.Context, req *SaveDesignRequest) (*Voice, error) {
	if req.GeneratedVoiceID == "" {
		return nil, fmt.Errorf("voice: generated_voice_id cannot be empty")
	}
	if req.VoiceName == "" {
		return nil, fmt.Errorf("voice: voice_name cannot be empty")
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

	resp, err := s.apiClient.CreateVoice(ctx, body, api.CreateVoiceParams{})
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
		return nil, fmt.Errorf("voice: unexpected response type")
	}
}
