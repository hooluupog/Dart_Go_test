import 'dart:html';
// import 'dart:typed_data';

// dart-enabled web worker does not be supported at present.
// https://github.com/dart-lang/sdk/issues/21573 
main() {
  window.onMessage.listen((e) {
    var view = e.data.asUint8List();
    var res = 0;
    for (var i in view) {
      if (i % 2 == 0) res += i;
    }
    window.postMessage(res,'*');
    window.close();
  });
}
