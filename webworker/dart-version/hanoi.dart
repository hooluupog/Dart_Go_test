import 'dart:isolate';

var replyTo;
var count = 0;
var text;

void main() {
  var sw = new Stopwatch();
  port.receive((msg, SendPort reply) {
    count = 0;
    replyTo = reply;
    if(msg == 'OK')
    {
      sw.stop();
      print('${sw.elapsedMilliseconds}ms.');
      text = '<li>time used: ${sw.elapsedMilliseconds}ms.</li>';
      reply.send(text,port.toSendPort());
    }
    else{
      sw.reset();
      sw.start();
      hanoi(msg,'A','B','C');
      reply.send('Over',port.toSendPort());
    }
  });
}

void hanoi(n,x,y,z) {
  if (n == 1){
    count++; 
    text = '<li>$count.move disk $n from $x to $z</li>';
    replyTo.send(text,port.toSendPort());
  }
  else{
    hanoi(n-1,x,z,y);
    count++;
    text = '<li>$count.move disk $n from $x to $z</li>';
    replyTo.send(text,port.toSendPort());
    hanoi(n-1,y,x,z);
  }
}