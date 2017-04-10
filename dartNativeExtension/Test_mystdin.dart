import 'dart:convert';
import 'mystdin/mystdin.dart';

void main() async {
  MyStdin mystdin = new MyStdin();
  var l = (await mystdin.readLine()).split(new RegExp("\\s+"));
  print(l);
}
