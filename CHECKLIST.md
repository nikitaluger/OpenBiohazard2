# Project Checklist

## Setup
- [x] Install Go and C compiler
- [x] Build project successfully
- [x] Run game data validation
- [x] Run application without panic

## Contributions
- [x] Update README task list (PR submitted)
- [x] Fix OpenGL context version (PR submitted)
- [x] Add README data setup notes
- [x] Add testing plan
- [x] Add CI workflow
- [x] Add contributing guide
- [x] Add hermetic PLD tests
- [x] Add hermetic streamreader extra tests
- [x] Add hermetic TIM tests
- [x] Add AOT sprite removal tests
- [x] Implement sprite kill opcode
- [x] Enable sprite rendering

## Next
- [x] Add hermetic tests for bitreader edge cases (mostly done)
- [x] Add valid-case TIM test
- [x] Add tests for EMD / MD1 / ESP
- [x] Support multiple sprites in renderer
- [x] Implement one PLC opcode (motion/dest/neck)
- [ ] Implement audio playback
  - [x] Add audio manager placeholder
  - [x] Wire BGM control opcode
  - [x] Integrate audio library
  - [x] Decode VAG ADPCM (verified)
  - [x] Create VAG test tool (vagtest)
  - [ ] Map BGM/SFX IDs to tracks
  - [ ] Parse VAB data
  - [ ] Looping and volume control
  - [ ] SFX triggers
- [ ] Implement enemy AI / models
  - [x] Add enemy AI state scaffold
  - [ ] Render enemy models
  - [ ] Add basic enemy movement
- [x] Implement door models and transitions
- [x] Implement inventory item use
- [ ] Implement puzzles
  - [ ] Add puzzle trigger handling
  - [ ] Implement puzzle state logic

## Notes
- Tests added so far are mostly scaffolding (error paths and crash guards).
- They do not yet verify full parsing correctness or game behavior.
