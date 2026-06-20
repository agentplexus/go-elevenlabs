package elevenlabs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ht "github.com/ogen-go/ogen/http"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// AgentsService handles Conversational AI agent operations.
type AgentsService struct {
	client *Client
}

// Agent represents a conversational AI agent.
type Agent struct {
	// AgentID is the unique identifier for the agent.
	AgentID string

	// Name is the display name of the agent.
	Name string

	// Tags are labels used to categorize the agent.
	Tags []string

	// CreatedAtUnixSecs is the creation time in unix seconds.
	CreatedAtUnixSecs int64

	// Archived indicates whether the agent is archived.
	Archived bool
}

// AgentSummary is a lightweight agent listing.
type AgentSummary struct {
	// AgentID is the unique identifier for the agent.
	AgentID string

	// Name is the display name of the agent.
	Name string

	// Tags are labels used to categorize the agent.
	Tags []string

	// Archived indicates whether the agent is archived.
	Archived bool

	// CreatedAtUnixSecs is the creation time in unix seconds.
	CreatedAtUnixSecs int64

	// LastCallTimeUnixSecs is the time of the most recent call, nil if no calls.
	LastCallTimeUnixSecs *int64
}

// AgentBranch represents a version branch for an agent.
type AgentBranch struct {
	// ID is the unique identifier for the branch.
	ID string

	// Name is the display name of the branch.
	Name string

	// Description is the description of the branch.
	Description string

	// AgentID is the ID of the agent this branch belongs to.
	AgentID string

	// CreatedAt is the creation time in unix seconds.
	CreatedAt int64

	// LastCommittedAt is the time of the last commit in unix seconds.
	LastCommittedAt int64

	// IsArchived indicates whether the branch is archived.
	IsArchived bool

	// CurrentLivePercentage is the percentage of traffic this branch receives.
	CurrentLivePercentage float64

	// ParentBranchID is the ID of the parent branch, nil if none.
	ParentBranchID *string

	// ParentBranchName is the name of the parent branch, nil if none.
	ParentBranchName *string
}

// AgentTopic represents conversation topic analytics.
type AgentTopic struct {
	// TopicID is the unique identifier for the topic.
	TopicID string

	// Label is the display label for the topic.
	Label string

	// Description is the description of the topic.
	Description string

	// ConversationCount is the number of conversations with this topic.
	ConversationCount int

	// ParentTopicID is the ID of the parent topic, nil if none.
	ParentTopicID *string
}

// AgentLink contains information for sharing an agent.
type AgentLink struct {
	// AgentID is the ID of the agent.
	AgentID string

	// SignedURL is the signed URL for accessing the agent.
	SignedURL string

	// ExpiresAtUnixMillis is when the link expires.
	ExpiresAtUnixMillis int64
}

// ListAgentsOptions configures the List operation.
type ListAgentsOptions struct {
	// PageSize is the maximum number of agents to return (max 100, default 30).
	PageSize int

	// Search filters by agent name.
	Search string

	// Archived filters by archived status.
	Archived *bool

	// SortBy specifies the field to sort by ("name", "created_at", etc.).
	SortBy string

	// SortDirection specifies the sort direction ("asc" or "desc").
	SortDirection string

	// Cursor is the pagination cursor from a previous response.
	Cursor string

	// CreatedByUserID filters to agents created by a specific user.
	// Use "@me" to filter to the authenticated user.
	CreatedByUserID string
}

// ListAgentsResponse contains paginated agent results.
type ListAgentsResponse struct {
	// Agents is the list of agent summaries.
	Agents []*AgentSummary

	// HasMore indicates if there are more results.
	HasMore bool

	// NextCursor is the cursor for the next page.
	NextCursor string
}

// CreateAgentRequest contains options for creating a new agent.
type CreateAgentRequest struct {
	// Name is the display name for the agent.
	Name string `json:"name,omitempty"`

	// Tags are labels for categorizing the agent.
	Tags []string `json:"tags,omitempty"`

	// ConversationConfig is the agent's conversation configuration.
	// This is a flexible map to accommodate the complex nested schema.
	ConversationConfig map[string]any `json:"conversation_config,omitempty"`

	// PlatformSettings are optional platform-specific settings.
	PlatformSettings map[string]any `json:"platform_settings,omitempty"`
}

