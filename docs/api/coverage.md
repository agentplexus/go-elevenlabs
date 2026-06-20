# API Coverage

This page documents the ElevenLabs API coverage in this Go SDK.

**Total API Methods:** 204
**SDK Status:** Core functionality covered, advanced features in progress

## Coverage Summary

| Category | API Methods | SDK Coverage |
|----------|-------------|--------------|
| Text-to-Speech | 4 | ✓ Full |
| WebSocket TTS | 1 | ✓ Full |
| Speech-to-Text | 2 | ✓ Full |
| WebSocket STT | 1 | ✓ Full |
| Speech-to-Speech | 2 | ✓ Full |
| Voices | 9 | ✓ Full |
| Models | 1 | ✓ Full |
| History | 5 | ✓ Full |
| User | 1 | ✓ Full |
| Sound Effects | 1 | ✓ Full |
| Forced Alignment | 1 | ✓ Full |
| Audio Isolation | 2 | ✓ Full |
| Text-to-Dialogue | 4 | ✓ Full |
| Voice Design | 8 | ✓ Partial |
| Music | 5 | ✓ Full |
| Pronunciation | 6 | ✓ Full |
| Projects (Studio) | 14 | ✓ Partial |
| Dubbing | 14 | ✓ Partial |
| Phone / Twilio | 7 | ✓ Partial |
| Professional Voice Cloning | 12 | ✗ Not covered |
| Voice Library | 5 | ✗ Not covered |
| Conversational AI | 26+ | ✓ Partial |
| Knowledge Base / RAG | 15 | ✗ Not covered |
| Workspace Management | 20 | ✗ Not covered |
| MCP / Tools | 5 | ✗ Not covered |
| Audio Native | 3 | ✗ Not covered |
| Transcription | 4 | ✗ Not covered |
| Miscellaneous | 6 | ✗ Not covered |

---

## Covered APIs

### Text-to-Speech (4 methods) ✓

Full coverage of text-to-speech functionality.

| Method | SDK Support |
|--------|-------------|
| `TextToSpeechFull` | ✓ `TextToSpeech().Generate()` |
| `TextToSpeechStream` | ✓ `TextToSpeech().GenerateStream()` |
| `TextToSpeechFullWithTimestamps` | ✓ `TextToSpeech().GenerateWithTimestamps()` |
| `TextToSpeechStreamWithTimestamps` | ✓ `TextToSpeech().GenerateStreamWithTimestamps()` |

### WebSocket TTS (1 method) ✓

Real-time text-to-speech streaming via WebSocket for low-latency voice synthesis.

| Method | SDK Support |
|--------|-------------|
| `TextToSpeechWebSocket` | ✓ `WebSocketTTS().Connect()` |

**Key Features:**

- Stream text to speech in real-time (ideal for LLM output)
- Low-latency audio generation with configurable optimization
- Word-level timing alignment
- SSML parsing support

### Speech-to-Text (2 methods) ✓

Full coverage of transcription functionality.

| Method | SDK Support |
|--------|-------------|
| `SpeechToText` | ✓ `SpeechToText().Transcribe()` |
| `Transcribe` | ✓ `SpeechToText().TranscribeURL()` |

### WebSocket STT (1 method) ✓

Real-time speech-to-text streaming via WebSocket for live transcription.

| Method | SDK Support |
|--------|-------------|
| `SpeechToTextWebSocket` | ✓ `WebSocketSTT().Connect()` |

**Key Features:**

- Stream audio for real-time transcription
- Partial (interim) results for responsive UIs
- Word-level timing with confidence scores
- Automatic language detection

### Speech-to-Speech (2 methods) ✓

Voice conversion - transform speech from one voice to another.

| Method | SDK Support |
|--------|-------------|
| `SpeechToSpeechFull` | ✓ `SpeechToSpeech().Convert()` |
| `SpeechToSpeechStream` | ✓ `SpeechToSpeech().ConvertStream()` |

**Key Features:**

- Convert voice while preserving speech content
- Background noise removal option
- Seed audio for consistent conversions
- Configurable voice settings

### Voices (9 methods) ✓

Full coverage of voice management.

| Method | SDK Support |
|--------|-------------|
| `GetVoices` | ✓ `Voices().List()` |
| `GetVoiceByID` | ✓ `Voices().Get()` |
| `AddVoice` | ✓ `Voices().Add()` |
| `EditVoice` | ✓ `Voices().Edit()` |
| `DeleteVoice` | ✓ `Voices().Delete()` |
| `GetVoiceSettings` | ✓ `Voices().GetSettings()` |
| `EditVoiceSettings` | ✓ `Voices().EditSettings()` |
| `GetVoiceSettingsDefault` | ✓ `Voices().GetDefaultSettings()` |
| `GetUserVoicesV2` | ✓ `Voices().ListUserVoices()` |

