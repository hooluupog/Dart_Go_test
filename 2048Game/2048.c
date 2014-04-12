// By hooluupog
// 2014/4/12
// Thanks to  Jay Chan's tiny 2048 algorithm.
// Thanks to  catull's description about keyboard escape sequences here
// -link:https://gist.github.com/anirudh-chhangani/9989292#comment-1206456
// Test passed on Linux
// The meaning of
//       read(0, &keyInput, 3);
//		 fillGrid(1, key[(keyInput >> 16) & 3]);
// The reason why 3 characters are read is that arrow keys generate
// 3-byte sequences. As follows,
// =========================================
// |    key      |    bytes     | Hex value|
// | arrow left  | 0x1b,'[','D' | 0x1b5b44 |
// | arrow right | 0x1b,'[','C' | 0x1b5b43 |
// | arrow up    | 0x1b,'[','A' | 0x1b5b41 |
// | arrow down  | 0x1b,'[','B' | 0x1b5b42 |
// =========================================
// 0x1b=27,represent 'ESC'. `0x1b,'[','x' ` is  keyboard escape sequence.
// When the characters are read, the 3-byte sequence actually is reversed.
// So for instance, in the last case (arrow down), keyInput is assigned 0x425b1b.
// keyInput is then used to decide which of the Arrow-array(key[]) value to use.
// keyInput is right-shifted by 16 bits, reducing it to 0x42, so we have
// key[(keyInput >> 16) & 3], which is equivalent to key[2], the 3rd element of Arrow-array.
// Thus, the statement `fillGrid(1,  key[(keyInput >> 16) & 3]);` evaluates to calling
// fillGrid(1, 1).

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>

#define SIDE 4
#define GRID_LENGTH (SIDE * SIDE)
#define SCORE  2048

enum {
	LOSE = 0, WIN, WAIT
} result;
int Grid[GRID_LENGTH];
int moveCount = 1; // Movement count of tiles alone one direction.
int key[] = { 2, 3, 1, 0 }; // Arrow array. 0x042 & 3, 0x043 & 3, 0x041 & 3,0x044 & 3.

// Get index of Grid from tile(i,j) according direction.
// When filling Grid according different move direction with the same code,
// for (row = 0; row < SIDE; row++){
//     for (col=0,col < SIDE; col++)
//     {
//     }
//  }
// you need to rotate grid to get the right value.
// direction = 2(down),3(right),1(up),0(left)
//            0
//      ------<------
//     |             |
//     |             |
//   2 ↓             ↑ 1
//     |             |
//     |             |
//      ------>------
//            3
// Suppose there is a two dimension array as follows,
//  1 2
//  3 4
// Initialized direction = 0,moveLeft,
//  1 2
//  3 4       No need rotation.
//  direction = 1,moveUp,
//  2 4
//  1 3       Rotate array 90 degree anticlock.
//  direction = 3, moveRight,
//  4 3
//  2 1       Rotate array 90 degree anticlock.
//  direction = 2, moveDown,
//  3 1
//  4 2       Rotate array 90 degree anticlock.
//  The course of fillGrid:
//  direction = 0, do loop: for row++    |
//                            for col++  |
//  direction = 1, do loop: for col--    |
//                            for row++  | after rotating    All of direction movement
//  direction = 3, do loop: for row--    | ===============>  do the same loop:
//                            for col--  |                   for row++
//  direction = 2, do loop: for col++    |                       for col++
//                            for row--  |
//  After doing the aboves,convert two dimension grid into one dimension
//  array. index = SIDE*i + j;
int tile(int direction, int i, int j) {
	if (direction <= 0) {
		return SIDE * i + j; //two-dimension to one-dimension
		// address conversion.
	}
	// Rotate grid according to direction.
	return tile(direction - 1, j, SIDE - 1 - i);
}

// Fill grid tiles with new data.
// algorithm course:
// Init: col = 0,present = 0, next = tile[cursor++],cursor = 0.
// do loop,
// when present == 0, present = next, next = tile[cursor++]
// when present != next, next != 0 && present != 0, present= next.
// when present == next, next != 0, merge two titles.
// when cursor >= SIDE, then do tile[col++] = present, until col >= SIDE.
// loop over.
// E.g
//      tile[4]=[0 2 0 4]             ||           tile[4]=[0 2 0 2]
//    tile   present next cursor      ||   //    tile   present next cursor
//  [0 2 0 4]   0     0      0        ||   //  [0 2 0 2]   0     0      0
//  [0 2 0 4]   0     0      1        ||   //  [0 2 0 2]   0     0      1
//  [0 2 0 4]   2     2      2        ||   //  [0 2 0 2]   2     2      2
//  [0 2 0 4]   2     0      3        ||   //  [0 2 0 2]   2     0      3
//  [2 2 0 4]   4     4      4        ||   //  [4 2 0 2]   0     2      4
//  [2 4 0 4]   0     4      4        ||   //  [4 0 0 2]   0     2      4
//  [2 4 0 4]   0     4      4        ||   //  [4 0 0 2]   0     2      4
//  [2 4 0 0]   0     4      4        ||   //  [4 0 0 0]   0     2      4

