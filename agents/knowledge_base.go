package agents

import (
	"bytes"
	"context"
	"errors"
	"io"

	ht "github.com/ogen-go/ogen/http"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Document represents a knowledge base document.
type Document struct {
	// ID is the unique identifier for the document.
	ID string

	// Name is the display name of the document.
	Name string

	// Type is the document type: "file", "text", "url", or "folder".
	Type string

	// FolderPath contains the path segment IDs to the document's folder.
	FolderPath []string
}

// DocumentSummary is a lightweight document listing.
type DocumentSummary struct {
	// ID is the unique identifier.
	ID string

	// Name is the display name.
	Name string

	// Type is the document type.
	Type string

	// FolderParentID is the ID of the parent folder.
	FolderParentID string
}

// DocumentChunk represents a chunk of a document.
type DocumentChunk struct {
	// ID is the unique identifier for the chunk.
	ID string

	// Name is the chunk name.
	Name string

	// Content is the chunk text content.
	Content string
}

// KnowledgeBaseFolder represents a folder for organizing documents.
type KnowledgeBaseFolder struct {
	// ID is the unique identifier for the folder.
	ID string

	// Name is the display name of the folder.
	Name string

	// FolderPath contains the path segment IDs to the folder.
	FolderPath []string
}

// DependentAgent represents an agent that depends on a knowledge base document.
type DependentAgent struct {
	// ID is the ID of the dependent agent.
	ID string

	// Name is the name of the agent.
	Name string

	// Type is the agent type.
	Type string
}

// CreateFileDocumentRequest contains options for creating a file document.
type CreateFileDocumentRequest struct {
	// Name is the display name for the document.
	Name string

	// File is the file content to upload.
	File io.Reader

	// Filename is the original filename.
	Filename string

	// FolderID is the optional parent folder.
	FolderID string
}

// CreateTextDocumentRequest contains options for creating a text document.
type CreateTextDocumentRequest struct {
	// Name is the display name for the document.
	Name string

	// Content is the text content.
	Content string

	// FolderID is the optional parent folder.
	FolderID string
}

// CreateURLDocumentRequest contains options for creating a URL document.
type CreateURLDocumentRequest struct {
	// Name is the display name for the document.
	Name string

	// URL is the source URL to crawl.
	URL string

	// FolderID is the optional parent folder.
	FolderID string

	// EnableAutoSync enables auto-sync for URL documents.
	EnableAutoSync bool

	// AutoRemove removes the document if URL becomes unavailable.
	AutoRemove bool
}

// CreateKnowledgeBaseFolderRequest contains options for creating a folder.
type CreateKnowledgeBaseFolderRequest struct {
	// Name is the display name for the folder.
	Name string

	// ParentFolderID is the optional parent folder.
	ParentFolderID string
}

// ListDocumentsOptions configures the ListDocuments operation.
type ListDocumentsOptions struct {
	// PageSize is the maximum number of documents to return (max 100).
	PageSize int

	// Search filters by document name prefix.
	Search string
}

// ListDocumentsResponse contains paginated document results.
type ListDocumentsResponse struct {
	// Documents is the list of document summaries.
	Documents []*DocumentSummary

	// HasMore indicates if there are more results.
	HasMore bool

	// NextCursor is the cursor for the next page.
	NextCursor string
}

// CreateFileDocument creates a document from a file upload.
func (s *Service) CreateFileDocument(ctx context.Context, req *CreateFileDocumentRequest) (*Document, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.File == nil {
		return nil, errors.New("file is required")
	}

	data, err := io.ReadAll(req.File)
	if err != nil {
		return nil, err
	}

	filename := req.Filename
	if filename == "" {
		filename = "document"
	}

	body := &api.BodyCreateFileDocumentV1ConvaiKnowledgeBaseFilePostMultipart{
		File: ht.MultipartFile{
			Name: filename,
			File: bytes.NewReader(data),
		},
	}
	if req.Name != "" {
		body.Name = api.OptNilString{Value: req.Name, Set: true}
	}
	if req.FolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.FolderID, Set: true}
	}

	resp, err := s.client.API().CreateFileDocumentRoute(ctx, body, api.CreateFileDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return documentFromAddResponse(result, "file"), nil
}

