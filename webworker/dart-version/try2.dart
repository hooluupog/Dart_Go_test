import 'dart:html';
import 'dart:isolate';


void main() {
  SendPort worker = spawnUri('hanoi.dart');
  port.receive((msg, SendPort reply) {
    if(msg == 'Over'){
      worker.send('OK', port.toSendPort());
    }
    else {
      query("#sample_text_id").innerHtml += msg;
    }
  });
  var disknum = 0;
  query("#title_id").text = 'Tower Of Hanoi disk numbers = $disknum';
  var runbutton = query("#run");
  runbutton.onClick.listen((e){
    disknum++;
    query("#title_id").text = 'Tower Of Hanoi disk numbers = $disknum';
    worker.send(disknum, port.toSendPort());
  });
  var clearbutton = query("#clear");
    clearbutton.onClick.listen((e){
    disknum = 0;
    query("#title_id").text = 'Tower Of Hanoi disk numbers = $disknum';
    query("#sample_text_id").innerHtml = '';
  });
}