self.onmessage = function(msg){
    var view = new Uint8Array(msg.data);
    var res = 0;
    for(var i = 0;i < view.length;i++){
        if (view[i] % 2 == 0) res += view[i];
    }
    postMessage(res);
    close();
}