### Models (1 method) ✓

| Method | SDK Support |
|--------|-------------|
| `GetModels` | ✓ `Models().List()` |

### History (5 methods) ✓

| Method | SDK Support |
|--------|-------------|
| `GetSpeechHistory` | ✓ `History().List()` |
| `GetSpeechHistoryItemByID` | ✓ `History().Get()` |
| `DeleteSpeechHistoryItem` | ✓ `History().Delete()` |
| `DownloadSpeechHistoryItems` | ✓ `History().Download()` |
| `GetAudioFullFromSpeechHistoryItem` | ✓ `History().GetAudio()` |

### User (1 method) ✓

| Method | SDK Support |
|--------|-------------|
| `GetUserInfo` | ✓ `User().Get()` |

### Sound Effects (1 method) ✓

| Method | SDK Support |
|--------|-------------|
| `SoundGeneration` | ✓ `SoundEffects().Generate()` |

### Forced Alignment (1 method) ✓

| Method | SDK Support |
|--------|-------------|
| `ForcedAlignment` | ✓ `ForcedAlignment().Align()` |

### Audio Isolation (2 methods) ✓

| Method | SDK Support |
|--------|-------------|
| `AudioIsolation` | ✓ `AudioIsolation().Isolate()` |
| `AudioIsolationStream` | ✓ `AudioIsolation().IsolateStream()` |

### Text-to-Dialogue (4 methods) ✓

| Method | SDK Support |
|--------|-------------|
| `TextToDialogue` | ✓ `TextToDialogue().Generate()` |
| `TextToDialogueStream` | ✓ `TextToDialogue().GenerateStream()` |
| `TextToDialogueFullWithTimestamps` | ✓ `TextToDialogue().GenerateWithTimestamps()` |
| `TextToDialogueStreamWithTimestamps` | ✓ Planned |

### Voice Design (8 methods) - Partial ✓

| Method | SDK Support |
|--------|-------------|
| `GenerateRandomVoice` | ✓ `VoiceDesign().GeneratePreview()` |
| `CreateVoiceOld` | ✓ `VoiceDesign().SaveVoice()` |
| `CreateVoice` | ✗ Not covered |
| `TextToVoice` | ✗ Not covered |
| `TextToVoiceDesign` | ✗ Not covered |
| `TextToVoicePreviewStream` | ✗ Not covered |
| `TextToVoiceRemix` | ✗ Not covered |
| `GetGenerateVoiceParameters` | ✗ Not covered |

### Music (5 methods) ✓

Full coverage of music generation and stem separation.

| Method | SDK Support |
|--------|-------------|
| `Generate` | ✓ `Music().Generate()` |
| `StreamCompose` | ✓ `Music().GenerateStream()` |
| `ComposeDetailed` | ✓ `Music().GenerateDetailed()` |
| `ComposePlan` | ✓ `Music().GeneratePlan()` |
| `SeparateSongStems` | ✓ `Music().SeparateStems()` |

### Pronunciation (6 methods) ✓

Full coverage of pronunciation dictionary management.

| Method | SDK Support |
|--------|-------------|
| `AddFromFile` | ✓ `Pronunciation().Create()` |
| `GetPronunciationDictionariesMetadata` | ✓ `Pronunciation().List()` |
| `GetPronunciationDictionaryMetadata` | ✓ `Pronunciation().Get()` |
| `PatchPronunciationDictionary` | ✓ `Pronunciation().Rename()`, `Pronunciation().Archive()` |
| `RemoveRules` | ✓ `Pronunciation().RemoveRules()` |
| `GetPronunciationDictionaryVersionPls` | ✓ `Pronunciation().GetVersionPLS()`, `Pronunciation().DownloadLatestPLS()` |

!!! note
    `UpdatePronunciationDictionaries` is a Projects API method that associates dictionaries with a project, not a pronunciation dictionary method.

### Projects / Studio (14 methods) - Partial ✓

