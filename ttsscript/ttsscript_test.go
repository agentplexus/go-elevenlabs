package ttsscript

import (
	"strings"
	"testing"
)

func TestParseScript(t *testing.T) {
	jsonData := `{
		"title": "Test Script",
		"default_voices": {"en": "voice-1"},
		"pronunciations": {"API": {"en": "A P I"}},
		"slides": [
			{
				"title": "Intro",
				"segments": [
					{"text": {"en": "Hello API world"}, "pause_after": "500ms"}
				]
			}
		]
	}`

	script, err := ParseScript([]byte(jsonData))
	if err != nil {
		t.Fatalf("ParseScript failed: %v", err)
	}

	if script.Title != "Test Script" {
		t.Errorf("expected title 'Test Script', got '%s'", script.Title)
	}

	if len(script.Slides) != 1 {
		t.Errorf("expected 1 slide, got %d", len(script.Slides))
	}

	if script.Slides[0].Title != "Intro" {
		t.Errorf("expected slide title 'Intro', got '%s'", script.Slides[0].Title)
	}
}

func TestCompiler(t *testing.T) {
	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Pronunciations: map[string]map[string]string{
			"API": {"en": "A P I"},
		},
		Slides: []Slide{
			{
				Title: "Slide 1",
				Segments: []Segment{
					{
						Text:       map[string]string{"en": "Hello API world"},
						PauseAfter: "500ms",
					},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	seg := segments[0]

	// Check pronunciation was applied
	if seg.Text != "Hello A P I world" {
		t.Errorf("expected 'Hello A P I world', got '%s'", seg.Text)
	}

	// Check voice was set
	if seg.VoiceID != "voice-1" {
		t.Errorf("expected voice 'voice-1', got '%s'", seg.VoiceID)
	}

	// Check pause was parsed (500ms + default slide pause 800ms = 800ms since slide pause > segment pause)
	if seg.PauseAfterMs != 800 {
		t.Errorf("expected pause 800ms, got %dms", seg.PauseAfterMs)
	}
}

func TestSSMLFormatter(t *testing.T) {
	segments := []CompiledSegment{
		{
			SlideIndex:   0,
			SegmentIndex: 0,
			SlideTitle:   "Intro",
			Text:         "Hello world",
			PauseAfterMs: 500,
		},
	}

	formatter := NewSSMLFormatter()
	ssml := formatter.Format(segments, "en")

	if !strings.Contains(ssml, "<speak") {
		t.Error("SSML should contain <speak> tag")
	}

	if !strings.Contains(ssml, "Hello world") {
		t.Error("SSML should contain the text")
	}

	if !strings.Contains(ssml, `<break time="500ms"`) {
		t.Error("SSML should contain break element")
	}

	if !strings.Contains(ssml, "<!-- Slide 1: Intro -->") {
		t.Error("SSML should contain slide comment")
	}
}

func TestElevenLabsFormatter(t *testing.T) {
	segments := []CompiledSegment{
		{
			SlideIndex:   0,
			SegmentIndex: 0,
			SlideTitle:   "Intro",
			Text:         "Hello world",
			VoiceID:      "voice-1",
			PauseAfterMs: 500,
		},
		{
			SlideIndex:   0,
			SegmentIndex: 1,
			Text:         "Second segment",
			VoiceID:      "voice-2",
		},
	}

	formatter := NewElevenLabsFormatter()
	result := formatter.Format(segments)

	if len(result) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(result))
	}

	if result[0].VoiceID != "voice-1" {
		t.Errorf("expected voice 'voice-1', got '%s'", result[0].VoiceID)
	}

	if result[0].PauseAfterMs != 500 {
		t.Errorf("expected pause 500ms, got %dms", result[0].PauseAfterMs)
	}

	// Test grouping by voice
	groups := formatter.GroupByVoice(result)
	if len(groups) != 2 {
		t.Errorf("expected 2 voice groups, got %d", len(groups))
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"500ms", 500},
		{"1s", 1000},
		{"1.5s", 1500},
		{"2s", 2000},
		{"", 0},
		{"100ms", 100},
	}

	for _, tt := range tests {
		result := ParseDuration(tt.input)
		if result != tt.expected {
			t.Errorf("ParseDuration(%q) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{500, "500ms"},
		{1000, "1s"},
		{2000, "2s"},
		{1500, "1500ms"},
		{0, ""},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("FormatDuration(%d) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestEscapeSSML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello & world", "Hello &amp; world"},
		{"<tag>", "&lt;tag&gt;"},
		{`Say "hello"`, `Say &quot;hello&quot;`},
		{"It's fine", "It&apos;s fine"},
	}

	for _, tt := range tests {
		result := EscapeSSML(tt.input)
		if result != tt.expected {
			t.Errorf("EscapeSSML(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestScriptLanguages(t *testing.T) {
	script := &Script{
		Slides: []Slide{
			{
				Segments: []Segment{
					{Text: map[string]string{"en": "Hello", "es": "Hola"}},
					{Text: map[string]string{"en": "World", "fr": "Monde"}},
				},
			},
		},
	}

	langs := script.Languages()
	if len(langs) != 3 {
		t.Errorf("expected 3 languages, got %d", len(langs))
	}

	langSet := make(map[string]bool)
	for _, l := range langs {
		langSet[l] = true
	}

	for _, expected := range []string{"en", "es", "fr"} {
		if !langSet[expected] {
			t.Errorf("expected language %q in result", expected)
		}
	}
}

func TestScriptValidate(t *testing.T) {
	// Valid script
	valid := &Script{
		Slides: []Slide{
			{Segments: []Segment{{Text: map[string]string{"en": "Hello"}}}},
		},
	}
	if issues := valid.Validate(); len(issues) != 0 {
		t.Errorf("valid script should have no issues, got: %v", issues)
	}

	// Empty script
	empty := &Script{}
	if issues := empty.Validate(); len(issues) == 0 {
		t.Error("empty script should have issues")
	}

	// Slide with no segments
	noSegs := &Script{
		Slides: []Slide{{}},
	}
	if issues := noSegs.Validate(); len(issues) == 0 {
		t.Error("slide with no segments should have issues")
	}
}

func TestShouldSpeakTitle(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }

	tests := []struct {
		name     string
		slide    Slide
		expected bool
	}{
		{
			name:     "regular slide without SpeakTitle",
			slide:    Slide{Title: "Slide 1"},
			expected: false,
		},
		{
			name:     "section header without SpeakTitle",
			slide:    Slide{Title: "Section 1", IsSectionHeader: true},
			expected: true,
		},
		{
			name:     "regular slide with SpeakTitle true",
			slide:    Slide{Title: "Slide 1", SpeakTitle: boolPtr(true)},
			expected: true,
		},
		{
			name:     "regular slide with SpeakTitle false",
			slide:    Slide{Title: "Slide 1", SpeakTitle: boolPtr(false)},
			expected: false,
		},
		{
			name:     "section header with SpeakTitle false",
			slide:    Slide{Title: "Section 1", IsSectionHeader: true, SpeakTitle: boolPtr(false)},
			expected: false,
		},
		{
			name:     "section header with SpeakTitle true",
			slide:    Slide{Title: "Section 1", IsSectionHeader: true, SpeakTitle: boolPtr(true)},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.slide.ShouldSpeakTitle()
			if result != tt.expected {
				t.Errorf("ShouldSpeakTitle() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCompilerSectionHeader(t *testing.T) {
	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Slides: []Slide{
			{
				Title:           "Introduction",
				IsSectionHeader: true,
				Segments: []Segment{
					{Text: map[string]string{"en": "Welcome to the course"}},
				},
			},
			{
				Title: "Regular Slide",
				Segments: []Segment{
					{Text: map[string]string{"en": "Some content"}},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: title segment + content segment for section header, + content segment for regular slide = 3
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}

	// First segment should be the title
	titleSeg := segments[0]
	if !titleSeg.IsTitleSegment {
		t.Error("first segment should be a title segment")
	}
	if titleSeg.Text != "Introduction" {
		t.Errorf("expected title text 'Introduction', got '%s'", titleSeg.Text)
	}
	if titleSeg.SegmentIndex != -1 {
		t.Errorf("title segment should have index -1, got %d", titleSeg.SegmentIndex)
	}
	if !titleSeg.IsSectionHeader {
		t.Error("title segment should be marked as section header")
	}
	// Default pause after section header title is 500ms
	if titleSeg.PauseAfterMs != 500 {
		t.Errorf("expected title pause 500ms, got %dms", titleSeg.PauseAfterMs)
	}

	// Second segment should be the section header content
	contentSeg := segments[1]
	if contentSeg.IsTitleSegment {
		t.Error("second segment should not be a title segment")
	}
	if contentSeg.Text != "Welcome to the course" {
		t.Errorf("expected content 'Welcome to the course', got '%s'", contentSeg.Text)
	}

	// Third segment should be from regular slide (no title segment since not section header)
	regularSeg := segments[2]
	if regularSeg.IsTitleSegment {
		t.Error("third segment should not be a title segment")
	}
	if regularSeg.Text != "Some content" {
		t.Errorf("expected content 'Some content', got '%s'", regularSeg.Text)
	}
}

func TestCompilerSpeakTitleExplicit(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }

	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Slides: []Slide{
			{
				Title:      "Spoken Title",
				SpeakTitle: boolPtr(true),
				Segments: []Segment{
					{Text: map[string]string{"en": "Content"}},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have title segment + content segment = 2
	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	titleSeg := segments[0]
	if !titleSeg.IsTitleSegment {
		t.Error("first segment should be a title segment")
	}
	if titleSeg.Text != "Spoken Title" {
		t.Errorf("expected title 'Spoken Title', got '%s'", titleSeg.Text)
	}
	// Default pause for regular slide title is 300ms
	if titleSeg.PauseAfterMs != 300 {
		t.Errorf("expected title pause 300ms, got %dms", titleSeg.PauseAfterMs)
	}
}

func TestCompilerTitleVoice(t *testing.T) {
	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "default-voice"},
		Slides: []Slide{
			{
				Title:           "Section Title",
				IsSectionHeader: true,
				TitleVoice:      map[string]string{"en": "title-voice"},
				Segments: []Segment{
					{
						Text:  map[string]string{"en": "Content"},
						Voice: map[string]string{"en": "segment-voice"},
					},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	// Title should use TitleVoice
	if segments[0].VoiceID != "title-voice" {
		t.Errorf("expected title voice 'title-voice', got '%s'", segments[0].VoiceID)
	}

	// Content should use segment voice
	if segments[1].VoiceID != "segment-voice" {
		t.Errorf("expected segment voice 'segment-voice', got '%s'", segments[1].VoiceID)
	}
}

func TestCompilerTitlePauseAfter(t *testing.T) {
	boolPtr := func(v bool) *bool { return &v }

	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Slides: []Slide{
			{
				Title:           "Custom Pause",
				SpeakTitle:      boolPtr(true),
				TitlePauseAfter: "1s",
				Segments: []Segment{
					{Text: map[string]string{"en": "Content"}},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	// Title should use custom pause
	if segments[0].PauseAfterMs != 1000 {
		t.Errorf("expected title pause 1000ms, got %dms", segments[0].PauseAfterMs)
	}
}

func TestCompilerSectionHeaderPauseBefore(t *testing.T) {
	script := &Script{
		Title:         "Test",
		DefaultVoices: map[string]string{"en": "voice-1"},
		Slides: []Slide{
			{
				Title: "First Slide",
				Segments: []Segment{
					{Text: map[string]string{"en": "First content"}},
				},
			},
			{
				Title:           "New Section",
				IsSectionHeader: true,
				Segments: []Segment{
					{Text: map[string]string{"en": "Section content"}},
				},
			},
		},
	}

	compiler := NewCompiler()
	segments, err := compiler.Compile(script, "en")
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Should have: first content + section title + section content = 3
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}

	// Section header title (index 1) should have pause before (1000ms)
	sectionTitle := segments[1]
	if !sectionTitle.IsTitleSegment {
		t.Error("second segment should be a title segment")
	}
	if sectionTitle.PauseBeforeMs != 1000 {
		t.Errorf("expected section title pause before 1000ms, got %dms", sectionTitle.PauseBeforeMs)
	}
}

func TestElevenLabsFormatterTitleSegment(t *testing.T) {
	segments := []CompiledSegment{
		{
			SlideIndex:     0,
			SegmentIndex:   -1,
			SlideTitle:     "Introduction",
			IsTitleSegment: true,
			Text:           "Introduction",
			VoiceID:        "voice-1",
			PauseAfterMs:   500,
		},
		{
			SlideIndex:   0,
			SegmentIndex: 0,
			SlideTitle:   "Introduction",
			Text:         "Welcome to the course",
			VoiceID:      "voice-1",
		},
	}

	formatter := NewElevenLabsFormatter()
	result := formatter.Format(segments)

	if len(result) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(result))
	}

	// Check title segment filename
	if result[0].SuggestedFilename != "slide01_title.mp3" {
		t.Errorf("expected title filename 'slide01_title.mp3', got '%s'", result[0].SuggestedFilename)
	}
	if !result[0].IsTitleSegment {
		t.Error("first segment should be marked as title segment")
	}

	// Check regular segment filename
	if result[1].SuggestedFilename != "slide01_seg01.mp3" {
		t.Errorf("expected segment filename 'slide01_seg01.mp3', got '%s'", result[1].SuggestedFilename)
	}
	if result[1].IsTitleSegment {
		t.Error("second segment should not be marked as title segment")
	}
}

func TestBatchConfigGenerateFilename(t *testing.T) {
	config := NewBatchConfig("./output")

	tests := []struct {
		name     string
		segment  ElevenLabsSegment
		language string
		expected string
	}{
		{
			name: "regular segment",
			segment: ElevenLabsSegment{
				SlideIndex:   0,
				SegmentIndex: 0,
			},
			language: "en",
			expected: "./output/slide01_seg01_en.mp3",
		},
		{
			name: "title segment",
			segment: ElevenLabsSegment{
				SlideIndex:     0,
				SegmentIndex:   -1,
				IsTitleSegment: true,
			},
			language: "en",
			expected: "./output/slide01_title_en.mp3",
		},
		{
			name: "second slide title",
			segment: ElevenLabsSegment{
				SlideIndex:     1,
				SegmentIndex:   -1,
				IsTitleSegment: true,
			},
			language: "es",
			expected: "./output/slide02_title_es.mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.GenerateFilename(tt.segment, tt.language)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
