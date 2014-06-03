var stopkeyboardListener = false;
document.onkeydown=function(e){
    var isIE = (document.all) ? true:false;
    var key;
    if(isIE){
        key = window.event.keyCode;
    }else{
        key = e.which;
    } 
    if(!stopkeyboardListener){
        if(key==37){
            move('left');
        }
        if(key==38){
            move('up');
        }
        if(key==39){
            move('right');
        }
        if(key==40){
            move('down');
        }
    }
};
function move(direction){
    ajax(direction);
}
