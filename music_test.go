package elevenlabs

import (
	"context"
	"errors"
	"testing"
)

// dummyReader is a minimal io.Reader for testing validation
type dummyReader struct{}

func (d dummyReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

func TestMusicCompositionUnavailable(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// All composition methods should return ErrMusicCompositionUnavailable
	t.Run("Generate", func(t *testing.T) {
		_, err := client.Music().Generate(ctx, &MusicRequest{Prompt: "test"})
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("Generate() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})

	t.Run("GenerateStream", func(t *testing.T) {
		_, err := client.Music().GenerateStream(ctx, &MusicRequest{Prompt: "test"})
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("GenerateStream() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})

	t.Run("Simple", func(t *testing.T) {
		_, err := client.Music().Simple(ctx, "test prompt")
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("Simple() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})

	t.Run("GenerateInstrumental", func(t *testing.T) {
		_, err := client.Music().GenerateInstrumental(ctx, "test prompt", 30000)
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("GenerateInstrumental() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})

	t.Run("GeneratePlan", func(t *testing.T) {
		_, err := client.Music().GeneratePlan(ctx, &CompositionPlanRequest{Prompt: "test"})
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("GeneratePlan() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})

	t.Run("GenerateDetailed", func(t *testing.T) {
		_, err := client.Music().GenerateDetailed(ctx, &MusicDetailedRequest{Prompt: "test"})
		if !errors.Is(err, ErrMusicCompositionUnavailable) {
			t.Errorf("GenerateDetailed() error = %v, want ErrMusicCompositionUnavailable", err)
		}
	})
}

func TestMusicService(t *testing.T) {
	client, _ := NewClient()

	// Test that service is accessible
	if client.Music() == nil {
		t.Error("Music() returned nil")
	}
}

func TestMusicRequest(t *testing.T) {
	// Test MusicRequest struct
	req := MusicRequest{
		Prompt:            "upbeat electronic music",
		DurationMs:        30000,
		ForceInstrumental: true,
		Seed:              12345,
	}

	if req.Prompt != "upbeat electronic music" {
		t.Errorf("Prompt = %s, want upbeat electronic music", req.Prompt)
	}
	if req.DurationMs != 30000 {
		t.Errorf("DurationMs = %d, want 30000", req.DurationMs)
	}
	if !req.ForceInstrumental {
		t.Error("ForceInstrumental should be true")
	}
	if req.Seed != 12345 {
		t.Errorf("Seed = %d, want 12345", req.Seed)
	}
}

func TestMusicResponse(t *testing.T) {
	// Test MusicResponse struct
	resp := MusicResponse{
		SongID: "song123",
	}

	if resp.SongID != "song123" {
		t.Errorf("SongID = %s, want song123", resp.SongID)
	}
}

func TestStemSeparationValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test nil file
	_, err := client.Music().SeparateStems(ctx, &StemSeparationRequest{
		Filename: "test.mp3",
	})
	if err == nil {
		t.Error("SeparateStems() with nil file should return error")
	}
	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "file" {
		t.Errorf("ValidationError field = %s, want file", valErr.Field)
	}

	// Test empty filename
	_, err = client.Music().SeparateStems(ctx, &StemSeparationRequest{
		File: dummyReader{},
	})
	if err == nil {
		t.Error("SeparateStems() with empty filename should return error")
	}
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "filename" {
		t.Errorf("ValidationError field = %s, want filename", valErr.Field)
	}
}

func TestVideoToMusicValidation(t *testing.T) {
	client, _ := NewClient()
	ctx := context.Background()

	// Test empty videos
	_, err := client.Music().VideoToMusic(ctx, &VideoToMusicRequest{})
	if err == nil {
		t.Error("VideoToMusic() with empty videos should return error")
	}
	var valErr *ValidationError
	if !isValidationError(err, &valErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if valErr.Field != "videos" {
		t.Errorf("ValidationError field = %s, want videos", valErr.Field)
	}

	// Test nil file in video
	_, err = client.Music().VideoToMusic(ctx, &VideoToMusicRequest{
		Videos: []VideoFile{{Name: "test.mp4"}},
	})
	if err == nil {
		t.Error("VideoToMusic() with nil video file should return error")
	}

	// Test empty filename in video
	_, err = client.Music().VideoToMusic(ctx, &VideoToMusicRequest{
		Videos: []VideoFile{{File: dummyReader{}}},
	})
	if err == nil {
		t.Error("VideoToMusic() with empty video filename should return error")
	}
}