// CreateTextDocument creates a document from text content.
func (s *Service) CreateTextDocument(ctx context.Context, req *CreateTextDocumentRequest) (*Document, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}

	body := &api.BodyCreateTextDocumentV1ConvaiKnowledgeBaseTextPost{
		Text: req.Content,
	}
	if req.Name != "" {
		body.Name = api.OptNilString{Value: req.Name, Set: true}
	}
	if req.FolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.FolderID, Set: true}
	}

	resp, err := s.client.API().CreateTextDocumentRoute(ctx, body, api.CreateTextDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return documentFromAddResponse(result, "text"), nil
}

// CreateURLDocument creates a document from a URL.
func (s *Service) CreateURLDocument(ctx context.Context, req *CreateURLDocumentRequest) (*Document, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.URL == "" {
		return nil, errors.New("url is required")
	}

	body := &api.BodyCreateURLDocumentV1ConvaiKnowledgeBaseURLPost{
		URL: req.URL,
	}
	if req.Name != "" {
		body.Name = api.OptNilString{Value: req.Name, Set: true}
	}
	if req.FolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.FolderID, Set: true}
	}
	if req.EnableAutoSync {
		body.EnableAutoSync = api.NewOptBool(true)
	}
	if req.AutoRemove {
		body.AutoRemove = api.NewOptBool(true)
	}

	resp, err := s.client.API().CreateURLDocumentRoute(ctx, body, api.CreateURLDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return documentFromAddResponse(result, "url"), nil
}

// CreateKnowledgeBaseFolder creates a folder for organizing documents.
func (s *Service) CreateKnowledgeBaseFolder(ctx context.Context, req *CreateKnowledgeBaseFolderRequest) (*KnowledgeBaseFolder, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	body := &api.BodyCreateFolderV1ConvaiKnowledgeBaseFolderPost{
		Name: req.Name,
	}
	if req.ParentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.ParentFolderID, Set: true}
	}

	resp, err := s.client.API().CreateFolderRoute(ctx, body, api.CreateFolderRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return folderFromAddResponse(result), nil
}

// GetDocument retrieves a document by ID.
func (s *Service) GetDocument(ctx context.Context, documentID string) (*Document, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.API().GetDocumentationFromKnowledgeBase(ctx, api.GetDocumentationFromKnowledgeBaseParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetDocumentationFromKnowledgeBaseOK:
		return documentFromUnion(r), nil
	default:
		return nil, errors.New("unexpected response type")
	}
}

// DeleteDocument removes a document.
func (s *Service) DeleteDocument(ctx context.Context, documentID string) error {
	if documentID == "" {
		return errors.New("document ID is required")
	}

	_, err := s.client.API().DeleteKnowledgeBaseDocument(ctx, api.DeleteKnowledgeBaseDocumentParams{
		DocumentationID: documentID,
	})
	return err
}

// ListDocuments returns a paginated list of documents.
func (s *Service) ListDocuments(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error) {
	params := api.GetKnowledgeBaseListRouteParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.Search != "" {
			params.Search = api.OptNilString{Value: opts.Search, Set: true}
		}
	}

	resp, err := s.client.API().GetKnowledgeBaseListRoute(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetKnowledgeBaseListResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	docs := make([]*DocumentSummary, 0, len(result.Documents))
	for _, d := range result.Documents {
		docs = append(docs, documentSummaryFromUnion(&d))
	}

	return &ListDocumentsResponse{
		Documents:  docs,
		HasMore:    result.HasMore,
		NextCursor: result.NextCursor.Value,
	}, nil
}

// GetDocumentContent retrieves the content of a document.
func (s *Service) GetDocumentContent(ctx context.Context, documentID string) (io.ReadCloser, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.API().GetKnowledgeBaseContent(ctx, api.GetKnowledgeBaseContentParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetKnowledgeBaseContentOK)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return io.NopCloser(result.Data), nil
}

