import 'dart:html';
import 'globals.dart' as globals;
import 'package:Game2048/grid.dart';

handle(String direction) {
  updateGrid(direction);
}

renderGrid(List tile) {
  var text = '';
  for (var row = 0; row < tile.length; row++) {
    for (var col = 0; col < tile.length; col++) {
      if (tile[row][col].Value != 0) {
        text += '<div class=\'tile tile-$row-$col tile-${tile[row][col].Value} newTile\'>' + 
            '<div class=\'tile-inner\'>${tile[row][col].Value}</div></div>';
      }
    }
  }
  querySelector('.tile-container').innerHtml = text;
}

initGrid() {
  var tile = InitGrid();
  renderGrid(tile);
}

var continueAfterwon = false;

showMessage() {
  var elem = querySelector('#result');
  var text = elem.text;
  if (text == 'lose') {
    var elem = querySelector('.message');
    globals.stopkeyboardListener = true;
    elem.className = 'message lose';
    var restart = querySelector('.retry-button');
    restart.onClick.listen((e) {
      elem.className = 'message';
      continueAfterwon = false;
      globals.stopkeyboardListener = false;
      initGrid();
    });
  }
  if (text == 'win' && !continueAfterwon) {
    var elem = querySelector('.message');
    globals.stopkeyboardListener = true;
    elem.className = 'message win';
    var continuePlay = querySelector('.keep-playing-button');
    continuePlay.onClick.listen((e) {
      elem.className = 'message';
      querySelector('#result').innerHtml = 'wait';
      continueAfterwon = true;
      globals.stopkeyboardListener = false;
    });
    var restart = querySelector('.retry-button');
    restart.onClick.listen((e) {
      elem.className = 'message';
      continueAfterwon = false;
      globals.stopkeyboardListener = false;
      initGrid();
    });
  }
}

class _TileStatus {
  Element _elem;
  int _row;
  int _col;
  int _value;
  _TileStatus(this._elem, this._row, this._col, this._value);
  Element get elem => _elem;
  int get row => _row;
  int get col => _col;
  int get value => _value;
}

clearstatus(List list) {
  var tilestatus;
  while (list.length != 0) {
    tilestatus = list.removeLast();
    tilestatus.elem.className = 'tile tile-${tilestatus.row}-${tilestatus.col} tile-${tilestatus.value}';
  }
}

var newTileList = [];
var mergedTileList = [];

updateGrid(String direction) {
  var object = move(direction);
  clearstatus(newTileList);
  clearstatus(mergedTileList);
  querySelector('#result').innerHtml = Win();
  for (var row = 0; row < object.length; row++) {
    for (var col = 0; col < object.length; col++) {
      var data = object[row][col];
      var pre = data.PrePosition.Value;
      if (data.Value != 0) {
        if (data.MergedFrom != null) {
          var ob1 = data.MergedFrom[0];
          var ob2 = data.MergedFrom[1];
          var elem1 = querySelector('.tile.tile-${ob1.Row}-${ob1.Col}.tile-${ob1.Value}');
          var elem2 = querySelector('.tile.tile-${ob2.Row}-${ob2.Col}.tile-${ob2.Value}');
          if (elem1 != null) {
            elem1.className = 'tile tile-$row-$col tile-${ob1.Value}';
          }
          if (elem2 != null) {
            elem2.className = 'tile tile-$row-$col tile-${ob2.Value}';
          }
          var node = querySelectorAll('.tile.tile-$row-$col.tile-$pre');
          var len = node.length;
          for (var i = len - 1; i >= 0; i--) {
            node[i].remove();
          }
          var newElem = document.createElement('div');
          newElem.className = 'tile tile-$row-$col tile-${data.Value} mergedTile';
          mergedTileList.add(new _TileStatus(newElem, row, col, data.Value));
          querySelector('.tile-container').append(newElem);
          newElem.innerHtml = '<div class=\'tile-inner\'>${data.Value}</div>';
        } else {
          if (data.IsNewTile) {
            // New generated random tile.
            var newElem = document.createElement('div');
            newElem.className = 'tile tile-$row-$col tile-${data.Value} newTile';
            newTileList.add(new _TileStatus(newElem, row, col, data.Value));
            querySelector('.tile-container').append(newElem);
            newElem.innerHtml = '<div class=\'tile-inner\'>${data.Value}</div>';
          } else {
            var elem = querySelector('.tile.tile-${data.PrePosition.Row}-${data.PrePosition.Col}.tile-$pre');
            if (elem != null) {
              elem.className = 'tile tile-$row-$col tile-${data.Value}';
            }
          }
        }
      }
    }
  }
  showMessage();
}
