package elevenlabs

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// AgentTestingService handles agent testing operations.
type AgentTestingService struct {
	client *Client
}

// AgentTestFolder represents a folder for organizing tests.
type AgentTestFolder struct {
	ID            string
	Name          string
	ChildrenCount int
	Path          []AgentTestFolderSegment
}

// AgentTestFolderSegment represents a segment in the folder path.
type AgentTestFolderSegment struct {
	ID   string
	Name string
}

// AgentResponseTest represents a response test for an agent.
type AgentResponseTest struct {
	ID       string
	Name     string
	Type     string // "llm", "tool", "simulation"
	FolderID string
}

// AgentResponseTestSummary represents a summary of a response test.
type AgentResponseTestSummary struct {
	ID       string
	Name     string
	FolderID string
}

// TestSummary represents a summary of a test (used for listing).
type TestSummary struct {
	ID            string
	Name          string
	FolderID      string
	CreatedAtUnix int
	EntityType    string
	Path          []AgentTestFolderSegment
}

// CreateTestFolderRequest contains options for creating a test folder.
type CreateTestFolderRequest struct {
	Name           string
	ParentFolderID string
}

// UpdateTestFolderRequest contains options for updating a test folder.
type UpdateTestFolderRequest struct {
	Name string
}

// ListTestsOptions contains options for listing tests.
type ListTestsOptions struct {
	PageSize       int
	Cursor         string
	ParentFolderID string
	Search         string
}

// ListTestsResponse contains the paginated list of tests.
type ListTestsResponse struct {
	Tests      []*TestSummary
	HasMore    bool
	NextCursor string
}

// BulkMoveTestsRequest contains options for bulk moving tests.
type BulkMoveTestsRequest struct {
	TestIDs        []string
	TargetFolderID string
}

