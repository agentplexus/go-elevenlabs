package account

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/plexusone/elevenlabs-go/internal/api"
)

// Service handles user account, models, history, and pronunciation dictionary services.
type Service struct {
	apiClient *api.Client
}

// New creates a new account service.
func New(apiClient *api.Client) *Service {
	return &Service{apiClient: apiClient}
}

// --- User ---

// User represents an ElevenLabs user.
type User struct {
	// UserID is the unique user identifier.
	UserID string

	// FirstName is the user's first name.
	FirstName string

	// Subscription contains the user's subscription details.
	Subscription *Subscription

	// CreatedAt is when the user was created.
	CreatedAt time.Time
}

// Subscription represents a user's subscription details.
type Subscription struct {
	// Tier is the subscription tier (e.g., "free", "starter", "creator").
	Tier string

	// Status is the subscription status.
	Status string

	// CharacterCount is the number of characters used.
	CharacterCount int

	// CharacterLimit is the maximum characters allowed.
	CharacterLimit int

	// VoiceLimit is the maximum number of voices allowed.
	VoiceLimit int

	// VoiceSlotsUsed is the number of voice slots used.
	VoiceSlotsUsed int

	// CanUseInstantVoiceCloning indicates if instant cloning is available.
	CanUseInstantVoiceCloning bool

	// CanUseProfessionalVoiceCloning indicates if pro cloning is available.
	CanUseProfessionalVoiceCloning bool

	// NextCharacterResetUnix is when characters reset (Unix timestamp).
	NextCharacterResetUnix int64
}

