library mystdin;

import 'dart:async';
import 'dart:io' show Platform;
import 'dart:io';
import 'asyn_extension.dart';

class MyStdin {
  Future<String> readLine(
      {Encoding encoding: SYSTEM_ENCODING, bool retainNewlines: false}) {
    const CR = 13;
    const LF = 10;
    final List<int> line = <int>[];
    List<int> _buf = [];
    int pos = 0;
    // On Windows, if lineMode is disabled, only CR is received.
    bool crIsNewline = Platform.isWindows &&
        (stdioType(stdin) == StdioType.TERMINAL) &&
        !lineMode;
    int byte;
    ReadBytes reader = new ReadBytes();
    return reader.readBytes().then((values) {
      if (values != null) {
        _buf.addAll(values);
        _buf.add(-1);
        if (retainNewlines) {
          byte = _buf[pos++];
          while (byte != LF && !(byte == CR && crIsNewline)) {
            if (byte < 0) {
              if (line.isEmpty) return null;
              break;
            }
            line.add(byte);
            byte = _buf[pos++];
          }
        } else if (crIsNewline) {
          // CR and LF are both line terminators, neither is retained.
          byte = _buf[pos++];
          while (byte != LF && byte != CR) {
            if (byte < 0) {
              if (line.isEmpty) return null;
              break;
            }
            line.add(byte);
            byte = _buf[pos++];
          }
        } else {
          // Case having to handel CR LF as a single unretained line terminator.
          byte = _buf[pos++];
          while (byte != LF) {
            if (byte < 0) {
              if (line.isEmpty) return null;
              break;
            }
            line.add(byte);
            byte = _buf[pos++];
          }
        }
      }
      return encoding.decode(line);
    });
  }
}
