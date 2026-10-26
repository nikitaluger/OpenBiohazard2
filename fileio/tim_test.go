package fileio

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func buildMinimalTIM4BPP() []byte {
	var buf bytes.Buffer

	// TIM header
	binary.Write(&buf, binary.LittleEndian, uint32(16)) // Magic
	binary.Write(&buf, binary.LittleEndian, uint32(8))  // BPP (4-bit)
	binary.Write(&buf, binary.LittleEndian, uint32(0))  // Offset (unused)
	binary.Write(&buf, binary.LittleEndian, uint16(0))  // OriginX
	binary.Write(&buf, binary.LittleEndian, uint16(0))  // OriginY
	binary.Write(&buf, binary.LittleEndian, uint16(16)) // NumColors
	binary.Write(&buf, binary.LittleEndian, uint16(1))  // NumCluts

	// One palette with 16 colors (all zero is acceptable for parsing)
	for i := 0; i < 16; i++ {
		binary.Write(&buf, binary.LittleEndian, uint16(0))
	}

	// Image header
	binary.Write(&buf, binary.LittleEndian, uint32(0)) // Size (unused)
	binary.Write(&buf, binary.LittleEndian, uint16(0))  // OriginX
	binary.Write(&buf, binary.LittleEndian, uint16(0))  // OriginY
	binary.Write(&buf, binary.LittleEndian, uint16(1))  // Width (pixels/4)
	binary.Write(&buf, binary.LittleEndian, uint16(1))  // Height

	// Image data: (4*1)/2 = 2 bytes
	buf.Write([]byte{0x00, 0x00})

	return buf.Bytes()
}

func TestLoadTIMStream_Valid4BPP(t *testing.T) {
	data := buildMinimalTIM4BPP()
	out, err := LoadTIMStream(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil output")
	}
	if out.ImageWidth != 4 {
		t.Errorf("ImageWidth = %d, want 4", out.ImageWidth)
	}
	if out.ImageHeight != 1 {
		t.Errorf("ImageHeight = %d, want 1", out.ImageHeight)
	}
	if out.NumPalettes != 1 {
		t.Errorf("NumPalettes = %d, want 1", out.NumPalettes)
	}
}

func TestLoadTIMStream_InvalidMagic(t *testing.T) {
	data := buildMinimalTIM4BPP()
	// Corrupt the magic number (first 4 bytes)
	data[0] = 0xFF
	data[1] = 0xFF
	data[2] = 0xFF
	data[3] = 0xFF

	_, err := LoadTIMStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for invalid TIM magic, got nil")
	}
}

func TestLoadTIMStream_EmptyReturnsError(t *testing.T) {
	_, err := LoadTIMStream(bytes.NewReader(nil), 0)
	if err == nil {
		t.Fatal("expected error for empty TIM stream, got nil")
	}
}
