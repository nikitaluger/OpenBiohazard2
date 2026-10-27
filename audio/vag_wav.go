package audio

import (
	"io"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/wav"
)

// pcmStreamer streams multi-channel PCM samples stored interleaved.
type pcmStreamer struct {
	samples  []int16
	channels int
	pos      int
}

// newPCMStreamer creates a channel-aware PCM streamer.
func newPCMStreamer(samples []int16, channels int) *pcmStreamer {
	if channels < 1 {
		channels = 1
	}
	if channels > 2 {
		channels = 2
	}
	return &pcmStreamer{samples: samples, channels: channels}
}

func (p *pcmStreamer) Stream(samples [][2]float64) (int, bool) {
	if p.channels < 1 {
		p.channels = 1
	}
	if p.channels > 2 {
		p.channels = 2
	}
	n := 0
	for i := range samples {
		if p.pos+p.channels > len(p.samples) {
			// Skip any trailing partial frame and stop.
			p.pos = len(p.samples)
			break
		}
		for c := 0; c < p.channels; c++ {
			v := float64(p.samples[p.pos+c]) / 32768.0
			if p.channels == 1 {
				samples[i][0] = v
				samples[i][1] = v
			} else {
				samples[i][c] = v
			}
		}
		p.pos += p.channels
		n++
	}
	return n, p.pos < len(p.samples)
}

func (p *pcmStreamer) Err() error { return nil }

// DecodeVAGToWAV writes decoded VAG data as a WAV file.
func DecodeVAGToWAV(data []byte, sampleRate int, w io.WriteSeeker) error {
	pcm := DecodeVAG(data)
	return writePCMToWAV(pcm, 1, sampleRate, w)
}

// DecodeVAGResetToWAV writes reset-state decoded VAG data as a WAV file.
func DecodeVAGResetToWAV(data []byte, sampleRate int, w io.WriteSeeker) error {
	pcm := DecodeVAGReset(data)
	return writePCMToWAV(pcm, 1, sampleRate, w)
}

func writePCMToWAV(pcm []int16, channels, sampleRate int, w io.WriteSeeker) error {
	sr := beep.SampleRate(sampleRate)

	streamer := newPCMStreamer(pcm, channels)
	format := beep.Format{
		SampleRate:  sr,
		NumChannels: channels,
		Precision:   2,
	}

	buf := beep.NewBuffer(format)
	buf.Append(streamer)

	return wav.Encode(w, buf.Streamer(0, buf.Len()), format)
}
