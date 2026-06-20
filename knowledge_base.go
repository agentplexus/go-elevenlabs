package elevenlabs

import (
	"context"
	"errors"
	"io"

	ht "github.com/ogen-go/ogen/http"
	"github.com/plexusone/elevenlabs-go/internal/api"
)

// KnowledgeBaseService handles knowledge base document operations.
type KnowledgeBaseService struct {
	client *Client
}

// KnowledgeBaseDocument represents a document in the knowledge base.
type KnowledgeBaseDocument struct {
	ID              string
	Name            string
	Type            string // "file", "text", "url", "folder"
	SourceURL       string
	Content         string
	ParentFolderID  string
	CreatedAtUnix   int
	LastUpdatedUnix int
	SizeBytes       int
}

// KnowledgeBaseChunk represents a chunk of a knowledge base document.
type KnowledgeBaseChunk struct {
	ID      string
	Name    string
	Content string
}

// RAGIndex represents a RAG index for a document.
type RAGIndex struct {
	ID            string
	DocumentID    string
	Status        string
	ChunkCount    int
	CreatedAtUnix int
}

// CreateFileDocumentRequest contains options for creating a file document.
type CreateFileDocumentRequest struct {
	// File to upload
	File io.Reader
	// File name
	FileName string
	// Optional name for the document
	Name string
	// Optional parent folder ID
	ParentFolderID string
}

// CreateTextDocumentRequest contains options for creating a text document.
type CreateTextDocumentRequest struct {
	// Text content
	Content string
	// Document name
	Name string
	// Optional parent folder ID
	ParentFolderID string
}

// CreateURLDocumentRequest contains options for creating a URL document.
type CreateURLDocumentRequest struct {
	// URL to index
	URL string
	// Optional document name
	Name string
	// Optional parent folder ID
	ParentFolderID string
}

// ListDocumentsOptions contains options for listing documents.
type ListDocumentsOptions struct {
	// Pagination cursor
	Cursor string
	// Maximum results
	PageSize int
	// Filter by parent folder
	ParentFolderID string
	// Search query
	Search string
}

// ListDocumentsResponse contains the paginated list of documents.
type ListDocumentsResponse struct {
	Documents  []*KnowledgeBaseDocument
	HasMore    bool
	NextCursor string
}

// CreateFileDocument uploads a file to the knowledge base.
func (s *KnowledgeBaseService) CreateFileDocument(ctx context.Context, req *CreateFileDocumentRequest) (*KnowledgeBaseDocument, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.File == nil {
		return nil, errors.New("file is required")
	}
	if req.FileName == "" {
		return nil, errors.New("file name is required")
	}

	body := &api.BodyCreateFileDocumentV1ConvaiKnowledgeBaseFilePostMultipart{
		File: ht.MultipartFile{
			Name: req.FileName,
			File: req.File,
		},
	}

	if req.Name != "" {
		body.Name = api.OptNilString{Value: req.Name, Set: true}
	}
	if req.ParentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.ParentFolderID, Set: true}
	}

	resp, err := s.client.apiClient.CreateFileDocumentRoute(ctx, body, api.CreateFileDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &KnowledgeBaseDocument{
		ID:   result.ID,
		Name: result.Name,
		Type: "file",
	}, nil
}

// CreateTextDocument creates a text document in the knowledge base.
func (s *KnowledgeBaseService) CreateTextDocument(ctx context.Context, req *CreateTextDocumentRequest) (*KnowledgeBaseDocument, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.Content == "" {
		return nil, errors.New("content is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	body := &api.BodyCreateTextDocumentV1ConvaiKnowledgeBaseTextPost{
		Name: api.OptNilString{Value: req.Name, Set: true},
		Text: req.Content,
	}

	if req.ParentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.ParentFolderID, Set: true}
	}

	resp, err := s.client.apiClient.CreateTextDocumentRoute(ctx, body, api.CreateTextDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &KnowledgeBaseDocument{
		ID:      result.ID,
		Name:    result.Name,
		Type:    "text",
		Content: req.Content,
	}, nil
}

// CreateURLDocument creates a URL document in the knowledge base.
func (s *KnowledgeBaseService) CreateURLDocument(ctx context.Context, req *CreateURLDocumentRequest) (*KnowledgeBaseDocument, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if req.URL == "" {
		return nil, errors.New("URL is required")
	}

	body := &api.BodyCreateURLDocumentV1ConvaiKnowledgeBaseURLPost{
		URL: req.URL,
	}

	if req.Name != "" {
		body.Name = api.OptNilString{Value: req.Name, Set: true}
	}
	if req.ParentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: req.ParentFolderID, Set: true}
	}

	resp, err := s.client.apiClient.CreateURLDocumentRoute(ctx, body, api.CreateURLDocumentRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &KnowledgeBaseDocument{
		ID:        result.ID,
		Name:      result.Name,
		Type:      "url",
		SourceURL: req.URL,
	}, nil
}