| Method | SDK Support |
|--------|-------------|
| `AddProject` | ✓ `Projects().Create()` |
| `GetProjects` | ✓ `Projects().List()` |
| `DeleteProject` | ✓ `Projects().Delete()` |
| `ConvertProjectEndpoint` | ✓ `Projects().Convert()` |
| `EditProject` | ✗ Not covered |
| `EditProjectContent` | ✗ Not covered |
| `GetChapters` | ✗ Not covered |
| `ConvertChapterEndpoint` | ✗ Not covered |
| `DeleteChapterEndpoint` | ✗ Not covered |
| `GetChapterSnapshots` | ✗ Not covered |
| `GetChapterSnapshotEndpoint` | ✗ Not covered |
| `GetProjectSnapshots` | ✗ Not covered |
| `GetProjectSnapshotEndpoint` | ✗ Not covered |
| `StreamProjectSnapshotArchiveEndpoint` | ✗ Not covered |

### Dubbing (14 methods) - Partial ✓

| Method | SDK Support |
|--------|-------------|
| `CreateDubbing` | ✓ `Dubbing().Create()` |
| `DeleteDubbing` | ✓ `Dubbing().Delete()` |
| `GetDubbingResource` | ✓ `Dubbing().GetStatus()` |
| `Dub` | ✗ Not covered |
| `AddLanguage` | ✗ Not covered |
| `CreateSpeaker` | ✗ Not covered |
| `UpdateSpeaker` | ✗ Not covered |
| `GetSpeakerAudio` | ✗ Not covered |
| `GetSimilarVoicesForSpeaker` | ✗ Not covered |
| `StartSpeakerSeparation` | ✗ Not covered |
| `CreateClip` | ✗ Not covered |
| `DeleteSegment` | ✗ Not covered |
| `UpdateSegmentLanguage` | ✗ Not covered |
| `MigrateSegments` | ✗ Not covered |

### Phone / Twilio (7 methods) - Partial ✓

Phone call and Twilio integration for conversational AI agents.

| Method | SDK Support |
|--------|-------------|
| `RegisterTwilioCall` | ✓ `Twilio().RegisterCall()` |
| `HandleTwilioOutboundCall` | ✓ `Twilio().OutboundCall()` |
| `HandleSipTrunkOutboundCall` | ✓ `Twilio().SIPOutboundCall()` |
| `ListPhoneNumbersRoute` | ✓ `PhoneNumbers().List()` |
| `GetPhoneNumberRoute` | ✓ `PhoneNumbers().Get()` |
| `UpdatePhoneNumberRoute` | ✓ `PhoneNumbers().Update()` |
| `DeletePhoneNumberRoute` | ✓ `PhoneNumbers().Delete()` |

**Key Features:**

- Register incoming Twilio calls with ElevenLabs agents
- Initiate outbound calls via Twilio or SIP trunks
- Manage phone numbers associated with agents
- Dynamic variables and prompt overrides per call

### Conversational AI Agents (26+ methods) - Partial ✓

AI agent management, branches, deployment, and testing.

| Method | SDK Support |
|--------|-------------|
| `GetAgentsRoute` | ✓ `Agents().List()` |
| `GetAgent` | ✓ `Agents().Get()` |
| `CreateAgent` | ✓ `Agents().Create()` |
| `UpdateAgent` | ✓ `Agents().Update()` |
| `DeleteAgentRoute` | ✓ `Agents().Delete()` |
| `DuplicateAgentRoute` | ✓ `Agents().Duplicate()` |
| `PostAgentAvatarRoute` | ✓ `Agents().UploadAvatar()` |
| `GetAgentLinkRoute` | ✓ `Agents().GetLink()` |
| `GetAgentTopicsRoute` | ✓ `Agents().GetTopics()` |
| `GetWidgetRoute` | ✓ `Agents().GetWidget()` |
| `GetBranchesRoute` | ✓ `Agents().ListBranches()` |
| `GetBranchRoute` | ✓ `Agents().GetBranch()` |
| `CreateBranch` | ✓ `Agents().CreateBranch()` |
| `UpdateBranchRoute` | ✓ `Agents().UpdateBranch()` |
| `MergeBranchIntoTarget` | ✓ `Agents().MergeBranch()` |
| `CreateAgentDeploymentRoute` | ✓ `Agents().Deploy()` |
| `SimulateConversation` | ✓ `Agents().SimulateConversation()` |
| `SimulateConversationStream` | ✓ `Agents().SimulateConversationStream()` |
| `CreateAgentTestFolderRoute` | ✓ `Agents().CreateTestFolder()` |
| `GetAgentTestFolderRoute` | ✓ `Agents().GetTestFolder()` |
| `UpdateAgentTestFolderRoute` | ✓ `Agents().UpdateTestFolder()` |
| `DeleteAgentTestFolderRoute` | ✓ `Agents().DeleteTestFolder()` |
| `BulkMoveAgentTestsRoute` | ✓ `Agents().BulkMoveTests()` |
| `GetAgentKnowledgeBaseSize` | ✓ `Agents().GetKnowledgeBaseSize()` |
| `GetConversationHistoriesRoute` | ✓ `Conversations().List()` |
| `GetConversationHistoryRoute` | ✓ `Conversations().Get()` |
| `DeleteConversationRoute` | ✓ `Conversations().Delete()` |
| `GetConversationAudioRoute` | ✓ `Conversations().GetAudio()` |
| `PostConversationFeedbackRoute` | ✓ `Conversations().SubmitFeedback()` |
| `RunConversationAnalysis` | ✓ `Conversations().RunAnalysis()` |
| `CreateBatchCall` | ✓ `BatchCalling().Create()` |
| `GetBatchCall` | ✓ `BatchCalling().Get()` |
| `CancelBatchCall` | ✓ `BatchCalling().Cancel()` |
| `RetryBatchCall` | ✓ `BatchCalling().Retry()` |
| `DeleteBatchCall` | ✓ `BatchCalling().Delete()` |
| `GetWorkspaceBatchCalls` | ✓ `BatchCalling().List()` |
| `GetLiveCount` | ✓ `Analytics().GetLiveCount()` |

