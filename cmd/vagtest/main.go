package main

import (
	"fmt"
	"io"
	"os"

	"github.com/OpenBiohazard2/OpenBiohazard2/audio"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: vagtest <door.do2>")
		os.Exit(1)
	}
	input := os.Args[1]
	file, _ := os.Open(input)
	defer file.Close()
	fi, _ := file.Stat()
	do2, _ := fileio.LoadDO2File(input)
	vh := do2.VABHeaderOutput
	vabDataReader := io.NewSectionReader(file, do2.DO2FileFormat.VBOffset, do2.DO2FileFormat.VBLength)
	vabData, _ := fileio.LoadVABDataStream(vabDataReader, do2.DO2FileFormat.VBLength, vh)

	fmt.Printf("VAB: programs=%d tones=%d waveforms=%d\n",
		vh.VABHeader.ProgramCount, vh.VABHeader.ToneCount, vh.VABHeader.WaveformCount)

	for i, w := range vabData.RawADPCMData {
		outName := fmt.Sprintf("vag_test_%d.wav", i)
		f, _ := os.Create(outName)
		audio.DecodeVAGToWAV(w, 22050, f)
		f.Close()
		fmt.Printf("Wrote %s (%d bytes raw)\n", outName, len(w))
	}
	_ = fi
}
