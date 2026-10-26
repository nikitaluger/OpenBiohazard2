package fileio

import (
	"bytes"
	"testing"
)

func TestLoadMD1Stream_EmptyReturnsError(t *testing.T) {
	_, err := LoadMD1Stream(bytes.NewReader(nil), 0)
	if err == nil {
		t.Fatal("expected error for empty MD1 stream, got nil")
	}
}

func TestLoadMD1Stream_TruncatedReturnsError(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	_, err := LoadMD1Stream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for truncated MD1 stream, got nil")
	}
}

func TestLoadMD1Stream_GarbageReturnsError(t *testing.T) {
	data := bytes.Repeat([]byte{0xFF}, 64)
	_, err := LoadMD1Stream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for garbage MD1 stream, got nil")
	}
}
