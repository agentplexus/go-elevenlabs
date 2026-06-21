package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// SimulationConnection represents a streaming conversation simulation.
type SimulationConnection struct {
	conn      *websocket.Conn
	agentID   string
	mu        sync.Mutex
	closed    bool
	msgOut    chan *SimulationMessage
	errChan   chan error
	closeChan chan struct{}
	closeOnce sync.Once
}

// SimulationMessage represents a message in the simulation.
type SimulationMessage struct {
	// Role indicates who sent the message ("user" or "agent").
	Role string

	// Content is the text content of the message.
	Content string

	// Audio contains optional audio data (base64 encoded).
	Audio []byte

	// Metadata contains additional message metadata.
	Metadata map[string]any
}

// SimulateConversationOptions configures the simulation.
type SimulateConversationOptions struct {
	// SimulationPersona describes the simulated user persona.
	SimulationPersona string

	// InitialMessage is the first message to send.
	InitialMessage string

	// MaxTurns limits the conversation length.
	MaxTurns int

	// IncludeAudio requests audio output.
	IncludeAudio bool
}

// SimulateConversation runs a non-streaming conversation simulation.
func (s *Service) SimulateConversation(ctx context.Context, agentID string, opts *SimulateConversationOptions) ([]*SimulationMessage, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	req := map[string]any{}
	if opts != nil {
		if opts.SimulationPersona != "" {
			req["simulation_persona"] = opts.SimulationPersona
		}
		if opts.InitialMessage != "" {
			req["initial_message"] = opts.InitialMessage
		}
		if opts.MaxTurns > 0 {
			req["max_turns"] = opts.MaxTurns
		}
		if opts.IncludeAudio {
			req["include_audio"] = true
		}
	}

	var result struct {
		Messages []struct {
			Role     string         `json:"role"`
			Content  string         `json:"content"`
			Audio    []byte         `json:"audio,omitempty"`
			Metadata map[string]any `json:"metadata,omitempty"`
		} `json:"messages"`
	}

	if err := s.doJSON(ctx, "POST", "/v1/convai/agents/"+agentID+"/simulate-conversation", req, &result); err != nil {
		return nil, err
	}

	messages := make([]*SimulationMessage, 0, len(result.Messages))
	for _, m := range result.Messages {
		messages = append(messages, &SimulationMessage{
			Role:     m.Role,
			Content:  m.Content,
			Audio:    m.Audio,
			Metadata: m.Metadata,
		})
	}

	return messages, nil
}

// SimulateConversationStream opens a streaming simulation connection.
func (s *Service) SimulateConversationStream(ctx context.Context, agentID string, opts *SimulateConversationOptions) (*SimulationConnection, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	wsURL, err := s.buildSimulationWebSocketURL(agentID)
	if err != nil {
		return nil, err
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 0,
	}

	headers := http.Header{}
	headers.Set("xi-api-key", s.client.APIKey())

	conn, _, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}

	sc := &SimulationConnection{
		conn:      conn,
		agentID:   agentID,
		msgOut:    make(chan *SimulationMessage, 100),
		errChan:   make(chan error, 1),
		closeChan: make(chan struct{}),
	}

	// Send initial configuration
	if opts != nil {
		initMsg := map[string]any{}
		if opts.SimulationPersona != "" {
			initMsg["simulation_persona"] = opts.SimulationPersona
		}
		if opts.InitialMessage != "" {
			initMsg["initial_message"] = opts.InitialMessage
		}
		if opts.MaxTurns > 0 {
			initMsg["max_turns"] = opts.MaxTurns
		}
		if opts.IncludeAudio {
			initMsg["include_audio"] = true
		}
		if len(initMsg) > 0 {
			if err := sc.sendJSON(initMsg); err != nil {
				conn.Close()
				return nil, err
			}
		}
	}

	go sc.readLoop()

	return sc, nil
}

func (s *Service) buildSimulationWebSocketURL(agentID string) (string, error) {
	baseURL := s.client.BaseURL()
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

	u.Path = fmt.Sprintf("/v1/convai/agents/%s/simulate-conversation-stream", agentID)

	return u.String(), nil
}

func (sc *SimulationConnection) sendJSON(msg any) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return fmt.Errorf("connection closed")
	}

	return sc.conn.WriteJSON(msg)
}

func (sc *SimulationConnection) readLoop() {
	defer sc.closeChannels()

	for {
		select {
		case <-sc.closeChan:
			return
		default:
		}

		_, message, err := sc.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				select {
				case sc.errChan <- err:
				default:
				}
			}
			return
		}

		var resp struct {
			Role     string         `json:"role"`
			Content  string         `json:"content"`
			Audio    []byte         `json:"audio,omitempty"`
			Metadata map[string]any `json:"metadata,omitempty"`
			Error    string         `json:"error,omitempty"`
		}

		if err := json.Unmarshal(message, &resp); err != nil {
			select {
			case sc.errChan <- fmt.Errorf("failed to parse response: %w", err):
			default:
			}
			continue
		}

		if resp.Error != "" {
			select {
			case sc.errChan <- fmt.Errorf("server error: %s", resp.Error):
			default:
			}
			continue
		}

		msg := &SimulationMessage{
			Role:     resp.Role,
			Content:  resp.Content,
			Audio:    resp.Audio,
			Metadata: resp.Metadata,
		}

		select {
		case sc.msgOut <- msg:
		case <-sc.closeChan:
			return
		}
	}
}

func (sc *SimulationConnection) closeChannels() {
	sc.closeOnce.Do(func() {
		close(sc.closeChan)
		close(sc.msgOut)
	})
}

// Messages returns a channel that receives simulation messages.
func (sc *SimulationConnection) Messages() <-chan *SimulationMessage {
	return sc.msgOut
}

// Errors returns a channel that receives errors from the connection.
func (sc *SimulationConnection) Errors() <-chan error {
	return sc.errChan
}

// SendMessage sends a user message in the simulation.
func (sc *SimulationConnection) SendMessage(content string) error {
	return sc.sendJSON(map[string]any{
		"role":    "user",
		"content": content,
	})
}

// Close closes the WebSocket connection.
func (sc *SimulationConnection) Close() error {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return nil
	}
	sc.closed = true
	sc.mu.Unlock()

	sc.closeChannels()
	return sc.conn.Close()
}
