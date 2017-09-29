void swap(List<int> a, int i, int j) {
  var tmp = a[i];
  a[i] = a[j];
  a[j] = tmp;
}

void all(List<int> s, int n) {
  if (n == 1) {
    //print(new String.fromCharCodes(s));
  } else {
    for (var i = 0; i < n; i++) {
      swap(s, i, n - 1);
      all(s, n - 1);
      swap(s, i, n - 1);
    }
  }
}

void main() {
  String a = "♠行到水穷处，坐看云起时";
  var b = a.runes.toList();
  all(b, b.length);
}
