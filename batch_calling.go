package elevenlabs

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// BatchCallingService handles batch outbound call operations.
type BatchCallingService struct {
	client *Client
}

// BatchCall represents a batch call job.
type BatchCall struct {
	ID                     string
	Name                   string
	AgentID                string
	AgentName              string
	BranchID               string
	BranchName             string
	Status                 string
	PhoneNumberID          string
	PhoneProvider          string
	Environment            string
	Timezone               string
	ScheduledTimeUnix      int
	CreatedAtUnix          int
	LastUpdatedAtUnix      int
	TotalCallsScheduled    int
	TotalCallsDispatched   int
	TotalCallsFinished     int
	RetryCount             int
	TargetConcurrencyLimit *int
}

// BatchCallRecipient represents a recipient for a batch call.
type BatchCallRecipient struct {
	// Phone number to call (E.164 format)
	PhoneNumber string
	// Optional unique ID for tracking
	ID string
	// Optional WhatsApp user ID for WhatsApp calls
	WhatsAppUserID string
	// Optional custom data to pass to the conversation
	CustomData map[string]any
}

// CreateBatchCallRequest contains options for creating a batch call.
type CreateBatchCallRequest struct {
	// Name for this batch call job
	Name string
	// Agent ID to use for calls
	AgentID string
	// Optional branch ID
	BranchID string
	// Phone number ID to use for outbound calls
	PhoneNumberID string
	// Recipients to call
	Recipients []BatchCallRecipient
	// Optional scheduled time (Unix timestamp)
	ScheduledTimeUnix *int
	// Optional timezone (e.g., "America/New_York")
	Timezone string
	// Optional environment (e.g., "production", "staging")
	Environment string
	// Optional maximum concurrent calls
	TargetConcurrencyLimit *int
}

// ListBatchCallsOptions contains options for listing batch calls.
type ListBatchCallsOptions struct {
	// Maximum number of results
	Limit int
	// Pagination cursor
	LastDoc string
	// Filter by agent ID
	AgentID string
}

// ListBatchCallsResponse contains the paginated list of batch calls.
type ListBatchCallsResponse struct {
	BatchCalls []*BatchCall
	LastDoc    string
}

