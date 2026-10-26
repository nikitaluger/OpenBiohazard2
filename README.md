# OpenBiohazard2
Open source re-implementation of the original Resident Evil 2 engine written in Go and OpenGL. You must own a copy of the original game.

<div style="display:inline-block;">
<img src="https://github.com/OpenBiohazard2/OpenBiohazard2/raw/master/screenshots/beginning.png" alt="beginning" width="400" height="300" />
<img src="https://github.com/OpenBiohazard2/OpenBiohazard2/raw/master/screenshots/inventory.png" alt="inventory" width="400" height="300" />
</div>

### Installation

1. Clone this project.
2. Get the game data from your installed location. Copy all the files to the `data/` folder in this repository.
3. Run `go build`.

### Data setup notes

The repository expects the original game data to be available under `data/`.
Some game installs (including many GOG copies) use slightly different folder
names. If file validation fails, check the following:

- `data/Pl0/RDP/` is expected. Some installs provide `data/Pl0/Rdu/`.
- `data/Common/DATP/` is expected. Some installs provide `data/Common/Datu/`.

On Windows, you can create junctions instead of copying or renaming folders:

```powershell
cd data
mklink /J "Pl0\RDP" "Pl0\Rdu"
mklink /J "Common\DATP" "Common\Datu"
```

Windows is case-insensitive, but folder names still need to match the expected names above.

### Task list

- [ ] Audio
  - [ ] Background music
  - [ ] Core sound
- [ ] Game
  - [x] Collision detection
  - [x] Event triggers
  - [ ] Inventory system
  - [ ] Enemy AI
  - [ ] Puzzles
  - [ ] Door transitions
- [ ] Renderer
  - [x] Animation
  - [x] Pre-rendered background
  - [x] Depth testing
  - [ ] Sprites
  - [ ] Shadows

### Controls

- W/S to move forward/backward.
- A/D to rotate left/right.
- Tab to access inventory.
- Enter is action button.

### Other tools

- [Bio2 script viewer](https://github.com/OpenBiohazard2/Bio2ScriptIde)

### References

- https://github.com/pmandin
- https://github.com/yanmingsohu/Biohazard2
- https://github.com/MeganGrass/BioScript
- https://github.com/mortician
