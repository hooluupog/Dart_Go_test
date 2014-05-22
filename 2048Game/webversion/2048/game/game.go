package game

import (
	"math/rand"
	"time"
)

const (
	SIZE = 4    // tile size*size.
	WIN  = 2048 // win scores.
)

var (
	tile       [][]uint32
	blankCount uint32 = SIZE * SIZE // count of empty cells.
	maxValue   uint32               //Highest score.
)

func InitGrid() [][]uint32 {
	blankCount = SIZE * SIZE
	tile = make([][]uint32, SIZE)
	for i := 0; i < SIZE; i++ {
		tile[i] = make([]uint32, SIZE)
	}
	fillRandom()
	fillRandom()
	return tile
}

func fillRandom() int {
	// Generate random numbers. Return 1: sucess; Return 0: fail.
	rand.Seed(time.Now().UnixNano())
	//random location
	row, col := rand.Intn(4), rand.Intn(4)
	if blankCount == 0 {
		return 0
	}
	for tile[row][col] != 0 {
		rand.Seed(time.Now().UnixNano())
		row, col = rand.Intn(4), rand.Intn(4)
	}
	//random value 2 or 4.
	tile[row][col] = ((rand.Uint32()*blankCount)%2 + 1) * 2
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
func rTile(direction int, row, col int) *uint32 {
	if direction <= 0 {
		return &tile[row][col]
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
	for row := 0; row < SIZE; row++ {
		merged := false // Each tile merged once.
		pos := 0        // Location of non-empty tile.
		flag := true    // The leftmost empty location markup.
		for col := 0; col < SIZE; col++ {
			if *rTile(direction, row, col) != 0 {
				if *rTile(direction, row, pos) == 0 { // Current tile[row][col] is
					// the leftmost non-empty tile.
					*rTile(direction, row, pos) = *rTile(direction, row, col)
					*rTile(direction, row, col) = 0
					flag = false
					moveCount++
					continue
				}

				if *rTile(direction, row, col) == *rTile(direction, row, pos) && !merged {
					if pos != col {
						// Merge equal tiles.
						*rTile(direction, row, pos) = *rTile(direction, row, col) + *rTile(direction, row, pos)
						*rTile(direction, row, col) = 0
						blankCount += 1
						moveCount++
						merged = true
						if maxValue < *rTile(direction, row, pos) {
							maxValue = *rTile(direction, row, pos)
						}
						continue
					}
				} else {
					if (pos + 1) != col {
						*rTile(direction, row, pos+1) = *rTile(direction, row, col)
						*rTile(direction, row, col) = 0
						pos += 1
						moveCount++
						merged = false
						continue
					}
				}
				pos = col
				merged = false
			} // if
			if flag { // Set the leftmost empty location.
				pos = col
				flag = false
			}
		} // for col
	} // for row
	return moveCount > 0
} // moveOneStep

func Move(direction string) [][]uint32 {
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
				if tile[row][col] == tile[row+1][col] {
					return true
				}
			}
		}
		for row := 0; row < SIZE; row++ {
			for col := 0; col < SIZE-1; col++ {
				if tile[row][col] == tile[row][col+1] {
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
