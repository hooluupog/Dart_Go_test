import 'dart:html';
import 'dart:typed_data';

void run(int n) {
  var sw = new Stopwatch()..start();
  var l = new Uint8List(n);
  for (var i = 0; i < l.length; i++) l[i] = i % 256;
  var res = l.where((i) => i % 2 == 0).reduce((a, b) => a + b);
  sw.stop();
  var text = '<li>sum = $res  time used: ${sw.elapsedMilliseconds}ms.</li>';
  querySelector('#sample_text_id').innerHtml = text;
  querySelector('#title_id').text = 'The sum of evens within $n';
}

void main() {
  var n = 1024 * 1024;
  querySelector('#run').onClick.listen((e) {
    run(n++);
  });
}
