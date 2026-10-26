package fileio

import (
	"os"
	"path/filepath"
	"testing"
)

const testPLDPath = "../data/Pl0/PLD/PL00.PLD"

func TestLoadPLDFile_MissingReturnsError(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "definitely-missing.PLD")
	if _, err := LoadPLDFile(missingPath); err == nil {
		t.Fatal("expected error for missing PLD file, got nil")
	}
}

func TestLoadPLDFile_Success(t *testing.T) {
	if _, err := os.Stat(testPLDPath); err != nil {
		if os.IsNotExist(err) {
			t.Skipf("game data not available at %s: %v", testPLDPath, err)
		}
		t.Fatalf("stat %s: %v", testPLDPath, err)
	}

	output, err := LoadPLDFile(testPLDPath)
	if err != nil {
		t.Fatalf("LoadPLDFile(%q) returned error: %v", testPLDPath, err)
	}
	if output == nil {
		t.Fatal("LoadPLDFile returned nil output")
	}

	if output.MeshData == nil {
		t.Error("expected MeshData to be loaded")
	}
	if output.TextureData == nil {
		t.Error("expected TextureData to be loaded")
	} else {
		if output.TextureData.ImageWidth <= 0 {
			t.Errorf("expected positive texture width, got %d", output.TextureData.ImageWidth)
		}
		if output.TextureData.ImageHeight <= 0 {
			t.Errorf("expected positive texture height, got %d", output.TextureData.ImageHeight)
		}
	}
	if output.AnimationData == nil {
		t.Error("expected AnimationData to be loaded")
	}
	if output.SkeletonData == nil {
		t.Error("expected SkeletonData to be loaded")
	}
}
