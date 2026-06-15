package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// formatValidationLocation formats the location array from validation errors.
func formatValidationLocation(loc []api.ValidationErrorLocItem) string {
	if len(loc) == 0 {
		return ""
	}
	parts := make([]string, 0, len(loc))
	for _, item := range loc {
		if item.IsString() {
			parts = append(parts, item.String)
		} else if item.IsInt() {
			parts = append(parts, string(rune('0'+item.Int)))
		}
	}
	return strings.Join(parts, ".")
}

// SpeechToTextService handles speech-to-text transcription.
type SpeechToTextService struct {
	client *Client
}

// TranscriptionRequest contains options for transcription.
type TranscriptionRequest struct {
	// FileURL is the HTTPS URL of the file to transcribe.
	FileURL string

	// FileBytes is raw audio file bytes for direct upload.
	// Use this for local file transcription.
	FileBytes []byte

	// FileName is the filename to use when uploading FileBytes.
	// Defaults to "audio.wav" if not specified.
	FileName string

	// LanguageCode is an ISO-639-1 or ISO-639-3 language code.
	// If not provided, language is auto-detected.
	LanguageCode string

	// Diarize enables speaker diarization (who said what).
	Diarize bool

	// NumSpeakers is the expected number of speakers (for diarization).
	NumSpeakers int

	// TagAudioEvents tags audio events like laughter, applause, etc.
	TagAudioEvents bool

	// ModelID is the transcription model to use (default: "scribe_v1").
	ModelID string
}

// TranscriptionResponse contains the transcription result.
type TranscriptionResponse struct {
	// Text is the full transcribed text.
	Text string

	// LanguageCode is the detected language.
	LanguageCode string

	// Words contains word-level details with timestamps.
	Words []TranscriptionWord

	// Utterances contains speaker-labeled segments (when diarization is enabled).
	Utterances []TranscriptionUtterance
}

// TranscriptionWord represents a single word with timing.
type TranscriptionWord struct {
	// Text is the word text.
	Text string

	// Start is the start time in seconds.
	Start float64

	// End is the end time in seconds.
	End float64

	// Confidence is the confidence score (0-1).
	Confidence float64

	// Speaker is the speaker ID (when diarization is enabled).
	Speaker string

	// Type is the word type (e.g., "word", "punctuation").
	Type string
}

// TranscriptionUtterance represents a speaker segment.
type TranscriptionUtterance struct {
	// Text is the utterance text.
	Text string

	// Start is the start time in seconds.
	Start float64

	// End is the end time in seconds.
	End float64

	// Speaker is the speaker ID.
	Speaker string
}

// Transcribe transcribes audio to text.
func (s *SpeechToTextService) Transcribe(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error) {
	if req.FileURL == "" && len(req.FileBytes) == 0 {
		return nil, &ValidationError{Field: "file", Message: "either FileURL or FileBytes must be provided"}
	}

	// Use handwritten multipart request for file bytes (ogen doesn't handle binary uploads correctly)
	if len(req.FileBytes) > 0 {
		return s.transcribeWithFileUpload(ctx, req)
	}

	// Use ogen-generated client for URL-based transcription
	return s.transcribeWithURL(ctx, req)
}

// transcribeWithURL uses the ogen-generated client for URL-based transcription.
func (s *SpeechToTextService) transcribeWithURL(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error) {
	body := &api.BodySpeechToTextV1SpeechToTextPostMultipart{}

	body.SourceURL = api.NewOptNilString(req.FileURL)

	if req.LanguageCode != "" {
		body.LanguageCode = api.NewOptNilString(req.LanguageCode)
	}
	if req.Diarize {
		body.Diarize = api.NewOptBool(true)
	}
	if req.NumSpeakers > 0 {
		body.NumSpeakers = api.NewOptNilInt(req.NumSpeakers)
	}
	if req.TagAudioEvents {
		body.TagAudioEvents = api.NewOptBool(true)
	}
	if req.ModelID != "" {
		body.ModelID = api.BodySpeechToTextV1SpeechToTextPostMultipartModelID(req.ModelID)
	} else {
		body.ModelID = api.BodySpeechToTextV1SpeechToTextPostMultipartModelIDScribeV1
	}

	resp, err := s.client.apiClient.SpeechToText(ctx, body, api.SpeechToTextParams{})
	if err != nil {
		return nil, err
	}

	return s.parseTranscriptionResponse(resp)
}

