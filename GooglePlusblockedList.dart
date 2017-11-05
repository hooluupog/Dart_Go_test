import 'dart:async';
import 'dart:io';

void main() async {
  var res = await new File('list.html').readAsString();
  var match = new RegExp(r'已屏蔽.*?</div>');
  print('已屏蔽 ${match.allMatches(res).length} 人：');
  match.allMatches(res).map((m) => m[0])
      .forEach((f) => print(f.substring(3, f.indexOf('</div>'))));
}
