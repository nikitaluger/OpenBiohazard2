package fileio

import (
	"encoding/binary"
	"testing"
)

func TestReadInt8Extra(t *testing.T) {
	data := []byte{0xFF} // -1 as int8
	reader := createTestStreamReader(data)

	value, err := reader.ReadInt8()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != -1 {
		t.Errorf("Expected -1, got %d", value)
	}
}

func TestReadInt32Extra(t *testing.T) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, 0xFFFFFFFE) // -2 as int32
	reader := createTestStreamReader(data)

	value, err := reader.ReadInt32()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != -2 {
		t.Errorf("Expected -2, got %d", value)
	}
}

func TestReadFloat64Extra(t *testing.T) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, 0x3FF0000000000000) // 1.0
	reader := createTestStreamReader(data)

	value, err := reader.ReadFloat64()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != 1.0 {
		t.Errorf("Expected 1.0, got %v", value)
	}
}

func TestReadBytesInvalidCountExtra(t *testing.T) {
	reader := createTestStreamReader([]byte{1, 2, 3})

	if _, err := reader.ReadBytes(0); err == nil {
		t.Error("Expected error for count 0")
	}
	if _, err := reader.ReadBytes(-1); err == nil {
		t.Error("Expected error for negative count")
	}
}
