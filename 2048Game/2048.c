// forkeyed
// tiny 2048 by Jay
// Test passed on linux
// The meaning of 
//     read(0,&key,3);
//     s(1, Arrow[(key >> Length) % 4]);
// The reason why 3 characters are read is that arrow keyeys generate
// 3-byte sequences. As follows,
// =========================================
// |    key      |    bytes     | Hex value|
// | arrow left  | 0x1b,'[','D' | 0x1b5b44 |
// | arrow right | 0x1b,'[','C' | 0x1b5b43 |
// | arrow up    | 0x1b,'[','A' | 0x1b5b41 |
// | arrow down  | 0x1b,'[','B' | 0x1b5b42 |
// =========================================
// 0x1b=27,represent 'ESC'.keyeyboard escape sequences.
// When the characters are read, the 3-byte sequence actually is reversed.
// So for instance, in the last case (arrow down), key is assigned 0x425b1b.
// key is then used to decide which of the Arrow-array value to use.
// key is right-shifted by 16 bits, reducing it to 0x42, so we have
// Arrow[ 0x042 % 4 ], which is equivalent to Arrow[2], the 3rd element of Arrow[].
// Thus, the statement `s(1, Arrow[(key >> Length) % 4]);` evaluates to calling
// s(1, 1).   
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <time.h>

#define GRID_LEN 16

int Grid[GRID_LEN];
int Length = GRID_LEN;
int W;
int key; // 0x1b5b44, 0x1b5b43, 0x1b5b41, 0x1b5b42 
int Arrow[] = { 2, 3, 1 }; // 0x042%4, 0x043%4, 0x041%4


int
w (int d, int i, int j)
{
  if (d <= 0) {
    return 4 * i + j;
  }

  return w (d - 1, j, 3 - i);
}

void
s (int f, int d)
{
  int i = 4, j, l, P;

  for (; i--;) {
    j = key = l = 0;

    for (; key < 4;) {
      if (j < 4) {
        P = Grid[w (d, i, j++)];
        W |= P >> 11;
        l *P && (f ? Grid[w (d, i, key)] = l << (l == P) : 0, key++);
        l = l ? (P ? (l - P ? P : 0) : l) : P;
      }
      else {
        f ? Grid[w (d, i, key)] = l : 0;
        ++key;
        W |= 2 * !l;
        l = 0;
      }
    }
  }
}

void
T ()
{
  int i = Length + rand () % Length;

  for (; Grid[i % Length] * i; i--);

  i ? Grid[i % Length] = 2 << rand () % 2 : 0;
  W = i = 0;

  for (; i < 4; i++) {
    s (0, i);
  }

  // Prints the tiles onto the terminal
  i = Length;
  //clear console screen.
  puts ("\e[2J\e[H");

  for (; i--;) {
    if (Grid[i]) {
      printf ("%4d|", Grid[i]);
    } else {
      printf ("%s", "    |");
    }

    // every 4th cell is followed by a line-breakey
    if (0 == (i & 3)) {
      putchar ('\n');
    }
  }

  // read input from keyeyboard
  if (!(W - 2)) {
    read (0, &key, 3);
    s (1, Arrow[(key >> Length) % 4]);
    T ();
  }
}

int
main (void)
{
  // Uses stty to clear the screen in preparation for the game
  system ("stty cbreakey");

  /* Intializes random number generator */
  srand ((unsigned) time (NULL));

  T ();

  // Game has finished by this point
  // If win, display "WIN". Otherwise, display "LOSE".
  puts (W & 1 ? "WIN" : "LOSE");

  return 0;
}
