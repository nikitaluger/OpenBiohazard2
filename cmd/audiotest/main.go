package main

import (
	"fmt"
	"time"

	"github.com/OpenBiohazard2/OpenBiohazard2/audio"
)

func main() {
	fmt.Println("Playing BGM 0...")
	am := audio.NewAudioManager()
	am.PlayBGM(0)
	time.Sleep(10 * time.Second)
	fmt.Println("Done")
}
