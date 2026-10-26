package fileio

import (
	"bytes"
	"testing"
)

func TestLoadEMDStream_EmptyReturnsError(t *testing.T) {
	_, err := LoadEMDStream(bytes.NewReader(nil), 0)
	if err == nil {
		t.Fatal("expected error for empty EMD stream, got nil")
	}
}

func TestLoadEMDStream_TruncatedReturnsError(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	_, err := LoadEMDStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for truncated EMD stream, got nil")
	}
}

func TestLoadEMDStream_GarbageReturnsError(t *testing.T) {
	data := bytes.Repeat([]byte{0xFF}, 64)
	_, err := LoadEMDStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for garbage EMD stream, got nil")
	}
}
