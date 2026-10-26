package main

import (
	"fmt"
	"io"
	"os"

	"github.com/OpenBiohazard2/OpenBiohazard2/audio"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	input := "data/Common/Door/DOOR00.DO2"
	file, _ := os.Open(input)
	defer file.Close()
	do2, _ := fileio.LoadDO2File(input)
	vh := do2.VABHeaderOutput
	vabDataReader := io.NewSectionReader(file, do2.DO2FileFormat.VBOffset, do2.DO2FileFormat.VBLength)
	_ = vabDataReader
	vabData, _ := fileio.LoadVABDataStream(vabDataReader, do2.DO2FileFormat.VBLength, vh)

	w := vabData.RawADPCMData[1]
	fmt.Printf("waveform1 len=%d\n", len(w))
	// skip 16-byte lead-in
	if len(w) > 16 {
		w = w[16:]
	}
	fmt.Printf("first block: % x\n", w[:16])

	pcm := audio.DecodeVAG(w)
	fmt.Printf("decoded samples: %d\n", len(pcm))
	min, max := int16(32767), int16(-32768)
	for _, s := range pcm {
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
	}
	fmt.Printf("sample range: %d to %d\n", min, max)

	f, _ := os.Create("vag_test_carry.wav")
	audio.DecodeVAGToWAV(w, 22050, f)
	f.Close()
	fmt.Println("Wrote vag_test_carry.wav")
}
