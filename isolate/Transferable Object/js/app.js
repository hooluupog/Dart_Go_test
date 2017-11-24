document.querySelector('button').addEventListener('click', startWorker);
function startWorker() {
    var arrayBuffer = new ArrayBuffer(1024*1024*320);
    var view = new Uint8Array(arrayBuffer);
    //var sab = new SharedArrayBuffer(1024*1024*320);
    //var view = new Uint8Array(sab);
    var start = performance.now();
    var w =new Worker('psum.js');
    w.onmessage = function (event) {
        //alert(view.length);
        var elapsed = performance.now() - start;
        document.getElementById("result").innerHTML=event.data + '  time used: ' +  elapsed + 'ms';
        w.terminate();
    };
    for (var i = 0; i < view.length; i++) {
        view[i] = i % 256;
    }
    w.postMessage(view.buffer,[view.buffer]);//transerable:ON
    //w.postMessage(view.buffer);// transerable:OFF
    //w.postMessage(sab);// SharedArrayBuffer
}

