package fileio

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func newStreamReader(data []byte) *StreamReader {
	return NewStreamReader(io.NewSectionReader(bytes.NewReader(data), 0, int64(len(data))))
}

func TestStreamReader_ReadUint8(t *testing.T) {
	sr := newStreamReader([]byte{0xAB})
	got, err := sr.ReadUint8()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0xAB {
		t.Errorf("got %d, want %d", got, 0xAB)
	}
}

func TestStreamReader_ReadUint16(t *testing.T) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, 0x1234)
	sr := newStreamReader(data)
	got, err := sr.ReadUint16()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0x1234 {
		t.Errorf("got %d, want %d", got, 0x1234)
	}
}

func TestStreamReader_ReadUint32(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 0x89ABCDEF)
	sr := newStreamReader(data)
	got, err := sr.ReadUint32()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0x89ABCDEF {
		t.Errorf("got %d, want %d", got, 0x89ABCDEF)
	}
}

func TestStreamReader_ReadInt16(t *testing.T) {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, 0xFFFB) // -5 as int16
	sr := newStreamReader(data)
	got, err := sr.ReadInt16()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != -5 {
		t.Errorf("got %d, want %d", got, -5)
	}
}

func TestStreamReader_ReadFloat32(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 0x3F800000) // 1.0
	sr := newStreamReader(data)
	got, err := sr.ReadFloat32()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 1.0 {
		t.Errorf("got %v, want 1.0", got)
	}
}

func TestStreamReader_ReadBytes(t *testing.T) {
	input := []byte{1, 2, 3, 4}
	sr := newStreamReader(input)
	got, err := sr.ReadBytes(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("got len %d, want 4", len(got))
	}
	for i := range input {
		if got[i] != input[i] {
			t.Errorf("byte %d: got %d, want %d", i, got[i], input[i])
		}
	}
}

func TestStreamReader_ReadBytes_InvalidCount(t *testing.T) {
	sr := newStreamReader([]byte{1, 2, 3})
	if _, err := sr.ReadBytes(0); err == nil {
		t.Fatal("expected error for count 0")
	}
	if _, err := sr.ReadBytes(-1); err == nil {
		t.Fatal("expected error for negative count")
	}
}

func TestStreamReader_ReadString_NullTerminated(t *testing.T) {
	input := []byte{'h', 'i', 0, 'x'}
	sr := newStreamReader(input)
	got, err := sr.ReadString(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hi" {
		t.Errorf("got %q, want %q", got, "hi")
	}
}

func TestStreamReader_ReadFixedString_TrimsPadding(t *testing.T) {
	input := []byte{'a', 'b', 0, 0}
	sr := newStreamReader(input)
	got, err := sr.ReadFixedString(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ab" {
		t.Errorf("got %q, want %q", got, "ab")
	}
}

func TestStreamReader_SetPosition(t *testing.T) {
	input := []byte{10, 20, 30}
	sr := newStreamReader(input)
	sr.SetPosition(1)
	got, err := sr.ReadUint8()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 20 {
		t.Errorf("got %d, want 20", got)
	}
}

func TestStreamReader_ReadBytes_Truncated(t *testing.T) {
	sr := newStreamReader([]byte{1, 2})
	if _, err := sr.ReadBytes(4); err == nil {
		t.Fatal("expected error reading more bytes than available")
	}
}

func TestStreamReader_ReadString_NoTerminator(t *testing.T) {
	input := []byte{'a', 'b', 'c'}
	sr := newStreamReader(input)
	got, err := sr.ReadString(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "abc" {
		t.Errorf("got %q, want %q", got, "abc")
	}
}

func TestStreamReader_ReadString_InvalidMaxLength(t *testing.T) {
	sr := newStreamReader([]byte{'a'})
	if _, err := sr.ReadString(0); err == nil {
		t.Fatal("expected error for maxLength 0")
	}
	if _, err := sr.ReadString(-1); err == nil {
		t.Fatal("expected error for negative maxLength")
	}
}

func TestStreamReader_ReadFixedString_NoPadding(t *testing.T) {
	input := []byte{'x', 'y', 'z', 'w'}
	sr := newStreamReader(input)
	got, err := sr.ReadFixedString(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "xyzw" {
		t.Errorf("got %q, want %q", got, "xyzw")
	}
}

func TestStreamReader_ReadFixedString_AllZeros(t *testing.T) {
	input := []byte{0, 0, 0, 0}
	sr := newStreamReader(input)
	got, err := sr.ReadFixedString(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestStreamReader_SetPosition_ResetToStart(t *testing.T) {
	input := []byte{10, 20, 30}
	sr := newStreamReader(input)
	if _, err := sr.ReadUint8(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sr.SetPosition(0)
	got, err := sr.ReadUint8()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 10 {
		t.Errorf("got %d, want 10", got)
	}
}

func TestStreamReader_SequentialReadsAdvance(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint16(data[0:2], 0x1234)
	binary.LittleEndian.PutUint16(data[2:4], 0x5678)
	sr := newStreamReader(data)
	first, err := sr.ReadUint16()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := sr.ReadUint16()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first != 0x1234 || second != 0x5678 {
		t.Errorf("got (%d, %d), want (0x1234, 0x5678)", first, second)
	}
}
