package audio

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/generators"
	"github.com/gopxl/beep/speaker"

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

// PlayBGM plays a BGM track by id using MS-ADPCM decoder (no ffmpeg).
func (a *AudioManager) PlayBGM(id int) {
	filename := fmt.Sprintf("data/Common/Sound/bgm/MAIN%02X.SAP", id)
	sap, err := fileio.LoadSAPFile(filename)
	if err != nil {
		fmt.Printf("PlayBGM: failed to read %s: %v\n", filename, err)
		a.playTone(440+float64(id*10), time.Second)
		return
	}

	wave, ok := DecodeSAP(sap.AudioData)
	if !ok {
		fmt.Printf("PlayBGM: failed to decode %s\n", filename)
		a.playTone(440+float64(id*10), time.Second)
		return
	}

	streamer := newPCMStreamer(wave.Samples, wave.Channels)
	buf := beep.NewBuffer(beep.Format{
		SampleRate:  beep.SampleRate(wave.SamplesPerSec),
		NumChannels: wave.Channels,
		Precision:   2,
	})
	buf.Append(streamer)

	fmt.Printf("PlayBGM: playing %s (%d Hz, %d ch, %d samples)\n", filename, wave.SamplesPerSec, wave.Channels, len(wave.Samples))
	speaker.Clear()
	speaker.Play(beep.Loop(-1, buf.Streamer(0, buf.Len())))
}

// StopBGM stops all currently playing audio.
func (a *AudioManager) StopBGM() {
	speaker.Clear()
}

// PlaySFX plays a sound effect using decoded VAG data from a door file.
func (a *AudioManager) PlaySFX(id int) {
	doorPath := "data/Common/Door/DOOR00.DO2"
	do2, err := fileio.LoadDO2File(doorPath)
	if err != nil {
		a.playTone(880, 200*time.Millisecond)
		return
	}
	file, err := os.Open(doorPath)
	if err != nil {
		a.playTone(880, 200*time.Millisecond)
		return
	}
	defer file.Close()
	vabData, err := fileio.LoadVABDataStream(
		io.NewSectionReader(file, do2.DO2FileFormat.VBOffset, do2.DO2FileFormat.VBLength),
		do2.DO2FileFormat.VBLength,
		do2.VABHeaderOutput,
	)
	if err != nil || len(vabData.RawADPCMData) < 2 {
		a.playTone(880, 200*time.Millisecond)
		return
	}

	pcm := DecodeVAG(vabData.RawADPCMData[1])
	streamer := newPCMStreamer(pcm, 1)
	buf := beep.NewBuffer(beep.Format{SampleRate: a.sampleRate, NumChannels: 1, Precision: 2})
	buf.Append(streamer)
	speaker.Play(buf.Streamer(0, buf.Len()))
}

// SFXNames maps sound effect ids observed in scripts to human-readable names.
var SFXNames = map[int]string{
	0x0A: "Unknown SFX 0x0A",
	0x0B: "Unknown SFX 0x0B",
	0x0C: "Unknown SFX 0x0C",
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
