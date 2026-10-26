package world

import (
	"testing"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
)

func TestAotManager_RemoveScriptSprite(t *testing.T) {
	manager := NewAotManager()

	sprite1 := fileio.ScriptInstrSceEsprOn{Id: 1}
	sprite2 := fileio.ScriptInstrSceEsprOn{Id: 2}

	manager.AddScriptSprite(sprite1)
	manager.AddScriptSprite(sprite2)

	if len(manager.Sprites) != 2 {
		t.Fatalf("expected 2 sprites, got %d", len(manager.Sprites))
	}

	manager.RemoveScriptSprite(1)

	if len(manager.Sprites) != 1 {
		t.Fatalf("expected 1 sprite after removal, got %d", len(manager.Sprites))
	}
	if manager.Sprites[0].Id != 2 {
		t.Errorf("expected remaining sprite id 2, got %d", manager.Sprites[0].Id)
	}
}

func TestAotManager_RemoveScriptSprite_MissingIsNoOp(t *testing.T) {
	manager := NewAotManager()
	manager.AddScriptSprite(fileio.ScriptInstrSceEsprOn{Id: 1})

	manager.RemoveScriptSprite(99) // not present

	if len(manager.Sprites) != 1 {
		t.Errorf("expected 1 sprite, got %d", len(manager.Sprites))
	}
}
