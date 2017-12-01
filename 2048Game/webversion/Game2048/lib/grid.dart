library grid;

import 'dart:math';

final int SIZE = 4; // tile size*size.
final int WIN = 2048; // win scores.

class _Position {
  int _x;
  int _y;
  int _value;
  _Position(this._x, this._y, this._value);
  _Position.position(this._x, this._y);
  int get Value => _value;
  int get Row => _x;
  int get Col => _y;
}

class Tile {
  int _value;
  List<_Position> _mergedFrom;
  _Position _prePosition;
  bool _isNewTile;
  Tile(this._value, this._mergedFrom, this._prePosition, this._isNewTile);

  int get Value => _value;
  List get MergedFrom => _mergedFrom;
  _Position get PrePosition => _prePosition;
  bool get IsNewTile => _isNewTile;

  set Value(int value) => _value = value;
  set MergedFrom(List mergedFrom) => _mergedFrom = mergedFrom;
  set PrePosition(_Position prePosition) => _prePosition = prePosition;
  set IsNewTile(bool isNewTile) => _isNewTile = isNewTile;

  @override
  String toString() => '$Value';
}

List tile;
int blankCount; // count of empty cells.
int maxValue = 0; //Highest score.

List InitGrid() {
  blankCount = SIZE * SIZE;
  tile = new List(SIZE).map((_) => new List(SIZE)).toList();
  for (var i = 0; i < SIZE; i++) {
    for (var j = 0; j < SIZE; j++) {
        tile[i][j] = new Tile(0,null,null,false);
    }
  }
  fillRandom();
  fillRandom();
  return tile;
}

clearStatus() {
  for (var i = 0; i < SIZE; i++) {
    for (var j = 0; j < SIZE; j++) {
      tile[i][j].MergedFrom = null;
      tile[i][j].IsNewTile = false;
      tile[i][j].PrePosition = new _Position(i, j, tile[i][j].Value);
    }
  }
}

int fillRandom() {
  // Generate random numbers. Return 1: sucess; Return 0: fail.
  var rand = new Random(new DateTime.now().millisecondsSinceEpoch);
  // random location
  var row = rand.nextInt(4);
  var col = rand.nextInt(4);
  if (blankCount == 0) {
    return 0;
  }
  while (tile[row][col].Value != 0) {
    rand = new Random(new DateTime.now().millisecondsSinceEpoch);
    row = rand.nextInt(4);
    col = rand.nextInt(4);
  }
  // random value 2 or 4.
  tile[row][col].Value = ((rand.nextInt(2 ^ 32 - 1) * blankCount) % 2 + 1) * 2;
  tile[row][col].IsNewTile = true;
  tile[row][col].PrePosition = new _Position(row, col, tile[row][col].Value);
  blankCount -= 1;
  return 1;
}

moveDirection(int direction) {
  bool canMove = moveOneStep(direction);
  if (canMove) {
    fillRandom();
  }
}

/// Get tile after rotation accoring given direction.
_Position rTile(int direction, int row, int col) {
  if (direction <= 0) {
    return new _Position.position(row, col);
  }
  return rTile(direction - 1, col, SIZE - 1 - row);
}