// CharactersRemaining returns the number of characters remaining.
func (s *Subscription) CharactersRemaining() int {
	remaining := s.CharacterLimit - s.CharacterCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetUser returns the current user's information including subscription.
func (s *Service) GetUser(ctx context.Context) (*User, error) {
	resp, err := s.apiClient.GetUserInfo(ctx, api.GetUserInfoParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.UserResponseModel:
		user := &User{
			UserID:    r.UserID,
			CreatedAt: time.Unix(int64(r.CreatedAt), 0),
		}

		if r.FirstName.Set && !r.FirstName.Null {
			user.FirstName = r.FirstName.Value
		}

		sub := r.Subscription
		user.Subscription = &Subscription{
			Tier:                           sub.Tier,
			Status:                         string(sub.Status),
			CharacterCount:                 sub.CharacterCount,
			CharacterLimit:                 sub.CharacterLimit,
			VoiceLimit:                     sub.VoiceLimit,
			CanUseInstantVoiceCloning:      sub.CanUseInstantVoiceCloning,
			CanUseProfessionalVoiceCloning: sub.CanUseProfessionalVoiceCloning,
		}

		if sub.NextCharacterCountResetUnix.Set && !sub.NextCharacterCountResetUnix.Null {
			user.Subscription.NextCharacterResetUnix = int64(sub.NextCharacterCountResetUnix.Value)
		}

		return user, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// GetSubscription returns the current user's subscription details.
func (s *Service) GetSubscription(ctx context.Context) (*Subscription, error) {
	user, err := s.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return user.Subscription, nil
}

// GetCharactersRemaining returns the number of characters remaining in the current period.
func (s *Service) GetCharactersRemaining(ctx context.Context) (int, error) {
	sub, err := s.GetSubscription(ctx)
	if err != nil {
		return 0, err
	}
	return sub.CharactersRemaining(), nil
}

// --- Models ---

// Language represents a language supported by a model.
type Language struct {
	// LanguageID is the unique identifier (ISO code).
	LanguageID string

	// Name is the display name of the language.
	Name string
}

// Model represents an ElevenLabs model.
type Model struct {
	// ModelID is the unique identifier for the model.
	ModelID string

	// Name is the display name of the model.
	Name string

	// Description is the model description.
	Description string

	// CanDoTextToSpeech indicates if the model supports TTS.
	CanDoTextToSpeech bool

	// CanDoVoiceConversion indicates if the model supports voice conversion.
	CanDoVoiceConversion bool

	// CanBeFinetuned indicates if the model can be fine-tuned.
	CanBeFinetuned bool

	// CanUseStyle indicates if the model supports style settings.
	CanUseStyle bool

	// CanUseSpeakerBoost indicates if the model supports speaker boost.
	CanUseSpeakerBoost bool

	// Languages is the list of supported languages.
	Languages []*Language

	// MaxCharactersFreeUser is the max characters for free users.
	MaxCharactersFreeUser int

	// MaxCharactersSubscribedUser is the max characters for subscribed users.
	MaxCharactersSubscribedUser int

	// TokenCostFactor is the cost factor for the model.
	TokenCostFactor float64
}

// ListModels returns all available models.
func (s *Service) ListModels(ctx context.Context) ([]*Model, error) {
	resp, err := s.apiClient.GetModels(ctx, api.GetModelsParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetModelsOKApplicationJSON:
		models := make([]*Model, 0, len(*r))
		for _, m := range *r {
			model := &Model{
				ModelID:                     m.ModelID,
				Name:                        m.Name,
				Description:                 m.Description,
				CanDoTextToSpeech:           m.CanDoTextToSpeech,
				CanDoVoiceConversion:        m.CanDoVoiceConversion,
				CanBeFinetuned:              m.CanBeFinetuned,
				CanUseStyle:                 m.CanUseStyle,
				CanUseSpeakerBoost:          m.CanUseSpeakerBoost,
				MaxCharactersFreeUser:       m.MaxCharactersRequestFreeUser,
				MaxCharactersSubscribedUser: m.MaxCharactersRequestSubscribedUser,
				TokenCostFactor:             m.TokenCostFactor,
				Languages:                   make([]*Language, 0, len(m.Languages)),
			}
			for _, lang := range m.Languages {
				model.Languages = append(model.Languages, &Language{
					LanguageID: lang.LanguageID,
					Name:       lang.Name,
				})
			}
			models = append(models, model)
		}
		return models, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// ListTTSModels returns only models that support text-to-speech.
func (s *Service) ListTTSModels(ctx context.Context) ([]*Model, error) {
	models, err := s.ListModels(ctx)
	if err != nil {
		return nil, err
	}

	ttsModels := make([]*Model, 0)
	for _, m := range models {
		if m.CanDoTextToSpeech {
			ttsModels = append(ttsModels, m)
		}
	}
	return ttsModels, nil
}

// --- History ---

// HistoryItem represents a speech generation history item.
type HistoryItem struct {
	// HistoryItemID is the unique identifier.
	HistoryItemID string

	// VoiceID is the ID of the voice used.
	VoiceID string

	// VoiceName is the name of the voice used.
	VoiceName string

	// VoiceCategory is the category of the voice.
	VoiceCategory string

	// ModelID is the ID of the model used.
	ModelID string

	// Text is the text that was converted to speech.
	Text string

	// State is the state of the history item.
	State string

	// Source is the source of the generation.
	Source string

	// ContentType is the content type of the audio.
	ContentType string

	// CharactersUsed is the number of characters used.
	CharactersUsed int

	// CreatedAt is when the item was created.
	CreatedAt time.Time
}

// HistoryListResponse contains the list of history items and pagination info.
type HistoryListResponse struct {
	// Items is the list of history items.
	Items []*HistoryItem

	// HasMore indicates if there are more items to fetch.
	HasMore bool

	// LastHistoryItemID is the ID of the last item (for pagination).
	LastHistoryItemID string
}

// HistoryListOptions contains options for listing history items.
type HistoryListOptions struct {
	// PageSize is the number of items per page.
	PageSize int

	// StartAfterHistoryItemID is for pagination (fetch items after this ID).
	StartAfterHistoryItemID string

	// VoiceID filters by voice ID.
	VoiceID string
}

// ListHistory returns a list of speech history items.
func (s *Service) ListHistory(ctx context.Context, opts *HistoryListOptions) (*HistoryListResponse, error) {
	params := api.GetSpeechHistoryParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.StartAfterHistoryItemID != "" {
			params.StartAfterHistoryItemID = api.NewOptNilString(opts.StartAfterHistoryItemID)
		}
		if opts.VoiceID != "" {
			params.VoiceID = api.NewOptNilString(opts.VoiceID)
		}
	}

	resp, err := s.apiClient.GetSpeechHistory(ctx, params)
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetSpeechHistoryResponseModel:
		result := &HistoryListResponse{
			HasMore: r.HasMore,
			Items:   make([]*HistoryItem, 0, len(r.History)),
		}

		if r.LastHistoryItemID.Set && !r.LastHistoryItemID.Null {
			result.LastHistoryItemID = r.LastHistoryItemID.Value
		}

		for _, h := range r.History {
			item := &HistoryItem{
				HistoryItemID:  h.HistoryItemID,
				State:          string(h.State),
				ContentType:    h.ContentType,
				CharactersUsed: h.CharacterCountChangeTo - h.CharacterCountChangeFrom,
				CreatedAt:      time.Unix(int64(h.DateUnix), 0),
			}

			if h.VoiceID.Set && !h.VoiceID.Null {
				item.VoiceID = h.VoiceID.Value
			}
			if h.VoiceName.Set && !h.VoiceName.Null {
				item.VoiceName = h.VoiceName.Value
			}
			if h.VoiceCategory.Set && !h.VoiceCategory.Null {
				item.VoiceCategory = string(h.VoiceCategory.Value)
			}
			if h.ModelID.Set && !h.ModelID.Null {
				item.ModelID = h.ModelID.Value
			}
			if h.Text.Set && !h.Text.Null {
				item.Text = h.Text.Value
			}
			if h.Source.Set && !h.Source.Null {
				item.Source = string(h.Source.Value)
			}

			result.Items = append(result.Items, item)
		}

		return result, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// GetHistoryItem returns a specific history item by ID.
func (s *Service) GetHistoryItem(ctx context.Context, historyItemID string) (*HistoryItem, error) {
	if historyItemID == "" {
		return nil, fmt.Errorf("account: history_item_id cannot be empty")
	}

	resp, err := s.apiClient.GetSpeechHistoryItemByID(ctx, api.GetSpeechHistoryItemByIDParams{
		HistoryItemID: historyItemID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.SpeechHistoryItemResponseModel:
		item := &HistoryItem{
			HistoryItemID:  r.HistoryItemID,
			State:          string(r.State),
			ContentType:    r.ContentType,
			CharactersUsed: r.CharacterCountChangeTo - r.CharacterCountChangeFrom,
			CreatedAt:      time.Unix(int64(r.DateUnix), 0),
		}

		if r.VoiceID.Set && !r.VoiceID.Null {
			item.VoiceID = r.VoiceID.Value
		}
		if r.VoiceName.Set && !r.VoiceName.Null {
			item.VoiceName = r.VoiceName.Value
		}
		if r.VoiceCategory.Set && !r.VoiceCategory.Null {
			item.VoiceCategory = string(r.VoiceCategory.Value)
		}
		if r.ModelID.Set && !r.ModelID.Null {
			item.ModelID = r.ModelID.Value
		}
		if r.Text.Set && !r.Text.Null {
			item.Text = r.Text.Value
		}
		if r.Source.Set && !r.Source.Null {
			item.Source = string(r.Source.Value)
		}

		return item, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// GetHistoryAudio returns the audio for a history item.
func (s *Service) GetHistoryAudio(ctx context.Context, historyItemID string) (io.Reader, error) {
	if historyItemID == "" {
		return nil, fmt.Errorf("account: history_item_id cannot be empty")
	}

	resp, err := s.apiClient.GetAudioFullFromSpeechHistoryItem(ctx, api.GetAudioFullFromSpeechHistoryItemParams{
		HistoryItemID: historyItemID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetAudioFullFromSpeechHistoryItemOK:
		return r.Data, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// DeleteHistoryItem deletes a history item by ID.
func (s *Service) DeleteHistoryItem(ctx context.Context, historyItemID string) error {
	if historyItemID == "" {
		return fmt.Errorf("account: history_item_id cannot be empty")
	}

	_, err := s.apiClient.DeleteSpeechHistoryItem(ctx, api.DeleteSpeechHistoryItemParams{
		HistoryItemID: historyItemID,
	})
	return err
}

// --- Pronunciation ---

// PronunciationDictionary represents a pronunciation dictionary.
type PronunciationDictionary struct {
	// ID is the unique identifier.
	ID string

	// Name is the display name.
	Name string

	// Description is the dictionary description.
	Description string

	// LatestVersionID is the ID of the latest version.
	LatestVersionID string

	// RulesCount is the number of rules in the latest version.
	RulesCount int

	// CreatedBy is the user ID who created the dictionary.
	CreatedBy string

	// CreatedAt is when the dictionary was created.
	CreatedAt time.Time
}

// PronunciationDictionaryListResponse contains the list result.
type PronunciationDictionaryListResponse struct {
	// Dictionaries is the list of pronunciation dictionaries.
	Dictionaries []*PronunciationDictionary

	// HasMore indicates if there are more items to fetch.
	HasMore bool

	// NextCursor is the cursor for pagination.
	NextCursor string
}

// PronunciationDictionaryListOptions contains options for listing.
type PronunciationDictionaryListOptions struct {
	// PageSize is the number of items per page (max 100).
	PageSize int

	// Cursor is the pagination cursor.
	Cursor string
}

// ListPronunciationDictionaries returns all pronunciation dictionaries.
func (s *Service) ListPronunciationDictionaries(ctx context.Context, opts *PronunciationDictionaryListOptions) (*PronunciationDictionaryListResponse, error) {
	params := api.GetPronunciationDictionariesMetadataParams{}

	if opts != nil {
		if opts.PageSize > 0 {
			params.PageSize = api.NewOptInt(opts.PageSize)
		}
		if opts.Cursor != "" {
			params.Cursor = api.NewOptNilString(opts.Cursor)
		}
	}

	resp, err := s.apiClient.GetPronunciationDictionariesMetadata(ctx, params)
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetPronunciationDictionariesMetadataResponseModel:
		result := &PronunciationDictionaryListResponse{
			HasMore:      r.HasMore,
			Dictionaries: make([]*PronunciationDictionary, 0, len(r.PronunciationDictionaries)),
		}

		if r.NextCursor.Set && !r.NextCursor.Null {
			result.NextCursor = r.NextCursor.Value
		}

		for _, d := range r.PronunciationDictionaries {
			dict := &PronunciationDictionary{
				ID:              d.ID,
				Name:            d.Name,
				LatestVersionID: d.LatestVersionID,
				RulesCount:      d.LatestVersionRulesNum,
				CreatedBy:       d.CreatedBy,
				CreatedAt:       time.Unix(int64(d.CreationTimeUnix), 0),
			}
			if d.Description.Set && !d.Description.Null {
				dict.Description = d.Description.Value
			}
			result.Dictionaries = append(result.Dictionaries, dict)
		}

		return result, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// CreatePronunciationDictionaryRequest contains options for creating a pronunciation dictionary.
type CreatePronunciationDictionaryRequest struct {
	// Name is the name of the dictionary (required).
	Name string

	// Description is an optional description.
	Description string

	// PLSContent is the PLS (Pronunciation Lexicon Specification) XML content.
	PLSContent string

	// Rules is a convenient alternative to PLSContent.
	Rules PronunciationRules

	// Language is the language code for the rules (default: "en-US").
	Language string
}

// CreatePronunciationDictionary creates a new pronunciation dictionary.
func (s *Service) CreatePronunciationDictionary(ctx context.Context, req *CreatePronunciationDictionaryRequest) (*PronunciationDictionary, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("account: name cannot be empty")
	}

	body := &api.BodyAddAPronunciationDictionaryV1PronunciationDictionariesAddFromFilePostMultipart{
		Name: req.Name,
	}

	if req.Description != "" {
		body.Description = api.NewOptNilString(req.Description)
	}

	// Handle PLS content
	plsContent := req.PLSContent
	if plsContent == "" && len(req.Rules) > 0 {
		lang := req.Language
		if lang == "" {
			lang = "en-US"
		}
		generated, err := req.Rules.ToPLSString(lang)
		if err != nil {
			return nil, err
		}
		plsContent = generated
	}

	if plsContent != "" {
		body.File = api.NewOptNilString(plsContent)
	}

	resp, err := s.apiClient.AddFromFile(ctx, body, api.AddFromFileParams{})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.AddPronunciationDictionaryResponseModel:
		dict := &PronunciationDictionary{
			ID:              r.ID,
			Name:            r.Name,
			LatestVersionID: r.VersionID,
			RulesCount:      r.VersionRulesNum,
			CreatedBy:       r.CreatedBy,
			CreatedAt:       time.Unix(int64(r.CreationTimeUnix), 0),
		}
		if r.Description.Set && !r.Description.Null {
			dict.Description = r.Description.Value
		}
		return dict, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// GetPronunciationDictionaryPLS returns the PLS file for a specific version.
func (s *Service) GetPronunciationDictionaryPLS(ctx context.Context, dictionaryID, versionID string) (io.Reader, error) {
	if dictionaryID == "" {
		return nil, fmt.Errorf("account: dictionary_id cannot be empty")
	}
	if versionID == "" {
		return nil, fmt.Errorf("account: version_id cannot be empty")
	}

	resp, err := s.apiClient.GetPronunciationDictionaryVersionPls(ctx, api.GetPronunciationDictionaryVersionPlsParams{
		DictionaryID: dictionaryID,
		VersionID:    versionID,
	})
	if err != nil {
		return nil, err
	}

	switch r := resp.(type) {
	case *api.GetPronunciationDictionaryVersionPlsOKHeaders:
		return r.Response.Data, nil
	default:
		return nil, fmt.Errorf("account: unexpected response type")
	}
}

// --- Pronunciation Rules ---

// PronunciationRule defines how a word or phrase should be pronounced.
type PronunciationRule struct {
	// Grapheme is the text to match (required).
	Grapheme string `json:"grapheme"`

	// Alias is the replacement text (mutually exclusive with Phoneme).
	Alias string `json:"alias,omitempty"`

	// Phoneme is the IPA pronunciation (mutually exclusive with Alias).
	Phoneme string `json:"phoneme,omitempty"`
}

// Validate checks that the rule is valid.
func (r *PronunciationRule) Validate() error {
	if r.Grapheme == "" {
		return fmt.Errorf("account: grapheme cannot be empty")
	}
	if r.Alias == "" && r.Phoneme == "" {
		return fmt.Errorf("account: either alias or phoneme must be specified")
	}
	if r.Alias != "" && r.Phoneme != "" {
		return fmt.Errorf("account: cannot specify both alias and phoneme")
	}
	return nil
}

// PronunciationRules is a collection of pronunciation rules.
type PronunciationRules []PronunciationRule

// LoadRulesFromJSON loads pronunciation rules from a JSON file.
func LoadRulesFromJSON(filename string) (PronunciationRules, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("account: reading pronunciation rules file: %w", err)
	}

	var rules PronunciationRules
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("account: parsing pronunciation rules JSON: %w", err)
	}

	for i, rule := range rules {
		if err := rule.Validate(); err != nil {
			return nil, fmt.Errorf("account: rule %d: %w", i, err)
		}
	}

	return rules, nil
}

// RulesFromMap creates pronunciation rules from a simple map.
func RulesFromMap(m map[string]string) PronunciationRules {
	rules := make(PronunciationRules, 0, len(m))
	for grapheme, alias := range m {
		rules = append(rules, PronunciationRule{
			Grapheme: grapheme,
			Alias:    alias,
		})
	}
	return rules
}

// ToPLS converts the rules to PLS (Pronunciation Lexicon Specification) XML format.
func (rules PronunciationRules) ToPLS(language string) ([]byte, error) {
	if language == "" {
		language = "en-US"
	}

	var lexemes []plsLexeme
	for _, rule := range rules {
		lexeme := plsLexeme{
			Grapheme: rule.Grapheme,
		}
		if rule.Alias != "" {
			lexeme.Alias = rule.Alias
		} else {
			lexeme.Phoneme = rule.Phoneme
		}
		lexemes = append(lexemes, lexeme)
	}

	lexicon := plsLexicon{
		Version:  "1.0",
		XMLNS:    "http://www.w3.org/2005/01/pronunciation-lexicon",
		Alphabet: "ipa",
		XMLLang:  language,
		Lexemes:  lexemes,
	}

	output, err := xml.MarshalIndent(lexicon, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("account: generating PLS XML: %w", err)
	}

	return []byte(xml.Header + string(output)), nil
}

// ToPLSString is a convenience method that returns the PLS as a string.
func (rules PronunciationRules) ToPLSString(language string) (string, error) {
	data, err := rules.ToPLS(language)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Graphemes returns a list of all graphemes (the words being defined).
func (rules PronunciationRules) Graphemes() []string {
	result := make([]string, len(rules))
	for i, rule := range rules {
		result[i] = rule.Grapheme
	}
	return result
}

// String returns a human-readable summary of the rules.
func (rules PronunciationRules) String() string {
	var sb strings.Builder
	for _, rule := range rules {
		if rule.Alias != "" {
			sb.WriteString(fmt.Sprintf("%s → %s\n", rule.Grapheme, rule.Alias))
		} else {
			sb.WriteString(fmt.Sprintf("%s → [%s]\n", rule.Grapheme, rule.Phoneme))
		}
	}
	return sb.String()
}

// PLS XML structures (internal)

type plsLexicon struct {
	XMLName  xml.Name    `xml:"lexicon"`
	Version  string      `xml:"version,attr"`
	XMLNS    string      `xml:"xmlns,attr"`
	Alphabet string      `xml:"alphabet,attr"`
	XMLLang  string      `xml:"xml:lang,attr"`
	Lexemes  []plsLexeme `xml:"lexeme"`
}

type plsLexeme struct {
	Grapheme string `xml:"grapheme"`
	Alias    string `xml:"alias,omitempty"`
	Phoneme  string `xml:"phoneme,omitempty"`
}
