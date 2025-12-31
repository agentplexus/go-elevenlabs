package elevenlabs

import (
	"bytes"
	"testing"
)

func TestPCMBytesToWAV(t *testing.T) {
	// Create some fake PCM data (silence)
	pcm := make([]byte, 1000)

	wav, err := PCMBytesToWAV(pcm, 44100)
	if err != nil {
		t.Fatalf("PCMBytesToWAV() error = %v", err)
	}

	// Check WAV header
	if string(wav[0:4]) != "RIFF" {
		t.Error("WAV should start with RIFF")
	}
	if string(wav[8:12]) != "WAVE" {
		t.Error("WAV should contain WAVE marker")
	}
	if string(wav[12:16]) != "fmt " {
		t.Error("WAV should contain fmt chunk")
	}
	if string(wav[36:40]) != "data" {
		t.Error("WAV should contain data chunk")
	}

	// Check total size (44 byte header + 1000 bytes PCM)
	if len(wav) != 44+1000 {
		t.Errorf("WAV size = %d, want %d", len(wav), 44+1000)
	}
}

func TestPCMToWAV(t *testing.T) {
	pcm := make([]byte, 500)
	reader := bytes.NewReader(pcm)

	wav, err := PCMToWAV(reader, 22050)
	if err != nil {
		t.Fatalf("PCMToWAV() error = %v", err)
	}

	if string(wav[0:4]) != "RIFF" {
		t.Error("WAV should start with RIFF")
	}
}

func TestParsePCMSampleRate(t *testing.T) {
	tests := []struct {
		format  string
		want    int
		wantErr bool
	}{
		{"pcm_44100", 44100, false},
		{"pcm_48000", 48000, false},
		{"pcm_16000", 16000, false},
		{"pcm_8000", 8000, false},
		{"mp3_44100_128", 0, true},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			got, err := ParsePCMSampleRate(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePCMSampleRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePCMSampleRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