**Key Features:**

- Full agent CRUD (create, read, update, delete, duplicate)
- Branch management with version control and merging
- Traffic splitting deployment for A/B testing
- Conversation simulation for agent testing (streaming and non-streaming)
- Test folder organization
- Conversation history with search and analytics
- Batch calling for outbound call campaigns
- Real-time analytics (live call count)

---

## Not Covered APIs

### Professional Voice Cloning - PVC (12 methods)

Professional-grade voice cloning with training.

| Method | Description |
|--------|-------------|
| `CreatePvcVoice` | Create a PVC voice |
| `EditPvcVoice` | Edit PVC voice settings |
| `AddPvcVoiceSamples` | Add training samples |
| `DeletePvcVoiceSample` | Delete a sample |
| `EditPvcVoiceSample` | Edit sample metadata |
| `GetPvcSampleAudio` | Get sample audio |
| `GetPvcSampleSpeakers` | Get detected speakers |
| `GetPvcSampleVisualWaveform` | Get waveform visualization |
| `GetPvcVoiceCaptcha` | Get verification captcha |
| `VerifyPvcVoiceCaptcha` | Verify captcha |
| `RequestPvcManualVerification` | Request manual verification |
| `RunPvcVoiceTraining` | Start voice training |

### Voice Library (5 methods)

Community voice discovery and sharing.

| Method | Description |
|--------|-------------|
| `GetLibraryVoices` | Browse community voices |
| `GetSimilarLibraryVoices` | Find similar voices |
| `AddSharingVoice` | Add a shared voice to your library |
| `ShareResourceEndpoint` | Share a resource |
| `UnshareResourceEndpoint` | Unshare a resource |

### Conversational AI - Remaining (Not Covered)

Methods not yet covered in the Agents service.

| Method | Description |
|--------|-------------|
| `GetAgentKnowledgeBaseSummariesRoute` | Get KB summaries |
| `GetAgentLlmExpectedCostCalculation` | Estimate LLM costs |
| `CreateAgentResponseTestRoute` | Create response test |
| `GetAgentResponseTestRoute` | Get test results |
| `UpdateAgentResponseTestRoute` | Update test |
| `GetAgentResponseTestsSummariesRoute` | Get test summaries |
| `GetConversationSignedLink` | Get signed link for agents |
| `GetLivekitToken` | Get LiveKit token |
| `ListChatResponseTestsRoute` | List chat tests |
| `DeleteChatResponseTestRoute` | Delete chat test |

### Knowledge Base / RAG (15 methods)

Document management and retrieval-augmented generation.

