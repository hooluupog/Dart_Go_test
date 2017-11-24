import 'dart:html';
import 'dart:typed_data';

void run() {
  var n = 1024 * 1024 * 320; // 320M buffer.
  var l = new Uint8List(n);
  for (var i = 0; i < l.length; i++) l[i] = i % 256;
  var buffer = l.buffer;
  var sw = new Stopwatch()..start();
  final worker = new Worker('psum.js');
  //final worker = new Worker('psum.dart');
  worker.onMessage.listen((e) {
    //window.alert(l.length);
    sw.stop();
    var text =
        '<li>sum = ${e.data}  time used: ${sw.elapsedMilliseconds}ms.</li>';
    querySelector('#sample_text_id').innerHtml = text;
    worker.terminate();
  });
  querySelector('#title_id').text = 'The sum of evens within $n';
  worker.postMessage(buffer, [buffer]);
  //worker.postMessage(buffer);
}

void main() {
  querySelector('#run').onClick.listen((e) {
    run();
  });
}
