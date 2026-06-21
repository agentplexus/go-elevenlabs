// Package elevenlabs provides a Go client for the ElevenLabs API.
//
// The client wraps the ogen-generated API client with a higher-level
// interface that handles authentication and provides convenient methods
// for common operations.
package elevenlabs

import (
	"net/http"
	"os"
	"time"

	"github.com/plexusone/elevenlabs-go/account"
	"github.com/plexusone/elevenlabs-go/agents"
	"github.com/plexusone/elevenlabs-go/audio"
	"github.com/plexusone/elevenlabs-go/content"
	"github.com/plexusone/elevenlabs-go/internal/api"
	"github.com/plexusone/elevenlabs-go/realtime"
	"github.com/plexusone/elevenlabs-go/stt"
	"github.com/plexusone/elevenlabs-go/telephony"
	"github.com/plexusone/elevenlabs-go/tts"
	"github.com/plexusone/elevenlabs-go/voice"
)

// Version is the SDK version.
const Version = "0.13.0"

// DefaultBaseURL is the default ElevenLabs API base URL.
const DefaultBaseURL = "https://api.elevenlabs.io"

// DefaultModelID is the recommended model for text-to-speech.
const DefaultModelID = "eleven_multilingual_v2"

// Client is the main ElevenLabs client for interacting with the API.
type Client struct {
	apiClient *api.Client
	apiKey    string
	baseURL   string

	// Domain-based service accessors
	ttsSvc       *tts.Service
	sttSvc       *stt.Service
	voiceSvc     *voice.Service
	audioSvc     *audio.Service
	contentSvc   *content.Service
	realtimeSvc  *realtime.Service
	telephonySvc *telephony.Service
	accountSvc   *account.Service
	agentsSvc    *agents.Service
}

// NewClient creates a new ElevenLabs client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	options := defaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Try environment variable if API key not set
	if options.apiKey == "" {
		options.apiKey = os.Getenv("ELEVENLABS_API_KEY")
	}

	// Create HTTP client with auth headers
	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.timeout,
		}
	}

	// Wrap with auth transport
	authClient := &authHTTPClient{
		client: httpClient,
		apiKey: options.apiKey,
	}

	// Create the ogen client
	apiClient, err := api.NewClient(
		options.baseURL,
		api.WithClient(authClient),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{
		apiClient: apiClient,
		apiKey:    options.apiKey,
		baseURL:   options.baseURL,
	}

	// Initialize domain-based services
	c.ttsSvc = tts.New(apiClient, options.apiKey, options.baseURL)
	c.sttSvc = stt.New(apiClient, options.apiKey, options.baseURL)
	c.voiceSvc = voice.New(apiClient)
	c.audioSvc = audio.New(apiClient)
	c.contentSvc = content.New(apiClient)
	c.realtimeSvc = realtime.New(options.apiKey, options.baseURL)
	c.telephonySvc = telephony.New(options.apiKey, options.baseURL)
	c.accountSvc = account.New(apiClient)
	c.agentsSvc = agents.New(c)

	return c, nil
}

// authHTTPClient wraps an http.Client to add authentication headers.
type authHTTPClient struct {
	client *http.Client
	apiKey string
}

// Do implements ht.Client interface.
func (c *authHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add authentication header
	if c.apiKey != "" {
		req.Header.Set("xi-api-key", c.apiKey)
	}

	// Add SDK version headers
	req.Header.Set("X-ElevenLabs-SDK-Version", Version)
	req.Header.Set("X-ElevenLabs-SDK-Lang", "go")

	return c.client.Do(req) //nolint:gosec // G704: HTTP client library, URL is caller-controlled by design
}

// API returns the underlying ogen-generated API client for advanced usage.
// Use this when you need access to API endpoints not covered by the
// high-level wrapper methods.
func (c *Client) API() *api.Client {
	return c.apiClient
}

// TTS returns the text-to-speech service.
// This includes TTS generation, speech-to-speech conversion, and dialogue generation.
func (c *Client) TTS() *tts.Service {
	return c.ttsSvc
}

// STT returns the speech-to-text service.
func (c *Client) STT() *stt.Service {
	return c.sttSvc
}

// Voice returns the voice management service.
// This includes listing voices, getting settings, and AI voice design.
func (c *Client) Voice() *voice.Service {
	return c.voiceSvc
}

// Audio returns the audio processing service.
// This includes audio isolation, forced alignment, and sound effects generation.
func (c *Client) Audio() *audio.Service {
	return c.audioSvc
}

// Content returns the content generation service.
// This includes music generation, dubbing, and Studio projects.
func (c *Client) Content() *content.Service {
	return c.contentSvc
}

// Realtime returns the real-time streaming service.
// This includes WebSocket-based TTS and STT connections.
func (c *Client) Realtime() *realtime.Service {
	return c.realtimeSvc
}

// Telephony returns the telephony integration service.
// This includes Twilio/SIP integration and phone number management.
func (c *Client) Telephony() *telephony.Service {
	return c.telephonySvc
}

// Account returns the account management service.
// This includes user info, models, history, and pronunciation dictionaries.
func (c *Client) Account() *account.Service {
	return c.accountSvc
}

// Agents returns the ElevenAgents (Conversational AI) service.
// This includes agent management, conversations, knowledge base, testing, and batch calling.
func (c *Client) Agents() *agents.Service {
	return c.agentsSvc
}

// APIKey returns the API key used by the client.
func (c *Client) APIKey() string {
	return c.apiKey
}

// BaseURL returns the base URL used by the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// clientOptions holds the options for creating a Client.
type clientOptions struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		baseURL: DefaultBaseURL,
		timeout: 120 * time.Second, // TTS can take a while
	}
}

// Option is a functional option for configuring the Client.
type Option func(*clientOptions)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(o *clientOptions) {
		o.apiKey = apiKey
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}