// UpdateAgentRequest contains options for updating an agent.
type UpdateAgentRequest struct {
	// Name is the new display name for the agent.
	Name string `json:"name,omitempty"`

	// Tags are the new labels for the agent.
	Tags []string `json:"tags,omitempty"`

	// ConversationConfig is the agent's conversation configuration.
	ConversationConfig map[string]any `json:"conversation_config,omitempty"`

	// PlatformSettings are optional platform-specific settings.
	PlatformSettings map[string]any `json:"platform_settings,omitempty"`
}

// CreateBranchRequest contains options for creating a branch.
type CreateBranchRequest struct {
	// Name is the display name for the branch.
	Name string `json:"name"`

	// Description is an optional description.
	Description string `json:"description,omitempty"`

	// BaseBranchID is the ID of the branch to fork from.
	BaseBranchID string `json:"base_branch_id,omitempty"`
}

// UpdateBranchRequest contains options for updating a branch.
type UpdateBranchRequest struct {
	// Name is the new display name for the branch.
	Name string `json:"name,omitempty"`

	// Description is the new description.
	Description string `json:"description,omitempty"`

	// IsArchived sets the archived status.
	IsArchived *bool `json:"is_archived,omitempty"`
}

// DeploymentRequest specifies branch deployment configuration.
type DeploymentRequest struct {
	// BranchID is the ID of the branch to deploy.
	BranchID string

	// Percentage is the traffic percentage for this branch (0-100).
	Percentage float64
}

// WidgetConfig contains agent widget configuration.
type WidgetConfig struct {
	// AgentID is the ID of the agent.
	AgentID string `json:"agent_id"`

	// WidgetConfig contains widget settings.
	WidgetConfig map[string]any `json:"widget_config,omitempty"`
}

// List returns a paginated list of agents.
func (s *AgentsService) List(ctx context.Context, opts *ListAgentsOptions) (*ListAgentsResponse, error) {
	params := api.GetAgentsRouteParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.Search != "" {
			params.Search = api.OptNilString{Value: opts.Search, Set: true}
		}
		if opts.Archived != nil {
			params.Archived = api.OptNilBool{Value: *opts.Archived, Set: true}
		}
		if opts.SortBy != "" {
			params.SortBy = api.NewOptAgentSortBy(api.AgentSortBy(opts.SortBy))
		}
		if opts.SortDirection != "" {
			params.SortDirection = api.NewOptSortDirection(api.SortDirection(opts.SortDirection))
		}
		if opts.Cursor != "" {
			params.Cursor = api.OptNilString{Value: opts.Cursor, Set: true}
		}
		if opts.CreatedByUserID != "" {
			params.CreatedByUserID = api.OptNilString{Value: opts.CreatedByUserID, Set: true}
		}
	}

	resp, err := s.client.apiClient.GetAgentsRoute(ctx, params)
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetAgentsPageResponseModel:
		result := &ListAgentsResponse{
			Agents:  make([]*AgentSummary, 0, len(r.Agents)),
			HasMore: r.HasMore,
		}
		if r.NextCursor.Set && !r.NextCursor.Null {
			result.NextCursor = r.NextCursor.Value
		}
		for _, a := range r.Agents {
			summary := &AgentSummary{
				AgentID:           a.AgentID,
				Name:              a.Name,
				Tags:              a.Tags,
				CreatedAtUnixSecs: int64(a.CreatedAtUnixSecs),
			}
			if a.Archived.Set {
				summary.Archived = a.Archived.Value
			}
			if a.LastCallTimeUnixSecs.Set && !a.LastCallTimeUnixSecs.Null {
				val := int64(a.LastCallTimeUnixSecs.Value)
				summary.LastCallTimeUnixSecs = &val
			}
			result.Agents = append(result.Agents, summary)
		}
		return result, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// Get retrieves an agent by ID.
// Note: Uses raw HTTP as ogen didn't generate this operation.
func (s *AgentsService) Get(ctx context.Context, agentID string) (*Agent, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	var result struct {
		AgentID           string   `json:"agent_id"`
		Name              string   `json:"name"`
		Tags              []string `json:"tags"`
		CreatedAtUnixSecs int64    `json:"created_at_unix_secs"`
	}

	if err := s.doJSON(ctx, "GET", "/v1/convai/agents/"+agentID, nil, &result); err != nil {
		return nil, err
	}

	return &Agent{
		AgentID:           result.AgentID,
		Name:              result.Name,
		Tags:              result.Tags,
		CreatedAtUnixSecs: result.CreatedAtUnixSecs,
	}, nil
}

