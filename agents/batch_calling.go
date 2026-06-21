package agents

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// BatchCall represents a batch calling job.
type BatchCall struct {
	// ID is the unique identifier for the batch call.
	ID string

	// Name is the display name of the batch call.
	Name string

	// AgentID is the ID of the agent used for the calls.
	AgentID string

	// AgentName is the name of the agent.
	AgentName string

	// BranchID is the branch ID if applicable.
	BranchID string

	// BranchName is the branch name if applicable.
	BranchName string

	// Status is the batch call status: "pending", "running", "completed", "failed", "cancelled".
	Status string

	// TotalCallsScheduled is the total number of calls scheduled.
	TotalCallsScheduled int

	// TotalCallsDispatched is the number of calls dispatched.
	TotalCallsDispatched int

	// TotalCallsFinished is the number of finished calls.
	TotalCallsFinished int

	// RetryCount is the number of retries performed.
	RetryCount int

	// CreatedAtUnix is the creation time in unix seconds.
	CreatedAtUnix int

	// LastUpdatedAtUnix is the last update time in unix seconds.
	LastUpdatedAtUnix int

	// ScheduledTimeUnix is the scheduled start time in unix seconds.
	ScheduledTimeUnix int
}

// BatchCallDetail contains detailed batch call information including recipients.
type BatchCallDetail struct {
	BatchCall

	// Recipients contains information about each call recipient.
	Recipients []*BatchCallRecipientStatus
}

// BatchCallRecipientStatus represents the status of a single recipient in a batch call.
type BatchCallRecipientStatus struct {
	// ID is the unique identifier for this recipient.
	ID string

	// PhoneNumber is the recipient's phone number.
	PhoneNumber string

	// ConversationID is the ID of the conversation created for this call.
	ConversationID string

	// Status is the recipient call status.
	Status string

	// CreatedAtUnix is the creation time in unix seconds.
	CreatedAtUnix int

	// UpdatedAtUnix is the last update time in unix seconds.
	UpdatedAtUnix int
}

// BatchCallRecipient represents a single call recipient for creating batch calls.
type BatchCallRecipient struct {
	// PhoneNumber is the recipient's phone number in E.164 format.
	PhoneNumber string

	// ID is an optional custom identifier for this recipient.
	ID string

	// DynamicVariables are optional context variables for this recipient.
	DynamicVariables map[string]any
}

// CreateBatchCallRequest contains options for creating a batch call.
type CreateBatchCallRequest struct {
	// AgentID is the ID of the agent to use for calls.
	AgentID string

	// Name is the display name for the batch call.
	Name string

	// Recipients is the list of call recipients.
	Recipients []BatchCallRecipient

	// AgentPhoneNumberID is the phone number ID to use for outbound calls.
	AgentPhoneNumberID string

	// BranchID specifies a specific branch to use.
	BranchID string

	// Environment specifies the environment for variable resolution.
	Environment string

	// ScheduledTimeUnix is an optional scheduled start time in unix seconds.
	ScheduledTimeUnix *int

	// TargetConcurrencyLimit limits the number of parallel calls.
	TargetConcurrencyLimit *int

	// Timezone for scheduling.
	Timezone string
}

// ListBatchCallsOptions configures the ListBatchCalls operation.
type ListBatchCallsOptions struct {
	// PageSize is the maximum number of batch calls to return (max 100).
	PageSize int

	// Cursor is the pagination cursor from a previous response.
	Cursor string

	// AgentID filters to a specific agent.
	AgentID string

	// Status filters by batch status.
	Status string
}

// ListBatchCallsResponse contains paginated batch call results.
type ListBatchCallsResponse struct {
	// BatchCalls is the list of batch calls.
	BatchCalls []*BatchCall

	// HasMore indicates if there are more results.
	HasMore bool

	// NextCursor is the cursor for the next page.
	NextCursor string
}

