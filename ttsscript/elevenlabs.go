package ttsscript

import (
	"fmt"
	"strings"
)

// ElevenLabsFormatter formats compiled segments for ElevenLabs TTS.
type ElevenLabsFormatter struct {
	// UsePauseMarkers includes [pause:Xms] markers in text output.
	// When false, pauses are tracked separately for post-processing.
	UsePauseMarkers bool

	// PauseMarkerFormat is the format for pause markers (default: "[pause:%s]").
	PauseMarkerFormat string
}

// NewElevenLabsFormatter creates a new ElevenLabs formatter.
func NewElevenLabsFormatter() *ElevenLabsFormatter {
	return &ElevenLabsFormatter{
		UsePauseMarkers:   false,
		PauseMarkerFormat: "[pause:%s]",
	}
}

// ElevenLabsSegment represents a segment ready for ElevenLabs TTS.
type ElevenLabsSegment struct {
	// Text is the text to generate speech for.
	Text string

	// VoiceID is the ElevenLabs voice ID.
	VoiceID string

	// SlideIndex is the source slide index.
	SlideIndex int

	// SegmentIndex is the source segment index (-1 for title segments).
	SegmentIndex int

	// SlideTitle is the slide title for reference.
	SlideTitle string

	// IsTitleSegment indicates this segment was generated from a slide title.
	IsTitleSegment bool

	// IsSectionHeader indicates this segment belongs to a section header slide.
	IsSectionHeader bool

	// PauseBeforeMs is silence to add before this segment.
	PauseBeforeMs int

	// PauseAfterMs is silence to add after this segment.
	PauseAfterMs int

	// SuggestedFilename is a suggested output filename.
	SuggestedFilename string
}

// Format formats compiled segments for ElevenLabs.
func (f *ElevenLabsFormatter) Format(segments []CompiledSegment) []ElevenLabsSegment {
	result := make([]ElevenLabsSegment, len(segments))

	for i, seg := range segments {
		text := seg.Text

		// Add pause markers if enabled
		if f.UsePauseMarkers {
			if seg.PauseBeforeMs > 0 {
				marker := fmt.Sprintf(f.PauseMarkerFormat, FormatDuration(seg.PauseBeforeMs))
				text = marker + " " + text
			}
			if seg.PauseAfterMs > 0 {
				marker := fmt.Sprintf(f.PauseMarkerFormat, FormatDuration(seg.PauseAfterMs))
				text = text + " " + marker
			}
		}

		// Generate appropriate filename
		var filename string
		if seg.IsTitleSegment {
			filename = fmt.Sprintf("slide%02d_title.mp3", seg.SlideIndex+1)
		} else {
			filename = fmt.Sprintf("slide%02d_seg%02d.mp3", seg.SlideIndex+1, seg.SegmentIndex+1)
		}

		result[i] = ElevenLabsSegment{
			Text:              text,
			VoiceID:           seg.VoiceID,
			SlideIndex:        seg.SlideIndex,
			SegmentIndex:      seg.SegmentIndex,
			SlideTitle:        seg.SlideTitle,
			IsTitleSegment:    seg.IsTitleSegment,
			IsSectionHeader:   seg.IsSectionHeader,
			PauseBeforeMs:     seg.PauseBeforeMs,
			PauseAfterMs:      seg.PauseAfterMs,
			SuggestedFilename: filename,
		}
	}

	return result
}

// FormatScript compiles and formats a script for ElevenLabs.
func (f *ElevenLabsFormatter) FormatScript(script *Script, language string) ([]ElevenLabsSegment, error) {
	compiler := NewCompiler()
	segments, err := compiler.Compile(script, language)
	if err != nil {
		return nil, err
	}
	return f.Format(segments), nil
}

// CombineForSingleRequest combines segments into a single text block.
// Useful when you want to generate all audio in one API call.
// Note: This loses per-segment voice control.
func (f *ElevenLabsFormatter) CombineForSingleRequest(segments []ElevenLabsSegment) string {
	var parts []string
	for _, seg := range segments {
		parts = append(parts, seg.Text)
	}
	return strings.Join(parts, " ")
}

