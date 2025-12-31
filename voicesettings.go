package elevenlabs

// Voice settings presets for different platforms and use cases.
//
// These presets are tuned for specific content types and platforms.
// Adjust as needed for your specific voice and content style.

// VoiceSettingsForUdemy returns settings tuned for Udemy courses.
// Neutral, clear, consistent, safe for long lectures.
func VoiceSettingsForUdemy() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Style:           0.05,
		Speed:           1.0,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForCoursera returns settings tuned for Coursera courses.
// Slightly expressive, engaging for mixed media content.
func VoiceSettingsForCoursera() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.7,
		SimilarityBoost: 0.85,
		Style:           0.2,
		Speed:           1.0,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForEdX returns settings tuned for edX courses.
// Very stable, highly intelligible, slightly faster for dense academic content.
func VoiceSettingsForEdX() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.8,
		SimilarityBoost: 0.9,
		Style:           0.15,
		Speed:           1.05,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForInstagram returns settings tuned for Instagram content.
// Energetic but polished, suitable for brand content.
func VoiceSettingsForInstagram() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.4,
		SimilarityBoost: 0.85,
		Style:           0.35,
		Speed:           1.1,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForTikTok returns settings tuned for TikTok content.
// Designed for immediate engagement in the first 1-3 seconds.
func VoiceSettingsForTikTok() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.3,
		SimilarityBoost: 0.85,
		Style:           0.45,
		Speed:           1.15,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForYouTube returns settings tuned for YouTube content.
// Designed to hold attention for 5-20 minutes without sounding robotic or theatrical.
func VoiceSettingsForYouTube() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.45,
		SimilarityBoost: 0.8,
		Style:           0.2,
		Speed:           1.05,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForPodcast returns settings tuned for podcast content.
// Natural conversational tone for long-form audio content.
func VoiceSettingsForPodcast() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.55,
		SimilarityBoost: 0.75,
		Style:           0.15,
		Speed:           1.0,
		UseSpeakerBoost: true,
	}
}

// VoiceSettingsForAudiobook returns settings tuned for audiobook narration.
// Clear, consistent, easy to listen to for extended periods.
func VoiceSettingsForAudiobook() *VoiceSettings {
	return &VoiceSettings{
		Stability:       0.65,
		SimilarityBoost: 0.8,
		Style:           0.1,
		Speed:           0.95,
		UseSpeakerBoost: true,
	}
}
