package stt

import (
	"os"
	"testing"

	"github.com/plexusone/omnivoice-core/stt/providertest"
)

// TestConformance runs the OmniVoice STT provider conformance tests.
func TestConformance(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping conformance tests")
	}

	p, err := New(WithAPIKey(apiKey))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	// Download test audio for conformance tests
	testAudio := downloadTestAudioBytes(t)

	providertest.RunAll(t, providertest.Config{
		Provider:          p,
		StreamingProvider: p,
		SkipIntegration:   false,
		TestAudio:         testAudio,
		TestAudioURL:      testAudioURL,
	})
}

// TestInterfaceConformance runs only interface tests (no API calls).
func TestInterfaceConformance(t *testing.T) {
	p, err := New(WithAPIKey("test-key"))
	if err != nil {
		t.Fatalf("failed to create provider: %v", err)
	}

	providertest.RunInterfaceTests(t, providertest.Config{
		Provider: p,
	})
}