// GroupByVoice groups segments by voice ID for batch processing.
func (f *ElevenLabsFormatter) GroupByVoice(segments []ElevenLabsSegment) map[string][]ElevenLabsSegment {
	groups := make(map[string][]ElevenLabsSegment)
	for _, seg := range segments {
		groups[seg.VoiceID] = append(groups[seg.VoiceID], seg)
	}
	return groups
}

// TTSRequest represents a request to the ElevenLabs TTS API.
// This is a simplified version for use with ttsscript.
type TTSRequest struct {
	VoiceID  string
	Text     string
	ModelID  string
	Segment  ElevenLabsSegment
	Language string
}

// GenerateTTSRequests creates TTS requests from formatted segments.
func GenerateTTSRequests(segments []ElevenLabsSegment, modelID, language string) []TTSRequest {
	requests := make([]TTSRequest, len(segments))
	for i, seg := range segments {
		requests[i] = TTSRequest{
			VoiceID:  seg.VoiceID,
			Text:     seg.Text,
			ModelID:  modelID,
			Segment:  seg,
			Language: language,
		}
	}
	return requests
}

// BatchConfig contains configuration for batch TTS processing.
type BatchConfig struct {
	// OutputDir is the directory for output files.
	OutputDir string

	// FilePrefix is added before each filename.
	FilePrefix string

	// FileSuffix is added after each filename (before extension).
	FileSuffix string

	// IncludeLanguageInFilename adds language code to filename.
	IncludeLanguageInFilename bool
}

// NewBatchConfig creates a batch config with defaults.
func NewBatchConfig(outputDir string) *BatchConfig {
	return &BatchConfig{
		OutputDir:                 outputDir,
		FilePrefix:                "",
		FileSuffix:                "",
		IncludeLanguageInFilename: true,
	}
}

// GenerateFilename generates an output filename for a segment.
func (c *BatchConfig) GenerateFilename(seg ElevenLabsSegment, language string) string {
	var name string
	if seg.IsTitleSegment {
		name = fmt.Sprintf("slide%02d_title", seg.SlideIndex+1)
	} else {
		name = fmt.Sprintf("slide%02d_seg%02d", seg.SlideIndex+1, seg.SegmentIndex+1)
	}

	if c.FilePrefix != "" {
		name = c.FilePrefix + "_" + name
	}

	if c.IncludeLanguageInFilename && language != "" {
		name = name + "_" + language
	}

	if c.FileSuffix != "" {
		name = name + "_" + c.FileSuffix
	}

	return fmt.Sprintf("%s/%s.mp3", c.OutputDir, name)
}

// ManifestEntry represents an entry in a generation manifest.
type ManifestEntry struct {
	SlideIndex      int    `json:"slide_index"`
	SegmentIndex    int    `json:"segment_index"`
	SlideTitle      string `json:"slide_title,omitempty"`
	IsTitleSegment  bool   `json:"is_title_segment,omitempty"`
	IsSectionHeader bool   `json:"is_section_header,omitempty"`
	Text            string `json:"text"`
	VoiceID         string `json:"voice_id"`
	Language        string `json:"language"`
	OutputFile      string `json:"output_file"`
	PauseBeforeMs   int    `json:"pause_before_ms,omitempty"`
	PauseAfterMs    int    `json:"pause_after_ms,omitempty"`
}

// GenerateManifest creates a manifest of all segments for tracking.
func GenerateManifest(segments []ElevenLabsSegment, config *BatchConfig, language string) []ManifestEntry {
	entries := make([]ManifestEntry, len(segments))
	for i, seg := range segments {
		entries[i] = ManifestEntry{
			SlideIndex:      seg.SlideIndex,
			SegmentIndex:    seg.SegmentIndex,
			SlideTitle:      seg.SlideTitle,
			IsTitleSegment:  seg.IsTitleSegment,
			IsSectionHeader: seg.IsSectionHeader,
			Text:            seg.Text,
			VoiceID:         seg.VoiceID,
			Language:        language,
			OutputFile:      config.GenerateFilename(seg, language),
			PauseBeforeMs:   seg.PauseBeforeMs,
			PauseAfterMs:    seg.PauseAfterMs,
		}
	}
	return entries
}
