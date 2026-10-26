package audio

// AudioManager is a placeholder for future audio playback support.
type AudioManager struct{}

// NewAudioManager creates a new audio manager.
func NewAudioManager() *AudioManager {
	return &AudioManager{}
}

// PlayBGM is a placeholder for background music playback.
func (a *AudioManager) PlayBGM(id int) {
	// TODO: implement using an audio library
}

// StopBGM is a placeholder for stopping background music.
func (a *AudioManager) StopBGM() {
	// TODO: implement
}

// PlaySFX is a placeholder for sound effect playback.
func (a *AudioManager) PlaySFX(id int) {
	// TODO: implement
}
