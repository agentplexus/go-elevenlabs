package content

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	ht "github.com/ogen-go/ogen/http"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Service handles content generation including music, dubbing, and projects.
type Service struct {
	apiClient *api.Client
}

// New creates a new content service.
func New(apiClient *api.Client) *Service {
	return &Service{apiClient: apiClient}
}

// --- Music ---

// ErrMusicCompositionUnavailable is returned when music composition endpoints
// are called but the underlying API client doesn't support them.
var ErrMusicCompositionUnavailable = errors.New("content: music composition endpoints are temporarily unavailable due to API schema complexity; use SeparateStems or VideoToMusic instead")

// MusicRequest contains options for music generation.
type MusicRequest struct {
	// Prompt is a simple text description of the music to generate.
	Prompt string

	// DurationMs is the length of the song in milliseconds (3000-600000).
	DurationMs int

	// ForceInstrumental ensures the song has no vocals.
	ForceInstrumental bool

	// Seed for deterministic generation (optional).
	Seed int
}

// GenerateMusic creates music from a text prompt.
//
// NOTE: This method is temporarily unavailable due to API schema complexity.
func (s *Service) GenerateMusic(ctx context.Context, req *MusicRequest) (io.Reader, error) {
	return nil, ErrMusicCompositionUnavailable
}

// StemSeparationRequest contains options for stem separation.
type StemSeparationRequest struct {
	// File is the audio file to separate.
	File io.Reader

	// Filename is the name of the file.
	Filename string

	// StemVariation specifies which stem variation to use.
	// Options: "two_stems_v1" (vocals + music), "six_stems_v1" (vocals, drums, bass, other - default)
	StemVariation string
}

