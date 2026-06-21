package elevenlabs

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// Test creating client without API key (uses environment)
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	// Verify domain-based services are initialized
	if client.TTS() == nil {
		t.Error("TTS() service is nil")
	}
	if client.STT() == nil {
		t.Error("STT() service is nil")
	}
	if client.Voice() == nil {
		t.Error("Voice() service is nil")
	}
	if client.Audio() == nil {
		t.Error("Audio() service is nil")
	}
	if client.Content() == nil {
		t.Error("Content() service is nil")
	}
	if client.Realtime() == nil {
		t.Error("Realtime() service is nil")
	}
	if client.Telephony() == nil {
		t.Error("Telephony() service is nil")
	}
	if client.Account() == nil {
		t.Error("Account() service is nil")
	}
	if client.API() == nil {
		t.Error("API() returned nil")
	}
}

func TestNewClientWithAPIKey(t *testing.T) {
	client, err := NewClient(WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("NewClient(WithAPIKey()) error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
}

func TestNewClientWithOptions(t *testing.T) {
	client, err := NewClient(
		WithAPIKey("test-api-key"),
		WithBaseURL("https://custom.api.com"),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}
	if client.baseURL != "https://custom.api.com" {
		t.Errorf("baseURL = %s, want https://custom.api.com", client.baseURL)
	}
}
