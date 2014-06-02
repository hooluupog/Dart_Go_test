package game

import (
	"math/rand"
	"time"
)

const (
	SIZE = 4    // tile size*size.
	WIN  = 2048 // win scores.
)

type Position struct {
	X     int
	Y     int
	Value uint32
}

type Tile struct {
	Value       uint32
	MergedFrom  []Position
	PrePosition Position
	IsNewTile   bool
}

var (
	tile       [][]Tile
	blankCount uint32 = SIZE * SIZE // count of empty cells.
	maxValue   uint32               //Highest score.
)

func InitGrid() [][]Tile {
	blankCount = SIZE * SIZE
	tile = make([][]Tile, SIZE)
	for i := 0; i < SIZE; i++ {
		tile[i] = make([]Tile, SIZE)
	}
	fillRandom()
	fillRandom()
	return tile
}

func clearStatus() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			tile[i][j].MergedFrom = nil
			tile[i][j].IsNewTile = false
			tile[i][j].PrePosition = Position{i, j, tile[i][j].Value}
		}
	}
}

func fillRandom() int {
	// Generate random numbers. Return 1: sucess; Return 0: fail.
	rand.Seed(time.Now().UnixNano())
	//random location
	row, col := rand.Intn(4), rand.Intn(4)
	if blankCount == 0 {
		return 0
	}
	for tile[row][col].Value != 0 {
		rand.Seed(time.Now().UnixNano())
		row, col = rand.Intn(4), rand.Intn(4)
	}
	//random value 2 or 4.
	tile[row][col].Value = ((rand.Uint32()*blankCount)%2 + 1) * 2
	tile[row][col].IsNewTile = true
	tile[row][col].PrePosition = Position{row, col, tile[row][col].Value}
	blankCount -= 1
	return 1
}

func moveDirection(direction int) {
	canMove := moveOneStep(direction)
	if canMove {
		fillRandom()
	}
}

//Get tile after rotation accoring given direction.
func rTile(direction int, row, col int) (*Tile, int, int) {
	if direction <= 0 {
		return &tile[row][col], row, col
	}
	return rTile(direction-1, col, SIZE-1-row)
}

// Make left direction movement as default to fill grid without need rotating grid.
// When other direction movement command input,rotate grid with new layout and then
// fill grid according left direction movement.
// double loop sequence: outer-row, inner-col.
// arrow key     direction  outer-inner loop  move  ||  after rotation  move
//   left            0        row++,col++      ←   ||   row++,col++     ←
//   right           2        row--,col--      →   ||   row++,col++     ←
//   up              1        col--,row++      ↑   ||   row++,col++     ←
//   down            3        col++,row--      ↓   ||   row++,col++     ←
// Do clock-anti rotation.
//                 left(0)       up(1)      right(2)    down(3)
//                     ←             ←         ←          ←
//                   1   2         2    4      4    3      3    1
//                 ↓      ↑
//                   3   4         1    3      2    1      4    2
//                     →
func moveOneStep(direction int) bool {
	moveCount := 0
	moveCountInSameCol := 0
	clearStatus()
	for row := 0; row < SIZE; row++ {
		pos := 0     // Location of non-empty tile.
		flag := true // The leftmost empty location markup.
		moveCountInSameCol = 0
		for col := 0; col < SIZE; col++ {
			currentTile, rrow, rcol := rTile(direction, row, col)
			posTile, prow, pcol := rTile(direction, row, pos)
			if currentTile.Value != 0 {
				if posTile.Value == 0 { // Current tile[row][col] is
					// the leftmost non-empty tile.
					posTile.Value = currentTile.Value
					posTile.PrePosition = Position{rrow, rcol, currentTile.Value}
					currentTile.PrePosition = Position{prow, pcol, currentTile.Value}
					currentTile.Value = 0
					flag = false
					moveCount++
					moveCountInSameCol++
					continue
				}
				if currentTile.Value == posTile.Value && posTile.MergedFrom == nil {
					if pos != col {
						// Merge equal tiles.
						posTile.MergedFrom = make([]Position, 2)
						posTile.MergedFrom[0] = posTile.PrePosition
						posTile.MergedFrom[1] = Position{rrow, rcol, currentTile.Value}
						posTile.Value = currentTile.Value + posTile.Value
						posTile.PrePosition = Position{rrow, rcol, currentTile.Value}
						currentTile.PrePosition = Position{prow, pcol, currentTile.Value}
						currentTile.Value = 0
						blankCount = blankCount + 1
						moveCount++
						moveCountInSameCol++
						if maxValue < posTile.Value {
							maxValue = posTile.Value
						}
						continue
					}
				} else {
					if (pos + 1) != col {
						posTile, prow, pcol = rTile(direction, row, pos+1)
						posTile.Value = currentTile.Value
						posTile.PrePosition = Position{rrow, rcol, currentTile.Value}
						currentTile.PrePosition = Position{prow, pcol, currentTile.Value}
						currentTile.Value = 0
						pos = pos + 1
						moveCount++
						moveCountInSameCol++
						continue
					}
				}
				if moveCountInSameCol == 0 {
					posTile.PrePosition = Position{prow, pcol, posTile.Value}
				}
				pos = col
			} // if
			if flag { // Set the leftmost empty location.
				pos = col
				flag = false
			}
		} // for col
	} // for row
	return moveCount > 0
} // moveOneStep

func Move(direction string) [][]Tile {
	switch direction {
	case "up":
		moveDirection(1)
	case "down":
		moveDirection(3)
	case "left":
		moveDirection(0)
	case "right":
		moveDirection(2)
	}
	return tile
}

func canMove() bool {
	if blankCount == 0 {
		for col := 0; col < SIZE; col++ {
			for row := 0; row < SIZE-1; row++ {
				if tile[row][col].Value == tile[row+1][col].Value {
					return true
				}
			}
		}
		for row := 0; row < SIZE; row++ {
			for col := 0; col < SIZE-1; col++ {
				if tile[row][col].Value == tile[row][col+1].Value {
					return true
				}
			}
		}
		return false
	}
	return true
}

func Win() string {
	if maxValue == WIN {
		return "win"
	}
	if !canMove() && maxValue != WIN {
		return "lose"
	}
	return "wait"
}

func GridSize() int {
	return SIZE
}

func Reset() {
	maxValue = 0
	blankCount = SIZE * SIZE
}