// Get retrieves a document from the knowledge base.
func (s *KnowledgeBaseService) Get(ctx context.Context, documentID string) (*KnowledgeBaseDocument, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.apiClient.GetDocumentationFromKnowledgeBase(ctx, api.GetDocumentationFromKnowledgeBaseParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	return documentFromResponse(resp)
}

// Delete removes a document from the knowledge base.
func (s *KnowledgeBaseService) Delete(ctx context.Context, documentID string) error {
	if documentID == "" {
		return errors.New("document ID is required")
	}

	_, err := s.client.apiClient.DeleteKnowledgeBaseDocument(ctx, api.DeleteKnowledgeBaseDocumentParams{
		DocumentationID: documentID,
	})
	return err
}

// List returns a paginated list of documents.
func (s *KnowledgeBaseService) List(ctx context.Context, opts *ListDocumentsOptions) (*ListDocumentsResponse, error) {
	params := api.GetKnowledgeBaseListRouteParams{}

	if opts != nil {
		if opts.Cursor != "" {
			params.Cursor = api.OptNilString{Value: opts.Cursor, Set: true}
		}
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.ParentFolderID != "" {
			params.ParentFolderID = api.OptNilString{Value: opts.ParentFolderID, Set: true}
		}
		if opts.Search != "" {
			params.Search = api.OptNilString{Value: opts.Search, Set: true}
		}
	}

	resp, err := s.client.apiClient.GetKnowledgeBaseListRoute(ctx, params)
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetKnowledgeBaseListResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	documents := make([]*KnowledgeBaseDocument, len(result.Documents))
	for i, doc := range result.Documents {
		documents[i] = documentItemToDocument(&doc)
	}

	return &ListDocumentsResponse{
		Documents:  documents,
		HasMore:    result.HasMore,
		NextCursor: result.NextCursor.Value,
	}, nil
}

// GetContent retrieves the content of a document as a reader.
func (s *KnowledgeBaseService) GetContent(ctx context.Context, documentID string) (io.ReadCloser, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.apiClient.GetKnowledgeBaseContent(ctx, api.GetKnowledgeBaseContentParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.GetKnowledgeBaseContentOK)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	// Wrap the reader in a ReadCloser
	if rc, ok := result.Data.(io.ReadCloser); ok {
		return rc, nil
	}
	return io.NopCloser(result.Data), nil
}

// GetChunks retrieves chunks of a document.
func (s *KnowledgeBaseService) GetChunks(ctx context.Context, documentID string) ([]*KnowledgeBaseChunk, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.apiClient.GetDocumentationChunksFromKnowledgeBase(ctx, api.GetDocumentationChunksFromKnowledgeBaseParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.KnowledgeBaseDocumentChunksResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	chunks := make([]*KnowledgeBaseChunk, len(result.Chunks))
	for i, c := range result.Chunks {
		chunks[i] = &KnowledgeBaseChunk{
			ID:      c.ID,
			Name:    c.Name,
			Content: c.Content,
		}
	}

	return chunks, nil
}

// GetChunk retrieves a specific chunk of a document.
func (s *KnowledgeBaseService) GetChunk(ctx context.Context, documentID, chunkID string) (*KnowledgeBaseChunk, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}
	if chunkID == "" {
		return nil, errors.New("chunk ID is required")
	}

	resp, err := s.client.apiClient.GetDocumentationChunkFromKnowledgeBase(ctx, api.GetDocumentationChunkFromKnowledgeBaseParams{
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

	return &KnowledgeBaseChunk{
		ID:      result.ID,
		Name:    result.Name,
		Content: result.Content,
	}, nil
}

// DependentAgent represents an agent that depends on a knowledge base document.
type DependentAgent struct {
	ID   string
	Name string
}

// GetDependentAgents returns agents that use this document.
func (s *KnowledgeBaseService) GetDependentAgents(ctx context.Context, documentID string) ([]*DependentAgent, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.apiClient.GetKnowledgeBaseDependentAgents(ctx, api.GetKnowledgeBaseDependentAgentsParams{
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
		switch a.Type {
		case api.DependentAvailableAgentIdentifierGetKnowledgeBaseDependentAgentsResponseModelAgentsItem:
			agents = append(agents, &DependentAgent{
				ID:   a.DependentAvailableAgentIdentifier.ID,
				Name: a.DependentAvailableAgentIdentifier.Name,
			})
		case api.DependentUnknownAgentIdentifierGetKnowledgeBaseDependentAgentsResponseModelAgentsItem:
			agents = append(agents, &DependentAgent{
				ID: a.DependentUnknownAgentIdentifier.ID,
			})
		}
	}

	return agents, nil
}

// GetRAGIndexes retrieves RAG indexes for a document.
func (s *KnowledgeBaseService) GetRAGIndexes(ctx context.Context, documentID string) ([]*RAGIndex, error) {
	if documentID == "" {
		return nil, errors.New("document ID is required")
	}

	resp, err := s.client.apiClient.GetRagIndexes(ctx, api.GetRagIndexesParams{
		DocumentationID: documentID,
	})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.RAGDocumentIndexesResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	indexes := make([]*RAGIndex, len(result.Indexes))
	for i, idx := range result.Indexes {
		indexes[i] = &RAGIndex{
			ID:         idx.ID,
			DocumentID: documentID,
			Status:     string(idx.Status),
		}
	}

	return indexes, nil
}

// DeleteRAGIndex deletes a RAG index.
func (s *KnowledgeBaseService) DeleteRAGIndex(ctx context.Context, documentID, ragIndexID string) error {
	if documentID == "" {
		return errors.New("document ID is required")
	}
	if ragIndexID == "" {
		return errors.New("RAG index ID is required")
	}

	_, err := s.client.apiClient.DeleteRagIndex(ctx, api.DeleteRagIndexParams{
		DocumentationID: documentID,
		RagIndexID:      ragIndexID,
	})
	return err
}

// CreateFolder creates a folder in the knowledge base.
func (s *KnowledgeBaseService) CreateFolder(ctx context.Context, name string, parentFolderID string) (*KnowledgeBaseDocument, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	body := &api.BodyCreateFolderV1ConvaiKnowledgeBaseFolderPost{
		Name: name,
	}

	if parentFolderID != "" {
		body.ParentFolderID = api.OptNilString{Value: parentFolderID, Set: true}
	}

	resp, err := s.client.apiClient.CreateFolderRoute(ctx, body, api.CreateFolderRouteParams{})
	if err != nil {
		return nil, err
	}

	result, ok := resp.(*api.AddKnowledgeBaseResponseModel)
	if !ok {
		return nil, errors.New("unexpected response type")
	}

	return &KnowledgeBaseDocument{
		ID:   result.ID,
		Name: result.Name,
		Type: "folder",
	}, nil
}

// documentFromResponse converts API response to KnowledgeBaseDocument.
func documentFromResponse(resp api.GetDocumentationFromKnowledgeBaseRes) (*KnowledgeBaseDocument, error) {
	switch r := resp.(type) {
	case *api.GetDocumentationFromKnowledgeBaseOK:
		switch r.Type {
		case api.GetKnowledgeBaseFileResponseModelGetDocumentationFromKnowledgeBaseOK:
			return &KnowledgeBaseDocument{
				ID:              r.GetKnowledgeBaseFileResponseModel.ID,
				Name:            r.GetKnowledgeBaseFileResponseModel.Name,
				Type:            "file",
				CreatedAtUnix:   r.GetKnowledgeBaseFileResponseModel.Metadata.CreatedAtUnixSecs,
				LastUpdatedUnix: r.GetKnowledgeBaseFileResponseModel.Metadata.LastUpdatedAtUnixSecs,
				SizeBytes:       r.GetKnowledgeBaseFileResponseModel.Metadata.SizeBytes,
			}, nil
		case api.GetKnowledgeBaseTextResponseModelGetDocumentationFromKnowledgeBaseOK:
			return &KnowledgeBaseDocument{
				ID:              r.GetKnowledgeBaseTextResponseModel.ID,
				Name:            r.GetKnowledgeBaseTextResponseModel.Name,
				Type:            "text",
				Content:         r.GetKnowledgeBaseTextResponseModel.ExtractedInnerHTML,
				CreatedAtUnix:   r.GetKnowledgeBaseTextResponseModel.Metadata.CreatedAtUnixSecs,
				LastUpdatedUnix: r.GetKnowledgeBaseTextResponseModel.Metadata.LastUpdatedAtUnixSecs,
				SizeBytes:       r.GetKnowledgeBaseTextResponseModel.Metadata.SizeBytes,
			}, nil
		case api.GetKnowledgeBaseURLResponseModelGetDocumentationFromKnowledgeBaseOK:
			return &KnowledgeBaseDocument{
				ID:              r.GetKnowledgeBaseURLResponseModel.ID,
				Name:            r.GetKnowledgeBaseURLResponseModel.Name,
				Type:            "url",
				SourceURL:       r.GetKnowledgeBaseURLResponseModel.URL,
				CreatedAtUnix:   r.GetKnowledgeBaseURLResponseModel.Metadata.CreatedAtUnixSecs,
				LastUpdatedUnix: r.GetKnowledgeBaseURLResponseModel.Metadata.LastUpdatedAtUnixSecs,
				SizeBytes:       r.GetKnowledgeBaseURLResponseModel.Metadata.SizeBytes,
			}, nil
		case api.GetKnowledgeBaseFolderResponseModelGetDocumentationFromKnowledgeBaseOK:
			return &KnowledgeBaseDocument{
				ID:              r.GetKnowledgeBaseFolderResponseModel.ID,
				Name:            r.GetKnowledgeBaseFolderResponseModel.Name,
				Type:            "folder",
				CreatedAtUnix:   r.GetKnowledgeBaseFolderResponseModel.Metadata.CreatedAtUnixSecs,
				LastUpdatedUnix: r.GetKnowledgeBaseFolderResponseModel.Metadata.LastUpdatedAtUnixSecs,
				SizeBytes:       r.GetKnowledgeBaseFolderResponseModel.Metadata.SizeBytes,
			}, nil
		}
	}
	return nil, errors.New("unexpected response type")
}

// documentItemToDocument converts list item to KnowledgeBaseDocument.
func documentItemToDocument(item *api.GetKnowledgeBaseListResponseModelDocumentsItem) *KnowledgeBaseDocument {
	doc := &KnowledgeBaseDocument{}

	switch item.Type {
	case api.GetKnowledgeBaseSummaryFileResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		doc.ID = item.GetKnowledgeBaseSummaryFileResponseModel.ID
		doc.Name = item.GetKnowledgeBaseSummaryFileResponseModel.Name
		doc.Type = "file"
		doc.CreatedAtUnix = item.GetKnowledgeBaseSummaryFileResponseModel.Metadata.CreatedAtUnixSecs
		doc.LastUpdatedUnix = item.GetKnowledgeBaseSummaryFileResponseModel.Metadata.LastUpdatedAtUnixSecs
		doc.SizeBytes = item.GetKnowledgeBaseSummaryFileResponseModel.Metadata.SizeBytes
	case api.GetKnowledgeBaseSummaryTextResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		doc.ID = item.GetKnowledgeBaseSummaryTextResponseModel.ID
		doc.Name = item.GetKnowledgeBaseSummaryTextResponseModel.Name
		doc.Type = "text"
		doc.CreatedAtUnix = item.GetKnowledgeBaseSummaryTextResponseModel.Metadata.CreatedAtUnixSecs
		doc.LastUpdatedUnix = item.GetKnowledgeBaseSummaryTextResponseModel.Metadata.LastUpdatedAtUnixSecs
		doc.SizeBytes = item.GetKnowledgeBaseSummaryTextResponseModel.Metadata.SizeBytes
	case api.GetKnowledgeBaseSummaryURLResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		doc.ID = item.GetKnowledgeBaseSummaryURLResponseModel.ID
		doc.Name = item.GetKnowledgeBaseSummaryURLResponseModel.Name
		doc.Type = "url"
		doc.SourceURL = item.GetKnowledgeBaseSummaryURLResponseModel.URL
		doc.CreatedAtUnix = item.GetKnowledgeBaseSummaryURLResponseModel.Metadata.CreatedAtUnixSecs
		doc.LastUpdatedUnix = item.GetKnowledgeBaseSummaryURLResponseModel.Metadata.LastUpdatedAtUnixSecs
		doc.SizeBytes = item.GetKnowledgeBaseSummaryURLResponseModel.Metadata.SizeBytes
	case api.GetKnowledgeBaseSummaryFolderResponseModelGetKnowledgeBaseListResponseModelDocumentsItem:
		doc.ID = item.GetKnowledgeBaseSummaryFolderResponseModel.ID
		doc.Name = item.GetKnowledgeBaseSummaryFolderResponseModel.Name
		doc.Type = "folder"
		doc.CreatedAtUnix = item.GetKnowledgeBaseSummaryFolderResponseModel.Metadata.CreatedAtUnixSecs
		doc.LastUpdatedUnix = item.GetKnowledgeBaseSummaryFolderResponseModel.Metadata.LastUpdatedAtUnixSecs
		doc.SizeBytes = item.GetKnowledgeBaseSummaryFolderResponseModel.Metadata.SizeBytes
	}

	return doc
}
