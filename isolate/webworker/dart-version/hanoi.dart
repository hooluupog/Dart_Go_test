import 'dart:isolate';

SendPort replyTo;
var count = 0;
var text;

void main(List<String> args, SendPort reply) {
  var sw = new Stopwatch();
  var port = new ReceivePort();
  // 向其他isolates通知自己的接收端口.
  reply.send(port.sendPort);
  replyTo = reply;
  port.listen((msg) {
    if (msg == 'OK') {
      sw.stop();
      print('${sw.elapsedMilliseconds}ms.');
      text = '<li>time used: ${sw.elapsedMilliseconds}ms.</li>';
      reply.send(text);
    } else {
      sw.reset();
      sw.start();
      count = 0;
      hanoi(msg, 'A', 'B', 'C');
      reply.send('OVER');
    }
  });
}

void hanoi(n, x, y, z) {
  if (n == 1) {
    count++;
    text = '<li>$count.move disk $n from $x to $z</li>';
    replyTo.send(text);
  } else {
    hanoi(n - 1, x, z, y);
    count++;
    text = '<li>$count.move disk $n from $x to $z</li>';
    replyTo.send(text);
    hanoi(n - 1, y, x, z);
  }
}
