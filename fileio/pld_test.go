package fileio

import (
	"bytes"
	"testing"
)

func TestLoadPLDStream_EmptyReturnsError(t *testing.T) {
	_, err := LoadPLDStream(bytes.NewReader(nil), 0)
	if err == nil {
		t.Fatal("expected error for empty PLD stream, got nil")
	}
}

func TestLoadPLDStream_TruncatedReturnsError(t *testing.T) {
	// A few bytes are not enough for a complete PLD header/offsets.
	data := []byte{0x01, 0x02, 0x03, 0x04}
	_, err := LoadPLDStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for truncated PLD stream, got nil")
	}
}

func TestLoadPLDStream_GarbageReturnsError(t *testing.T) {
	data := bytes.Repeat([]byte{0xFF}, 64)
	_, err := LoadPLDStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for garbage PLD stream, got nil")
	}
}
