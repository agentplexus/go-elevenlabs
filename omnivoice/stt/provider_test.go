package stt

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/omnivoice-core/stt"
)

// testAudioURL is Deepgram's public test audio file.
const testAudioURL = "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"

func TestTranscribeURL(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := p.TranscribeURL(ctx, testAudioURL, stt.TranscriptionConfig{})
	if err != nil {
		t.Fatalf("TranscribeURL failed: %v", err)
	}

	if result.Text == "" {
		t.Error("TranscribeURL returned empty text")
	}

	// Check for expected content
	lower := strings.ToLower(result.Text)
	if !strings.Contains(lower, "life") && !strings.Contains(lower, "fast") {
		t.Errorf("TranscribeURL text doesn't contain expected content: %q", result.Text)
	}

	t.Logf("TranscribeURL result: %q", result.Text)
}

func TestTranscribeStream(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	writer, events, err := p.TranscribeStream(ctx, stt.TranscriptionConfig{})
	if err != nil {
		t.Fatalf("TranscribeStream failed: %v", err)
	}
	defer writer.Close()

	// Send minimal audio to trigger connection (empty audio is fine for connection test)
	// The WebSocket connection is the main thing we're testing here
	t.Log("TranscribeStream connection established successfully")

	// Close the writer to signal end of stream
	if err := writer.Close(); err != nil {
		t.Errorf("writer.Close() failed: %v", err)
	}

	// Drain events
	for range events {
		// Just drain
	}
}

func TestTranscribe(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Download test audio for this test
	audioData := downloadTestAudioBytes(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := p.Transcribe(ctx, audioData, stt.TranscriptionConfig{})
	if err != nil {
		t.Fatalf("Transcribe failed: %v", err)
	}

	if result.Text == "" {
		t.Error("Transcribe returned empty text")
	}

	// Check for expected content
	lower := strings.ToLower(result.Text)
	if !strings.Contains(lower, "life") && !strings.Contains(lower, "fast") {
		t.Errorf("Transcribe text doesn't contain expected content: %q", result.Text)
	}

	t.Logf("Transcribe result: %q", result.Text)
}

func TestTranscribeFile(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Download test audio to a file
	audioFile := downloadTestAudioToFile(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := p.TranscribeFile(ctx, audioFile, stt.TranscriptionConfig{})
	if err != nil {
		t.Fatalf("TranscribeFile failed: %v", err)
	}

	if result.Text == "" {
		t.Error("TranscribeFile returned empty text")
	}

	// Check for expected content
	lower := strings.ToLower(result.Text)
	if !strings.Contains(lower, "life") && !strings.Contains(lower, "fast") {
		t.Errorf("TranscribeFile text doesn't contain expected content: %q", result.Text)
	}

	t.Logf("TranscribeFile result: %q", result.Text)
}

func TestProviderName(t *testing.T) {
	p, err := New(WithAPIKey("test"))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	name := p.Name()
	if name == "" {
		t.Error("Provider.Name() should not be empty")
	}
	if name != "elevenlabs" {
		t.Errorf("Provider.Name() = %q, want %q", name, "elevenlabs")
	}
}

// downloadTestAudioBytes downloads the test audio file and returns its bytes.
func downloadTestAudioBytes(t *testing.T) []byte {
	t.Helper()

	resp, err := http.Get(testAudioURL)
	if err != nil {
		t.Fatalf("failed to download test audio: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to download test audio: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read test audio: %v", err)
	}

	return data
}

// downloadTestAudioToFile downloads the test audio file to a temp file and returns its path.
func downloadTestAudioToFile(t *testing.T) string {
	t.Helper()

	data := downloadTestAudioBytes(t)

	tmpDir := t.TempDir()
	audioFile := filepath.Join(tmpDir, "test-audio.wav")

	if err := os.WriteFile(audioFile, data, 0644); err != nil {
		t.Fatalf("failed to write test audio file: %v", err)
	}

	return audioFile
}