void fillGrid(int mode, int direction) {
	//fill Grid tiles according to every step.
	// mode = 1 OR 0
	// if mode == 0, do win or lose test.
	// if mode == 1, fill grid tiles according to every step.
	// direction, arrow array value. direction = key[index].
	// direction, arrow array value. direction = key[index].
	int row = SIDE, col, cursor, present, next;
	if (mode != 0) {
		moveCount = 0;
	}
	for (; row--;) {
		cursor = col = present = 0;
		for (; col < SIDE;) {
			if (cursor < SIDE) {
				next = Grid[tile(direction, row, cursor++)];
				result |= next / SCORE; //if next/score == 1,result = WIN (1).
				if (present != 0 && next != 0) {
					// Merge two tiles value if they equals.
					if (mode != 0) {
						if (present == next) {
							Grid[tile(direction, row, col)] = present << 1;
							moveCount++;
						} else {
							if (Grid[tile(direction, row, col)] != present) {
								Grid[tile(direction, row, col)] = present;
								moveCount++;
							}
						}
					}
					col++;
				}
				if (present == 0) {
					present = next;
				} else {
					if (next != 0) {
						present = (present == next) ? 0 : next;
					}
				}
			} else {
				// When code run into this block,it turns out that there is at least
				// an empty tile in this row. present always != 0.
				if (mode != 0) {
					if (Grid[tile(direction, row, col)] != present) {
						Grid[tile(direction, row, col)] = present;
						moveCount++;
					}
				}
				col++;
				// if present = 0,
				//     then there is at least one empty tile
				//     in this row. !present = 1, WAIT * 1 = WAIT.
				//     result = result | WAIT.
				// if present != 0,
				//     then there is none empty tile in this row. !present = 0,
				//     WAIT * 0 = LOSE. result = result | LOSE.
				result |= WAIT * (!present);
				present = 0;// Tile movement finished and the rest tiles are empty.
				// set present = 0, means there are empty tiles in the row.
			}
		}
	}
}

//Generate random nums for tiles.
void cellRand() {
	int randNum, randPosition;
	randNum = GRID_LENGTH + (rand() & (GRID_LENGTH - 1));  // randNum = 16~31
	randPosition = randNum & (GRID_LENGTH - 1);  //randPosition = 0~15
	//Find a non zero tile and fill the generated random number into it.
	for (; (Grid[randPosition] * randNum) != 0; randNum--) {
		randPosition = randNum & (GRID_LENGTH - 1);  //randPosition = 0~15
	}
	if (randNum != 0) { // All tiles are full if randNum == 0.
		Grid[randPosition] = 2 << (rand() % 2); // 2 or 4
	}
}

// Get user keyboard input and run game step by step.
void step() {
	//keyInput : 0x1b5b44, 0x1b5b43, 0x1b5b41, 0x1b5b42
	int keyInput = 0, input = 0;
	int i = SIDE;
	result = LOSE;
	for (; i--;) {
		// In verify mode(mode=0),test will win or lose.
		fillGrid(0, i);
	}
	if (moveCount != 0) { // The movement count of tiles alone one direction
		// during last step fillGrid(1,direction) running.
		cellRand();
		// clear console screen and reset cursor
		puts("\e[2J\e[H");
		// print the tiles onto the terminal.
		for (i = GRID_LENGTH; i--;) {
			if (Grid[i] != 0) {
				printf("%4d|", Grid[i]);
			} else {
				printf("%s", "    |");
			}
			if (0 == (i % SIDE)) {
				putchar('\n');
			}
		}
	}
	if (result == WAIT) {
		// read input from keyboard
		read(0, &keyInput, 3);
		//(0x1b5b44, 0x1b5b43, 0x1b5b41, 0x1b5b42) >> 16 => (0x44,0x43,0x41,0x42) &3  =>  0,3,1,2
		input = keyInput >> 16;
		//printf("%d  %d\n",keyInput,input);
		if (input == 0x44 | input == 0x43 | input == 0x41 | input == 0x42) {
			fillGrid(1, key[input & 3]);
		} else {
			moveCount = 0;
		}
		step();
	}
}

int main(void) {
	// clear the screen in preparation for the game
	system("stty cbreak -echo");
	/* intialize random number generator */
	srand((unsigned) time(NULL));
	// Init Grid.
	int i;
	for (i = 0; i < SIDE; i++) {
		Grid[i] = 0;
	}
	cellRand();
	step();
	// game has finished by this point
	puts(result & WIN ? "You win" : "You lose");
	system("stty cbreak echo");
	return 0;
}
//[2048]