| Method | Description |
|--------|-------------|
| `AddDocumentationToKnowledgeBase` | Add documentation |
| `CreateFileDocumentRoute` | Create from file |
| `CreateTextDocumentRoute` | Create from text |
| `CreateURLDocumentRoute` | Create from URL |
| `DeleteKnowledgeBaseDocument` | Delete document |
| `GetDocumentationFromKnowledgeBase` | Get documentation |
| `GetDocumentationChunkFromKnowledgeBase` | Get chunk |
| `GetKnowledgeBaseContent` | Get KB content |
| `GetKnowledgeBaseListRoute` | List knowledge bases |
| `UpdateDocumentRoute` | Update document |
| `GetOrCreateRagIndexes` | Get/create RAG indexes |
| `GetRagIndexes` | List RAG indexes |
| `GetRagIndexOverview` | Get index overview |
| `DeleteRagIndex` | Delete index |
| `RagIndexStatus` | Get index status |

### WhatsApp Integration (6 methods)

WhatsApp call and messaging integration.

| Method | Description |
|--------|-------------|
| `GetWhatsappAccount` | Get WhatsApp account |
| `ListWhatsappAccounts` | List accounts |
| `ImportWhatsappAccount` | Import account |
| `UpdateWhatsappAccount` | Update account |
| `DeleteWhatsappAccount` | Delete account |
| `WhatsappOutboundCall` | Make WhatsApp call |

### Workspace Management (20 methods)

Team, permissions, and workspace settings.

| Method | Description |
|--------|-------------|
| `SearchGroups` | Search user groups |
| `AddMember` | Add group member |
| `RemoveMember` | Remove member |
| `InviteUser` | Invite user |
| `InviteUsersBulk` | Bulk invite |
| `DeleteInvite` | Delete invitation |
| `UpdateWorkspaceMember` | Update member |
| `CreateSecretRoute` | Create secret |
| `GetSecretsRoute` | List secrets |
| `UpdateSecretRoute` | Update secret |
| `DeleteSecretRoute` | Delete secret |
| `GetWorkspaceServiceAccounts` | List service accounts |
| `CreateServiceAccountAPIKey` | Create API key |
| `EditServiceAccountAPIKey` | Edit API key |
| `DeleteServiceAccountAPIKey` | Delete API key |
| `GetServiceAccountAPIKeysRoute` | List API keys |
| `GetSettingsRoute` | Get settings |
| `UpdateSettingsRoute` | Update settings |
| `GetDashboardSettingsRoute` | Get dashboard settings |
| `UpdateDashboardSettingsRoute` | Update dashboard |

### Webhooks (4 methods)

| Method | Description |
|--------|-------------|
| `CreateWorkspaceWebhookRoute` | Create webhook |
| `GetWorkspaceWebhooksRoute` | List webhooks |
| `EditWorkspaceWebhookRoute` | Edit webhook |
| `DeleteWorkspaceWebhookRoute` | Delete webhook |

### MCP / Tools (5 methods)

Model Context Protocol and tool integrations.

| Method | Description |
|--------|-------------|
| `DeleteMcpServerRoute` | Delete MCP server |
| `ListMcpServerToolsRoute` | List MCP tools |
| `GetMcpToolConfigOverrideRoute` | Get tool config |
| `DeleteToolRoute` | Delete tool |
| `GetToolDependentAgentsRoute` | Get dependent agents |

### Audio Native (3 methods)

Audio Native project embedding.

| Method | Description |
|--------|-------------|
| `CreateAudioNativeProject` | Create project |
| `AudioNativeProjectUpdateContentEndpoint` | Update content |
| `GetAudioNativeProjectSettingsEndpoint` | Get settings |

### Transcription (4 methods)

Advanced transcription features.

| Method | Description |
|--------|-------------|
| `GetTranscriptByID` | Get transcript |
| `DeleteTranscriptByID` | Delete transcript |
| `GetDubbedTranscriptFile` | Get dubbed transcript |
| `Translate` | Translate content |

### Miscellaneous (6 methods)

| Method | Description |
|--------|-------------|
| `UsageCharacters` | Get character usage stats |
| `GetSingleUseToken` | Get single-use token |
| `GetResourceMetadata` | Get resource metadata |
| `GetSignedURLDeprecated` | Get signed URL (deprecated) |
| `GetPublicLlmExpectedCostCalculation` | Public LLM cost calc |
| `RedirectToMintlify` | Redirect to docs |

---

## Contributing

Want to help expand SDK coverage? Contributions are welcome! Priority areas:

1. **Knowledge Base / RAG** - Document management for agent context
2. **Voice Library** - Community voice discovery and sharing
3. **Professional Voice Cloning** - Premium voice training features
4. **Conversation History** - Conversation management and analytics
5. **Batch Calls** - Batch call management for agents

See the [Contributing Guide](https://github.com/plexusone/elevenlabs-go/blob/main/CONTRIBUTING.md) for details.