// CreateBatchCall creates a new batch calling job.
func (s *Service) CreateBatchCall(ctx context.Context, req *CreateBatchCallRequest) (*BatchCall, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.AgentID == "" {
		return nil, errors.New("agent ID is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if len(req.Recipients) == 0 {
		return nil, errors.New("at least one recipient is required")
	}

	recipients := make([]api.OutboundCallRecipient, 0, len(req.Recipients))
	for _, r := range req.Recipients {
		recipient := api.OutboundCallRecipient{}
		if r.PhoneNumber != "" {
			recipient.PhoneNumber = api.OptNilString{Value: r.PhoneNumber, Set: true}
		}
		if r.ID != "" {
			recipient.ID = api.OptNilString{Value: r.ID, Set: true}
		}
		// Note: DynamicVariables would need to be converted to api.NilDynamicVariableValueTypeInput
		// format. For simplicity, we support basic string variables via ConversationInitiationClientData.
		recipients = append(recipients, recipient)
	}

	body := &api.BodySubmitABatchCallRequestV1ConvaiBatchCallingSubmitPost{
		AgentID:    req.AgentID,
		CallName:   req.Name,
		Recipients: recipients,
	}

	if req.AgentPhoneNumberID != "" {
		body.AgentPhoneNumberID = api.OptNilString{Value: req.AgentPhoneNumberID, Set: true}
	}
	if req.BranchID != "" {
		body.BranchID = api.OptNilString{Value: req.BranchID, Set: true}
	}
	if req.Environment != "" {
		body.Environment = api.OptNilString{Value: req.Environment, Set: true}
	}
	if req.ScheduledTimeUnix != nil {
		body.ScheduledTimeUnix = api.OptNilInt{Value: *req.ScheduledTimeUnix, Set: true}
	}
	if req.TargetConcurrencyLimit != nil {
		body.TargetConcurrencyLimit = api.OptNilInt{Value: *req.TargetConcurrencyLimit, Set: true}
	}
	if req.Timezone != "" {
		body.Timezone = api.OptNilString{Value: req.Timezone, Set: true}
	}

	resp, err := s.client.API().CreateBatchCall(ctx, body, api.CreateBatchCallParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.BatchCallResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return batchCallFromResponse(result), nil
}

// GetBatchCall retrieves a batch call by ID.
func (s *Service) GetBatchCall(ctx context.Context, batchCallID string) (*BatchCallDetail, error) {
	if batchCallID == "" {
		return nil, errors.New("batch call ID is required")
	}

	resp, err := s.client.API().GetBatchCall(ctx, api.GetBatchCallParams{
		BatchID: batchCallID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.BatchCallDetailedResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return batchCallDetailFromResponse(result), nil
}

// CancelBatchCall cancels a running batch call.
func (s *Service) CancelBatchCall(ctx context.Context, batchCallID string) error {
	if batchCallID == "" {
		return errors.New("batch call ID is required")
	}

	_, err := s.client.API().CancelBatchCall(ctx, api.CancelBatchCallParams{
		BatchID: batchCallID,
	})
	return err
}

// RetryBatchCall retries failed calls in a batch.
func (s *Service) RetryBatchCall(ctx context.Context, batchCallID string) (*BatchCall, error) {
	if batchCallID == "" {
		return nil, errors.New("batch call ID is required")
	}

	resp, err := s.client.API().RetryBatchCall(ctx, api.RetryBatchCallParams{
		BatchID: batchCallID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.BatchCallResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return batchCallFromResponse(result), nil
}

// DeleteBatchCall deletes a batch call.
func (s *Service) DeleteBatchCall(ctx context.Context, batchCallID string) error {
	if batchCallID == "" {
		return errors.New("batch call ID is required")
	}

	_, err := s.client.API().DeleteBatchCall(ctx, api.DeleteBatchCallParams{
		BatchID: batchCallID,
	})
	return err
}

// ListBatchCalls returns a paginated list of batch calls.
func (s *Service) ListBatchCalls(ctx context.Context, opts *ListBatchCallsOptions) (*ListBatchCallsResponse, error) {
	params := api.GetWorkspaceBatchCallsParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.Limit = api.NewOptInt(opts.PageSize)
		}
		if opts.Cursor != "" {
			params.LastDoc = api.OptNilString{Value: opts.Cursor, Set: true}
		}
		if opts.AgentID != "" {
			params.AgentID = api.OptNilString{Value: opts.AgentID, Set: true}
		}
		// Note: Status filtering is not supported by the current API
	}

	resp, err := s.client.API().GetWorkspaceBatchCalls(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.WorkspaceBatchCallsResponse)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	batchCalls := make([]*BatchCall, 0, len(result.BatchCalls))
	for i := range result.BatchCalls {
		batchCalls = append(batchCalls, batchCallFromResponse(&result.BatchCalls[i]))
	}

	var nextCursor string
	if result.NextDoc.Set && !result.NextDoc.Null {
		nextCursor = result.NextDoc.Value
	}

	var hasMore bool
	if result.HasMore.Set {
		hasMore = result.HasMore.Value
	}

	return &ListBatchCallsResponse{
		BatchCalls: batchCalls,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

// batchCallFromResponse converts API response to BatchCall.
func batchCallFromResponse(r *api.BatchCallResponse) *BatchCall {
	bc := &BatchCall{
		ID:                   r.ID,
		Name:                 r.Name,
		AgentID:              r.AgentID,
		AgentName:            r.AgentName,
		Status:               string(r.Status),
		TotalCallsScheduled:  r.TotalCallsScheduled,
		TotalCallsDispatched: r.TotalCallsDispatched,
		TotalCallsFinished:   r.TotalCallsFinished,
		RetryCount:           r.RetryCount,
		CreatedAtUnix:        r.CreatedAtUnix,
		LastUpdatedAtUnix:    r.LastUpdatedAtUnix,
		ScheduledTimeUnix:    r.ScheduledTimeUnix,
	}
	if !r.BranchID.Null {
		bc.BranchID = r.BranchID.Value
	}
	if !r.BranchName.Null {
		bc.BranchName = r.BranchName.Value
	}
	return bc
}

// batchCallDetailFromResponse converts API detailed response to BatchCallDetail.
func batchCallDetailFromResponse(r *api.BatchCallDetailedResponse) *BatchCallDetail {
	bc := &BatchCallDetail{
		BatchCall: BatchCall{
			ID:                   r.ID,
			Name:                 r.Name,
			AgentID:              r.AgentID,
			AgentName:            r.AgentName,
			Status:               string(r.Status),
			TotalCallsScheduled:  r.TotalCallsScheduled,
			TotalCallsDispatched: r.TotalCallsDispatched,
			TotalCallsFinished:   r.TotalCallsFinished,
			RetryCount:           r.RetryCount,
			CreatedAtUnix:        r.CreatedAtUnix,
			LastUpdatedAtUnix:    r.LastUpdatedAtUnix,
			ScheduledTimeUnix:    r.ScheduledTimeUnix,
		},
	}
	if !r.BranchID.Null {
		bc.BranchID = r.BranchID.Value
	}
	if !r.BranchName.Null {
		bc.BranchName = r.BranchName.Value
	}

	// Convert recipients
	if len(r.Recipients) > 0 {
		bc.Recipients = make([]*BatchCallRecipientStatus, 0, len(r.Recipients))
		for _, rec := range r.Recipients {
			recipientStatus := &BatchCallRecipientStatus{
				ID:            rec.ID,
				Status:        string(rec.Status),
				CreatedAtUnix: rec.CreatedAtUnix,
				UpdatedAtUnix: rec.UpdatedAtUnix,
			}
			if rec.PhoneNumber.Set && !rec.PhoneNumber.Null {
				recipientStatus.PhoneNumber = rec.PhoneNumber.Value
			}
			if !rec.ConversationID.Null {
				recipientStatus.ConversationID = rec.ConversationID.Value
			}
			bc.Recipients = append(bc.Recipients, recipientStatus)
		}
	}
	return bc
}