// SeparateStems separates a song into individual stems (vocals, instruments, etc.).
//
// Example:
//
//	f, _ := os.Open("song.mp3")
//	stems, err := client.Content().SeparateStems(ctx, &content.StemSeparationRequest{
//	    File:     f,
//	    Filename: "song.mp3",
//	})
//	// Save the separated stems (returned as a zip file)
//	output, _ := os.Create("stems.zip")
//	io.Copy(output, stems)
func (s *Service) SeparateStems(ctx context.Context, req *StemSeparationRequest) (io.Reader, error) {
	if req.File == nil {
		return nil, fmt.Errorf("content: file cannot be nil")
	}
	if req.Filename == "" {
		return nil, fmt.Errorf("content: filename cannot be empty")
	}

	body := &api.BodyStemSeparationV1MusicStemSeparationPostMultipart{
		File: ht.MultipartFile{
			Name: req.Filename,
			File: req.File,
		},
	}

	if req.StemVariation != "" {
		body.StemVariationID = api.NewOptBodyStemSeparationV1MusicStemSeparationPostMultipartStemVariationID(
			api.BodyStemSeparationV1MusicStemSeparationPostMultipartStemVariationID(req.StemVariation))
	}

	resp, err := s.apiClient.SeparateSongStems(ctx, body, api.SeparateSongStemsParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.SeparateSongStemsOKHeaders:
		return r.Response.Data, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// VideoToMusicRequest contains options for converting a video to music.
type VideoToMusicRequest struct {
	// Videos is one or more video files to convert.
	Videos []VideoFile

	// Description is an optional text description of the music.
	Description string

	// Tags are optional style tags (e.g., ["upbeat", "cinematic"]).
	Tags []string

	// ModelID specifies the model to use (optional).
	ModelID string
}

// VideoFile represents a video file for VideoToMusic.
type VideoFile struct {
	// Name is the filename.
	Name string
	// File is the video content.
	File io.Reader
}

// VideoToMusic converts a video to music.
func (s *Service) VideoToMusic(ctx context.Context, req *VideoToMusicRequest) (io.Reader, error) {
	if len(req.Videos) == 0 {
		return nil, fmt.Errorf("content: at least one video is required")
	}

	var videos []ht.MultipartFile
	for _, v := range req.Videos {
		if v.File == nil {
			return nil, fmt.Errorf("content: video file cannot be nil")
		}
		if v.Name == "" {
			return nil, fmt.Errorf("content: video filename cannot be empty")
		}
		videos = append(videos, ht.MultipartFile{
			Name: v.Name,
			File: v.File,
		})
	}

	body := &api.BodyVideoToMusicV1MusicVideoToMusicPostMultipart{
		Videos: videos,
		Tags:   req.Tags,
	}

	if req.Description != "" {
		body.Description = api.NewOptNilString(req.Description)
	}

	if req.ModelID != "" {
		body.ModelID = api.NewOptBodyVideoToMusicV1MusicVideoToMusicPostMultipartModelID(
			api.BodyVideoToMusicV1MusicVideoToMusicPostMultipartModelID(req.ModelID))
	}

	resp, err := s.apiClient.VideoToMusic(ctx, body, api.VideoToMusicParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.VideoToMusicOKHeaders:
		return r.Response.Data, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// --- Dubbing ---

// DubbingProject represents a dubbing project.
type DubbingProject struct {
	// DubbingID is the unique identifier.
	DubbingID string

	// Name is the project name.
	Name string

	// Status is the current status (dubbed, dubbing, failed, cloning).
	Status string

	// TargetLanguages are the target languages for dubbing.
	TargetLanguages []string

	// SourceLanguage is the source language.
	SourceLanguage string

	// Error contains any error message if the project failed.
	Error string

	// CreatedAt is when the project was created.
	CreatedAt time.Time
}

// IsComplete checks if a dubbing project is complete.
func (p *DubbingProject) IsComplete() bool {
	return p.Status == "dubbed"
}

// IsFailed checks if a dubbing project has failed.
func (p *DubbingProject) IsFailed() bool {
	return p.Status == "failed"
}

// IsProcessing checks if a dubbing project is still processing.
func (p *DubbingProject) IsProcessing() bool {
	return p.Status == "dubbing" || p.Status == "cloning"
}

// DubbingResponse contains the result of creating a dubbing project.
type DubbingResponse struct {
	// DubbingID is the ID of the created project.
	DubbingID string

	// ExpectedDurationSeconds is the expected duration.
	ExpectedDurationSeconds float64
}

// DubbingRequest contains options for creating a dubbing project.
type DubbingRequest struct {
	// Name is the name of the dubbing project.
	Name string

	// SourceURL is the URL of the source media (alternative to file upload).
	SourceURL string

	// File is the source media file (alternative to SourceURL).
	File io.Reader

	// SourceLanguage is the source language code (ISO 639-1).
	SourceLanguage string

	// TargetLanguage is the target language code (ISO 639-1).
	TargetLanguage string

	// NumSpeakers is the number of speakers (0 for auto-detection).
	NumSpeakers int

	// Watermark enables watermark (for free tier).
	Watermark bool

	// StartTime is the start time in seconds for dubbing.
	StartTime int

	// EndTime is the end time in seconds for dubbing.
	EndTime int

	// HighestResolution requests highest resolution output.
	HighestResolution bool

	// DropBackgroundAudio removes background audio.
	DropBackgroundAudio bool
}

// CreateDubbing creates a dubbing project from a URL source.
func (s *Service) CreateDubbing(ctx context.Context, req *DubbingRequest) (*DubbingResponse, error) {
	if req.SourceURL == "" {
		return nil, fmt.Errorf("content: source_url cannot be empty")
	}
	if req.TargetLanguage == "" {
		return nil, fmt.Errorf("content: target_language cannot be empty")
	}

	// Build request body
	body := api.BodyDubAVideoOrAnAudioFileV1DubbingPostMultipart{}
	body.SourceURL = api.NewOptNilString(req.SourceURL)
	body.TargetLang = api.NewOptNilString(req.TargetLanguage)

	if req.Name != "" {
		body.Name = api.NewOptNilString(req.Name)
	}
	if req.SourceLanguage != "" {
		body.SourceLang = api.NewOptString(req.SourceLanguage)
	}
	if req.NumSpeakers != 0 {
		body.NumSpeakers = api.NewOptInt(req.NumSpeakers)
	}
	if req.Watermark {
		body.Watermark = api.NewOptBool(true)
	}
	if req.StartTime > 0 {
		body.StartTime = api.NewOptNilInt(req.StartTime)
	}
	if req.EndTime > 0 {
		body.EndTime = api.NewOptNilInt(req.EndTime)
	}
	if req.HighestResolution {
		body.HighestResolution = api.NewOptBool(true)
	}
	if req.DropBackgroundAudio {
		body.DropBackgroundAudio = api.NewOptBool(true)
	}

	resp, err := s.apiClient.CreateDubbing(ctx, api.NewOptBodyDubAVideoOrAnAudioFileV1DubbingPostMultipart(body), api.CreateDubbingParams{})
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.DoDubbingResponseModel:
		return &DubbingResponse{
			DubbingID:               r.DubbingID,
			ExpectedDurationSeconds: r.ExpectedDurationSec,
		}, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// GetDubbing returns a dubbing project metadata by ID.
func (s *Service) GetDubbing(ctx context.Context, dubbingID string) (*DubbingProject, error) {
	if dubbingID == "" {
		return nil, fmt.Errorf("content: dubbing_id cannot be empty")
	}

	resp, err := s.apiClient.GetDubbedMetadata(ctx, api.GetDubbedMetadataParams{
		DubbingID: dubbingID,
	})
	if err != nil {
		return nil, err
	}

	// Handle response type
	switch r := resp.(type) {
	case *api.DubbingMetadataResponse:
		project := &DubbingProject{
			DubbingID:       r.DubbingID,
			Name:            r.Name,
			Status:          r.Status,
			TargetLanguages: r.TargetLanguages,
			CreatedAt:       r.CreatedAt,
		}

		if r.Error.Set && !r.Error.Null {
			project.Error = r.Error.Value
		}

		return project, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// DeleteDubbing deletes a dubbing project by ID.
func (s *Service) DeleteDubbing(ctx context.Context, dubbingID string) error {
	if dubbingID == "" {
		return fmt.Errorf("content: dubbing_id cannot be empty")
	}

	_, err := s.apiClient.DeleteDubbing(ctx, api.DeleteDubbingParams{
		DubbingID: dubbingID,
	})
	return err
}

// GetDubbedFile returns the dubbed audio/video file for a specific language.
func (s *Service) GetDubbedFile(ctx context.Context, dubbingID, languageCode string) (io.Reader, error) {
	if dubbingID == "" {
		return nil, fmt.Errorf("content: dubbing_id cannot be empty")
	}
	if languageCode == "" {
		return nil, fmt.Errorf("content: language_code cannot be empty")
	}

	resp, err := s.apiClient.GetDubbedFile(ctx, api.GetDubbedFileParams{
		DubbingID:    dubbingID,
		LanguageCode: languageCode,
	})
	if err != nil {
		return nil, err
	}

	// Handle response type - can be audio or video
	switch r := resp.(type) {
	case *api.GetDubbedFileOKAudioMpeg:
		return r.Data, nil
	case *api.GetDubbedFileOKVideoMP4:
		return r.Data, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// --- Projects (Studio) ---

// Project represents a Studio project.
type Project struct {
	// ProjectID is the unique identifier.
	ProjectID string

	// Name is the project name.
	Name string

	// Description is the project description.
	Description string

	// Author is the project author.
	Author string

	// Language is the two-letter language code (ISO 639-1).
	Language string

	// DefaultModelID is the default model for TTS.
	DefaultModelID string

	// DefaultParagraphVoiceID is the default voice for paragraphs.
	DefaultParagraphVoiceID string

	// DefaultTitleVoiceID is the default voice for titles.
	DefaultTitleVoiceID string

	// ContentType is the content type (e.g., "Novel", "Short Story").
	ContentType string

	// CoverImageURL is the cover image URL.
	CoverImageURL string

	// CreatedAt is the creation timestamp.
	CreatedAt time.Time

	// CanBeDownloaded indicates if the project can be downloaded.
	CanBeDownloaded bool

	// AccessLevel is the access level of the project.
	AccessLevel string
}

// Chapter represents a chapter within a project.
type Chapter struct {
	// ChapterID is the unique identifier.
	ChapterID string

	// Name is the chapter name.
	Name string

	// ConversionProgress is the conversion progress percentage.
	ConversionProgress float64

	// State is the current state.
	State string

	// LastConversionError is the last conversion error if any.
	LastConversionError string
}

// CreateProjectRequest contains options for creating a project.
type CreateProjectRequest struct {
	// Name is the project name (required).
	Name string

	// Description is an optional description.
	Description string

	// Author is an optional author name.
	Author string

	// Language is the two-letter language code (ISO 639-1).
	Language string

	// DefaultModelID is the model to use for TTS.
	DefaultModelID string

	// DefaultParagraphVoiceID is the default voice for paragraphs.
	DefaultParagraphVoiceID string

	// DefaultTitleVoiceID is the default voice for titles.
	DefaultTitleVoiceID string

	// FromURL is a URL to extract content from.
	FromURL string

	// ContentType is the content type (e.g., "Novel", "Short Story").
	ContentType string

	// Genres is a list of genres.
	Genres []string

	// QualityPreset is the output quality: "standard", "high", "ultra", "ultra lossless".
	QualityPreset string

	// AutoConvert automatically converts the project to audio.
	AutoConvert bool
}

// ListProjects returns all projects.
func (s *Service) ListProjects(ctx context.Context) ([]*Project, error) {
	resp, err := s.apiClient.GetProjects(ctx, api.GetProjectsParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetProjectsResponseModel:
		projects := make([]*Project, 0, len(r.Projects))
		for _, p := range r.Projects {
			proj := projectFromAPI(&p)
			projects = append(projects, proj)
		}
		return projects, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// CreateProject creates a new project.
func (s *Service) CreateProject(ctx context.Context, req *CreateProjectRequest) (*Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("content: name cannot be empty")
	}

	body := &api.BodyCreateStudioProjectV1StudioProjectsPostMultipart{
		Name: req.Name,
	}

	if req.Description != "" {
		body.Description = api.NewOptNilString(req.Description)
	}
	if req.Author != "" {
		body.Author = api.NewOptNilString(req.Author)
	}
	if req.Language != "" {
		body.Language = api.NewOptNilString(req.Language)
	}
	if req.DefaultModelID != "" {
		body.DefaultModelID = api.NewOptNilString(req.DefaultModelID)
	}
	if req.DefaultParagraphVoiceID != "" {
		body.DefaultParagraphVoiceID = api.NewOptNilString(req.DefaultParagraphVoiceID)
	}
	if req.DefaultTitleVoiceID != "" {
		body.DefaultTitleVoiceID = api.NewOptNilString(req.DefaultTitleVoiceID)
	}
	if req.FromURL != "" {
		body.FromURL = api.NewOptNilString(req.FromURL)
	}
	if req.ContentType != "" {
		body.ContentType = api.NewOptNilString(req.ContentType)
	}
	if len(req.Genres) > 0 {
		body.Genres = req.Genres
	}
	if req.QualityPreset != "" {
		body.QualityPreset = api.NewOptQualityPresetType(api.QualityPresetType(req.QualityPreset))
	}
	if req.AutoConvert {
		body.AutoConvert = api.NewOptBool(true)
	}

	resp, err := s.apiClient.AddProject(ctx, body, api.AddProjectParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AddProjectResponseModel:
		return projectFromAPI(&r.Project), nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// DeleteProject deletes a project.
func (s *Service) DeleteProject(ctx context.Context, projectID string) error {
	if projectID == "" {
		return fmt.Errorf("content: project_id cannot be empty")
	}

	_, err := s.apiClient.DeleteProject(ctx, api.DeleteProjectParams{
		ProjectID: projectID,
	})
	return err
}

// ConvertProject initiates conversion of a project to audio.
func (s *Service) ConvertProject(ctx context.Context, projectID string) error {
	if projectID == "" {
		return fmt.Errorf("content: project_id cannot be empty")
	}

	_, err := s.apiClient.ConvertProjectEndpoint(ctx, api.ConvertProjectEndpointParams{
		ProjectID: projectID,
	})
	return err
}

// ListChapters returns all chapters in a project.
func (s *Service) ListChapters(ctx context.Context, projectID string) ([]*Chapter, error) {
	if projectID == "" {
		return nil, fmt.Errorf("content: project_id cannot be empty")
	}

	resp, err := s.apiClient.GetChapters(ctx, api.GetChaptersParams{
		ProjectID: projectID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetChaptersResponseModel:
		chapters := make([]*Chapter, 0, len(r.Chapters))
		for _, c := range r.Chapters {
			ch := &Chapter{
				ChapterID: c.ChapterID,
				Name:      c.Name,
				State:     string(c.State),
			}
			if c.ConversionProgress.Set && !c.ConversionProgress.Null {
				ch.ConversionProgress = c.ConversionProgress.Value
			}
			if c.LastConversionError.Set && !c.LastConversionError.Null {
				ch.LastConversionError = c.LastConversionError.Value
			}
			chapters = append(chapters, ch)
		}
		return chapters, nil
	default:
		return nil, fmt.Errorf("content: unexpected response type")
	}
}

// projectFromAPI converts an API ProjectResponseModel to our Project type.
func projectFromAPI(p *api.ProjectResponseModel) *Project {
	proj := &Project{
		ProjectID:               p.ProjectID,
		Name:                    p.Name,
		DefaultModelID:          p.DefaultModelID,
		DefaultParagraphVoiceID: p.DefaultParagraphVoiceID,
		DefaultTitleVoiceID:     p.DefaultTitleVoiceID,
		CreatedAt:               time.Unix(int64(p.CreateDateUnix), 0),
		CanBeDownloaded:         p.CanBeDownloaded,
		AccessLevel:             string(p.AccessLevel),
	}

	if p.Description.Set && !p.Description.Null {
		proj.Description = p.Description.Value
	}
	if p.Author.Set && !p.Author.Null {
		proj.Author = p.Author.Value
	}
	if p.Language.Set && !p.Language.Null {
		proj.Language = p.Language.Value
	}
	if p.ContentType.Set && !p.ContentType.Null {
		proj.ContentType = p.ContentType.Value
	}
	if p.CoverImageURL.Set && !p.CoverImageURL.Null {
		proj.CoverImageURL = p.CoverImageURL.Value
	}

	return proj
}
