package script

import (
	"bytes"
	"encoding/binary"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/OpenBiohazard2/OpenBiohazard2/game"
)

// PLC commands are used for 3D model animation

func (scriptDef *ScriptDef) ScriptPlcMotion(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcMotion{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if gameDef != nil && gameDef.Player != nil {
		gameDef.Player.Motion = int(instruction.MoveNumber)
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcNeck(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcNeck{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if gameDef != nil && gameDef.Player != nil {
		gameDef.Player.NeckX = instruction.NeckX
		gameDef.Player.NeckY = instruction.NeckY
		gameDef.Player.NeckZ = instruction.NeckZ
	}

	return 1
}

func (scriptDef *ScriptDef) ScriptPlcDest(lineData []byte, gameDef *game.GameDef) int {
	byteArr := bytes.NewBuffer(lineData)
	instruction := fileio.ScriptInstrPlcDest{}
	binary.Read(byteArr, binary.LittleEndian, &instruction)

	if gameDef != nil && gameDef.Player != nil {
		gameDef.Player.DestX = instruction.DestX
		gameDef.Player.DestZ = instruction.DestZ
	}

	return 1
}
