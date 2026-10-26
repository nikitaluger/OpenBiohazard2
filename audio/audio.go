package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	sr := beep.SampleRate(22050)
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

// PlayBGM plays a BGM track by id, converting SAP (ADPCM) to PCM via ffmpeg.
func (a *AudioManager) PlayBGM(id int) {
	filename := fmt.Sprintf("data/Common/Sound/bgm/MAIN%02X.SAP", id)
	streamer, err := a.loadSAPStreamer(filename)
	if err != nil {
		fmt.Printf("PlayBGM: failed to load %s: %v\n", filename, err)
		a.playTone(440+float64(id*10), time.Second)
		return
	}

	fmt.Printf("PlayBGM: playing %s\n", filename)
	speaker.Clear()
	speaker.Play(beep.Loop(-1, streamer))
}

// StopBGM stops all currently playing audio.
func (a *AudioManager) StopBGM() {
	speaker.Clear()
}

// PlaySFX plays a short placeholder tone for a sound effect id.
func (a *AudioManager) PlaySFX(id int) {
	a.playTone(880, 200*time.Millisecond)
}

// SFXNames maps sound effect ids observed in scripts to human-readable names.
var SFXNames = map[int]string{
	0x0A: "Unknown SFX 0x0A",
	0x0B: "Unknown SFX 0x0B",
	0x21: "Unknown SFX 0x21",
}

// SFXName returns a readable name for a sound effect id.
func SFXName(id int) string {
	if name, ok := SFXNames[id]; ok {
		return name
	}
	return "Unknown SFX"
}

func (a *AudioManager) playTone(freq float64, duration time.Duration) {
	sine, err := generators.SineTone(a.sampleRate, freq)
	if err != nil {
		return
	}
	speaker.Play(beep.Take(a.sampleRate.N(duration), sine))
}

func (a *AudioManager) loadSAPStreamer(filename string) (beep.StreamSeeker, error) {
	sap, err := fileio.LoadSAPFile(filename)
	if err != nil {
		return nil, err
	}

	tmpDir := os.TempDir()
	inPath := filepath.Join(tmpDir, "sap_input.wav")
	outPath := filepath.Join(tmpDir, "sap_output.wav")

	if err := os.WriteFile(inPath, sap.AudioData, 0644); err != nil {
		return nil, err
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", inPath, outPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err := os.Open(outPath)
	if err != nil {
		return nil, err
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	buf := beep.NewBuffer(format)
	buf.Append(streamer)
	streamer.Close()
	f.Close()

	return buf.Streamer(0, buf.Len()), nil
}
