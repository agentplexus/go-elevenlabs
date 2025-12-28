package elevenlabs

import (
	"context"
	"io"

	"github.com/agentplexus/go-elevenlabs/internal/api"
)

// TextToSpeechService handles text-to-speech operations.
type TextToSpeechService struct {
	client *Client
}

// VoiceSettings contains the voice configuration for text-to-speech.
type VoiceSettings struct {
	// Stability determines how stable the voice is (0.0 to 1.0).
	// Lower values introduce broader emotional range.
	Stability float64

	// SimilarityBoost determines how closely the AI should adhere to
	// the original voice (0.0 to 1.0).
	SimilarityBoost float64

	// Style determines the style exaggeration (0.0 to 1.0).
	// Higher values amplify the original speaker's style.
	Style float64

	// Speed adjusts the speed of the voice (0.25 to 4.0).
	// 1.0 is the default speed.
	Speed float64

	// UseSpeakerBoost boosts similarity to the original speaker.
	UseSpeakerBoost bool
}

// Validate validates the voice settings.
func (vs *VoiceSettings) Validate() error {
	if vs.Stability < 0 || vs.Stability > 1 {
		return ErrInvalidStability
	}
	if vs.SimilarityBoost < 0 || vs.SimilarityBoost > 1 {
		return ErrInvalidSimilarityBoost
	}
	if vs.Style < 0 || vs.Style > 1 {
		return ErrInvalidStyle
	}
	if vs.Speed != 0 && (vs.Speed < 0.25 || vs.Speed > 4.0) {
		return ErrInvalidSpeed
	}
	return nil
}

// DefaultVoiceSettings returns sensible default voice settings.
func DefaultVoiceSettings() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Style:           0.0,
		Speed:           1.0,
		UseSpeakerBoost: true,
	}
}

// TTSRequest is a request to generate speech from text.
type TTSRequest struct {
	// VoiceID is the voice to use for generation.
	VoiceID string

	// Text is the text to convert to speech.
	Text string

	// ModelID is the model to use. Defaults to DefaultModelID.
	ModelID string

	// VoiceSettings configures the voice parameters.
	// If nil, default settings will be used.
	VoiceSettings *VoiceSettings

	// OutputFormat specifies the audio output format.
	// Examples: "mp3_44100_128", "pcm_16000", "pcm_22050"
	OutputFormat string

	// LanguageCode is the ISO 639-1 language code for text normalization.
	LanguageCode string
}

// Validate validates the TTS request.
func (r *TTSRequest) Validate() error {
	if r.VoiceID == "" {
		return ErrEmptyVoiceID
	}
	if r.Text == "" {
		return ErrEmptyText
	}
	if r.VoiceSettings != nil {
		if err := r.VoiceSettings.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TTSResponse contains the generated audio from text-to-speech.
type TTSResponse struct {
	// Audio is the generated audio data.
	Audio io.Reader
}

// Generate generates speech from text.
func (s *TextToSpeechService) Generate(ctx context.Context, req *TTSRequest) (*TTSResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Build request body
	body := &api.BodyTextToSpeechFull{
		Text: req.Text,
	}

	// Set model ID
	modelID := req.ModelID
	if modelID == "" {
		modelID = DefaultModelID
	}
	body.ModelID = api.NewOptString(modelID)

	// Set voice settings if provided
	if req.VoiceSettings != nil {
		vs := api.VoiceSettingsResponseModel{
			Stability:       api.NewOptNilFloat64(req.VoiceSettings.Stability),
			SimilarityBoost: api.NewOptNilFloat64(req.VoiceSettings.SimilarityBoost),
			Style:           api.NewOptNilFloat64(req.VoiceSettings.Style),
		}
		if req.VoiceSettings.Speed != 0 {
			vs.Speed = api.NewOptNilFloat64(req.VoiceSettings.Speed)
		}
		body.VoiceSettings = api.NewOptVoiceSettingsResponseModel(vs)
	}

	// Set language code if provided
	if req.LanguageCode != "" {
		body.LanguageCode = api.NewOptNilString(req.LanguageCode)
	}

	// Build params
	params := api.TextToSpeechFullParams{
		VoiceID: req.VoiceID,
	}

	// Set output format if provided
	if req.OutputFormat != "" {
		params.OutputFormat = api.NewOptTextToSpeechFullOutputFormat(
			api.TextToSpeechFullOutputFormat(req.OutputFormat),
		)
	}

	// Make the API call
	resp, err := s.client.apiClient.TextToSpeechFull(ctx, body, params)
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.TextToSpeechFullOK:
		return &TTSResponse{Audio: r.Data}, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// GenerateToWriter generates speech and writes it to a writer.
func (s *TextToSpeechService) GenerateToWriter(ctx context.Context, req *TTSRequest, w io.Writer) error {
	resp, err := s.Generate(ctx, req)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, resp.Audio)
	return err
}

// Simple is a convenience method that generates speech with minimal parameters.
func (s *TextToSpeechService) Simple(ctx context.Context, voiceID, text string) (io.Reader, error) {
	resp, err := s.Generate(ctx, &TTSRequest{
		VoiceID:       voiceID,
		Text:          text,
		VoiceSettings: DefaultVoiceSettings(),
	})
	if err != nil {
		return nil, err
	}
	return resp.Audio, nil
}