// Create creates a new agent.
// Note: Uses raw HTTP as ogen didn't generate this operation.
func (s *AgentsService) Create(ctx context.Context, req *CreateAgentRequest) (*Agent, error) {
	if req == nil {
		req = &CreateAgentRequest{}
	}

	var result struct {
		AgentID string `json:"agent_id"`
	}

	if err := s.doJSON(ctx, "POST", "/v1/convai/agents/create", req, &result); err != nil {
		return nil, err
	}

	return &Agent{
		AgentID: result.AgentID,
		Name:    req.Name,
		Tags:    req.Tags,
	}, nil
}

// Update updates an existing agent.
// Note: Uses raw HTTP as ogen didn't generate this operation.
func (s *AgentsService) Update(ctx context.Context, agentID string, req *UpdateAgentRequest) (*Agent, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}
	if req == nil {
		req = &UpdateAgentRequest{}
	}

	var result struct {
		AgentID string `json:"agent_id"`
	}

	if err := s.doJSON(ctx, "PATCH", "/v1/convai/agents/"+agentID, req, &result); err != nil {
		return nil, err
	}

	return &Agent{
		AgentID: result.AgentID,
		Name:    req.Name,
		Tags:    req.Tags,
	}, nil
}

// Delete deletes an agent.
func (s *AgentsService) Delete(ctx context.Context, agentID string) error {
	if agentID == "" {
		return &APIError{Message: "agent_id is required"}
	}

	_, err := s.client.apiClient.DeleteAgentRoute(ctx, api.DeleteAgentRouteParams{
		AgentID: agentID,
	})
	return err
}

