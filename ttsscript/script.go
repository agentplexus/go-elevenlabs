// Package ttsscript provides a structured format for authoring multilingual
// TTS scripts that can be compiled to various output formats (SSML, ElevenLabs, etc.).
//
// This package is engine-agnostic and can be used with any TTS provider.
package ttsscript

import (
	"encoding/json"
	"fmt"
	"os"
)

// Script represents a multilingual TTS script with slides/segments.
// This is the canonical format for authoring TTS content that can be
// compiled to SSML (Google TTS, Amazon Polly) or ElevenLabs-compatible text.
type Script struct {
	// Title is the script title.
	Title string `json:"title,omitempty"`

	// Description is an optional description.
	Description string `json:"description,omitempty"`

	// DefaultLanguage is the primary language code (e.g., "en-US").
	DefaultLanguage string `json:"default_language,omitempty"`

	// DefaultVoices maps language codes to default voice IDs.
	DefaultVoices map[string]string `json:"default_voices,omitempty"`

	// Pronunciations maps terms to their pronunciation by language.
	// Example: {"ADK": {"en": "A D K", "es": "A D K"}}
	Pronunciations map[string]map[string]string `json:"pronunciations,omitempty"`

	// Slides contains the ordered list of slides/sections.
	Slides []Slide `json:"slides"`
}

// Slide represents a slide or section of the script.
type Slide struct {
	// Title is the slide title (optional).
	Title string `json:"title,omitempty"`

	// Notes are speaker notes or comments (not rendered to audio).
	Notes string `json:"notes,omitempty"`

	// IsSectionHeader marks this slide as the start of a new section.
	// Section headers can have their titles spoken and use longer transition pauses.
	IsSectionHeader bool `json:"is_section_header,omitempty"`

	// SpeakTitle causes the slide title to be spoken before the segments.
	// If true, the title is converted to a segment. Defaults to true for section headers.
	SpeakTitle *bool `json:"speak_title,omitempty"`

	// TitleVoice overrides the voice used for speaking the title, by language.
	// If not set, uses the segment voice or default voice.
	TitleVoice map[string]string `json:"title_voice,omitempty"`

	// TitlePauseAfter is the pause after the spoken title (e.g., "500ms").
	// Defaults to "500ms" for section headers, "300ms" for regular slides.
	TitlePauseAfter string `json:"title_pause_after,omitempty"`

	// Segments are the audio segments for this slide.
	Segments []Segment `json:"segments"`
}

// Segment represents a single audio segment within a slide.
type Segment struct {
	// Text contains the text content by language code.
	// Example: {"en": "Hello world", "es": "Hola mundo"}
	Text map[string]string `json:"text"`

	// Voice overrides the default voice for this segment by language.
	// Example: {"en": "voice-id-1", "es": "voice-id-2"}
	Voice map[string]string `json:"voice,omitempty"`

	// PauseBefore is the pause duration before this segment (e.g., "500ms", "1s").
	PauseBefore string `json:"pause_before,omitempty"`

	// PauseAfter is the pause duration after this segment (e.g., "500ms", "1s").
	PauseAfter string `json:"pause_after,omitempty"`

	// Emphasis indicates the emphasis level ("strong", "moderate", "reduced").
	Emphasis string `json:"emphasis,omitempty"`

	// Rate is the speaking rate ("slow", "medium", "fast", or percentage like "80%").
	Rate string `json:"rate,omitempty"`

	// Pitch adjusts the pitch ("low", "medium", "high", or percentage like "+10%").
	Pitch string `json:"pitch,omitempty"`

	// Pronunciations are segment-specific pronunciation overrides.
	Pronunciations map[string]map[string]string `json:"pronunciations,omitempty"`
}

// LoadScript loads a script from a JSON file.
func LoadScript(filePath string) (*Script, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading script file: %w", err)
	}
	return ParseScript(data)
}

// ParseScript parses a script from JSON data.
func ParseScript(data []byte) (*Script, error) {
	var script Script
	if err := json.Unmarshal(data, &script); err != nil {
		return nil, fmt.Errorf("parsing script JSON: %w", err)
	}
	return &script, nil
}

// Save saves a script to a JSON file.
func (s *Script) Save(filePath string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling script: %w", err)
	}
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return fmt.Errorf("writing script file: %w", err)
	}
	return nil
}

// Languages returns all language codes used in the script.
func (s *Script) Languages() []string {
	langs := make(map[string]bool)
	for _, slide := range s.Slides {
		for _, seg := range slide.Segments {
			for lang := range seg.Text {
				langs[lang] = true
			}
		}
	}
	result := make([]string, 0, len(langs))
	for lang := range langs {
		result = append(result, lang)
	}
	return result
}

// SlideCount returns the number of slides.
func (s *Script) SlideCount() int {
	return len(s.Slides)
}

// SegmentCount returns the total number of segments across all slides.
func (s *Script) SegmentCount() int {
	count := 0
	for _, slide := range s.Slides {
		count += len(slide.Segments)
	}
	return count
}

// ShouldSpeakTitle returns true if the slide title should be spoken.
// Returns true if SpeakTitle is explicitly true, or if the slide is a
// section header and SpeakTitle is not explicitly false.
func (s *Slide) ShouldSpeakTitle() bool {
	if s.SpeakTitle != nil {
		return *s.SpeakTitle
	}
	// Default: speak title for section headers only
	return s.IsSectionHeader
}

// Validate checks the script for common issues.
func (s *Script) Validate() []string {
	var issues []string

	if len(s.Slides) == 0 {
		issues = append(issues, "script has no slides")
	}

	for i, slide := range s.Slides {
		if len(slide.Segments) == 0 {
			issues = append(issues, fmt.Sprintf("slide %d has no segments", i+1))
		}
		for j, seg := range slide.Segments {
			if len(seg.Text) == 0 {
				issues = append(issues, fmt.Sprintf("slide %d, segment %d has no text", i+1, j+1))
			}
		}
	}

	return issues
}
