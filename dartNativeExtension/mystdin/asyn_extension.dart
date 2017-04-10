library asyn_extension;

import 'dart:async';
import 'dart:isolate';
import 'dart-ext:mystdin';

// A class caches the native port used to call an asynchronous extension.
class ReadBytes {
  static SendPort _port;

  Future<int> readBytes() {
    var completer = new Completer();
    var replyPort = new RawReceivePort();
    var args = new List(1);
    args[0] = replyPort.sendPort;
    _servicePort.send(args);
    replyPort.handler = (result) {
      replyPort.close();
      if (result != null) {
        completer.complete(result);
      } else {
        completer.completeError(new Exception("Random array creation failed"));
      }
    };
    return completer.future;
  }

  SendPort get _servicePort {
    if (_port == null) {
      _port = _newServicePort();
    }
    return _port;
  }

  SendPort _newServicePort() native "ReadBytes_ServicePort";
}
