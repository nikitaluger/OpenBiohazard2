package fileio

import (
	"bytes"
	"testing"
)

func TestLoadESPStream_EmptyReturnsError(t *testing.T) {
	_, err := LoadESPStream(bytes.NewReader(nil), 0, -4)
	if err == nil {
		t.Fatal("expected error for empty ESP stream, got nil")
	}
}

func TestLoadESPStream_TruncatedReturnsError(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x04}
	_, err := LoadESPStream(bytes.NewReader(data), int64(len(data)), int64(len(data))-4)
	if err == nil {
		t.Fatal("expected error for truncated ESP stream, got nil")
	}
}
