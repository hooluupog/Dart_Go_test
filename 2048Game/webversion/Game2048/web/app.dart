import 'dart:html';
import 'globals.dart' as globals;
import 'game_manager.dart' show handle, initGrid;

move(direction) {
  handle(direction);
}

void main() {
  initGrid();
  window.onKeyUp.listen((e) {
    if (!globals.stopkeyboardListener) {
      switch (e.keyCode) {
        case KeyCode.UP:
          move('up');
          break;
        case KeyCode.DOWN:
          move('down');
          break;
        case KeyCode.LEFT:
          move('left');
          break;
        case KeyCode.RIGHT:
          move('right');
          break;
      }
    }
  });
}
