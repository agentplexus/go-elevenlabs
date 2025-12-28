package ttsscript

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Compiler compiles scripts to various output formats.
type Compiler struct {
	// AdditionalPronunciations are extra pronunciations to apply.
	AdditionalPronunciations map[string]map[string]string

	// DefaultPauseAfterSlide is the pause after each slide if not specified.
	DefaultPauseAfterSlide string

	// DefaultPauseAfterSegment is the pause after each segment if not specified.
	DefaultPauseAfterSegment string
}

// NewCompiler creates a new script compiler with default settings.
func NewCompiler() *Compiler {
	return &Compiler{
		AdditionalPronunciations: make(map[string]map[string]string),
		DefaultPauseAfterSlide:   "800ms",
		DefaultPauseAfterSegment: "",
	}
}

// CompiledSegment represents a compiled segment ready for TTS.
type CompiledSegment struct {
	// SlideIndex is the 0-based slide index.
	SlideIndex int

	// SegmentIndex is the 0-based segment index within the slide.
	// For title segments, this is -1.
	SegmentIndex int

	// SlideTitle is the slide title (if any).
	SlideTitle string

	// IsTitleSegment indicates this segment was generated from a slide title.
	IsTitleSegment bool

	// IsSectionHeader indicates this segment belongs to a section header slide.
	IsSectionHeader bool

	// Text is the processed text with pronunciations applied.
	Text string

	// OriginalText is the text before pronunciation substitutions.
	OriginalText string

	// VoiceID is the voice to use for this segment.
	VoiceID string

	// Language is the language code.
	Language string

	// PauseBeforeMs is the pause before in milliseconds.
	PauseBeforeMs int

	// PauseAfterMs is the pause after in milliseconds.
	PauseAfterMs int

	// Emphasis is the emphasis level.
	Emphasis string

	// Rate is the speaking rate.
	Rate string

	// Pitch is the pitch adjustment.
	Pitch string
}

// Compile compiles the script for the specified language.
// Returns a slice of compiled segments ready for TTS processing.
func (c *Compiler) Compile(script *Script, language string) ([]CompiledSegment, error) {
	var segments []CompiledSegment

	for slideIdx, slide := range script.Slides {
		// Check if we should speak the title
		if slide.ShouldSpeakTitle() && slide.Title != "" {
			titleText := slide.Title

			// Apply pronunciations to title
			titleText = c.applyPronunciations(titleText, language, script.Pronunciations, nil)

			// Determine voice for title
			voiceID := ""
			if v, ok := slide.TitleVoice[language]; ok {
				voiceID = v
			} else if len(slide.Segments) > 0 {
				// Fall back to first segment's voice
				if v, ok := slide.Segments[0].Voice[language]; ok {
					voiceID = v
				}
			}
			if voiceID == "" {
				if v, ok := script.DefaultVoices[language]; ok {
					voiceID = v
				}
			}

			// Determine pause after title
			titlePauseAfter := ParseDuration(slide.TitlePauseAfter)
			if titlePauseAfter == 0 {
				// Apply defaults based on slide type
				if slide.IsSectionHeader {
					titlePauseAfter = 500 // 500ms for section headers
				} else {
					titlePauseAfter = 300 // 300ms for regular slides
				}
			}

			// Add pause before section headers
			pauseBefore := 0
			if slide.IsSectionHeader && slideIdx > 0 {
				pauseBefore = 1000 // 1s pause before section headers
			}

			segments = append(segments, CompiledSegment{
				SlideIndex:      slideIdx,
				SegmentIndex:    -1, // Title segments use -1
				SlideTitle:      slide.Title,
				IsTitleSegment:  true,
				IsSectionHeader: slide.IsSectionHeader,
				Text:            titleText,
				OriginalText:    slide.Title,
				VoiceID:         voiceID,
				Language:        language,
				PauseBeforeMs:   pauseBefore,
				PauseAfterMs:    titlePauseAfter,
			})
		}

		for segIdx, seg := range slide.Segments {
			text, ok := seg.Text[language]
			if !ok {
				continue // Skip segments without this language
			}

			originalText := text

			// Apply pronunciations
			text = c.applyPronunciations(text, language, script.Pronunciations, seg.Pronunciations)

			// Determine voice
			voiceID := ""
			if v, ok := seg.Voice[language]; ok {
				voiceID = v
			} else if v, ok := script.DefaultVoices[language]; ok {
				voiceID = v
			}

			// Parse pauses
			pauseBefore := ParseDuration(seg.PauseBefore)
			pauseAfter := ParseDuration(seg.PauseAfter)

			// Apply default segment pause
			if pauseAfter == 0 && c.DefaultPauseAfterSegment != "" {
				pauseAfter = ParseDuration(c.DefaultPauseAfterSegment)
			}

			// Add default slide pause after last segment
			if segIdx == len(slide.Segments)-1 && c.DefaultPauseAfterSlide != "" {
				slidePause := ParseDuration(c.DefaultPauseAfterSlide)
				if slidePause > pauseAfter {
					pauseAfter = slidePause
				}
			}

			segments = append(segments, CompiledSegment{
				SlideIndex:      slideIdx,
				SegmentIndex:    segIdx,
				SlideTitle:      slide.Title,
				IsSectionHeader: slide.IsSectionHeader,
				Text:            text,
				OriginalText:    originalText,
				VoiceID:         voiceID,
				Language:        language,
				PauseBeforeMs:   pauseBefore,
				PauseAfterMs:    pauseAfter,
				Emphasis:        seg.Emphasis,
				Rate:            seg.Rate,
				Pitch:           seg.Pitch,
			})
		}
	}

	return segments, nil
}

