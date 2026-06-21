package realtime

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// Service handles real-time WebSocket TTS and STT.
type Service struct {
	apiKey  string
	baseURL string
}

// New creates a new realtime service.
func New(apiKey, baseURL string) *Service {
	return &Service{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// --- WebSocket TTS ---

// TTSOptions configures the WebSocket TTS connection.
type TTSOptions struct {
	// ModelID is the model to use. Defaults to "eleven_turbo_v2_5" for low latency.
	ModelID string

	// OutputFormat specifies the audio output format.
	// Recommended for real-time: "pcm_16000", "pcm_22050", "pcm_24000", "pcm_44100"
	OutputFormat string

	// VoiceSettings configures the voice parameters.
	VoiceSettings *VoiceSettings

	// OptimizeStreamingLatency reduces latency at the cost of quality (0-4).
	OptimizeStreamingLatency int

	// EnableSSMLParsing enables SSML parsing for the input text.
	EnableSSMLParsing bool

	// LanguageCode is the ISO language code (e.g., "en", "es").
	LanguageCode string

	// ChunkLengthSchedule controls text chunking for audio generation.
	ChunkLengthSchedule []int

	// InactivityTimeout is the context timeout in seconds (default 20).
	InactivityTimeout int

	// PronunciationDictionaryIDs is a list of pronunciation dictionary IDs to use.
	PronunciationDictionaryIDs []string
}

// VoiceSettings contains voice configuration for TTS.
type VoiceSettings struct {
	Stability       float64
	SimilarityBoost float64
	Style           float64
	Speed           float64
	UseSpeakerBoost bool
}

// DefaultVoiceSettings returns sensible default voice settings for real-time TTS.
func DefaultVoiceSettings() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Style:           0.0,
		UseSpeakerBoost: true,
	}
}

// DefaultTTSOptions returns default options optimized for low latency.
func DefaultTTSOptions() *TTSOptions {
	return &TTSOptions{
		ModelID:                  "eleven_turbo_v2_5",
		OutputFormat:             "pcm_16000",
		OptimizeStreamingLatency: 3,
	}
}

// TTSConnection represents an active WebSocket TTS connection.
type TTSConnection struct {
	conn    *websocket.Conn
	voiceID string
	options *TTSOptions
	mu      sync.Mutex
	closed  bool
	flushed bool

	audioOut  chan []byte
	alignOut  chan *TTSAlignment
	errChan   chan error
	doneChan  chan struct{}
	closeChan chan struct{}
	closeOnce sync.Once
	doneOnce  sync.Once
}

// TTSAlignment contains word-level timing information.
type TTSAlignment struct {
	Characters     []string  `json:"characters"`
	CharacterStart []float64 `json:"character_start_times_seconds"`
	CharacterEnd   []float64 `json:"character_end_times_seconds"`
}

type ttsWSMessage struct {
	Text                       string           `json:"text,omitempty"`
	VoiceSettings              *wsVoiceSettings `json:"voice_settings,omitempty"`
	GenerationConfig           *wsGenConfig     `json:"generation_config,omitempty"`
	XIAPIKey                   string           `json:"xi_api_key,omitempty"`
	TryTriggerGeneration       bool             `json:"try_trigger_generation,omitempty"`
	Flush                      bool             `json:"flush,omitempty"`
	CloseConnection            bool             `json:"close_connection,omitempty"`
	ContextID                  string           `json:"context_id,omitempty"`
	PronunciationDictionaryIDs []string         `json:"pronunciation_dictionary_locators,omitempty"`
}

type wsVoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           float64 `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

type wsGenConfig struct {
	ChunkLengthSchedule []int `json:"chunk_length_schedule,omitempty"`
}

type ttsWSResponse struct {
	Audio               string        `json:"audio,omitempty"`
	IsFinal             bool          `json:"isFinal,omitempty"`
	NormalizedAlignment *TTSAlignment `json:"normalizedAlignment,omitempty"`
	Alignment           *TTSAlignment `json:"alignment,omitempty"`
	Error               string        `json:"error,omitempty"`
	Message             string        `json:"message,omitempty"`
	Code                int           `json:"code,omitempty"`
}

// ConnectTTS establishes a WebSocket connection for real-time TTS.
func (s *Service) ConnectTTS(ctx context.Context, voiceID string, opts *TTSOptions) (*TTSConnection, error) {
	if voiceID == "" {
		return nil, fmt.Errorf("realtime: voice ID is required")
	}

	if opts == nil {
		opts = DefaultTTSOptions()
	}

	// Build WebSocket URL
	wsURL, err := s.buildTTSWebSocketURL(voiceID, opts)
	if err != nil {
		return nil, err
	}

	// Create dialer with context
	dialer := websocket.Dialer{
		HandshakeTimeout: 0,
	}

	// Add headers
	headers := http.Header{}
	headers.Set("xi-api-key", s.apiKey)

	// Connect
	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("realtime: websocket dial failed: %w", err)
	}

	wsc := &TTSConnection{
		conn:      conn,
		voiceID:   voiceID,
		options:   opts,
		audioOut:  make(chan []byte, 100),
		alignOut:  make(chan *TTSAlignment, 100),
		errChan:   make(chan error, 1),
		doneChan:  make(chan struct{}),
		closeChan: make(chan struct{}),
	}

	// Send initial configuration
	if err := wsc.sendInit(); err != nil {
		conn.Close()
		return nil, err
	}

	// Start reading responses
	go wsc.readLoop()

	return wsc, nil
}

func (s *Service) buildTTSWebSocketURL(voiceID string, opts *TTSOptions) (string, error) {
	baseURL := s.baseURL
	if baseURL == "" {
		baseURL = "https://api.elevenlabs.io"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = fmt.Sprintf("/v1/text-to-speech/%s/stream-input", voiceID)

	q := u.Query()
	if opts.ModelID != "" {
		q.Set("model_id", opts.ModelID)
	}
	if opts.OutputFormat != "" {
		q.Set("output_format", opts.OutputFormat)
	}
	if opts.OptimizeStreamingLatency > 0 {
		q.Set("optimize_streaming_latency", fmt.Sprintf("%d", opts.OptimizeStreamingLatency))
	}
	if opts.EnableSSMLParsing {
		q.Set("enable_ssml_parsing", "true")
	}
	if opts.LanguageCode != "" {
		q.Set("language_code", opts.LanguageCode)
	}
	if opts.InactivityTimeout > 0 {
		q.Set("inactivity_timeout", fmt.Sprintf("%d", opts.InactivityTimeout))
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (wsc *TTSConnection) sendInit() error {
	msg := ttsWSMessage{
		Text: " ",
	}

	if wsc.options.VoiceSettings != nil {
		msg.VoiceSettings = &wsVoiceSettings{
			Stability:       wsc.options.VoiceSettings.Stability,
			SimilarityBoost: wsc.options.VoiceSettings.SimilarityBoost,
			Style:           wsc.options.VoiceSettings.Style,
			UseSpeakerBoost: wsc.options.VoiceSettings.UseSpeakerBoost,
		}
	}

	if len(wsc.options.ChunkLengthSchedule) > 0 {
		msg.GenerationConfig = &wsGenConfig{
			ChunkLengthSchedule: wsc.options.ChunkLengthSchedule,
		}
	}

	if len(wsc.options.PronunciationDictionaryIDs) > 0 {
		msg.PronunciationDictionaryIDs = wsc.options.PronunciationDictionaryIDs
	}

	return wsc.sendJSON(msg)
}

func (wsc *TTSConnection) sendJSON(msg any) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return fmt.Errorf("realtime: connection closed")
	}

	return wsc.conn.WriteJSON(msg)
}

func (wsc *TTSConnection) readLoop() {
	defer wsc.closeChannels()

	for {
		select {
		case <-wsc.closeChan:
			return
		default:
		}

		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				select {
				case wsc.errChan <- err:
				default:
				}
			}
			return
		}

		var resp ttsWSResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			select {
			case wsc.errChan <- fmt.Errorf("realtime: failed to parse response: %w", err):
			default:
			}
			continue
		}

		if resp.Error != "" || resp.Message != "" {
			errMsg := resp.Error
			if errMsg == "" {
				errMsg = resp.Message
			}
			select {
			case wsc.errChan <- fmt.Errorf("realtime: server error: %s", errMsg):
			default:
			}
			continue
		}

		if resp.Audio != "" {
			audioBytes, err := base64.StdEncoding.DecodeString(resp.Audio)
			if err != nil {
				select {
				case wsc.errChan <- fmt.Errorf("realtime: failed to decode audio: %w", err):
				default:
				}
				continue
			}
			if len(audioBytes) > 0 {
				select {
				case wsc.audioOut <- audioBytes:
				case <-wsc.closeChan:
					return
				}
			}
		}

		if resp.NormalizedAlignment != nil {
			select {
			case wsc.alignOut <- resp.NormalizedAlignment:
			default:
			}
		} else if resp.Alignment != nil {
			select {
			case wsc.alignOut <- resp.Alignment:
			default:
			}
		}

		if resp.IsFinal {
			wsc.mu.Lock()
			flushed := wsc.flushed
			wsc.mu.Unlock()

			if flushed {
				wsc.doneOnce.Do(func() {
					close(wsc.doneChan)
				})
			}
		}
	}
}

func (wsc *TTSConnection) closeChannels() {
	wsc.closeOnce.Do(func() {
		close(wsc.closeChan)
		close(wsc.audioOut)
		close(wsc.alignOut)
	})
}

// SendText sends text to be converted to speech.
func (wsc *TTSConnection) SendText(text string) error {
	if text == "" {
		return nil
	}
	return wsc.sendJSON(ttsWSMessage{Text: text})
}

// TriggerGeneration forces audio generation for buffered text.
func (wsc *TTSConnection) TriggerGeneration() error {
	return wsc.sendJSON(ttsWSMessage{Text: " ", TryTriggerGeneration: true})
}

// Flush signals that no more text will be sent and flushes remaining audio.
func (wsc *TTSConnection) Flush() error {
	wsc.mu.Lock()
	wsc.flushed = true
	wsc.mu.Unlock()

	return wsc.sendJSON(ttsWSMessage{Text: "", Flush: true})
}

// Audio returns a channel that receives audio chunks.
func (wsc *TTSConnection) Audio() <-chan []byte {
	return wsc.audioOut
}

// Alignments returns a channel that receives word alignment information.
func (wsc *TTSConnection) Alignments() <-chan *TTSAlignment {
	return wsc.alignOut
}

// Errors returns a channel that receives errors.
func (wsc *TTSConnection) Errors() <-chan error {
	return wsc.errChan
}

// Done returns a channel that is closed when all audio has been received after Flush().
func (wsc *TTSConnection) Done() <-chan struct{} {
	return wsc.doneChan
}

// Close closes the WebSocket connection gracefully.
func (wsc *TTSConnection) Close() error {
	wsc.mu.Lock()
	if wsc.closed {
		wsc.mu.Unlock()
		return nil
	}
	wsc.closed = true
	wsc.mu.Unlock()

	_ = wsc.sendJSON(ttsWSMessage{CloseConnection: true})
	wsc.closeChannels()
	return wsc.conn.Close()
}

// StreamText is a convenience method that sends all text from a channel and returns audio.
// It handles flushing automatically when the input channel closes.
func (wsc *TTSConnection) StreamText(ctx context.Context, textStream <-chan string) (<-chan []byte, <-chan error) {
	audioOut := make(chan []byte, 100)
	errOut := make(chan error, 1)

	go func() {
		defer close(audioOut)
		defer close(errOut)

		// Forward audio from connection
		done := make(chan struct{})
		go func() {
			defer close(done)
			for audio := range wsc.Audio() {
				select {
				case audioOut <- audio:
				case <-ctx.Done():
					return
				}
			}
		}()

		// Send text as it arrives
		for {
			select {
			case text, ok := <-textStream:
				if !ok {
					// Input stream closed, flush and wait for remaining audio
					if err := wsc.Flush(); err != nil {
						errOut <- err
						return
					}
					<-done
					return
				}
				if err := wsc.SendText(text); err != nil {
					errOut <- err
					return
				}
			case err := <-wsc.Errors():
				errOut <- err
				return
			case <-ctx.Done():
				errOut <- ctx.Err()
				return
			}
		}
	}()

	return audioOut, errOut
}

// --- WebSocket STT ---

// STTOptions configures the WebSocket STT connection.
type STTOptions struct {
	// ModelID is the transcription model to use.
	ModelID string

	// AudioFormat specifies the audio encoding format.
	AudioFormat string

	// LanguageCode is the expected language.
	LanguageCode string

	// IncludeTimestamps enables word-level timing information.
	IncludeTimestamps bool

	// IncludeLanguageDetection includes detected language in responses.
	IncludeLanguageDetection bool

	// CommitStrategy determines how transcripts are committed.
	CommitStrategy string

	// VAD settings
	VADSilenceThresholdSecs float64
	VADThreshold            float64
	MinSpeechDurationMs     int
	MinSilenceDurationMs    int
}

// DefaultSTTOptions returns default options for real-time STT.
func DefaultSTTOptions() *STTOptions {
	return &STTOptions{
		ModelID:           "scribe_v2_realtime",
		AudioFormat:       "pcm_16000",
		CommitStrategy:    "manual",
		IncludeTimestamps: true,
	}
}

// STTConnection represents an active WebSocket STT connection.
type STTConnection struct {
	conn      *websocket.Conn
	options   *STTOptions
	sessionID string
	mu        sync.Mutex
	closed    bool

	transcriptOut chan *STTTranscript
	errChan       chan error
	closeChan     chan struct{}
	closeOnce     sync.Once
}

// STTTranscript represents a transcription result.
type STTTranscript struct {
	Text         string    `json:"text"`
	IsFinal      bool      `json:"is_final"`
	Words        []STTWord `json:"words,omitempty"`
	LanguageCode string    `json:"language_code,omitempty"`
}

// STTWord represents a single word with timing.
type STTWord struct {
	Text      string  `json:"text"`
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	Type      string  `json:"type,omitempty"`
	SpeakerID string  `json:"speaker_id,omitempty"`
}

type sttInputAudioChunk struct {
	MessageType  string `json:"message_type"`
	AudioBase64  string `json:"audio_base_64"`
	Commit       bool   `json:"commit,omitempty"`
	SampleRate   int    `json:"sample_rate,omitempty"`
	PreviousText string `json:"previous_text,omitempty"`
}

type sttResponse struct {
	MessageType  string    `json:"message_type"`
	Text         string    `json:"text,omitempty"`
	LanguageCode string    `json:"language_code,omitempty"`
	Words        []STTWord `json:"words,omitempty"`
	Error        string    `json:"error,omitempty"`
	SessionID    string    `json:"session_id,omitempty"`
}

// ConnectSTT establishes a WebSocket connection for real-time STT.
func (s *Service) ConnectSTT(ctx context.Context, opts *STTOptions) (*STTConnection, error) {
	if opts == nil {
		opts = DefaultSTTOptions()
	}

	wsURL, err := s.buildSTTWebSocketURL(opts)
	if err != nil {
		return nil, err
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 0,
	}

	headers := http.Header{}
	headers.Set("xi-api-key", s.apiKey)

	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("realtime: websocket dial failed: %w", err)
	}

	wsc := &STTConnection{
		conn:          conn,
		options:       opts,
		transcriptOut: make(chan *STTTranscript, 100),
		errChan:       make(chan error, 1),
		closeChan:     make(chan struct{}),
	}

	go wsc.readLoop()

	return wsc, nil
}

func (s *Service) buildSTTWebSocketURL(opts *STTOptions) (string, error) {
	baseURL := s.baseURL
	if baseURL == "" {
		baseURL = "https://api.elevenlabs.io"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	u.Path = "/v1/speech-to-text/realtime"

	q := u.Query()
	if opts.ModelID != "" {
		q.Set("model_id", opts.ModelID)
	}
	if opts.AudioFormat != "" {
		q.Set("audio_format", opts.AudioFormat)
	}
	if opts.LanguageCode != "" {
		q.Set("language_code", opts.LanguageCode)
	}
	if opts.IncludeTimestamps {
		q.Set("include_timestamps", "true")
	}
	if opts.IncludeLanguageDetection {
		q.Set("include_language_detection", "true")
	}
	if opts.CommitStrategy != "" {
		q.Set("commit_strategy", opts.CommitStrategy)
	}

	if opts.CommitStrategy == "vad" {
		if opts.VADSilenceThresholdSecs > 0 {
			q.Set("vad_silence_threshold_secs", fmt.Sprintf("%.2f", opts.VADSilenceThresholdSecs))
		}
		if opts.VADThreshold > 0 {
			q.Set("vad_threshold", fmt.Sprintf("%.2f", opts.VADThreshold))
		}
		if opts.MinSpeechDurationMs > 0 {
			q.Set("min_speech_duration_ms", fmt.Sprintf("%d", opts.MinSpeechDurationMs))
		}
		if opts.MinSilenceDurationMs > 0 {
			q.Set("min_silence_duration_ms", fmt.Sprintf("%d", opts.MinSilenceDurationMs))
		}
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (wsc *STTConnection) sendJSON(msg any) error {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()

	if wsc.closed {
		return fmt.Errorf("realtime: connection closed")
	}

	return wsc.conn.WriteJSON(msg)
}

func (wsc *STTConnection) readLoop() {
	defer wsc.closeChannels()

	for {
		select {
		case <-wsc.closeChan:
			return
		default:
		}

		_, message, err := wsc.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				select {
				case wsc.errChan <- err:
				default:
				}
			}
			return
		}

		var resp sttResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			select {
			case wsc.errChan <- fmt.Errorf("realtime: failed to parse response: %w", err):
			default:
			}
			continue
		}

		switch resp.MessageType {
		case "session_started":
			wsc.mu.Lock()
			wsc.sessionID = resp.SessionID
			wsc.mu.Unlock()

		case "partial_transcript":
			transcript := &STTTranscript{Text: resp.Text, IsFinal: false}
			select {
			case wsc.transcriptOut <- transcript:
			case <-wsc.closeChan:
				return
			}

		case "committed_transcript":
			transcript := &STTTranscript{Text: resp.Text, IsFinal: true}
			select {
			case wsc.transcriptOut <- transcript:
			case <-wsc.closeChan:
				return
			}

		case "committed_transcript_with_timestamps":
			transcript := &STTTranscript{
				Text:         resp.Text,
				IsFinal:      true,
				LanguageCode: resp.LanguageCode,
				Words:        resp.Words,
			}
			select {
			case wsc.transcriptOut <- transcript:
			case <-wsc.closeChan:
				return
			}

		case "error", "auth_error", "quota_exceeded", "rate_limited",
			"input_error", "transcriber_error", "chunk_size_exceeded",
			"insufficient_audio_activity", "session_time_limit_exceeded",
			"resource_exhausted", "queue_overflow", "commit_throttled",
			"unaccepted_terms":
			errMsg := resp.Error
			if errMsg == "" {
				errMsg = resp.MessageType
			}
			select {
			case wsc.errChan <- fmt.Errorf("realtime: server error (%s): %s", resp.MessageType, errMsg):
			default:
			}
		}
	}
}

func (wsc *STTConnection) closeChannels() {
	wsc.closeOnce.Do(func() {
		close(wsc.closeChan)
		close(wsc.transcriptOut)
	})
}

// SendAudio sends audio data for transcription.
func (wsc *STTConnection) SendAudio(audio []byte) error {
	if len(audio) == 0 {
		return nil
	}

	return wsc.sendJSON(sttInputAudioChunk{
		MessageType: "input_audio_chunk",
		AudioBase64: base64.StdEncoding.EncodeToString(audio),
	})
}

// Commit forces a commit of the current transcript segment.
func (wsc *STTConnection) Commit() error {
	return wsc.sendJSON(sttInputAudioChunk{
		MessageType: "input_audio_chunk",
		AudioBase64: "",
		Commit:      true,
	})
}

// Transcripts returns a channel that receives transcription results.
func (wsc *STTConnection) Transcripts() <-chan *STTTranscript {
	return wsc.transcriptOut
}

// Errors returns a channel that receives errors.
func (wsc *STTConnection) Errors() <-chan error {
	return wsc.errChan
}

// SessionID returns the session ID assigned by the server.
func (wsc *STTConnection) SessionID() string {
	wsc.mu.Lock()
	defer wsc.mu.Unlock()
	return wsc.sessionID
}

// Close closes the WebSocket connection gracefully.
func (wsc *STTConnection) Close() error {
	wsc.mu.Lock()
	if wsc.closed {
		wsc.mu.Unlock()
		return nil
	}
	wsc.closed = true
	wsc.mu.Unlock()

	wsc.closeChannels()
	return wsc.conn.Close()
}

// StreamAudio is a convenience method that streams audio from a channel.
// It handles committing automatically when the input channel closes.
func (wsc *STTConnection) StreamAudio(ctx context.Context, audioStream <-chan []byte) (<-chan *STTTranscript, <-chan error) {
	transcriptOut := make(chan *STTTranscript, 100)
	errOut := make(chan error, 1)

	go func() {
		defer close(transcriptOut)
		defer close(errOut)

		// Forward transcripts from connection
		done := make(chan struct{})
		go func() {
			defer close(done)
			for transcript := range wsc.Transcripts() {
				select {
				case transcriptOut <- transcript:
				case <-ctx.Done():
					return
				}
			}
		}()

		// Send audio as it arrives
		for {
			select {
			case audio, ok := <-audioStream:
				if !ok {
					// Input stream closed, commit final transcript and wait
					if err := wsc.Commit(); err != nil {
						// Commit error is non-fatal, connection might be closing
					}
					<-done
					return
				}
				if err := wsc.SendAudio(audio); err != nil {
					errOut <- err
					return
				}
			case err := <-wsc.Errors():
				errOut <- err
				return
			case <-ctx.Done():
				errOut <- ctx.Err()
				return
			}
		}
	}()

	return transcriptOut, errOut
}
