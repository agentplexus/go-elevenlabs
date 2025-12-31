package elevenlabs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// PCMToWAV wraps raw PCM audio data in a WAV header.
// ElevenLabs PCM is 16-bit signed little-endian mono.
//
// Usage:
//
//	resp, _ := client.TextToSpeech().Generate(ctx, &TTSRequest{
//	    VoiceID:      voiceID,
//	    Text:         "Hello",
//	    OutputFormat: "pcm_44100",
//	})
//	wavData, _ := elevenlabs.PCMToWAV(resp.Audio, 44100)
func PCMToWAV(pcmData io.Reader, sampleRate int) ([]byte, error) {
	pcm, err := io.ReadAll(pcmData)
	if err != nil {
		return nil, fmt.Errorf("read PCM data: %w", err)
	}

	return PCMBytesToWAV(pcm, sampleRate)
}

// PCMBytesToWAV wraps raw PCM bytes in a WAV header.
func PCMBytesToWAV(pcm []byte, sampleRate int) ([]byte, error) {
	const (
		numChannels   = 1  // mono
		bitsPerSample = 16 // 16-bit
		maxUint32     = 1<<32 - 1
	)

	if sampleRate <= 0 || sampleRate > maxUint32 {
		return nil, fmt.Errorf("invalid sample rate: %d", sampleRate)
	}

	byteRate := sampleRate * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	dataSize := len(pcm)
	fileSize := 36 + dataSize // 44 byte header - 8 bytes for RIFF header

	if fileSize > maxUint32 || byteRate > maxUint32 || dataSize > maxUint32 {
		return nil, fmt.Errorf("audio data too large for WAV format")
	}

	// Convert to uint32 after validation (safe due to bounds checks above)
	fileSize32 := uint32(fileSize)     //nolint:gosec // validated above
	sampleRate32 := uint32(sampleRate) //nolint:gosec // validated above
	byteRate32 := uint32(byteRate)     //nolint:gosec // validated above
	dataSize32 := uint32(dataSize)     //nolint:gosec // validated above

	buf := new(bytes.Buffer)

	// writeBinary panics on error; bytes.Buffer writes cannot fail
	writeBinary := func(data any) {
		if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
			panic(fmt.Sprintf("binary.Write to bytes.Buffer failed: %v", err))
		}
	}

	// RIFF header
	buf.WriteString("RIFF")
	writeBinary(fileSize32)
	buf.WriteString("WAVE")

	// fmt subchunk
	buf.WriteString("fmt ")
	writeBinary(uint32(16)) // subchunk size
	writeBinary(uint16(1))  // audio format (1 = PCM)
	writeBinary(uint16(numChannels))
	writeBinary(sampleRate32)
	writeBinary(byteRate32)
	writeBinary(uint16(blockAlign))
	writeBinary(uint16(bitsPerSample))

	// data subchunk
	buf.WriteString("data")
	writeBinary(dataSize32)
	buf.Write(pcm)

	return buf.Bytes(), nil
}

// ParsePCMSampleRate extracts the sample rate from a PCM format string.
// Example: "pcm_44100" returns 44100.
func ParsePCMSampleRate(format string) (int, error) {
	if !strings.HasPrefix(format, "pcm_") {
		return 0, fmt.Errorf("not a PCM format: %s", format)
	}
	rateStr := strings.TrimPrefix(format, "pcm_")
	rate, err := strconv.Atoi(rateStr)
	if err != nil {
		return 0, fmt.Errorf("invalid sample rate in format %s: %w", format, err)
	}
	return rate, nil
}
