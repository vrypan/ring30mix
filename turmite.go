package main

// Turmite represents an 8x8 grid with 4 colors (2 bits per cell) and a moving cell
type Turmite struct {
	grid [16]byte // 64 cells packed: 4 cells per byte (2 bits each)
	x, y int      // position [0-7]
	dir  int      // direction: 0=N, 1=E, 2=S, 3=W
}

// Movement deltas: [N, E, S, W] → [dx, dy]
var deltas = [4][2]int{
	{0, -1}, // North
	{1, 0},  // East
	{0, 1},  // South
	{-1, 0}, // West
}

// NewTurmite creates a new turmite with given position and direction
func NewTurmite(x, y, dir int) *Turmite {
	return &Turmite{
		x:   x & 7, // ensure [0-7]
		y:   y & 7,
		dir: dir & 3, // ensure [0-3]
	}
}

// Get returns the color value at cell (x, y)
func (t *Turmite) Get(x, y int) int {
	// Calculate bit position: (y*8 + x) * 2
	bitPos := ((y & 7) * 8 + (x & 7)) * 2
	byteIdx := bitPos >> 3     // divide by 8
	bitOffset := bitPos & 7    // modulo 8

	return int((t.grid[byteIdx] >> bitOffset) & 0x03)
}

// Set sets the color value at cell (x, y)
func (t *Turmite) Set(x, y, val int) {
	bitPos := ((y & 7) * 8 + (x & 7)) * 2
	byteIdx := bitPos >> 3
	bitOffset := bitPos & 7

	// Clear 2 bits then set new value
	t.grid[byteIdx] &= ^(0x03 << bitOffset)
	t.grid[byteIdx] |= byte(val&0x03) << bitOffset
}

// Step executes one iteration of the LLLR turmite
// DNA: "1L.2L.3L.0R"
// State 0 → paint 1, turn left
// State 1 → paint 2, turn left
// State 2 → paint 3, turn left
// State 3 → paint 0, turn right
func (t *Turmite) Step() {
	state := t.Get(t.x, t.y)

	// Inline DNA logic for performance
	switch state {
	case 0:
		t.Set(t.x, t.y, 1)
		t.turnLeft()
	case 1:
		t.Set(t.x, t.y, 2)
		t.turnLeft()
	case 2:
		t.Set(t.x, t.y, 3)
		t.turnLeft()
	case 3:
		t.Set(t.x, t.y, 0)
		t.turnRight()
	}

	t.moveForward()
}

// turnLeft rotates the turmite 90° counter-clockwise
func (t *Turmite) turnLeft() {
	t.dir = (t.dir + 3) & 3 // (dir - 1) mod 4
}

// turnRight rotates the turmite 90° clockwise
func (t *Turmite) turnRight() {
	t.dir = (t.dir + 1) & 3
}

// moveForward moves the turmite one cell forward (with wrapping)
func (t *Turmite) moveForward() {
	delta := deltas[t.dir]
	t.x = (t.x + delta[0]) & 7 // wrap at edges
	t.y = (t.y + delta[1]) & 7
}

// CopyGrid returns a copy of the current grid state
func (t *Turmite) CopyGrid() [16]byte {
	var result [16]byte
	copy(result[:], t.grid[:])
	return result
}
