package agents

import (
	"context"
	"errors"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// TestFolder represents a folder for organizing tests.
type TestFolder struct {
	ID            string
	Name          string
	ChildrenCount int
	Path          []TestFolderSegment
}

// TestFolderSegment represents a segment in the folder path.
type TestFolderSegment struct {
	ID   string
	Name string
}

// ResponseTest represents a response test for an agent.
type ResponseTest struct {
	ID       string
	Name     string
	Type     string // "llm", "tool", "simulation"
	FolderID string
}

// ResponseTestSummary represents a summary of a response test.
type ResponseTestSummary struct {
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
	Path          []TestFolderSegment
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

// CreateTestFolder creates a new test folder.
func (s *Service) CreateTestFolder(ctx context.Context, req *CreateTestFolderRequest) (*TestFolder, error) {
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

	resp, err := s.client.API().CreateAgentTestFolderRoute(ctx, body, api.CreateAgentTestFolderRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.CreateAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &TestFolder{
		ID:   result.ID,
		Name: result.Name,
	}, nil
}

// GetTestFolder retrieves a test folder by ID.
func (s *Service) GetTestFolder(ctx context.Context, folderID string) (*TestFolder, error) {
	if folderID == "" {
		return nil, errors.New("folder ID is required")
	}

	resp, err := s.client.API().GetAgentTestFolderRoute(ctx, api.GetAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return testFolderFromResponse(result), nil
}

// UpdateTestFolder updates a test folder.
func (s *Service) UpdateTestFolder(ctx context.Context, folderID string, req *UpdateTestFolderRequest) (*TestFolder, error) {
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

	resp, err := s.client.API().UpdateAgentTestFolderRoute(ctx, body, api.UpdateAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetAgentTestFolderResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return testFolderFromResponse(result), nil
}

// DeleteTestFolder deletes a test folder.
func (s *Service) DeleteTestFolder(ctx context.Context, folderID string) error {
	if folderID == "" {
		return errors.New("folder ID is required")
	}

	_, err := s.client.API().DeleteAgentTestFolderRoute(ctx, api.DeleteAgentTestFolderRouteParams{
		FolderID: folderID,
	})
	return err
}

// GetResponseTest retrieves a response test by ID.
func (s *Service) GetResponseTest(ctx context.Context, testID string) (*ResponseTest, error) {
	if testID == "" {
		return nil, errors.New("test ID is required")
	}

	resp, err := s.client.API().GetAgentResponseTestRoute(ctx, api.GetAgentResponseTestRouteParams{
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
func (s *Service) GetResponseTestSummaries(ctx context.Context, testIDs []string) (map[string]*ResponseTestSummary, error) {
	if len(testIDs) == 0 {
		return nil, errors.New("test IDs are required")
	}

	body := &api.ListTestsByIdsRequestModel{
		TestIds: testIDs,
	}

	resp, err := s.client.API().GetAgentResponseTestsSummariesRoute(ctx, body, api.GetAgentResponseTestsSummariesRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetTestsSummariesByIdsResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	summaries := make(map[string]*ResponseTestSummary, len(result.Tests))
	for id, t := range result.Tests {
		summaries[id] = &ResponseTestSummary{
			ID:       t.ID,
			Name:     t.Name,
			FolderID: t.FolderParentID.Value,
		}
	}

	return summaries, nil
}

// ListTests lists tests and folders.
func (s *Service) ListTests(ctx context.Context, opts *ListTestsOptions) (*ListTestsResponse, error) {
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

	resp, err := s.client.API().ListChatResponseTestsRoute(ctx, params)
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
func (s *Service) DeleteTest(ctx context.Context, testID string) error {
	if testID == "" {
		return errors.New("test ID is required")
	}

	_, err := s.client.API().DeleteChatResponseTestRoute(ctx, api.DeleteChatResponseTestRouteParams{
		TestID: testID,
	})
	return err
}

// BulkMoveTests moves multiple tests to a folder.
func (s *Service) BulkMoveTests(ctx context.Context, req *BulkMoveTestsRequest) error {
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

	_, err := s.client.API().AgentTestingBulkMoveRoute(ctx, body, api.AgentTestingBulkMoveRouteParams{})
	return err
}

// testFolderFromResponse converts API response to TestFolder.
func testFolderFromResponse(r *api.GetAgentTestFolderResponseModel) *TestFolder {
	path := make([]TestFolderSegment, len(r.FolderPath))
	for i, p := range r.FolderPath {
		path[i] = TestFolderSegment{
			ID:   p.ID,
			Name: p.Name.Value,
		}
	}

	childrenCount := 0
	if r.ChildrenCount.Set {
		childrenCount = r.ChildrenCount.Value
	}

	return &TestFolder{
		ID:            r.ID,
		Name:          r.Name,
		ChildrenCount: childrenCount,
		Path:          path,
	}
}

// responseTestFromOK converts API response to ResponseTest.
func responseTestFromOK(r *api.GetAgentResponseTestRouteOK) *ResponseTest {
	test := &ResponseTest{
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
	path := make([]TestFolderSegment, len(t.FolderPath))
	for i, p := range t.FolderPath {
		path[i] = TestFolderSegment{
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