// GetDocumentChunks retrieves all chunks for a document.
func (s *Service) GetDocumentChunks(ctx context.Context, documentID string) ([]*DocumentChunk, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.API().GetDocumentationChunksFromKnowledgeBase(ctx, api.GetDocumentationChunksFromKnowledgeBaseParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.KnowledgeBaseDocumentChunksResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	chunks := make([]*DocumentChunk, 0, len(result.Chunks))
	for _, c := range result.Chunks {
		chunks = append(chunks, &DocumentChunk{
			ID:      c.ID,
			Name:    c.Name,
			Content: c.Content,
		})
	}

	return chunks, nil
}

// GetDocumentChunk retrieves a specific chunk from a document by chunk ID.
func (s *Service) GetDocumentChunk(ctx context.Context, documentID, chunkID string) (*DocumentChunk, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}
	if chunkID == "" {
		return nil, errors.New("chunk ID is required")
	}

	resp, err := s.client.API().GetDocumentationChunkFromKnowledgeBase(ctx, api.GetDocumentationChunkFromKnowledgeBaseParams{
		DocumentationID: documentID,
		ChunkID:         chunkID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.KnowledgeBaseDocumentChunkResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &DocumentChunk{
		ID:      result.ID,
		Name:    result.Name,
		Content: result.Content,
	}, nil
}

// GetDocumentDependentAgents retrieves agents that depend on a document.
func (s *Service) GetDocumentDependentAgents(ctx context.Context, documentID string) ([]*DependentAgent, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.API().GetKnowledgeBaseDependentAgents(ctx, api.GetKnowledgeBaseDependentAgentsParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetKnowledgeBaseDependentAgentsResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	agents := make([]*DependentAgent, 0, len(result.Agents))
	for _, a := range result.Agents {
		agent := dependentAgentFromUnion(&a)
		if agent != nil {
			agents = append(agents, agent)
		}
	}

	return agents, nil
}

// MoveDocument moves a document to a different folder.
func (s *Service) MoveDocument(ctx context.Context, documentID, targetFolderID string) error {
	if documentID == "" {
		return errors.New("document ID is required")
	}

	var body api.OptBodyMoveEntityToFolderV1ConvaiKnowledgeBaseDocumentIDMovePost
	if targetFolderID != "" {
		body = api.NewOptBodyMoveEntityToFolderV1ConvaiKnowledgeBaseDocumentIDMovePost(
			api.BodyMoveEntityToFolderV1ConvaiKnowledgeBaseDocumentIDMovePost{
				MoveTo: api.OptNilString{Value: targetFolderID, Set: true},
			},
		)
	}

	_, err := s.client.API().PostKnowledgeBaseMoveRoute(ctx, body, api.PostKnowledgeBaseMoveRouteParams{
		DocumentID: documentID,
	})
	return err
}

// BulkMoveDocuments moves multiple documents to a folder.
func (s *Service) BulkMoveDocuments(ctx context.Context, documentIDs []string, targetFolderID string) error {
	if len(documentIDs) == 0 {
		return errors.New("document IDs are required")
	}

	body := &api.BodyBulkMoveEntitiesToFolderV1ConvaiKnowledgeBaseBulkMovePost{
		DocumentIds: documentIDs,
	}
	if targetFolderID != "" {
		body.MoveTo = api.OptNilString{Value: targetFolderID, Set: true}
	}

	_, err := s.client.API().PostKnowledgeBaseBulkMoveRoute(ctx, body, api.PostKnowledgeBaseBulkMoveRouteParams{})
	return err
}

// RefreshURLDocument refreshes the content of a URL document.
func (s *Service) RefreshURLDocument(ctx context.Context, documentID string) error {
	if documentID == "" {
		return errors.New("document ID is required")
	}

	_, err := s.client.API().RefreshURLDocumentRoute(ctx, api.RefreshURLDocumentRouteParams{
		DocumentationID: documentID,
	})
	return err
}

// documentFromAddResponse converts API add response to Document.
func documentFromAddResponse(r *api.AddKnowledgeBaseResponseModel, docType string) *Document {
	doc := &Document{
		ID:   r.ID,
		Name: r.Name,
		Type: docType,
	}

	if len(r.FolderPath) > 0 {
		doc.FolderPath = make([]string, len(r.FolderPath))
		for i, p := range r.FolderPath {
			doc.FolderPath[i] = p.ID
		}
	}

	return doc
}

// folderFromAddResponse converts API add response to KnowledgeBaseFolder.
func folderFromAddResponse(r *api.AddKnowledgeBaseResponseModel) *KnowledgeBaseFolder {
	folder := &KnowledgeBaseFolder{
		ID:   r.ID,
		Name: r.Name,
	}

	if len(r.FolderPath) > 0 {
		folder.FolderPath = make([]string, len(r.FolderPath))
		for i, p := range r.FolderPath {
			folder.FolderPath[i] = p.ID
		}
	}

	return folder
}

// documentFromUnion converts API union response to Document.
func documentFromUnion(r *api.GetDocumentationFromKnowledgeBaseOK) *Document {
	doc := &Document{}

	switch r.Type {
	case api.GetKnowledgeBaseFileResponseModelGetDocumentationFromKnowledgeBaseOK:
		f := r.GetKnowledgeBaseFileResponseModel
		doc.ID = f.ID
		doc.Name = f.Name
		doc.Type = string(f.Type)
	case api.GetKnowledgeBaseTextResponseModelGetDocumentationFromKnowledgeBaseOK:
		t := r.GetKnowledgeBaseTextResponseModel
		doc.ID = t.ID
		doc.Name = t.Name
		doc.Type = string(t.Type)
	case api.GetKnowledgeBaseURLResponseModelGetDocumentationFromKnowledgeBaseOK:
		u := r.GetKnowledgeBaseURLResponseModel
		doc.ID = u.ID
		doc.Name = u.Name
		doc.Type = string(u.Type)
	case api.GetKnowledgeBaseFolderResponseModelGetDocumentationFromKnowledgeBaseOK:
		fo := r.GetKnowledgeBaseFolderResponseModel
		doc.ID = fo.ID
		doc.Name = fo.Name
		doc.Type = string(fo.Type)
	}

	return doc
}

// documentSummaryFromUnion converts API union list item to DocumentSummary.
func documentSummaryFromUnion(d *api.GetKnowledgeBaseListResponseModelDocumentsItem) *DocumentSummary {
	summary := &DocumentSummary{}

	switch d.Type {
	case api.GetKnowledgeBaseSummaryFileResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		f := d.GetKnowledgeBaseSummaryFileResponseModel
		summary.ID = f.ID
		summary.Name = f.Name
		summary.Type = string(f.Type)
		if f.FolderParentID.Set && !f.FolderParentID.Null {
			summary.FolderParentID = f.FolderParentID.Value
		}
	case api.GetKnowledgeBaseSummaryTextResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		t := d.GetKnowledgeBaseSummaryTextResponseModel
		summary.ID = t.ID
		summary.Name = t.Name
		summary.Type = string(t.Type)
		if t.FolderParentID.Set && !t.FolderParentID.Null {
			summary.FolderParentID = t.FolderParentID.Value
		}
	case api.GetKnowledgeBaseSummaryURLResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		u := d.GetKnowledgeBaseSummaryURLResponseModel
		summary.ID = u.ID
		summary.Name = u.Name
		summary.Type = string(u.Type)
		if u.FolderParentID.Set && !u.FolderParentID.Null {
			summary.FolderParentID = u.FolderParentID.Value
		}
	case api.GetKnowledgeBaseSummaryFolderResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		fo := d.GetKnowledgeBaseSummaryFolderResponseModel
		summary.ID = fo.ID
		summary.Name = fo.Name
		summary.Type = string(fo.Type)
		if fo.FolderParentID.Set && !fo.FolderParentID.Null {
			summary.FolderParentID = fo.FolderParentID.Value
		}
	}

	return summary
}

// dependentAgentFromUnion converts API union agent to DependentAgent.
func dependentAgentFromUnion(a *api.GetKnowledgeBaseDependentAgentsResponseModelAgentsItem) *DependentAgent {
	switch a.Type {
	case api.DependentAvailableAgentIdentifierGetKnowledgeBaseDependentAgentsResponseModelAgentsItem:
		agent := a.DependentAvailableAgentIdentifier
		agentType := ""
		if agent.Type.Set {
			agentType = string(agent.Type.Value)
		}
		return &DependentAgent{
			ID:   agent.ID,
			Name: agent.Name,
			Type: agentType,
		}
	case api.DependentUnknownAgentIdentifierGetKnowledgeBaseDependentAgentsResponseModelAgentsItem:
		// Unknown agents don't have accessible details
		return nil
	}
	return nil
}