// CreateFolder creates a new test folder.
func (s *AgentTestingService) CreateFolder(ctx context.Context, req *CreateTestFolderRequest) (*AgentTestFolder, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	body := &api.BodyCreateAgentTestFolderV1ConvaiAgentTestingFoldersPost{
		Name: req.Name,
	}

	if req.ParentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.ParentFolderID, Set: true}
	}

	resp, err := s.client.apiClient.CreateAgentTestFolderRoute(ctx, body, api.CreateAgentTestFolderRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.CreateAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &AgentTestFolder{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

// GetFolder retrieves a test folder by ID.
func (s *AgentTestingService) GetFolder(ctx context.Context, folderID string) (*AgentTestFolder, error) {
	if folderID == "" {
		return nil, errors.New("folder ID is required")
	}

	resp, err := s.client.apiClient.GetAgentTestFolderRoute(ctx, api.GetAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return getFolderFromResponse(result), nil
}

// UpdateFolder updates a test folder.
func (s *AgentTestingService) UpdateFolder(ctx context.Context, folderID string, req *UpdateTestFolderRequest) (*AgentTestFolder, error) {
	if folderID == "" {
		return nil, errors.New("folder ID is required")
	}
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	body := &api.BodyUpdateAgentTestFolderV1ConvaiAgentTestingFoldersFolderIDPatch{
		Name: req.Name,
	}

	resp, err := s.client.apiClient.UpdateAgentTestFolderRoute(ctx, body, api.UpdateAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return getFolderFromResponse(result), nil
}

// DeleteFolder deletes a test folder.
func (s *AgentTestingService) DeleteFolder(ctx context.Context, folderID string) error {
	if folderID == "" {
		return errors.New("folder ID is required")
	}

	_, err := s.client.apiClient.DeleteAgentTestFolderRoute(ctx, api.DeleteAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	return err
}

// GetResponseTest retrieves a response test by ID.
func (s *AgentTestingService) GetResponseTest(ctx context.Context, testID string) (*AgentResponseTest, error) {
	if testID == "" {
		return nil, errors.New("test ID is required")
	}

	resp, err := s.client.apiClient.GetAgentResponseTestRoute(ctx, api.GetAgentResponseTestRouteParams{
		TestID: testID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetAgentResponseTestRouteOK)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return responseTestFromOK(result), nil
}

// GetResponseTestSummaries retrieves summaries for multiple response tests.
func (s *AgentTestingService) GetResponseTestSummaries(ctx context.Context, testIDs []string) (map[string]*AgentResponseTestSummary, error) {
	if len(testIDs) == 0 {
		return nil, errors.New("test IDs are required")
	}

	body := &api.ListTestsByIdsRequestModel{
		TestIds: testIDs,
	}

	resp, err := s.client.apiClient.GetAgentResponseTestsSummariesRoute(ctx, body, api.GetAgentResponseTestsSummariesRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetTestsSummariesByIdsResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	summaries := make(map[string]*AgentResponseTestSummary, len(result.Tests))
	for id, t := range result.Tests {
		summaries[id] = &AgentResponseTestSummary{
			ID:       t.ID,
			Name:     t.Name,
			FolderID: t.FolderParentID.Value,
		}
	}

	return summaries, nil
}

// ListTests lists tests and folders.
func (s *AgentTestingService) ListTests(ctx context.Context, opts *ListTestsOptions) (*ListTestsResponse, error) {
	params := api.ListChatResponseTestsRouteParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.Cursor != "" {
			params.Cursor = api.OptNilString{Value: opts.Cursor, Set: true}
		}
		if opts.ParentFolderID != "" {
			params.ParentFolderID = api.OptNilString{Value: opts.ParentFolderID, Set: true}
		}
		if opts.Search != "" {
			params.Search = api.OptNilString{Value: opts.Search, Set: true}
		}
	}

	resp, err := s.client.apiClient.ListChatResponseTestsRoute(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetTestsPageResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	tests := make([]*TestSummary, len(result.Tests))
	for i, t := range result.Tests {
		tests[i] = testSummaryFromResponse(&t)
	}

	return &ListTestsResponse{
		Tests:      tests,
		HasMore:    result.HasMore,
		NextCursor: result.NextCursor.Value,
	}, nil
}

// DeleteTest deletes a test.
func (s *AgentTestingService) DeleteTest(ctx context.Context, testID string) error {
	if testID == "" {
		return errors.New("test ID is required")
	}

	_, err := s.client.apiClient.DeleteChatResponseTestRoute(ctx, api.DeleteChatResponseTestRouteParams{
		TestID: testID,
	})
	return err
}

// BulkMoveTests moves multiple tests to a folder.
func (s *AgentTestingService) BulkMoveTests(ctx context.Context, req *BulkMoveTestsRequest) error {
	if req == nil {
		return errors.New("request is required")
	}
	if len(req.TestIDs) == 0 {
		return errors.New("test IDs are required")
	}

	body := &api.BodyBulkMoveTestsToFolderV1ConvaiAgentTestingBulkMovePost{
		EntityIds: req.TestIDs,
	}

	if req.TargetFolderID != "" {
		body.MoveTo = api.OptNilString{Value: req.TargetFolderID, Set: true}
	}

	_, err := s.client.apiClient.AgentTestingBulkMoveRoute(ctx, body, api.AgentTestingBulkMoveRouteParams{})
	return err
}

// getFolderFromResponse converts API response to AgentTestFolder.
func getFolderFromResponse(r *api.GetAgentTestFolderResponseModel) *AgentTestFolder {
	path := make([]AgentTestFolderSegment, len(r.FolderPath))
	for i, p := range r.FolderPath {
		path[i] = AgentTestFolderSegment{
			ID:   p.ID,
			Name: p.Name.Value,
		}
	}

	childrenCount := 0
	if r.ChildrenCount.Set {
		childrenCount = r.ChildrenCount.Value
	}

	return &AgentTestFolder{
		ID:            r.ID,
		Name:          r.Name,
		ChildrenCount: childrenCount,
		Path:          path,
	}
}

// responseTestFromOK converts API response to AgentResponseTest.
func responseTestFromOK(r *api.GetAgentResponseTestRouteOK) *AgentResponseTest {
	test := &AgentResponseTest{
		Type: string(r.Type),
	}

	switch r.Type {
	case api.GetResponseUnitTestResponseModelGetAgentResponseTestRouteOK:
		test.ID = r.GetResponseUnitTestResponseModel.ID
		test.Name = r.GetResponseUnitTestResponseModel.Name
	case api.GetToolCallUnitTestResponseModelGetAgentResponseTestRouteOK:
		test.ID = r.GetToolCallUnitTestResponseModel.ID
		test.Name = r.GetToolCallUnitTestResponseModel.Name
	case api.GetSimulationTestResponseModelGetAgentResponseTestRouteOK:
		test.ID = r.GetSimulationTestResponseModel.ID
		test.Name = r.GetSimulationTestResponseModel.Name
	}

	return test
}

// testSummaryFromResponse converts API response to TestSummary.
func testSummaryFromResponse(t *api.UnitTestSummaryResponseModel) *TestSummary {
	path := make([]AgentTestFolderSegment, len(t.FolderPath))
	for i, p := range t.FolderPath {
		path[i] = AgentTestFolderSegment{
			ID:   p.ID,
			Name: p.Name.Value,
		}
	}

	entityType := ""
	if t.EntityType.Set {
		entityType = string(t.EntityType.Value)
	}

	return &TestSummary{
		ID:            t.ID,
		Name:          t.Name,
		FolderID:      t.FolderParentID.Value,
		CreatedAtUnix: t.CreatedAtUnixSecs,
		EntityType:    entityType,
		Path:          path,
	}
}
