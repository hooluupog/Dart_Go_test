import 'dart:isolate';
import 'dart:async';
import 'dart:io';
import 'dart:convert';

void main() async {
  var port = new ReceivePort();
  //创建isolate并通过Uri发送自己的接收端口.
  await Isolate.spawnUri(Uri.parse('echo.dart'), ['start'], port.sendPort);
  SendPort sendPort;
  port.listen((msg) {
    if (sendPort == null) {
      sendPort = msg;
      run(sendPort);
    } else {
      if (msg == 'OVER') {
        sendPort.send('OK');
        run(sendPort);
      } else {
        print(msg);
      }
    }
  });
}

void run(SendPort sendPort) {
  var disknum = int.parse(stdin.readLineSync().trim());
  if (disknum < 1) exit(0);
  print('Tower Of Hanoi disk numbers = $disknum');
  sendPort.send(disknum);
}
