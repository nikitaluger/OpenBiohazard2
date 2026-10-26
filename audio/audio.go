package audio

import (
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/speaker"
)

// AudioManager wraps basic audio playback using beep.
type AudioManager struct {
	sampleRate beep.SampleRate
}

// NewAudioManager initializes the audio system.
func NewAudioManager() *AudioManager {
	sr := beep.SampleRate(44100)
	_ = speaker.Init(sr, 1024)
	return &AudioManager{sampleRate: sr}
}

// BGMNames maps BGM ids observed in scripts to human-readable names.
var BGMNames = map[int]string{
	0: "Main",
	1: "Track 1 (unknown)",
}

// BGMTypeName describes the volume target type.
func BGMTypeName(t uint8) string {
	switch t {
	case 0:
		return "MAIN_VOL"
	case 1:
		return "PROG0_VOL"
	case 2:
		return "PROG1_VOL"
	case 3:
		return "PROG2_VOL"
	default:
		return "UNKNOWN"
	}
}
// PlayBGM plays a short placeholder tone for a BGM id.
func (a *AudioManager) PlayBGM(id int) {
	freq := 440.0 + float64(id*10)
	sine, err := generators.SineTone(a.sampleRate, freq)
	if err != nil {
		return
	}
	speaker.Play(beep.Take(a.sampleRate.N(time.Second), sine))
}

// StopBGM stops all currently playing audio.
func (a *AudioManager) StopBGM() {
	speaker.Clear()
}

// PlaySFX plays a short placeholder tone for a sound effect id.
func (a *AudioManager) PlaySFX(id int) {
	sine, err := generators.SineTone(a.sampleRate, 880)
	if err != nil {
		return
	}
	speaker.Play(beep.Take(a.sampleRate.N(200*time.Millisecond), sine))
}
