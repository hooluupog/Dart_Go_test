var count = 0;
function hanoi(n,x,y,z) {
  if (n == 1){
    count++;
    postMessage('<li>' + count + '.move disk ' + n + ' from ' + x + ' to ' + z + '</li>');
  }
  else{
    hanoi(n-1,x,z,y);
    count++;
    postMessage('<li>' + count + '.move disk ' + n + ' from ' + x + ' to ' + z + '</li>');
    hanoi(n-1,y,x,z);
  }
}
hanoi(12,'A','B','C');