// Create creates a new batch call job.
func (s *BatchCallingService) Create(ctx context.Context, req *CreateBatchCallRequest) (*BatchCall, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.AgentID == "" {
		return nil, errors.New("agent ID is required")
	}
	if len(req.Recipients) == 0 {
		return nil, errors.New("at least one recipient is required")
	}

	// Convert recipients
	recipients := make([]api.OutboundCallRecipient, len(req.Recipients))
	for i, r := range req.Recipients {
		rec := api.OutboundCallRecipient{}
		if r.PhoneNumber != "" {
			rec.PhoneNumber = api.OptNilString{Value: r.PhoneNumber, Set: true}
		}
		if r.ID != "" {
			rec.ID = api.OptNilString{Value: r.ID, Set: true}
		}
		if r.WhatsAppUserID != "" {
			rec.WhatsappUserID = api.OptNilString{Value: r.WhatsAppUserID, Set: true}
		}
		recipients[i] = rec
	}

	body := &api.BodySubmitABatchCallRequestV1ConvaiBatchCallingSubmitPost{
		CallName:   req.Name,
		AgentID:    req.AgentID,
		Recipients: recipients,
	}

	if req.BranchID != "" {
		body.BranchID = api.OptNilString{Value: req.BranchID, Set: true}
	}
	if req.PhoneNumberID != "" {
		body.AgentPhoneNumberID = api.OptNilString{Value: req.PhoneNumberID, Set: true}
	}
	if req.ScheduledTimeUnix != nil {
		body.ScheduledTimeUnix = api.OptNilInt{Value: *req.ScheduledTimeUnix, Set: true}
	}
	if req.Timezone != "" {
		body.Timezone = api.OptNilString{Value: req.Timezone, Set: true}
	}
	if req.Environment != "" {
		body.Environment = api.OptNilString{Value: req.Environment, Set: true}
	}
	if req.TargetConcurrencyLimit != nil {
		body.TargetConcurrencyLimit = api.OptNilInt{Value: *req.TargetConcurrencyLimit, Set: true}
	}

	resp, err := s.client.apiClient.CreateBatchCall(ctx, body, api.CreateBatchCallParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.BatchCallResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return batchCallFromResponse(result), nil
}

// Get retrieves a batch call by ID.
func (s *BatchCallingService) Get(ctx context.Context, batchID string) (*BatchCall, error) {
	if batchID == "" {
		return nil, errors.New("batch ID is required")
	}

	resp, err := s.client.apiClient.GetBatchCall(ctx, api.GetBatchCallParams{
		BatchID: batchID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.BatchCallDetailedResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return batchCallDetailedFromResponse(result), nil
}

// Cancel cancels a batch call job.
func (s *BatchCallingService) Cancel(ctx context.Context, batchID string) error {
	if batchID == "" {
		return errors.New("batch ID is required")
	}

	_, err := s.client.apiClient.CancelBatchCall(ctx, api.CancelBatchCallParams{
		BatchID: batchID,
	})
	return err
}

// Retry retries failed calls in a batch.
func (s *BatchCallingService) Retry(ctx context.Context, batchID string) error {
	if batchID == "" {
		return errors.New("batch ID is required")
	}

	_, err := s.client.apiClient.RetryBatchCall(ctx, api.RetryBatchCallParams{
		BatchID: batchID,
	})
	return err
}

// Delete deletes a batch call job.
func (s *BatchCallingService) Delete(ctx context.Context, batchID string) error {
	if batchID == "" {
		return errors.New("batch ID is required")
	}

	_, err := s.client.apiClient.DeleteBatchCall(ctx, api.DeleteBatchCallParams{
		BatchID: batchID,
	})
	return err
}

// List returns a paginated list of batch calls in the workspace.
func (s *BatchCallingService) List(ctx context.Context, opts *ListBatchCallsOptions) (*ListBatchCallsResponse, error) {
	params := api.GetWorkspaceBatchCallsParams{}

	if opts != nil {
		if opts.Limit > 0 {
			params.Limit = api.NewOptInt(opts.Limit)
		}
		if opts.LastDoc != "" {
			params.LastDoc = api.OptNilString{Value: opts.LastDoc, Set: true}
		}
		if opts.AgentID != "" {
			params.AgentID = api.OptNilString{Value: opts.AgentID, Set: true}
		}
	}

	resp, err := s.client.apiClient.GetWorkspaceBatchCalls(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.WorkspaceBatchCallsResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	batchCalls := make([]*BatchCall, len(result.BatchCalls))
	for i, bc := range result.BatchCalls {
		batchCalls[i] = batchCallFromResponse(&bc)
	}

	return &ListBatchCallsResponse{
		BatchCalls: batchCalls,
		LastDoc:    result.NextDoc.Value,
	}, nil
}

// batchCallFromResponse converts API response to BatchCall.
//
//nolint:dupl // Similar to batchCallDetailedFromResponse but different input type
func batchCallFromResponse(r *api.BatchCallResponse) *BatchCall {
	bc := &BatchCall{
		ID:                   r.ID,
		Name:                 r.Name,
		AgentID:              r.AgentID,
		AgentName:            r.AgentName,
		Status:               string(r.Status),
		PhoneProvider:        string(r.PhoneProvider),
		ScheduledTimeUnix:    r.ScheduledTimeUnix,
		CreatedAtUnix:        r.CreatedAtUnix,
		LastUpdatedAtUnix:    r.LastUpdatedAtUnix,
		TotalCallsScheduled:  r.TotalCallsScheduled,
		TotalCallsDispatched: r.TotalCallsDispatched,
		TotalCallsFinished:   r.TotalCallsFinished,
		RetryCount:           r.RetryCount,
	}

	if !r.BranchID.Null {
		bc.BranchID = r.BranchID.Value
	}
	if !r.BranchName.Null {
		bc.BranchName = r.BranchName.Value
	}
	if !r.PhoneNumberID.Null {
		bc.PhoneNumberID = r.PhoneNumberID.Value
	}
	if !r.Environment.Null {
		bc.Environment = r.Environment.Value
	}
	if !r.Timezone.Null {
		bc.Timezone = r.Timezone.Value
	}
	if !r.TargetConcurrencyLimit.Null {
		bc.TargetConcurrencyLimit = &r.TargetConcurrencyLimit.Value
	}

	return bc
}

// batchCallDetailedFromResponse converts detailed API response to BatchCall.
//
//nolint:dupl // Similar to batchCallFromResponse but different input type
func batchCallDetailedFromResponse(r *api.BatchCallDetailedResponse) *BatchCall {
	bc := &BatchCall{
		ID:                   r.ID,
		Name:                 r.Name,
		AgentID:              r.AgentID,
		AgentName:            r.AgentName,
		Status:               string(r.Status),
		PhoneProvider:        string(r.PhoneProvider),
		ScheduledTimeUnix:    r.ScheduledTimeUnix,
		CreatedAtUnix:        r.CreatedAtUnix,
		LastUpdatedAtUnix:    r.LastUpdatedAtUnix,
		TotalCallsScheduled:  r.TotalCallsScheduled,
		TotalCallsDispatched: r.TotalCallsDispatched,
		TotalCallsFinished:   r.TotalCallsFinished,
		RetryCount:           r.RetryCount,
	}

	if !r.BranchID.Null {
		bc.BranchID = r.BranchID.Value
	}
	if !r.BranchName.Null {
		bc.BranchName = r.BranchName.Value
	}
	if !r.PhoneNumberID.Null {
		bc.PhoneNumberID = r.PhoneNumberID.Value
	}
	if !r.Environment.Null {
		bc.Environment = r.Environment.Value
	}
	if !r.Timezone.Null {
		bc.Timezone = r.Timezone.Value
	}
	if !r.TargetConcurrencyLimit.Null {
		bc.TargetConcurrencyLimit = &r.TargetConcurrencyLimit.Value
	}

	return bc
}
