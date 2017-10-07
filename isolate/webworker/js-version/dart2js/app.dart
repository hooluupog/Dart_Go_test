import 'dart:html';
import 'dart:isolate';
import 'dart:async';

void main() async {
  var port = new ReceivePort();
  await Isolate.spawnUri(Uri.parse('hanoi.dart.js'), [], port.sendPort);
  SendPort sendPort;
  port.listen((msg) {
    if (sendPort == null) {
      sendPort = msg;
    } else {
      if (msg == 'OVER') {
        sendPort.send('OK');
      } else {
        querySelector('#sample_text_id').innerHtml += msg;
      }
    }
  });
  var disknum = 0;
  querySelector('#title_id').text = 'Tower Of Hanoi disk numbers = $disknum';
  var runbutton = querySelector('#run');
  runbutton.onClick.listen((e) {
    disknum++;
    querySelector('#title_id').text = 'Tower Of Hanoi disk numbers = $disknum';
    sendPort.send(disknum);
  });
  var clearbutton = querySelector('#clear');
  clearbutton.onClick.listen((e) {
    disknum = 0;
    querySelector('#title_id').text = 'Tower Of Hanoi disk numbers = $disknum';
    querySelector('#sample_text_id').innerHtml = '';
  });
}
