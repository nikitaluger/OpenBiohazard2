package fileio

import (
	"bytes"
	"testing"
)

func buildMinimalVABHeader() []byte {
	data := make([]byte, 32+128*16+0*16*32+256*2)
	// Magic "pBAV"
	copy(data[0:4], []byte("pBAV"))
	// Version, Id, Filesize, Reserved0 can stay zero.
	// ProgramCount, ToneCount, WaveformCount = 0
	// AudioSizes array (256+1 uint16) stays zero.
	return data
}

func TestLoadVABHeaderStream_ValidMinimal(t *testing.T) {
	data := buildMinimalVABHeader()
	out, err := LoadVABHeaderStream(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil {
		t.Fatal("expected non-nil output")
	}
	if string(out.VABHeader.Magic[:]) != "pBAV" {
		t.Errorf("magic = %s, want pBAV", string(out.VABHeader.Magic[:]))
	}
}

func TestLoadVABHeaderStream_InvalidMagic(t *testing.T) {
	data := buildMinimalVABHeader()
	data[0] = 'X'
	_, err := LoadVABHeaderStream(bytes.NewReader(data), int64(len(data)))
	if err == nil {
		t.Fatal("expected error for invalid VAB magic, got nil")
	}
}

func TestLoadVABHeaderStream_EmptyReturnsError(t *testing.T) {
	_, err := LoadVABHeaderStream(bytes.NewReader(nil), 0)
	if err == nil {
		t.Fatal("expected error for empty VAB stream, got nil")
	}
}
