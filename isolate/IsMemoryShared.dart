import 'dart:isolate';

echo(msg) {
  SendPort replyTo = msg[1];
  replyTo.send(msg[0]);
}

main() {
  var a = [1, 2, 3];
  var c = [];
  var d = [];
  var count = 0;
  for (var i = 0; i < 2; i++) c.add(a);
  print(identical(c[0], a));

  var port = new ReceivePort();
  for (var i = 0; i < 2; i++) Isolate.spawn(echo, [a, port.sendPort]);
  port.listen((msg) {
    d.add(msg);
    count++;
    if (count == 2) {
      port.close();
      print(identical(d[0], d[1]));
      print(identical(d[0], a));
      print(d.any((i) => identical(i, a)));
    }
  });
}
