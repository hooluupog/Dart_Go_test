var stopkeyboardListener = false;
var map = {
    37:1,
    38:1,
    39:1,
    40:1
};
var sum = 0;
document.onkeydown=function(e){
    var isIE = (document.all) ? true:false;
    var key;
    if(isIE){
        key = window.event.keyCode;
    }else{
        key = e.which;
    } 
    if (key in map){
        sum += map[key];
        if(sum > 1){
            //FIRE EVENT
            // Detect multiple keys pressed at the same time and ignore this case.
        }else{
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
        }
    }
};
document.onkeyup=function(e){
    var isIE = (document.all) ? true:false;
    var key;
    if(isIE){
        key = window.event.keyCode;
    }else{
        key = e.which;
    } 
    if (key in map){
        sum -= map[key]; 
    }
};
function move(direction){
    ajax(direction);
}