// applyPronunciations applies pronunciation substitutions to the text.
func (c *Compiler) applyPronunciations(text, language string, scriptProns, segmentProns map[string]map[string]string) string {
	// Build combined pronunciation map
	// Priority: additional > segment > script
	prons := make(map[string]string)

	// Script-level pronunciations
	for term, langMap := range scriptProns {
		if replacement, ok := langMap[language]; ok {
			prons[term] = replacement
		}
	}

	// Segment-level pronunciations (override script-level)
	for term, langMap := range segmentProns {
		if replacement, ok := langMap[language]; ok {
			prons[term] = replacement
		}
	}

	// Additional pronunciations from compiler (override all)
	for term, langMap := range c.AdditionalPronunciations {
		if replacement, ok := langMap[language]; ok {
			prons[term] = replacement
		}
	}

	// Apply substitutions (case-insensitive word boundary matching)
	result := text
	for term, replacement := range prons {
		pattern := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(term) + `\b`)
		result = pattern.ReplaceAllString(result, replacement)
	}

	return result
}

// AddPronunciation adds a pronunciation rule.
func (c *Compiler) AddPronunciation(term, language, replacement string) {
	if c.AdditionalPronunciations[term] == nil {
		c.AdditionalPronunciations[term] = make(map[string]string)
	}
	c.AdditionalPronunciations[term][language] = replacement
}

// AddPronunciations adds multiple pronunciation rules for a language.
func (c *Compiler) AddPronunciations(language string, rules map[string]string) {
	for term, replacement := range rules {
		c.AddPronunciation(term, language, replacement)
	}
}

// ParseDuration parses a duration string like "500ms" or "1s" to milliseconds.
func ParseDuration(s string) int {
	if s == "" {
		return 0
	}

	s = strings.TrimSpace(strings.ToLower(s))

	if strings.HasSuffix(s, "ms") {
		numStr := strings.TrimSuffix(s, "ms")
		if ms, err := strconv.Atoi(numStr); err == nil {
			return ms
		}
		return 0
	}

	if strings.HasSuffix(s, "s") {
		numStr := strings.TrimSuffix(s, "s")
		if sec, err := strconv.ParseFloat(numStr, 64); err == nil {
			return int(sec * 1000)
		}
		return 0
	}

	return 0
}

// FormatDuration formats milliseconds as a duration string.
func FormatDuration(ms int) string {
	if ms == 0 {
		return ""
	}
	if ms%1000 == 0 {
		return fmt.Sprintf("%ds", ms/1000)
	}
	return fmt.Sprintf("%dms", ms)
}

// GroupByVoice groups compiled segments by voice ID.
// Useful for batch processing with the same voice.
func GroupByVoice(segments []CompiledSegment) map[string][]CompiledSegment {
	groups := make(map[string][]CompiledSegment)
	for _, seg := range segments {
		groups[seg.VoiceID] = append(groups[seg.VoiceID], seg)
	}
	return groups
}

// GroupBySlide groups compiled segments by slide index.
func GroupBySlide(segments []CompiledSegment) map[int][]CompiledSegment {
	groups := make(map[int][]CompiledSegment)
	for _, seg := range segments {
		groups[seg.SlideIndex] = append(groups[seg.SlideIndex], seg)
	}
	return groups
}

// CombineText combines all segment texts into a single string with pause markers.
func CombineText(segments []CompiledSegment) string {
	var sb strings.Builder
	for i, seg := range segments {
		if i > 0 && seg.PauseBeforeMs > 0 {
			sb.WriteString(fmt.Sprintf(" [pause:%s] ", FormatDuration(seg.PauseBeforeMs)))
		}
		sb.WriteString(seg.Text)
		if seg.PauseAfterMs > 0 {
			sb.WriteString(fmt.Sprintf(" [pause:%s]", FormatDuration(seg.PauseAfterMs)))
		}
		if i < len(segments)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}
