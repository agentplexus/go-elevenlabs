package elevenlabs

import (
	"context"
	"io"
	"testing"
)

func TestVoiceSettingsValidate(t *testing.T) {
	tests := []struct {
		name    string
		vs      *VoiceSettings
		wantErr error
	}{
		{
			name:    "valid settings",
			vs:      DefaultVoiceSettings(),
			wantErr: nil,
		},
		{
			name:    "stability too low",
			vs:      &VoiceSettings{Stability: -0.1, SimilarityBoost: 0.5},
			wantErr: ErrInvalidStability,
		},
		{
			name:    "stability too high",
			vs:      &VoiceSettings{Stability: 1.1, SimilarityBoost: 0.5},
			wantErr: ErrInvalidStability,
		},
		{
			name:    "similarity_boost too low",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: -0.1},
			wantErr: ErrInvalidSimilarityBoost,
		},
		{
			name:    "similarity_boost too high",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 1.1},
			wantErr: ErrInvalidSimilarityBoost,
		},
		{
			name:    "style too low",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 0.5, Style: -0.1},
			wantErr: ErrInvalidStyle,
		},
		{
			name:    "style too high",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 0.5, Style: 1.1},
			wantErr: ErrInvalidStyle,
		},
		{
			name:    "speed too low",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 0.5, Speed: 0.1},
			wantErr: ErrInvalidSpeed,
		},
		{
			name:    "speed too high",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 0.5, Speed: 5.0},
			wantErr: ErrInvalidSpeed,
		},
		{
			name:    "speed zero is valid (default)",
			vs:      &VoiceSettings{Stability: 0.5, SimilarityBoost: 0.5, Speed: 0},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vs.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTTSRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *TTSRequest
		wantErr error
	}{
		{
			name:    "valid request",
			req:     &TTSRequest{VoiceID: "test-voice", Text: "Hello"},
			wantErr: nil,
		},
		{
			name:    "empty voice ID",
			req:     &TTSRequest{VoiceID: "", Text: "Hello"},
			wantErr: ErrEmptyVoiceID,
		},
		{
			name:    "empty text",
			req:     &TTSRequest{VoiceID: "test-voice", Text: ""},
			wantErr: ErrEmptyText,
		},
		{
			name: "invalid voice settings",
			req: &TTSRequest{
				VoiceID:       "test-voice",
				Text:          "Hello",
				VoiceSettings: &VoiceSettings{Stability: 2.0},
			},
			wantErr: ErrInvalidStability,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTTSRequestValidate_OutputFormat(t *testing.T) {
	tests := []struct {
		name       string
		format     string
		shouldPass bool
	}{
		{"empty (default)", "", true},
		{"mp3_44100_128", "mp3_44100_128", true},
		{"mp3_44100_192 (highest mp3)", "mp3_44100_192", true},
		{"mp3_22050_32", "mp3_22050_32", true},
		{"pcm_16000", "pcm_16000", true},
		{"pcm_48000 (highest quality)", "pcm_48000", true},
		{"ulaw_8000", "ulaw_8000", true},
		{"alaw_8000", "alaw_8000", true},
		{"opus_48000_192", "opus_48000_192", true},
		{"invalid mp3", "mp3", false},
		{"invalid wav", "wav", false},
		{"invalid random", "foo_bar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &TTSRequest{
				VoiceID:      "test-voice",
				Text:         "Hello",
				OutputFormat: tt.format,
			}
			err := req.Validate()
			if tt.shouldPass && err != nil {
				t.Errorf("Validate() error = %v, want nil", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Validate() error = nil, want error for format %q", tt.format)
			}
		})
	}
}

func TestDefaultVoiceSettings(t *testing.T) {
	vs := DefaultVoiceSettings()
	if vs == nil {
		t.Fatal("DefaultVoiceSettings() returned nil")
	}

	// Validate that defaults are within valid ranges
	if err := vs.Validate(); err != nil {
		t.Errorf("DefaultVoiceSettings().Validate() error = %v", err)
	}

	// Check expected default values
	if vs.Stability != 0.5 {
		t.Errorf("Stability = %v, want 0.5", vs.Stability)
	}
	if vs.SimilarityBoost != 0.75 {
		t.Errorf("SimilarityBoost = %v, want 0.75", vs.SimilarityBoost)
	}
	if vs.Speed != 1.0 {
		t.Errorf("Speed = %v, want 1.0", vs.Speed)
	}
}

// Live API tests - only run when ELEVENLABS_API_KEY is set
func TestTextToSpeechGenerate_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Get list of voices first to find a valid voice ID
	voices, err := client.Voices().List(context.Background())
	if err != nil {
		t.Fatalf("Voices().List() error = %v", err)
	}
	if len(voices) == 0 {
		t.Skip("No voices available")
	}

	voiceID := voices[0].VoiceID

	resp, err := client.TextToSpeech().Generate(context.Background(), &TTSRequest{
		VoiceID:       voiceID,
		Text:          "Hello, this is a test.",
		VoiceSettings: DefaultVoiceSettings(),
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp == nil {
		t.Fatal("Generate() returned nil response")
	}
	if resp.Audio == nil {
		t.Fatal("Generate() returned nil audio")
	}

	// Read some bytes to ensure we got audio data
	buf := make([]byte, 1024)
	n, err := resp.Audio.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Audio.Read() error = %v", err)
	}
	if n == 0 {
		t.Error("Audio.Read() returned 0 bytes")
	}
}

func TestTextToSpeechSimple_Live(t *testing.T) {
	apiKey := getAPIKey(t)

	client, err := NewClient(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Get list of voices first
	voices, err := client.Voices().List(context.Background())
	if err != nil {
		t.Fatalf("Voices().List() error = %v", err)
	}
	if len(voices) == 0 {
		t.Skip("No voices available")
	}

	voiceID := voices[0].VoiceID

	audio, err := client.TextToSpeech().Simple(context.Background(), voiceID, "Hello world")
	if err != nil {
		t.Fatalf("Simple() error = %v", err)
	}
	if audio == nil {
		t.Fatal("Simple() returned nil audio")
	}

	// Read some bytes to ensure we got audio data
	buf := make([]byte, 1024)
	n, err := audio.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Audio.Read() error = %v", err)
	}
	if n == 0 {
		t.Error("Audio.Read() returned 0 bytes")
	}
}
