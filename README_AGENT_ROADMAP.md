# Real-Time Voice Agent Roadmap

This document outlines the API priorities and implementation roadmap for supporting real-time voice agents that participate in phone calls or virtual meetings via services like Twilio or Zoom, powered by LLMs or services like Google Dialogflow.

## Current Coverage vs. Real-Time Agent Needs

| Capability | Current Status | Priority for Agents |
|------------|----------------|---------------------|
| HTTP TTS | ✅ Covered | Medium (too slow for real-time) |
| HTTP STT | ✅ Covered | Medium (too slow for real-time) |
| **WebSocket TTS** | ✅ Covered | **Critical** |
| **WebSocket STT** | ✅ Covered | **Critical** |
| **Speech-to-Speech** | ✅ Covered | High |
| **Twilio Integration** | ✅ Covered | **Critical** for phone |
| Conversational AI Agents | ❌ Not covered | Optional* |

*Optional if building your own orchestration with Dialogflow/LLM

## Priority APIs

### 1. WebSocket Text-to-Speech (Critical)

Real-time TTS with sub-100ms latency for conversational agents.

- **Endpoint**: `wss://api.elevenlabs.io/v1/text-to-speech/{voice_id}/stream-input`
- **Features**:
  - Input streaming (send text chunks, receive audio in real-time)
  - Multi-context sessions (multiple conversation contexts per connection)
  - Word-level timestamps/alignment
  - Pronunciation dictionary support
- **Docs**: https://elevenlabs.io/docs/api-reference/text-to-speech/v-1-text-to-speech-voice-id-stream-input

### 2. WebSocket Speech-to-Text (Critical)

Real-time transcription for capturing caller speech.

- **Endpoint**: `wss://api.elevenlabs.io/v1/speech-to-text/realtime`
- **Features**:
  - Streaming audio input
  - Real-time text output with word timestamps
  - Speaker diarization
- **Docs**: https://elevenlabs.io/docs/api-reference/speech-to-text/v-1-speech-to-text-realtime

### 3. Twilio/Phone Integration (Critical for Phone)

Connect agents to phone calls.

| Endpoint | Purpose |
|----------|---------|
| `POST /v1/convai/twilio/register-call` | Register incoming Twilio call |
| `POST /v1/convai/twilio/outbound-call` | Make outbound call via Twilio |
| `POST /v1/convai/sip-trunk/outbound-call` | Make call via SIP trunk |
| Phone number management | List, create, update, delete phone numbers |

**Docs**: https://elevenlabs.io/docs/agents-platform/phone-numbers/twilio-integration/native-integration

### 4. Speech-to-Speech (High Priority)

Voice conversion in real-time - useful for consistent agent voice.

| Method | Purpose |
|--------|---------|
| `SpeechToSpeechFull` | Convert audio to target voice |
| `SpeechToSpeechStream` | Streaming voice conversion |

## Architecture Options

### Option A: Use ElevenLabs Conversational AI Platform

ElevenLabs provides full orchestration including STT, LLM routing, and TTS:

```
Caller → Twilio → ElevenLabs Agent → (STT → LLM → TTS) → Caller
```

**Requirements**: Agents Platform APIs (26 methods)

### Option B: Build Your Own Orchestration (Recommended)

You control the logic using Dialogflow, OpenAI, or your own LLM. ElevenLabs provides the voice layer:

```
Caller → Twilio → Your Server → Dialogflow/LLM
                      ↓
              ElevenLabs WebSocket
              (STT + TTS streams)
```

**Requirements**:
- WebSocket TTS
- WebSocket STT
- Twilio call registration
- Speech-to-Speech (optional)

## Implementation Status

### Phase 1: Core Real-Time APIs ✅ Complete

- [x] `websockettts.go` - WebSocket TTS streaming
- [x] `websocketstt.go` - WebSocket STT streaming

### Phase 2: Phone Integration ✅ Complete

- [x] `twilio.go` - Twilio call registration, outbound calls, and SIP trunk integration
- [x] Phone number management (List, Get, Update, Delete)

### Phase 3: Voice Processing ✅ Complete

- [x] `speechtospeech.go` - Voice conversion (HTTP + streaming)

### Phase 4: Full Agent Platform (Future)

- [ ] `agents.go` - Agent CRUD operations
- [ ] `conversations.go` - Conversation management
- [ ] `knowledgebase.go` - RAG/Knowledge base

## Usage Examples

### WebSocket TTS for Real-Time Response

```go
// Connect to WebSocket TTS
ws, err := client.WebSocketTTS().Connect(ctx, voiceID, &elevenlabs.WebSocketTTSOptions{
    ModelID: "eleven_turbo_v2_5",
    OutputFormat: "pcm_16000",
})
if err != nil {
    log.Fatal(err)
}
defer ws.Close()

// Stream text from LLM, receive audio chunks
for chunk := range llmResponseStream {
    audioChunks, err := ws.SendText(ctx, chunk)
    if err != nil {
        log.Fatal(err)
    }
    for audio := range audioChunks {
        // Send audio to Twilio/caller
        twilioStream.Write(audio)
    }
}

// Flush remaining audio
ws.Flush(ctx)
```

### WebSocket STT for Real-Time Transcription

```go
// Connect to WebSocket STT
ws, err := client.WebSocketSTT().Connect(ctx, &elevenlabs.WebSocketSTTOptions{
    ModelID: "scribe_v1",
    SampleRate: 16000,
})
if err != nil {
    log.Fatal(err)
}
defer ws.Close()

// Stream audio from caller, receive transcription
go func() {
    for audioChunk := range callerAudioStream {
        ws.SendAudio(ctx, audioChunk)
    }
    ws.EndStream(ctx)
}()

for transcript := range ws.Transcripts() {
    // Send to Dialogflow/LLM for processing
    llmResponse := dialogflow.DetectIntent(transcript.Text)
    // Generate voice response...
}
```

### Twilio Call Registration

```go
// Register incoming Twilio call with ElevenLabs
twiml, err := client.Twilio().RegisterCall(ctx, &elevenlabs.TwilioRegisterCallRequest{
    AgentID: "your-agent-id",
    // or use custom handling
})

// Return TwiML to Twilio webhook
w.Header().Set("Content-Type", "application/xml")
w.Write([]byte(twiml))
```

## References

- [ElevenLabs Conversational AI](https://elevenlabs.io/conversational-ai)
- [WebSocket TTS Documentation](https://elevenlabs.io/docs/developers/websockets)
- [WebSocket STT Documentation](https://elevenlabs.io/docs/api-reference/speech-to-text/v-1-speech-to-text-realtime)
- [Twilio Integration](https://elevenlabs.io/docs/agents-platform/phone-numbers/twilio-integration/native-integration)
- [Agents Platform Overview](https://elevenlabs.io/docs/agents-platform/overview)
- [Multi-Context WebSocket](https://elevenlabs.io/docs/cookbooks/multi-context-web-socket)
