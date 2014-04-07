package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
#include <curses.h>
#include <stdio.h>
#cgo LDFLAGS: -lncurses

//Get single keystroke.
int getKey(){
    return getch();
}

//Clear screen.
void cls(){
    system("clear");
}

//initialize curse scheme
void initCurse(){
    initscr();
    noecho();
   // use curse stdscr to handle arrow key press.
      keypad(stdscr,TRUE);
}

//end curses scheme
void endCurse(){
    endwin();
}
*/
import "C" //Embed c function int Go.

const (
	KEY_EXIT  = 27   // ESC
	KEY_UP    = 259  // press up arrow key
	KEY_DOWN  = 258  // press down arrow key
	KEY_LEFT  = 260  // press left arrow key
	KEY_RIGHT = 261  // press right arrow key
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
	// Plot tile and fill tile cell data.
	fmt.Println("=============================")
	for row := 0; row < SIZE; row++ {
		fmt.Print("\r|")
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
	fmt.Println("\r=============================")
}

// Up command.
func Up() {
	canMoveUp := moveUp()
	if canMoveUp {
		fillRandom()
		C.cls()
		plot()
	}
}

// Down command.
func Down() {
	canMoveDown := moveDown()
	if canMoveDown {
		fillRandom()
		C.cls()
		plot()
	}
}

// Left command.
func Left() {
	canMoveLeft := moveLeft()
	if canMoveLeft {
		fillRandom()
		C.cls()
		plot()
	}
}

// Right command.
func Right() {
	canMoveRight := moveRight()
	if canMoveRight {
		fillRandom()
		C.cls()
		plot()
	}
}

func moveUp() bool {
	moveCount := 0
	for col := 0; col < SIZE; col++ {
		merged := false // Each cell merged once.
		pos := 0        // Location of non-empty cells.
		flag := true    // The upmost empty location markup.
		for row := 0; row < SIZE; row++ {
			if tile[row][col] != 0 {
				if tile[pos][col] == 0 { // Current tile[rowl][col] is the
					//upmost non-empty cell.
					tile[pos][col] = tile[row][col]
					tile[row][col] = 0
					flag = false
					moveCount++
					continue
				}
				if tile[row][col] == tile[pos][col] && !merged {
					if pos != row {
						// Merge equal tiles.
						tile[pos][col] = tile[row][col] + tile[pos][col]
						tile[row][col] = 0
						blankCount += 1
						moveCount++
						merged = true
						if maxValue < tile[pos][col] {
							maxValue = tile[pos][col]
						}
						continue
					}
				} else {
					if (pos + 1) != row { // Non adjacent cells with non-empty value.
						tile[pos+1][col] = tile[row][col]
						tile[row][col] = 0
						pos += 1
						merged = false // tile[pos][col] is new unmerged cell.
						moveCount++
						continue
					}
				}
				pos = row
				merged = false
			} // if
			if flag { // Set the upmost empty location.
				pos = row
				flag = false
			}
		} // for row
	} // for col
	return moveCount > 0
} // moveUp

func moveDown() bool {
	moveCount := 0
	for col := 0; col < SIZE; col++ {
		merged := false // Each cell merged once.
		pos := SIZE - 1 // Location of non-empty cells.
		flag := true    // The downmost empty location markup.
		for row := SIZE - 1; row >= 0; row-- {
			if tile[row][col] != 0 {
				if tile[pos][col] == 0 { // Current tile[row][col] is
					// the downmost non-empty cell.
					tile[pos][col] = tile[row][col]
					tile[row][col] = 0
					flag = false
					moveCount++
					continue
				}

				if tile[row][col] == tile[pos][col] && !merged {
					if pos != row {
						// Merge equal tile cells.
						tile[pos][col] = tile[row][col] + tile[pos][col]
						tile[row][col] = 0
						blankCount += 1
						moveCount++
						merged = true
						if maxValue < tile[pos][col] {
							maxValue = tile[pos][col]
						}
						continue
					}
				} else {
					if (pos - 1) != row {
						tile[pos-1][col] = tile[row][col]
						tile[row][col] = 0
						pos -= 1
						moveCount++
						merged = false
						continue
					}
				}
				pos = row
				merged = false
			} // if
			if flag { // Set the downmost empty location.
				pos = row
				flag = false
			}
		} // for row
	} // for col
	return moveCount > 0
} // moveDown

func moveLeft() bool {
	moveCount := 0
	for row := 0; row < SIZE; row++ {
		merged := false // Each cell merged once.
		pos := 0        // Location of non-empty cell.
		flag := true    // The leftmost empty location markup.
		for col := 0; col < SIZE; col++ {
			if tile[row][col] != 0 {
				if tile[row][pos] == 0 { // Current tile[row][col] is
					// the leftmost non-empty cell.
					tile[row][pos] = tile[row][col]
					tile[row][col] = 0
					flag = false
					moveCount++
					continue
				}

				if tile[row][col] == tile[row][pos] && !merged {
					if pos != col {
						// Merge equal tiles.
						tile[row][pos] = tile[row][col] + tile[row][pos]
						tile[row][col] = 0
						blankCount += 1
						moveCount++
						merged = true
						if maxValue < tile[row][pos] {
							maxValue = tile[row][pos]
						}
						continue
					}
				} else {
					if (pos + 1) != col {
						tile[row][pos+1] = tile[row][col]
						tile[row][col] = 0
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
} // moveLeft

func moveRight() bool {
	moveCount := 0
	for row := 0; row < SIZE; row++ {
		merged := false // Each cell merged once.
		pos := SIZE - 1 // Location of non-empty cells.
		flag := true    // The rightmost empty location markup.
		for col := SIZE - 1; col >= 0; col-- {
			if tile[row][col] != 0 {
				if tile[row][pos] == 0 { // Current tile[row][col] is
					// the rightmost non-empty cell.
					tile[row][pos] = tile[row][col]
					tile[row][col] = 0
					flag = false
					moveCount++
					continue
				}

				if tile[row][col] == tile[row][pos] && !merged {
					if pos != col {
						// Merge equal tiles.
						tile[row][pos] = tile[row][col] + tile[row][pos]
						tile[row][col] = 0
						blankCount += 1
						moveCount++
						merged = true
						if maxValue < tile[row][pos] {
							maxValue = tile[row][pos]
						}
						continue
					}
				} else {
					if (pos - 1) != col {
						tile[row][pos-1] = tile[row][col]
						tile[row][col] = 0
						pos -= 1
						merged = false
						moveCount++
						continue
					}
				}
				pos = col
				merged = false
			} // if
			if flag { // Set the rightmost empty location.
				pos = col
				flag = false
			}
		} // fow col
	} // for row
	return moveCount > 0
} // moveRight

func getUserCmd() bool {
	cmd := C.getKey()
	if cmd == KEY_EXIT {
		return false
	}
	switch cmd {
	case KEY_UP:
		Up()
	case KEY_DOWN:
		Down()
	case KEY_LEFT:
		Left()
	case KEY_RIGHT:
		Right()
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
	C.initCurse()
	isRunning := true
	for isRunning {
		if !canMove() {
			if maxValue == WIN {
				fmt.Println("\rYou win!")
				break
			} else {
				fmt.Println("\rYou lose.")
				break
			}
		}
		isRunning = getUserCmd()
	}
	fmt.Println("\r Press any key to exit.")
	C.getKey()
	C.endCurse()
}
