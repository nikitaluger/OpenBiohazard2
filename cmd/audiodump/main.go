package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenBiohazard2/OpenBiohazard2/audio"
	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func main() {
	root := "data/Common/Sound"
	outDir := "audio_dump"
	os.MkdirAll(outDir, 0755)

	count := 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".sap") {
			return nil
		}
		sap, err := fileio.LoadSAPFile(path)
		if err != nil {
			fmt.Printf("skip %s: %v\n", path, err)
			return nil
		}
		outName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".wav"
		outPath := filepath.Join(outDir, outName)
		f, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("skip %s: %v\n", path, err)
			return nil
		}
		defer f.Close()
		if err := audio.DecodeSAPToWAV(sap.AudioData, f); err != nil {
			fmt.Printf("skip %s: %v\n", path, err)
			return nil
		}
		count++
		fmt.Printf("wrote %s\n", outPath)
		return nil
	})
	fmt.Printf("Done: %d files\n", count)
}
