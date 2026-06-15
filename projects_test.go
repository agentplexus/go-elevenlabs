package elevenlabs

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestCreateProjectRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateProjectRequest
		wantErr bool
	}{
		{
			name:    "empty name",
			req:     &CreateProjectRequest{Name: ""},
			wantErr: true,
		},
		{
			name:    "valid request",
			req:     &CreateProjectRequest{Name: "My Project"},
			wantErr: false,
		},
		{
			name: "valid request with options",
			req: &CreateProjectRequest{
				Name:        "My Course",
				Description: "A comprehensive course",
				Author:      "John Doe",
				Language:    "en",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectsService(t *testing.T) {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		t.Skip("ELEVENLABS_API_KEY not set, skipping live test")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.Background()

	// Test that service is accessible
	if client.Projects() == nil {
		t.Error("Projects() returned nil")
	}

	t.Run("List", func(t *testing.T) {
		projects, err := client.Projects().List(ctx)
		if err != nil {
			// Projects API requires specific subscription tier
			if strings.Contains(err.Error(), "403") {
				t.Skip("Projects API requires Studio subscription tier")
			}
			t.Errorf("List() error = %v", err)
			return
		}
		// Projects may be empty, that's OK
		t.Logf("Found %d projects", len(projects))
	})

	t.Run("Update with empty project ID", func(t *testing.T) {
		err := client.Projects().Update(ctx, "", &UpdateProjectRequest{
			Name:                    "Test",
			DefaultParagraphVoiceID: "voice1",
			DefaultTitleVoiceID:     "voice2",
		})
		if err == nil {
			t.Error("Update() with empty project ID should return error")
		}
		if _, ok := err.(*ValidationError); !ok {
			t.Errorf("Update() with empty project ID should return ValidationError, got %T", err)
		}
	})

	t.Run("Update with empty name", func(t *testing.T) {
		err := client.Projects().Update(ctx, "test-id", &UpdateProjectRequest{
			Name:                    "",
			DefaultParagraphVoiceID: "voice1",
			DefaultTitleVoiceID:     "voice2",
		})
		if err == nil {
			t.Error("Update() with empty name should return error")
		}
	})

	t.Run("Delete with empty ID", func(t *testing.T) {
		err := client.Projects().Delete(ctx, "")
		if err == nil {
			t.Error("Delete() with empty ID should return error")
		}
	})

	t.Run("Convert with empty ID", func(t *testing.T) {
		err := client.Projects().Convert(ctx, "")
		if err == nil {
			t.Error("Convert() with empty ID should return error")
		}
	})

	t.Run("ListChapters with empty ID", func(t *testing.T) {
		_, err := client.Projects().ListChapters(ctx, "")
		if err == nil {
			t.Error("ListChapters() with empty ID should return error")
		}
	})

	t.Run("ConvertChapter with empty project ID", func(t *testing.T) {
		err := client.Projects().ConvertChapter(ctx, "", "chapter-id")
		if err == nil {
			t.Error("ConvertChapter() with empty project ID should return error")
		}
	})

	t.Run("ConvertChapter with empty chapter ID", func(t *testing.T) {
		err := client.Projects().ConvertChapter(ctx, "project-id", "")
		if err == nil {
			t.Error("ConvertChapter() with empty chapter ID should return error")
		}
	})

	t.Run("DeleteChapter with empty IDs", func(t *testing.T) {
		err := client.Projects().DeleteChapter(ctx, "", "chapter-id")
		if err == nil {
			t.Error("DeleteChapter() with empty project ID should return error")
		}

		err = client.Projects().DeleteChapter(ctx, "project-id", "")
		if err == nil {
			t.Error("DeleteChapter() with empty chapter ID should return error")
		}
	})

	t.Run("ListSnapshots with empty ID", func(t *testing.T) {
		_, err := client.Projects().ListSnapshots(ctx, "")
		if err == nil {
			t.Error("ListSnapshots() with empty ID should return error")
		}
	})

	t.Run("DownloadSnapshotArchive with empty IDs", func(t *testing.T) {
		_, err := client.Projects().DownloadSnapshotArchive(ctx, "", "snapshot-id")
		if err == nil {
			t.Error("DownloadSnapshotArchive() with empty project ID should return error")
		}

		_, err = client.Projects().DownloadSnapshotArchive(ctx, "project-id", "")
		if err == nil {
			t.Error("DownloadSnapshotArchive() with empty snapshot ID should return error")
		}
	})

	t.Run("ListChapterSnapshots with empty IDs", func(t *testing.T) {
		_, err := client.Projects().ListChapterSnapshots(ctx, "", "chapter-id")
		if err == nil {
			t.Error("ListChapterSnapshots() with empty project ID should return error")
		}

		_, err = client.Projects().ListChapterSnapshots(ctx, "project-id", "")
		if err == nil {
			t.Error("ListChapterSnapshots() with empty chapter ID should return error")
		}
	})

	t.Run("StreamChapterAudio with empty IDs", func(t *testing.T) {
		_, err := client.Projects().StreamChapterAudio(ctx, "", "chapter", "snapshot")
		if err == nil {
			t.Error("StreamChapterAudio() with empty project ID should return error")
		}

		_, err = client.Projects().StreamChapterAudio(ctx, "project", "", "snapshot")
		if err == nil {
			t.Error("StreamChapterAudio() with empty chapter ID should return error")
		}

		_, err = client.Projects().StreamChapterAudio(ctx, "project", "chapter", "")
		if err == nil {
			t.Error("StreamChapterAudio() with empty snapshot ID should return error")
		}
	})
}
