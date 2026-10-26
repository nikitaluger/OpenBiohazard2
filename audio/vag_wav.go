package audio

import (
	"io"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/wav"
)

// pcmStreamer streams 16-bit mono PCM samples.
type pcmStreamer struct {
	samples []int16
	pos     int
}

func (p *pcmStreamer) Stream(samples [][2]float64) (int, bool) {
	n := 0
	for i := range samples {
		if p.pos+1 >= len(p.samples) {
			break
		}
		vL := float64(p.samples[p.pos]) / 32768.0
		vR := float64(p.samples[p.pos+1]) / 32768.0
		samples[i][0] = vL
		samples[i][1] = vR
		p.pos += 2
		n++
	}
	return n, p.pos < len(p.samples)
}

func (p *pcmStreamer) Err() error { return nil }

// DecodeVAGToWAV writes decoded VAG data as a WAV file.
func DecodeVAGToWAV(data []byte, sampleRate int, w io.WriteSeeker) error {
	pcm := DecodeVAG(data)
	return writePCMToWAV(pcm, sampleRate, w)
}

// DecodeVAGResetToWAV writes reset-state decoded VAG data as a WAV file.
func DecodeVAGResetToWAV(data []byte, sampleRate int, w io.WriteSeeker) error {
	pcm := DecodeVAGReset(data)
	return writePCMToWAV(pcm, sampleRate, w)
}

func writePCMToWAV(pcm []int16, sampleRate int, w io.WriteSeeker) error {
	sr := beep.SampleRate(sampleRate)

	streamer := &pcmStreamer{samples: pcm}
	format := beep.Format{
		SampleRate:  sr,
		NumChannels: 1,
		Precision:   2,
	}

	buf := beep.NewBuffer(format)
	buf.Append(streamer)

	return wav.Encode(w, buf.Streamer(0, buf.Len()), format)
}
