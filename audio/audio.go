package audio

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

// AudioManager handles basic audio playback.
type AudioManager struct {
	sampleRate beep.SampleRate

	// Volume levels observed from BGM control commands.
	mainVolume  float64
	prog0Volume float64
	prog1Volume float64
	prog2Volume float64
}

// NewAudioManager initializes the audio system.
func NewAudioManager() *AudioManager {
	sr := beep.SampleRate(44100)
	_ = speaker.Init(sr, 1024)
	return &AudioManager{
		sampleRate:  sr,
		mainVolume:  1.0,
		prog0Volume: 1.0,
		prog1Volume: 1.0,
		prog2Volume: 1.0,
	}
}

// BGMNames maps BGM ids observed in scripts to human-readable names.
var BGMNames = map[int]string{
	0: "Main",
	1: "Track 1",
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

// SetVolume stores a volume level for a BGM volume target.
func (a *AudioManager) SetVolume(target uint8, left, right uint8) {
	level := (float64(left) + float64(right)) / 2.0 / 127.0
	switch target {
	case 0:
		a.mainVolume = level
	case 1:
		a.prog0Volume = level
	case 2:
		a.prog1Volume = level
	case 3:
		a.prog2Volume = level
	}
}

// PlayBGM plays a BGM track by id using SAP files.
func (a *AudioManager) PlayBGM(id int) {
	filename := fmt.Sprintf("data/Common/Sound/bgm/MAIN%02X.SAP", id)
	sap, err := fileio.LoadSAPFile(filename)
	if err != nil {
		// Fall back to a placeholder tone if the file is missing.
		a.playTone(440+float64(id*10), time.Second)
		return
	}

	streamer, format, err := wav.Decode(bytes.NewReader(sap.AudioData))
	if err != nil {
		a.playTone(440+float64(id*10), time.Second)
		return
	}

	speaker.Clear()
	speaker.Play(beep.Loop(-1, streamer))
	_ = format
}

// StopBGM stops all currently playing audio.
func (a *AudioManager) StopBGM() {
	speaker.Clear()
}

// PlaySFX plays a short placeholder tone for a sound effect id.
func (a *AudioManager) PlaySFX(id int) {
	a.playTone(880, 200*time.Millisecond)
}

func (a *AudioManager) playTone(freq float64, duration time.Duration) {
	sine, err := generators.SineTone(a.sampleRate, freq)
	if err != nil {
		return
	}
	speaker.Play(beep.Take(a.sampleRate.N(duration), sine))
}