// transcribeWithFileUpload uses a handwritten multipart request for file uploads.
// This bypasses ogen because it doesn't handle nullable binary fields correctly.
func (s *SpeechToTextService) transcribeWithFileUpload(ctx context.Context, req *TranscriptionRequest) (*TranscriptionResponse, error) {
	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	fileName := req.FileName
	if fileName == "" {
		fileName = "audio.wav"
	}
	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := filePart.Write(req.FileBytes); err != nil {
		return nil, fmt.Errorf("write file data: %w", err)
	}

	// Add model_id (required)
	modelID := req.ModelID
	if modelID == "" {
		modelID = "scribe_v1"
	}
	if err := writer.WriteField("model_id", modelID); err != nil {
		return nil, fmt.Errorf("write model_id: %w", err)
	}

	// Add optional fields
	if req.LanguageCode != "" {
		if err := writer.WriteField("language_code", req.LanguageCode); err != nil {
			return nil, fmt.Errorf("write language_code: %w", err)
		}
	}
	if req.Diarize {
		if err := writer.WriteField("diarize", "true"); err != nil {
			return nil, fmt.Errorf("write diarize: %w", err)
		}
	}
	if req.NumSpeakers > 0 {
		if err := writer.WriteField("num_speakers", fmt.Sprintf("%d", req.NumSpeakers)); err != nil {
			return nil, fmt.Errorf("write num_speakers: %w", err)
		}
	}
	if req.TagAudioEvents {
		if err := writer.WriteField("tag_audio_events", "true"); err != nil {
			return nil, fmt.Errorf("write tag_audio_events: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.client.baseURL+"/v1/speech-to-text", &buf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("xi-api-key", s.client.apiKey)

	// Execute request using default HTTP client
	// (we set auth header manually, don't need authHTTPClient wrapper)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("API request failed: %s", string(respBody)),
		}
	}

	// Parse response
	var result sttAPIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return result.toTranscriptionResponse(), nil
}

// sttAPIResponse is the raw API response structure.
type sttAPIResponse struct {
	Text         string        `json:"text"`
	LanguageCode string        `json:"language_code"`
	Words        []sttWordResp `json:"words"`
}

type sttWordResp struct {
	Text      string   `json:"text"`
	Start     *float64 `json:"start"`
	End       *float64 `json:"end"`
	SpeakerID *string  `json:"speaker_id"`
	Type      string   `json:"type"`
}

func (r *sttAPIResponse) toTranscriptionResponse() *TranscriptionResponse {
	result := &TranscriptionResponse{
		Text:         r.Text,
		LanguageCode: r.LanguageCode,
	}

	for _, w := range r.Words {
		word := TranscriptionWord{
			Text: w.Text,
			Type: w.Type,
		}
		if w.Start != nil {
			word.Start = *w.Start
		}
		if w.End != nil {
			word.End = *w.End
		}
		if w.SpeakerID != nil {
			word.Speaker = *w.SpeakerID
		}
		result.Words = append(result.Words, word)
	}

	return result
}

// parseTranscriptionResponse converts the ogen response to our response type.
func (s *SpeechToTextService) parseTranscriptionResponse(resp api.SpeechToTextRes) (*TranscriptionResponse, error) {
	switch r := resp.(type) {
	case *api.SpeechToTextOK:
		if !r.IsSpeechToTextChunkResponseModel() {
			return nil, &APIError{Message: "unexpected response format"}
		}
		chunk := r.SpeechToTextChunkResponseModel

		result := &TranscriptionResponse{
			Text:         chunk.Text,
			LanguageCode: chunk.LanguageCode,
		}

		for _, w := range chunk.Words {
			word := TranscriptionWord{
				Text: w.Text,
				Type: string(w.Type),
			}
			if w.Start.Set && !w.Start.Null {
				word.Start = w.Start.Value
			}
			if w.End.Set && !w.End.Null {
				word.End = w.End.Value
			}
			if w.SpeakerID.Set && !w.SpeakerID.Null {
				word.Speaker = w.SpeakerID.Value
			}
			result.Words = append(result.Words, word)
		}

		return result, nil
	case *api.HTTPValidationError:
		if len(r.Detail) > 0 {
			return nil, &ValidationError{
				Field:   formatValidationLocation(r.Detail[0].Loc),
				Message: r.Detail[0].Msg,
			}
		}
		return nil, &APIError{StatusCode: 422, Message: "validation error"}
	case *api.SpeechToTextAcceptedApplicationJSON:
		return nil, &APIError{StatusCode: 202, Message: "async transcription not supported"}
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// TranscribeURL transcribes audio from a URL.
func (s *SpeechToTextService) TranscribeURL(ctx context.Context, url string) (*TranscriptionResponse, error) {
	return s.Transcribe(ctx, &TranscriptionRequest{FileURL: url})
}

// TranscribeWithDiarization transcribes audio with speaker identification.
func (s *SpeechToTextService) TranscribeWithDiarization(ctx context.Context, url string) (*TranscriptionResponse, error) {
	return s.Transcribe(ctx, &TranscriptionRequest{
		FileURL: url,
		Diarize: true,
	})
}

// TranscribeFile transcribes audio from raw file bytes.
func (s *SpeechToTextService) TranscribeFile(ctx context.Context, fileBytes []byte) (*TranscriptionResponse, error) {
	if len(fileBytes) == 0 {
		return nil, &ValidationError{Field: "file", Message: "file bytes cannot be empty"}
	}

	return s.Transcribe(ctx, &TranscriptionRequest{
		FileBytes: fileBytes,
	})
}
