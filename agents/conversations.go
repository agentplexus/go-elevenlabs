package agents

import (
	"context"
	"errors"
	"io"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Conversation represents a conversation with an AI agent.
type Conversation struct {
	ConversationID         string
	AgentID                string
	AgentName              string
	BranchID               string
	Status                 string
	StartTimeUnixSecs      int
	CallDurationSecs       int
	MessageCount           int
	CallSuccessful         string
	CallSummaryTitle       string
	TranscriptSummary      string
	Rating                 *float64
	MainLanguage           string
	TerminationReason      string
	Direction              string
	ConversationInitiation string
	ToolNames              []string
	VersionID              string
}

// ConversationDetail represents detailed conversation information including transcript.
type ConversationDetail struct {
	Conversation
	Transcript []TranscriptMessage
	Metadata   map[string]any
	Analysis   *ConversationAnalysis
	HasAudio   bool
	TagIDs     []string
}

// TranscriptMessage represents a message in the conversation transcript.
type TranscriptMessage struct {
	Role           string
	Message        string
	TimeInCallSecs int
	Interrupted    bool
	ToolCalls      []ToolCall
}

// ToolCall represents a tool invocation in a conversation.
type ToolCall struct {
	ToolName  string
	RequestID string
	Params    map[string]any
}

// ConversationAnalysis represents analysis results for a conversation.
type ConversationAnalysis struct {
	CallSuccessful    string
	TranscriptSummary string
	CallSummaryTitle  string
}

// ListConversationsOptions contains options for listing conversations.
type ListConversationsOptions struct {
	// Cursor for pagination
	Cursor string
	// Filter by agent ID
	AgentID string
	// Filter by call success result: "success", "failure", "unknown"
	CallSuccessful string
	// Filter conversations before this Unix timestamp
	CallStartBeforeUnix *int
	// Filter conversations after this Unix timestamp
	CallStartAfterUnix *int
	// Minimum call duration in seconds
	CallDurationMinSecs *int
	// Maximum call duration in seconds
	CallDurationMaxSecs *int
	// Minimum rating (1-5)
	RatingMin *int
	// Maximum rating (1-5)
	RatingMax *int
	// Filter by presence of feedback comments
	HasFeedbackComment *bool
	// Filter by user ID
	UserID string
	// Filter by branch ID
	BranchID string
	// Filter by topic IDs
	TopicIDs []string
	// Filter by tool names used
	ToolNames []string
	// Filter by main languages
	MainLanguages []string
	// Exclude conversations with these statuses
	ExcludeStatuses []string
	// Maximum conversations to return (max 100, default 30)
	PageSize int
	// Whether to include transcript summaries
	IncludeSummary bool
}

// ListConversationsResponse contains the paginated list of conversations.
type ListConversationsResponse struct {
	Conversations []*Conversation
	HasMore       bool
	NextCursor    string
}

// ListConversations returns a paginated list of conversations.
func (s *Service) ListConversations(ctx context.Context, opts *ListConversationsOptions) (*ListConversationsResponse, error) {
	params := api.GetConversationHistoriesRouteParams{
		PageSize: api.NewOptInt(30),
	}

	if opts != nil {
		if opts.Cursor != "" {
			params.Cursor = api.OptNilString{Value: opts.Cursor, Set: true}
		}
		if opts.AgentID != "" {
			params.AgentID = api.OptNilString{Value: opts.AgentID, Set: true}
		}
		if opts.CallSuccessful != "" {
			params.CallSuccessful = api.NewOptEvaluationSuccessResult(api.EvaluationSuccessResult(opts.CallSuccessful))
		}
		if opts.CallStartBeforeUnix != nil {
			params.CallStartBeforeUnix = api.OptNilInt{Value: *opts.CallStartBeforeUnix, Set: true}
		}
		if opts.CallStartAfterUnix != nil {
			params.CallStartAfterUnix = api.OptNilInt{Value: *opts.CallStartAfterUnix, Set: true}
		}
		if opts.CallDurationMinSecs != nil {
			params.CallDurationMinSecs = api.OptNilInt{Value: *opts.CallDurationMinSecs, Set: true}
		}
		if opts.CallDurationMaxSecs != nil {
			params.CallDurationMaxSecs = api.OptNilInt{Value: *opts.CallDurationMaxSecs, Set: true}
		}
		if opts.RatingMin != nil {
			params.RatingMin = api.OptNilInt{Value: *opts.RatingMin, Set: true}
		}
		if opts.RatingMax != nil {
			params.RatingMax = api.OptNilInt{Value: *opts.RatingMax, Set: true}
		}
		if opts.HasFeedbackComment != nil {
			params.HasFeedbackComment = api.OptNilBool{Value: *opts.HasFeedbackComment, Set: true}
		}
		if opts.UserID != "" {
			params.UserID = api.OptNilString{Value: opts.UserID, Set: true}
		}
		if opts.BranchID != "" {
			params.BranchID = api.OptNilString{Value: opts.BranchID, Set: true}
		}
		if len(opts.TopicIDs) > 0 {
			params.TopicIds = api.OptNilStringArray{Value: opts.TopicIDs, Set: true}
		}
		if len(opts.ToolNames) > 0 {
			params.ToolNames = api.OptNilStringArray{Value: opts.ToolNames, Set: true}
		}
		if len(opts.MainLanguages) > 0 {
			params.MainLanguages = api.OptNilStringArray{Value: opts.MainLanguages, Set: true}
		}
		if len(opts.ExcludeStatuses) > 0 {
			statuses := make([]api.GetConversationHistoriesRouteExcludeStatusesItem, len(opts.ExcludeStatuses))
			for i, st := range opts.ExcludeStatuses {
				statuses[i] = api.GetConversationHistoriesRouteExcludeStatusesItem(st)
			}
			params.ExcludeStatuses = api.OptNilGetConversationHistoriesRouteExcludeStatusesItemArray{Value: statuses, Set: true}
		}
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.IncludeSummary {
			params.SummaryMode = api.NewOptGetConversationHistoriesRouteSummaryMode(api.GetConversationHistoriesRouteSummaryModeInclude)
		}
	}

	resp, err := s.client.API().GetConversationHistoriesRoute(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetConversationsPageResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	conversations := make([]*Conversation, len(result.Conversations))
	for i, c := range result.Conversations {
		conversations[i] = conversationFromSummary(&c)
	}

	return &ListConversationsResponse{
		Conversations: conversations,
		HasMore:       result.HasMore,
		NextCursor:    result.NextCursor.Value,
	}, nil
}

// GetConversation retrieves a specific conversation by ID.
func (s *Service) GetConversation(ctx context.Context, conversationID string) (*ConversationDetail, error) {
	if conversationID == "" {
		return nil, errors.New("conversation ID is required")
	}

	resp, err := s.client.API().GetConversationHistoryRoute(ctx, api.GetConversationHistoryRouteParams{
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetConversationResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return conversationDetailFromResponse(result), nil
}

// DeleteConversation removes a conversation.
func (s *Service) DeleteConversation(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}

	_, err := s.client.API().DeleteConversationRoute(ctx, api.DeleteConversationRouteParams{
		ConversationID: conversationID,
	})
	return err
}

// GetConversationAudio retrieves the audio recording for a conversation.
func (s *Service) GetConversationAudio(ctx context.Context, conversationID string) (io.ReadCloser, error) {
	if conversationID == "" {
		return nil, errors.New("conversation ID is required")
	}

	resp, err := s.client.API().GetConversationAudioRoute(ctx, api.GetConversationAudioRouteParams{
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, err
	}

	stream, ok := resp.(*api.GetConversationAudioRouteOK)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return io.NopCloser(stream.Data), nil
}

// SubmitConversationFeedback submits user feedback for a conversation.
// Feedback must be "like" or "dislike".
func (s *Service) SubmitConversationFeedback(ctx context.Context, conversationID string, feedback string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}
	if feedback != "like" && feedback != "dislike" {
		return errors.New("feedback must be 'like' or 'dislike'")
	}

	req := &api.ConversationFeedbackRequestModel{
		Feedback: api.NewOptUserFeedbackScore(api.UserFeedbackScore(feedback)),
	}

	_, err := s.client.API().PostConversationFeedbackRoute(ctx, req, api.PostConversationFeedbackRouteParams{
		ConversationID: conversationID,
	})
	return err
}

// RunConversationAnalysis triggers analysis on a conversation.
func (s *Service) RunConversationAnalysis(ctx context.Context, conversationID string) error {
	if conversationID == "" {
		return errors.New("conversation ID is required")
	}

	_, err := s.client.API().RunConversationAnalysis(ctx, api.RunConversationAnalysisParams{
		ConversationID: conversationID,
	})
	return err
}

// SearchConversations performs a text search across conversation transcripts.
func (s *Service) SearchConversations(ctx context.Context, query string, opts *ListConversationsOptions) (*ListConversationsResponse, error) {
	if query == "" {
		return nil, errors.New("query is required")
	}

	path := "/v1/convai/conversations/messages/smart-search"

	body := map[string]any{
		"query": query,
	}
	if opts != nil {
		if opts.AgentID != "" {
			body["agent_id"] = opts.AgentID
		}
		if opts.PageSize > 0 {
			body["page_size"] = opts.PageSize
		}
	}

	var result struct {
		Conversations []api.ConversationSummaryResponseModel `json:"conversations"`
		HasMore       bool                                   `json:"has_more"`
		NextCursor    string                                 `json:"next_cursor"`
	}

	if err := s.doJSON(ctx, "POST", path, body, &result); err != nil {
		return nil, err
	}

	conversations := make([]*Conversation, len(result.Conversations))
	for i, c := range result.Conversations {
		conversations[i] = conversationFromSummary(&c)
	}

	return &ListConversationsResponse{
		Conversations: conversations,
		HasMore:       result.HasMore,
		NextCursor:    result.NextCursor,
	}, nil
}

// conversationFromSummary converts API summary model to Conversation.
func conversationFromSummary(c *api.ConversationSummaryResponseModel) *Conversation {
	conv := &Conversation{
		ConversationID:    c.ConversationID,
		AgentID:           c.AgentID,
		Status:            string(c.Status),
		StartTimeUnixSecs: c.StartTimeUnixSecs,
		CallDurationSecs:  c.CallDurationSecs,
		MessageCount:      c.MessageCount,
		CallSuccessful:    string(c.CallSuccessful),
	}

	if c.AgentName.Set {
		conv.AgentName = c.AgentName.Value
	}
	if c.BranchID.Set {
		conv.BranchID = c.BranchID.Value
	}
	if c.CallSummaryTitle.Set {
		conv.CallSummaryTitle = c.CallSummaryTitle.Value
	}
	if c.TranscriptSummary.Set {
		conv.TranscriptSummary = c.TranscriptSummary.Value
	}
	if c.Rating.Set {
		conv.Rating = &c.Rating.Value
	}
	if c.MainLanguage.Set {
		conv.MainLanguage = c.MainLanguage.Value
	}
	if c.TerminationReason.Set {
		conv.TerminationReason = c.TerminationReason.Value
	}
	if c.Direction.Set {
		conv.Direction = string(c.Direction.Value)
	}
	if c.ConversationInitiationSource.Set {
		conv.ConversationInitiation = string(c.ConversationInitiationSource.Value)
	}
	if c.ToolNames.Set {
		conv.ToolNames = c.ToolNames.Value
	}
	if c.VersionID.Set {
		conv.VersionID = c.VersionID.Value
	}

	return conv
}

// conversationDetailFromResponse converts API response to ConversationDetail.
func conversationDetailFromResponse(r *api.GetConversationResponseModel) *ConversationDetail {
	detail := &ConversationDetail{
		Conversation: Conversation{
			ConversationID: r.ConversationID,
			AgentID:        r.AgentID,
			Status:         string(r.Status),
		},
		HasAudio: r.HasAudio,
		TagIDs:   r.TagIds,
	}

	if r.AgentName.Set {
		detail.AgentName = r.AgentName.Value
	}
	if r.BranchID.Set {
		detail.BranchID = r.BranchID.Value
	}
	if r.VersionID.Set {
		detail.VersionID = r.VersionID.Value
	}

	// Convert transcript
	if len(r.Transcript) > 0 {
		detail.Transcript = make([]TranscriptMessage, len(r.Transcript))
		for i, m := range r.Transcript {
			msg := TranscriptMessage{
				Role:           string(m.Role),
				TimeInCallSecs: m.TimeInCallSecs,
			}
			if m.Message.Set {
				msg.Message = m.Message.Value
			}
			if m.Interrupted.Set {
				msg.Interrupted = m.Interrupted.Value
			}
			if len(m.ToolCalls) > 0 {
				msg.ToolCalls = make([]ToolCall, len(m.ToolCalls))
				for j, tc := range m.ToolCalls {
					msg.ToolCalls[j] = ToolCall{
						ToolName:  tc.ToolName,
						RequestID: tc.RequestID,
					}
				}
			}
			detail.Transcript[i] = msg
		}
	}

	// Convert analysis if present
	if r.Analysis.Set {
		detail.Analysis = &ConversationAnalysis{
			CallSuccessful:    string(r.Analysis.Value.CallSuccessful),
			TranscriptSummary: r.Analysis.Value.TranscriptSummary,
		}
		if r.Analysis.Value.CallSummaryTitle.Set {
			detail.Analysis.CallSummaryTitle = r.Analysis.Value.CallSummaryTitle.Value
		}
	}

	return detail
}