/// Make left direction movement as default to fill grid without need rotating grid.
/// When other direction movement command input,rotate grid with new layout and then
/// fill grid according left direction movement.
/// double loop sequence: outer-row, inner-col.
/// arrow key     direction  outer-inner loop  move  ||  after rotation  move
///   left            0        row++,col++      ←   ||   row++,col++     ←
///   right           2        row--,col--      →   ||   row++,col++     ←
///   up              1        col--,row++      ↑   ||   row++,col++     ←
///   down            3        col++,row--      ↓   ||   row++,col++     ←
/// Do clock-anti rotation.
///                 left(0)       up(1)      right(2)    down(3)
///                     ←             ←         ←          ←
///                   1   2         2    4      4    3      3    1
///                 ↓      ↑
///                   3   4         1    3      2    1      4    2
///                     →
bool moveOneStep(int direction) {
  var moveCount = 0;
  var moveCountInSameCol = 0;
  clearStatus();
  for (var row = 0; row < SIZE; row++) {
    var pos = 0; // Location of non-empty tile.
    var flag = true; // The leftmost empty location markup.
    moveCountInSameCol = 0;
    for (var col = 0; col < SIZE; col++) {
      var position = rTile(direction, row, col);
      var rrow = position.Row;
      var rcol = position.Col;
      var currentTile = tile[rrow][rcol];
      position = rTile(direction, row, pos);
      var prow = position.Row;
      var pcol = position.Col;
      var posTile = tile[prow][pcol];
      if (currentTile.Value != 0) {
        if (posTile.Value == 0) {
          // Current tile[row][col] is the leftmost non-empty tile.
          posTile.Value = currentTile.Value;
          posTile.PrePosition = new _Position(rrow, rcol, currentTile.Value);
          currentTile.PrePosition =
              new _Position(prow, pcol, currentTile.Value);
          currentTile.Value = 0;
          flag = false;
          moveCount++;
          moveCountInSameCol++;
          continue;
        }
        if (currentTile.Value == posTile.Value && posTile.MergedFrom == null) {
          if (pos != col) {
            // Merge equal tiles.
            posTile.MergedFrom = new List<_Position>(2);
            posTile.MergedFrom[0] = posTile.PrePosition;
            posTile.MergedFrom[1] =
                new _Position(rrow, rcol, currentTile.Value);
            posTile.Value = currentTile.Value + posTile.Value;
            posTile.PrePosition = new _Position(rrow, rcol, currentTile.Value);
            currentTile.PrePosition =
                new _Position(prow, pcol, currentTile.Value);
            currentTile.Value = 0;
            blankCount += 1;
            moveCount++;
            moveCountInSameCol++;
            if (maxValue < posTile.Value) {
              maxValue = posTile.Value;
            }
            continue;
          }
        } else {
          if ((pos + 1) != col) {
            position = rTile(direction, row, pos + 1);
            prow = position.Row;
            pcol = position.Col;
            posTile = tile[prow][pcol];
            posTile.Value = currentTile.Value;
            posTile.PrePosition = new _Position(rrow, rcol, currentTile.Value);
            currentTile.PrePosition =
                new _Position(prow, pcol, currentTile.Value);
            currentTile.Value = 0;
            pos = pos + 1;
            moveCount++;
            moveCountInSameCol++;
            continue;
          }
        }
        if (moveCountInSameCol == 0) {
          posTile.PrePosition = new _Position(prow, pcol, posTile.Value);
        }
        pos = col;
      } // if
      if (flag) {
        // Set the leftmost empty location.
        pos = col;
        flag = false;
      }
    } // for col
  } // for row
  return moveCount > 0;
} // moveOneStep

List move(String direction) {
  switch (direction) {
    case "up":
      moveDirection(1);
      break;
    case "down":
      moveDirection(3);
      break;
    case "left":
      moveDirection(0);
      break;
    case "right":
      moveDirection(2);
      break;
  }
  return tile;
}

bool canMove() {
  if (blankCount == 0) {
    for (var col = 0; col < SIZE; col++) {
      for (var row = 0; row < SIZE - 1; row++) {
        if (tile[row][col].Value == tile[row + 1][col].Value) {
          return true;
        }
      }
    }
    for (var row = 0; row < SIZE; row++) {
      for (var col = 0; col < SIZE - 1; col++) {
        if (tile[row][col].Value == tile[row][col + 1].Value) {
          return true;
        }
      }
    }
    return false;
  }
  return true;
}

String Win() {
  if (maxValue == WIN) {
    return 'win';
  }
  if (!canMove() && maxValue != WIN) {
    return 'lose';
  }
  return 'wait';
}

Reset() {
  maxValue = 0;
  blankCount = SIZE * SIZE;
}
