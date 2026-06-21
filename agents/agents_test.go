// Package agents tests for ElevenAgents (Conversational AI) services.
//
// Run unit tests:
//
//	go test -v ./agents/...
//
// Run integration tests (requires ELEVENLABS_API_KEY):
//
//	go test -v -tags=integration ./agents/...
package agents

import (
	"testing"
)

func TestNewService(t *testing.T) {
	// Service requires a client, so we just test that nil is handled
	svc := New(nil)
	if svc == nil {
		t.Fatal("New() returned nil")
	}
}

func TestCreateBatchCallRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateBatchCallRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "missing agent ID",
			req:     &CreateBatchCallRequest{Name: "Test"},
			wantErr: "agent ID is required",
		},
		{
			name:    "missing name",
			req:     &CreateBatchCallRequest{AgentID: "agent-123"},
			wantErr: "name is required",
		},
		{
			name: "missing recipients",
			req: &CreateBatchCallRequest{
				AgentID: "agent-123",
				Name:    "Test",
			},
			wantErr: "at least one recipient is required",
		},
	}

	svc := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateBatchCall(nil, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateFileDocumentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateFileDocumentRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "missing file",
			req:     &CreateFileDocumentRequest{Filename: "test.pdf"},
			wantErr: "file is required",
		},
		{
			name:    "missing filename",
			req:     &CreateFileDocumentRequest{},
			wantErr: "file is required",
		},
	}

	svc := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateFileDocument(nil, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateTextDocumentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateTextDocumentRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "missing content",
			req:     &CreateTextDocumentRequest{Name: "Test"},
			wantErr: "content is required",
		},
	}

	svc := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateTextDocument(nil, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateURLDocumentRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateURLDocumentRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "missing URL",
			req:     &CreateURLDocumentRequest{Name: "Test"},
			wantErr: "url is required",
		},
	}

	svc := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateURLDocument(nil, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestCreateKnowledgeBaseFolderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateKnowledgeBaseFolderRequest
		wantErr string
	}{
		{
			name:    "nil request",
			req:     nil,
			wantErr: "request is required",
		},
		{
			name:    "missing name",
			req:     &CreateKnowledgeBaseFolderRequest{},
			wantErr: "name is required",
		},
	}

	svc := New(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.CreateKnowledgeBaseFolder(nil, tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("error = %q, want %q", err.Error(), tt.wantErr)
			}
		})
	}
}

func TestBatchCallRecipient(t *testing.T) {
	r := BatchCallRecipient{
		PhoneNumber: "+14155551234",
		ID:          "customer-001",
		DynamicVariables: map[string]any{
			"name": "John",
		},
	}

	if r.PhoneNumber != "+14155551234" {
		t.Errorf("PhoneNumber = %q, want +14155551234", r.PhoneNumber)
	}
	if r.ID != "customer-001" {
		t.Errorf("ID = %q, want customer-001", r.ID)
	}
	if r.DynamicVariables["name"] != "John" {
		t.Errorf("DynamicVariables[name] = %v, want John", r.DynamicVariables["name"])
	}
}

func TestDocumentTypes(t *testing.T) {
	// Test Document struct
	doc := &Document{
		ID:   "doc-123",
		Name: "Test Document",
		Type: "file",
	}

	if doc.ID != "doc-123" {
		t.Errorf("ID = %q, want doc-123", doc.ID)
	}
	if doc.Type != "file" {
		t.Errorf("Type = %q, want file", doc.Type)
	}

	// Test DocumentSummary struct
	summary := &DocumentSummary{
		ID:   "doc-456",
		Name: "Summary Doc",
		Type: "text",
	}

	if summary.ID != "doc-456" {
		t.Errorf("ID = %q, want doc-456", summary.ID)
	}

	// Test DocumentChunk struct
	chunk := &DocumentChunk{
		ID:      "chunk-1",
		Name:    "Chunk 1",
		Content: "Some content",
	}

	if chunk.ID != "chunk-1" {
		t.Errorf("ID = %q, want chunk-1", chunk.ID)
	}
}

func TestBatchCallTypes(t *testing.T) {
	// Test BatchCall struct
	bc := &BatchCall{
		ID:                   "batch-123",
		Name:                 "Test Campaign",
		AgentID:              "agent-456",
		Status:               "running",
		TotalCallsScheduled:  100,
		TotalCallsDispatched: 50,
		TotalCallsFinished:   25,
	}

	if bc.ID != "batch-123" {
		t.Errorf("ID = %q, want batch-123", bc.ID)
	}
	if bc.Status != "running" {
		t.Errorf("Status = %q, want running", bc.Status)
	}
	if bc.TotalCallsScheduled != 100 {
		t.Errorf("TotalCallsScheduled = %d, want 100", bc.TotalCallsScheduled)
	}

	// Test BatchCallDetail struct
	detail := &BatchCallDetail{
		BatchCall: *bc,
		Recipients: []*BatchCallRecipientStatus{
			{
				ID:          "r-1",
				PhoneNumber: "+14155551234",
				Status:      "completed",
			},
		},
	}

	if len(detail.Recipients) != 1 {
		t.Errorf("len(Recipients) = %d, want 1", len(detail.Recipients))
	}
	if detail.Recipients[0].Status != "completed" {
		t.Errorf("Recipients[0].Status = %q, want completed", detail.Recipients[0].Status)
	}
}
