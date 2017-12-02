import 'dart:io';
import 'dart:convert';

const String ANSI_CODE = "\x1b";
const String ANSI_ESCAPE = "\x1b[";

main() {
  stdin.echoMode = false;
  stdin.lineMode = false;
  stdin.asBroadcastStream().map((bytes) {
    var it = ASCII.decode(bytes);
    var original = bytes;
    var code = it.replaceAll(ANSI_CODE, "");

    if (code.isNotEmpty) {
      code = code.substring(1);
    }
    return [original, it, code];
  }).listen((List<dynamic> m) {
    print('${m[0]}, ${m[1]}, ${m[2]}');
  });
}
