package audio

// DecodeVAG decodes PlayStation SPU ADPCM (VAG) data into 16-bit PCM samples.
// This variant uses low nibble first and carries predictor state across blocks.
func DecodeVAG(data []byte) []int16 {
	// Some VAB waveforms appear to have a 16-byte lead-in.
	if len(data) > 16 {
		data = data[16:]
	}
	return decodeVAG(data, false)
}

// DecodeVAGReset is an experimental variant that resets predictor state per block.
func DecodeVAGReset(data []byte) []int16 {
	return decodeVAG(data, true)
}

func decodeVAG(data []byte, resetPerBlock bool) []int16 {
	samples := make([]int16, 0, (len(data)/16)*28)

	filterCoeff := [5][2]int{
		{0, 0},
		{60, 0},
		{115, -52},
		{98, -55},
		{122, -60},
	}

	prev1, prev2 := 0, 0

	for offset := 0; offset+16 <= len(data); offset += 16 {
		block := data[offset : offset+16]

		if block[1] == 0x07 {
			break
		}

		shift := int(block[0] & 0x0F)
		filter := int(block[0] >> 4)
		if filter > 4 {
			filter = 0
		}
		f0 := filterCoeff[filter][0]
		f1 := filterCoeff[filter][1]

		if resetPerBlock {
			prev1, prev2 = 0, 0
		}

		for i := 0; i < 28; i++ {
			var nibble int
			if i%2 == 0 {
				nibble = int(block[2+i/2] & 0x0F)
			} else {
				nibble = int(block[2+i/2] >> 4)
			}
			if nibble >= 8 {
				nibble -= 16
			}

			sample := nibble << 12
			if shift <= 12 {
				sample >>= shift
			} else {
				sample <<= (shift - 12)
			}

			sample += (prev1*f0 + prev2*f1 + 32) / 64

			if sample > 32767 {
				sample = 32767
			}
			if sample < -32768 {
				sample = -32768
			}

			samples = append(samples, int16(sample))
			prev2 = prev1
			prev1 = sample
		}
	}

	return samples
}