// Duplicate creates a copy of an agent.
func (s *AgentsService) Duplicate(ctx context.Context, agentID string, name string) (*Agent, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	var req api.OptBodyDuplicateAgentV1ConvaiAgentsAgentIDDuplicatePost
	if name != "" {
		req = api.NewOptBodyDuplicateAgentV1ConvaiAgentsAgentIDDuplicatePost(
			api.BodyDuplicateAgentV1ConvaiAgentsAgentIDDuplicatePost{
				Name: api.OptNilString{Value: name, Set: true},
			},
		)
	}

	resp, err := s.client.apiClient.DuplicateAgentRoute(ctx, req, api.DuplicateAgentRouteParams{
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.CreateAgentResponseModel:
		return &Agent{
			AgentID: r.AgentID,
			Name:    name,
		}, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// GetLink returns the shareable link for an agent.
func (s *AgentsService) GetLink(ctx context.Context, agentID string) (*AgentLink, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	resp, err := s.client.apiClient.GetAgentLinkRoute(ctx, api.GetAgentLinkRouteParams{
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetAgentLinkResponseModel:
		link := &AgentLink{
			AgentID: r.AgentID,
		}
		if r.Token.Set {
			token := r.Token.Value
			// Use conversation token as the signed URL
			link.SignedURL = token.ConversationToken
			if token.ExpirationTimeUnixSecs.Set && !token.ExpirationTimeUnixSecs.Null {
				// Convert to milliseconds
				link.ExpiresAtUnixMillis = int64(token.ExpirationTimeUnixSecs.Value) * 1000
			}
		}
		return link, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// UploadAvatar uploads an avatar image for an agent.
func (s *AgentsService) UploadAvatar(ctx context.Context, agentID string, avatarData io.Reader) error {
	if agentID == "" {
		return &APIError{Message: "agent_id is required"}
	}

	// Read all data for the multipart file
	data, err := io.ReadAll(avatarData)
	if err != nil {
		return fmt.Errorf("failed to read avatar data: %w", err)
	}

	// Create multipart file
	file := api.BodyPostAgentAvatarV1ConvaiAgentsAgentIDAvatarPostMultipart{
		AvatarFile: ht.MultipartFile{
			Name: "avatar",
			File: bytes.NewReader(data),
		},
	}

	_, err = s.client.apiClient.PostAgentAvatarRoute(ctx, &file, api.PostAgentAvatarRouteParams{
		AgentID: agentID,
	})
	return err
}

// GetTopics returns conversation topic analytics for an agent.
func (s *AgentsService) GetTopics(ctx context.Context, agentID string) ([]*AgentTopic, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	resp, err := s.client.apiClient.GetAgentTopicsRoute(ctx, api.GetAgentTopicsRouteParams{
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetAgentTopicsResponseModel:
		topics := make([]*AgentTopic, 0, len(r.Topics))
		for _, t := range r.Topics {
			topic := &AgentTopic{
				TopicID:           t.TopicID,
				Label:             t.Label,
				Description:       t.Description,
				ConversationCount: t.ConversationCount,
			}
			if t.ParentTopicID.Set && !t.ParentTopicID.Null {
				topic.ParentTopicID = &t.ParentTopicID.Value
			}
			topics = append(topics, topic)
		}
		return topics, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// ListBranches returns all branches for an agent.
func (s *AgentsService) ListBranches(ctx context.Context, agentID string) ([]*AgentBranch, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	resp, err := s.client.apiClient.GetBranchesRoute(ctx, api.GetBranchesRouteParams{
		AgentID: agentID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.ListResponseAgentBranchSummary:
		branches := make([]*AgentBranch, 0, len(r.Results))
		for i := range r.Results {
			branch := branchSummaryToAgentBranch(&r.Results[i])
			branches = append(branches, branch)
		}
		return branches, nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// GetBranch retrieves a specific branch.
func (s *AgentsService) GetBranch(ctx context.Context, agentID, branchID string) (*AgentBranch, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}
	if branchID == "" {
		return nil, &APIError{Message: "branch_id is required"}
	}

	resp, err := s.client.apiClient.GetBranchRoute(ctx, api.GetBranchRouteParams{
		AgentID:  agentID,
		BranchID: branchID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AgentBranchResponse:
		return branchResponseToAgentBranch(r), nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// CreateBranch creates a new branch for an agent.
// Note: Uses raw HTTP as ogen didn't generate this operation.
func (s *AgentsService) CreateBranch(ctx context.Context, agentID string, req *CreateBranchRequest) (*AgentBranch, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}
	if req == nil {
		return nil, &APIError{Message: "request is required"}
	}
	if req.Name == "" {
		return nil, &APIError{Message: "name is required"}
	}

	var result struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		AgentID     string `json:"agent_id"`
		CreatedAt   int64  `json:"created_at"`
	}

	if err := s.doJSON(ctx, "POST", "/v1/convai/agents/"+agentID+"/branches", req, &result); err != nil {
		return nil, err
	}

	return &AgentBranch{
		ID:          result.ID,
		Name:        result.Name,
		Description: result.Description,
		AgentID:     result.AgentID,
		CreatedAt:   result.CreatedAt,
	}, nil
}

// UpdateBranch updates a branch.
func (s *AgentsService) UpdateBranch(ctx context.Context, agentID, branchID string, req *UpdateBranchRequest) (*AgentBranch, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}
	if branchID == "" {
		return nil, &APIError{Message: "branch_id is required"}
	}

	var body api.BodyUpdateAgentBranchV1ConvaiAgentsAgentIDBranchesBranchIDPatch
	if req != nil {
		if req.Name != "" {
			body.Name = api.OptNilString{Value: req.Name, Set: true}
		}
		// Note: Description field is not supported in the API
		if req.IsArchived != nil {
			body.IsArchived = api.OptNilBool{Value: *req.IsArchived, Set: true}
		}
	}

	resp, err := s.client.apiClient.UpdateBranchRoute(ctx,
		api.NewOptBodyUpdateAgentBranchV1ConvaiAgentsAgentIDBranchesBranchIDPatch(body),
		api.UpdateBranchRouteParams{
			AgentID:  agentID,
			BranchID: branchID,
		},
	)
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AgentBranchResponse:
		return branchResponseToAgentBranch(r), nil
	default:
		return nil, &APIError{Message: "unexpected response type"}
	}
}

// MergeBranch merges a source branch into a target branch.
func (s *AgentsService) MergeBranch(ctx context.Context, agentID, sourceBranchID, targetBranchID string) error {
	if agentID == "" {
		return &APIError{Message: "agent_id is required"}
	}
	if sourceBranchID == "" {
		return &APIError{Message: "source_branch_id is required"}
	}
	if targetBranchID == "" {
		return &APIError{Message: "target_branch_id is required"}
	}

	// Body is optional, contains force and archive_source_branch flags
	var body api.OptBodyMergeABranchIntoATargetBranchV1ConvaiAgentsAgentIDBranchesSourceBranchIDMergePost

	_, err := s.client.apiClient.MergeBranchIntoTarget(ctx,
		body,
		api.MergeBranchIntoTargetParams{
			AgentID:        agentID,
			SourceBranchID: sourceBranchID,
			TargetBranchID: targetBranchID,
		},
	)
	return err
}

// Deploy deploys branches with specified traffic percentages.
func (s *AgentsService) Deploy(ctx context.Context, agentID string, deployments []DeploymentRequest) error {
	if agentID == "" {
		return &APIError{Message: "agent_id is required"}
	}
	if len(deployments) == 0 {
		return &APIError{Message: "at least one deployment is required"}
	}

	items := make([]api.AgentDeploymentRequestItem, 0, len(deployments))
	for _, d := range deployments {
		items = append(items, api.AgentDeploymentRequestItem{
			BranchID: d.BranchID,
			DeploymentStrategy: api.AgentDeploymentPercentageStrategy{
				TrafficPercentage: d.Percentage,
			},
		})
	}

	_, err := s.client.apiClient.CreateAgentDeploymentRoute(ctx,
		&api.BodyCreateOrUpdateDeploymentsV1ConvaiAgentsAgentIDDeploymentsPost{
			DeploymentRequest: api.AgentDeploymentRequest{
				Requests: items,
			},
		},
		api.CreateAgentDeploymentRouteParams{
			AgentID: agentID,
		},
	)
	return err
}

// GetKnowledgeBaseSize retrieves the knowledge base size for an agent.
func (s *AgentsService) GetKnowledgeBaseSize(ctx context.Context, agentID string) (float64, error) {
	if agentID == "" {
		return 0, &APIError{Message: "agent_id is required"}
	}

	resp, err := s.client.apiClient.GetAgentKnowledgeBaseSize(ctx, api.GetAgentKnowledgeBaseSizeParams{
		AgentID: agentID,
	})
	if err != nil {
		return 0, err
	}

	result, ok := resp.(*api.GetAgentKnowledgebaseSizeResponseModel)
	if !ok {
		return 0, &APIError{Message: "unexpected response type"}
	}

	return result.NumberOfPages, nil
}

// GetWidget retrieves the widget configuration for an agent.
// Note: Uses raw HTTP as ogen didn't generate this operation.
func (s *AgentsService) GetWidget(ctx context.Context, agentID string) (*WidgetConfig, error) {
	if agentID == "" {
		return nil, &APIError{Message: "agent_id is required"}
	}

	var result WidgetConfig
	if err := s.doJSON(ctx, "GET", "/v1/convai/agents/"+agentID+"/widget", nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// doJSON is a helper for making JSON HTTP requests.
//
//nolint:dupl // Intentional duplicate - service-specific HTTP helper
func (s *AgentsService) doJSON(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, s.client.baseURL+path, bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("xi-api-key", s.client.apiKey)

	resp, err := http.DefaultClient.Do(httpReq) //nolint:gosec // G704: API client, URL is fixed ElevenLabs endpoint
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(respBody),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// branchSummaryToAgentBranch converts an API branch summary to AgentBranch.
func branchSummaryToAgentBranch(b *api.AgentBranchSummary) *AgentBranch {
	branch := &AgentBranch{
		ID:              b.ID,
		Name:            b.Name,
		Description:     b.Description,
		AgentID:         b.AgentID,
		CreatedAt:       int64(b.CreatedAt),
		LastCommittedAt: int64(b.LastCommittedAt),
		IsArchived:      b.IsArchived,
	}
	if b.CurrentLivePercentage.Set {
		branch.CurrentLivePercentage = b.CurrentLivePercentage.Value
	}
	if b.ParentBranchID.Set && !b.ParentBranchID.Null {
		branch.ParentBranchID = &b.ParentBranchID.Value
	}
	return branch
}

// branchResponseToAgentBranch converts an API branch response to AgentBranch.
func branchResponseToAgentBranch(b *api.AgentBranchResponse) *AgentBranch {
	branch := &AgentBranch{
		ID:              b.ID,
		Name:            b.Name,
		Description:     b.Description,
		AgentID:         b.AgentID,
		CreatedAt:       int64(b.CreatedAt),
		LastCommittedAt: int64(b.LastCommittedAt),
		IsArchived:      b.IsArchived,
	}
	if b.CurrentLivePercentage.Set {
		branch.CurrentLivePercentage = b.CurrentLivePercentage.Value
	}
	if b.ParentBranch.Set {
		branch.ParentBranchID = &b.ParentBranch.Value.ID
		branch.ParentBranchName = &b.ParentBranch.Value.Name
	}
	return branch
}
