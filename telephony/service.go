package telephony

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Service handles Twilio and SIP phone integration for conversational AI.
type Service struct {
	apiKey  string
	baseURL string
}

// New creates a new telephony service.
func New(apiKey, baseURL string) *Service {
	return &Service{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// postJSON is a helper for making JSON POST requests.
func (s *Service) postJSON(ctx context.Context, path string, req any, result any) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("telephony: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		s.baseURL+path,
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("xi-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return fmt.Errorf("telephony: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telephony: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("telephony: failed to decode response: %w", err)
	}

	return nil
}

// --- Twilio Integration ---

// RegisterCallRequest is the request to register an incoming Twilio call.
type RegisterCallRequest struct {
	// AgentID is the ElevenLabs agent ID to handle the call.
	AgentID string `json:"agent_id"`

	// AgentPhoneNumberID is the ElevenLabs phone number ID (if using imported number).
	AgentPhoneNumberID string `json:"agent_phone_number_id,omitempty"`

	// CustomLLMExtraBody is additional data to pass to the LLM.
	CustomLLMExtraBody map[string]any `json:"custom_llm_extra_body,omitempty"`

	// DynamicVariables are variables to inject into the agent prompt.
	DynamicVariables map[string]string `json:"dynamic_variables,omitempty"`

	// FirstMessage overrides the agent's default first message.
	FirstMessage string `json:"first_message,omitempty"`

	// SystemPrompt overrides the agent's system prompt.
	SystemPrompt string `json:"system_prompt,omitempty"`
}

// RegisterCallResponse is the response from registering a call.
type RegisterCallResponse struct {
	// TwiML is the TwiML response to return to Twilio.
	TwiML string `json:"twiml"`

	// ConversationID is the ElevenLabs conversation ID for this call.
	ConversationID string `json:"conversation_id,omitempty"`
}

// OutboundCallRequest is the request to make an outbound call via Twilio.
type OutboundCallRequest struct {
	// AgentID is the ElevenLabs agent ID to handle the call.
	AgentID string `json:"agent_id"`

	// AgentPhoneNumberID is the ElevenLabs phone number ID to call from.
	AgentPhoneNumberID string `json:"agent_phone_number_id"`

	// ToNumber is the phone number to call (E.164 format).
	ToNumber string `json:"to_number"`

	// CustomLLMExtraBody is additional data to pass to the LLM.
	CustomLLMExtraBody map[string]any `json:"custom_llm_extra_body,omitempty"`

	// DynamicVariables are variables to inject into the agent prompt.
	DynamicVariables map[string]string `json:"dynamic_variables,omitempty"`

	// FirstMessage overrides the agent's default first message.
	FirstMessage string `json:"first_message,omitempty"`

	// SystemPrompt overrides the agent's system prompt.
	SystemPrompt string `json:"system_prompt,omitempty"`
}

// OutboundCallResponse is the response from making an outbound call.
type OutboundCallResponse struct {
	// CallSID is the Twilio call SID.
	CallSID string `json:"call_sid"`

	// ConversationID is the ElevenLabs conversation ID for this call.
	ConversationID string `json:"conversation_id"`

	// Status is the initial call status.
	Status string `json:"status"`
}

// SIPOutboundCallRequest is the request to make an outbound call via SIP trunk.
type SIPOutboundCallRequest struct {
	// AgentID is the ElevenLabs agent ID to handle the call.
	AgentID string `json:"agent_id"`

	// ToNumber is the phone number to call (E.164 format).
	ToNumber string `json:"to_number"`

	// SIPTrunkID is the SIP trunk ID to use.
	SIPTrunkID string `json:"sip_trunk_id"`

	// FromNumber is the caller ID to display (must be verified).
	FromNumber string `json:"from_number,omitempty"`

	// CustomLLMExtraBody is additional data to pass to the LLM.
	CustomLLMExtraBody map[string]any `json:"custom_llm_extra_body,omitempty"`

	// DynamicVariables are variables to inject into the agent prompt.
	DynamicVariables map[string]string `json:"dynamic_variables,omitempty"`

	// FirstMessage overrides the agent's default first message.
	FirstMessage string `json:"first_message,omitempty"`

	// SystemPrompt overrides the agent's system prompt.
	SystemPrompt string `json:"system_prompt,omitempty"`
}

// SIPOutboundCallResponse is the response from making a SIP outbound call.
type SIPOutboundCallResponse struct {
	// ConversationID is the ElevenLabs conversation ID for this call.
	ConversationID string `json:"conversation_id"`

	// Status is the initial call status.
	Status string `json:"status"`
}

// RegisterCall registers an incoming Twilio call with ElevenLabs.
// Returns TwiML that should be returned to Twilio's webhook.
func (s *Service) RegisterCall(ctx context.Context, req *RegisterCallRequest) (*RegisterCallResponse, error) {
	if req.AgentID == "" {
		return nil, fmt.Errorf("telephony: agent_id is required")
	}

	var result RegisterCallResponse
	if err := s.postJSON(ctx, "/v1/convai/twilio/register-call", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// OutboundCall initiates an outbound call via Twilio.
func (s *Service) OutboundCall(ctx context.Context, req *OutboundCallRequest) (*OutboundCallResponse, error) {
	if req.AgentID == "" {
		return nil, fmt.Errorf("telephony: agent_id is required")
	}
	if req.AgentPhoneNumberID == "" {
		return nil, fmt.Errorf("telephony: agent_phone_number_id is required")
	}
	if req.ToNumber == "" {
		return nil, fmt.Errorf("telephony: to_number is required")
	}

	var result OutboundCallResponse
	if err := s.postJSON(ctx, "/v1/convai/twilio/outbound-call", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SIPOutboundCall initiates an outbound call via SIP trunk.
func (s *Service) SIPOutboundCall(ctx context.Context, req *SIPOutboundCallRequest) (*SIPOutboundCallResponse, error) {
	if req.AgentID == "" {
		return nil, fmt.Errorf("telephony: agent_id is required")
	}
	if req.SIPTrunkID == "" {
		return nil, fmt.Errorf("telephony: sip_trunk_id is required")
	}
	if req.ToNumber == "" {
		return nil, fmt.Errorf("telephony: to_number is required")
	}

	var result SIPOutboundCallResponse
	if err := s.postJSON(ctx, "/v1/convai/sip-trunk/outbound-call", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// --- Phone Numbers ---

// PhoneNumber represents an ElevenLabs phone number.
type PhoneNumber struct {
	ID          string `json:"phone_number_id"`
	PhoneNumber string `json:"phone_number"`
	Label       string `json:"label"`
	AgentID     string `json:"agent_id,omitempty"`
	Provider    string `json:"provider"` // "twilio", "sip"
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// ListPhoneNumbersResponse is the response from listing phone numbers.
type ListPhoneNumbersResponse struct {
	PhoneNumbers []PhoneNumber `json:"phone_numbers"`
}

// ListPhoneNumbers lists all phone numbers in the workspace.
func (s *Service) ListPhoneNumbers(ctx context.Context) ([]PhoneNumber, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "GET",
		s.baseURL+"/v1/convai/phone-numbers",
		nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("xi-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return nil, fmt.Errorf("telephony: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("telephony: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result ListPhoneNumbersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("telephony: failed to decode response: %w", err)
	}

	return result.PhoneNumbers, nil
}

// GetPhoneNumber retrieves a specific phone number by ID.
func (s *Service) GetPhoneNumber(ctx context.Context, phoneNumberID string) (*PhoneNumber, error) {
	if phoneNumberID == "" {
		return nil, fmt.Errorf("telephony: phone_number_id is required")
	}

	httpReq, err := http.NewRequestWithContext(ctx, "GET",
		s.baseURL+"/v1/convai/phone-numbers/"+phoneNumberID,
		nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("xi-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return nil, fmt.Errorf("telephony: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("telephony: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result PhoneNumber
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("telephony: failed to decode response: %w", err)
	}

	return &result, nil
}

// UpdatePhoneNumberRequest is the request to update a phone number.
type UpdatePhoneNumberRequest struct {
	// Label is a descriptive label for the phone number.
	Label string `json:"label,omitempty"`

	// AgentID is the agent to associate with this phone number.
	AgentID string `json:"agent_id,omitempty"`
}

// UpdatePhoneNumber updates a phone number's settings.
func (s *Service) UpdatePhoneNumber(ctx context.Context, phoneNumberID string, req *UpdatePhoneNumberRequest) (*PhoneNumber, error) {
	if phoneNumberID == "" {
		return nil, fmt.Errorf("telephony: phone_number_id is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("telephony: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "PATCH",
		s.baseURL+"/v1/convai/phone-numbers/"+phoneNumberID,
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("xi-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return nil, fmt.Errorf("telephony: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("telephony: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result PhoneNumber
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("telephony: failed to decode response: %w", err)
	}

	return &result, nil
}

// DeletePhoneNumber removes a phone number from the workspace.
func (s *Service) DeletePhoneNumber(ctx context.Context, phoneNumberID string) error {
	if phoneNumberID == "" {
		return fmt.Errorf("telephony: phone_number_id is required")
	}

	httpReq, err := http.NewRequestWithContext(ctx, "DELETE",
		s.baseURL+"/v1/convai/phone-numbers/"+phoneNumberID,
		nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("xi-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return fmt.Errorf("telephony: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telephony: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
