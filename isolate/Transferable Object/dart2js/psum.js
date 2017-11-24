self.onmessage = function(e) {
    var view = new Uint8Array(e.data);
    var res = 0;
    for (var i = 0;i < view.length;i++) {
        if (view[i] % 2 == 0) res += view[i];
    }
    postMessage(res);
    close();
}; 

