// 2014/4/12
// by hooluupog

package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
#include <conio.h>

//Get single keystroke.
int getKey(){
    return _getch();
}

//Clear screen.
void cls(){
    system("cls");
}
*/
import "C" //Embed c function int Go.

const (
	KEY_EXIT  = 27   // ESC
	KEY_CTRL  = 224  // make code E0
	KEY_UP    = 72   // press up arrow key
	KEY_DOWN  = 80   // press down arrow key
	KEY_LEFT  = 75   // press left arrow key
	KEY_RIGHT = 77   // press right arrow key
	SIZE      = 4    // tile size*size.
	WIN       = 2048 // win scores.
)

var (
	tile       [SIZE][SIZE]uint32
	blankCount uint32 = SIZE * SIZE // count of empty cells.
	maxValue   uint32               //Highest score.
)

func init() {
	fillRandom()
	fillRandom()
	C.cls()
	plot()
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

func plot() {
	// Print grid on the console.
	fmt.Println("=============================")
	for row := 0; row < SIZE; row++ {
		fmt.Print("|")
		for col := 0; col < SIZE; col++ {
			data := tile[row][col]
			if data == 0 {
				fmt.Print("      |")
			} else {
				fmt.Printf(" %4d |", data)
			}

		}
		fmt.Print("\n")
	}
	fmt.Println("=============================")
}

func moveDirection(direction int) {
	canMove := moveOneStep(direction)
	if canMove {
		fillRandom()
		C.cls()
		plot()
	}
}

//Get tile after rotation accoring given direction.
func rTile(direction int, row, col int) *uint32 {
	if direction <= 0 {
		return &tile[row][col]
	}
	return rTile(direction-1, col, SIZE-1-row)
}

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

func getCmd() bool {
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

	cmd := C.getKey()
	if cmd == KEY_EXIT {
		return false
	}
	if cmd == KEY_CTRL {
		cmd = C.getKey()
		switch cmd {
		case KEY_UP:
			moveDirection(1)
		case KEY_DOWN:
			moveDirection(3)
		case KEY_LEFT:
			moveDirection(0)
		case KEY_RIGHT:
			moveDirection(2)
		}
	}
	return true
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

func main() {
	isRunning := true
	for isRunning {
		if maxValue == WIN {
			fmt.Println("You win!")
			break
		}
		if !canMove() && maxValue != WIN {
			fmt.Println("You lose.")
			break
		}
		isRunning = getCmd()
	}
	fmt.Println("Press any key to exit.")
	C.getKey()
}
