import 'dart:io';
import 'dart:async';
import 'dart:isolate';
import 'package:http/http.dart' as http;
import 'package:args/args.dart';

void main(List<String> args) {
  var parsedArgs = parseArgs(args);
  // download(parsedArgs[0], parsedArgs[1]);
  progressiveDownload(parsedArgs[0], parsedArgs[1]);
}

/// Parse the command line args
List<String> parseArgs(List<String> args) {
  var parser = new ArgParser()
    ..addOption('url', abbr: 'u', help: 'The url to download', valueHelp: 'url')
    ..addOption('filename', abbr: 'f', help: 'save filename', valueHelp: 'filename')
    ..addFlag('help', abbr: 'h', help: 'help info');
  var parsedArgs = parser.parse(args);
  if (parsedArgs['help']) {
    print(parser.usage);
    exit(0);
  }

  String url;
  if (parsedArgs['url'] == null) {
    stderr.writeln('A url must be provided');
    exit(1);
  } else {
    url = parsedArgs['url'];
  }

  var filename = parsedArgs['filename'];
  if (filename == null) {
    var uri = Uri.parse(url);
    filename = uri.pathSegments.isNotEmpty ? uri.pathSegments.last : 'output.html';
  }
  return [url, filename];
}

/// Downloads and saves the specified URI
void download(String url, String filename) async {
  var res = await http.get(Uri.parse(url));
  switch (res.statusCode) {
    case 200:
      print('Downloading content of length: ${res.contentLength}');
      new File(filename).writeAsBytes(res.bodyBytes);
      break;
    default:
      stderr.writeln('Unexpected response code : ${res.statusCode}');
  }
}

/// Format file size Bytes show GB,MB,KB,B.
String sizeFormat(num size) {
  String fmtSize = "";
  var KB = size / 1024;
  var MB = KB / 1024;
  var GB = MB / 1024;
  if (KB <= 1)
    fmtSize = size.toStringAsFixed(0).padLeft(7, ' ') + 'B';
  else if (MB <= 1)
    fmtSize = KB.toStringAsFixed(1).padLeft(6, ' ') + 'KB';
  else if (GB <= 1)
    fmtSize = MB.toStringAsFixed(1).padLeft(6, ' ') + 'MB';
  else
    fmtSize = GB.toStringAsFixed(1).padLeft(6, ' ') + 'GB';
  return fmtSize;
}

/// Format total seconds into h:m:s
String timeFormat(int seconds) {
  if (seconds == -1) return ' > 1 day';
  String fmt = '';
  var date = new Duration(seconds: seconds).toString().split(':');
  var hour = int.parse(date[0]);
  var minute = int.parse(date[1]);
  var second = int.parse(date[2].split('.').first);
  if (hour > 24) return ' > 1 day';
  if (hour < 1) fmt = '00' + ':';
  if (hour >= 1 && hour < 10) fmt = '0' + hour.toString() + ':';
  if (hour >= 10) fmt = hour.toString() + ':';
  if (minute < 1) fmt += '00' + ':';
  if (minute >= 1 && minute < 10) fmt += '0' + minute.toString() + ':';
  if (minute >= 10) fmt += minute.toString() + ':';
  if (second < 1) fmt += '00';
  if (second >= 1 && second < 10) fmt += '0' + second.toString();
  if (second >= 10) fmt += second.toString();
  return fmt;
}

/// Progressively downloads and saves the specified URI
void progressiveDownload(String url, String filename) async {
  var receivePort = new ReceivePort();
  await Isolate.spawn(progress, receivePort.sendPort);
  var sendPort = await receivePort.first;

  var req = new http.Request('GET', Uri.parse(url));
  var res = await req.send();
  if (res.statusCode == 200) {
    var fileSize = res.contentLength;
    int currSize = 0;
    print('Downloading file size: ${sizeFormat(fileSize)}');
    sendPort.send(fileSize);
    res.stream.map((chunk) {
      sendPort.send(chunk.length);
      currSize += chunk.length;
      if (currSize == fileSize) {
        sendPort.send('OVER');
      }
      return chunk;
    }).pipe(new File(filename).openWrite());
  } else {
    stderr.writeln('Unexpected response code : ${res.statusCode}');
    res.stream.drain();
  }
}

/// the entry point for the isolate
progress(SendPort sendPort) {
  const oneSec = const Duration(seconds: 1);
  int fileSize = -1;
  int currSize = 0;
  int preSize = 0;
  int chunk = 0;
  int elapsed = 0;
  var port = new ReceivePort();
  sendPort.send(port.sendPort);
  Timer timer;
  port.listen((msg) {
    if (msg == 'OVER') {
      update(chunk, elapsed++, preSize, currSize, fileSize);
      print('\nDownloading Finished!');
      timer.cancel();
      port.close();
    } else if (fileSize == -1) {
      fileSize = msg;
    } else {
      chunk = msg;
      currSize += chunk;
    }
  });
  timer = new Timer.periodic( oneSec, (Timer t) =>
          preSize = update(chunk, elapsed++, preSize, currSize, fileSize));
}

/// print downloading rate of progress
int update(int chunk, int elapsed, int preSize, int currSize, int fileSize) {
  if (fileSize != -1) {
    var present = (currSize / fileSize * 100).toStringAsFixed(2);
    var speed = currSize - preSize;
    int rTime = speed == 0 ? -1 : (fileSize - currSize) ~/ speed;
    stdout.write('\r$present%  ${sizeFormat(currSize)}/${sizeFormat(fileSize).trim()}   ${sizeFormat(speed)}/S  Elapsed[${timeFormat(elapsed)}]  Remain[${timeFormat(rTime)}]');
    preSize = currSize;
  }
  return preSize;
}
