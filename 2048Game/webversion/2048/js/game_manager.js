var xmlreq;
function ajax(direction){
    var req = new XMLHttpRequest();
    xmlreq = req;
    if (direction != null) {
        req.onreadystatechange = updateGrid;
        req.open("POST","/"+ direction,true);
        req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
        req.send(direction);
    }else{
        req.onreadystatechange = function(){
            var req = xmlreq;
            if(!req || req.readyState != 4) {
                return;
            }
            document.getElementById("output").innerHTML = req.responseText;
        };
        req.open("POST","/grid",true); 
        req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
        req.send(null);
    }
}

var continueAfterwon = false;
function showMessage(){
    var elem = document.getElementById("result");
    if (navigator.appName.indexOf("Explorer")>-1){
        vartext = elem.innerText;
    }else{
        vartext = elem.textContent;
    }
    if(vartext === "lose"){
        var elem = document.getElementsByClassName("message")[0];
        stopkeyboardListener = true;
        elem.className = "message lose"; 
        var restart = document.getElementsByClassName("retry-button")[0];
        restart.onclick = function(){
            elem.className = "message"; 
            elem.style.display='none';
            continueAfterwon = false;
            stopkeyboardListener = false;
            ajax();
        };
    }
    if(vartext === "win" && !continueAfterwon){
        var elem = document.getElementsByClassName("message")[0];
        stopkeyboardListener = true;
        elem.className = "message win"; 
        var continuePlay = document.getElementsByClassName("keep-playing-button")[0];
        continuePlay.onclick = function(){
            elem.className = "message"; 
            document.getElementById("result").innerHTML = "wait";
            continueAfterwon = true;
            stopkeyboardListener = false;
        };
        var restart = document.getElementsByClassName("retry-button")[0];
        restart.onclick = function(){
            elem.className = "message"; 
            continueAfterwon = false;
            stopkeyboardListener = false;
            ajax();
        };
    }
}
function TileStatus(elem,row,col,value){
    this.elem = elem;
    this.row = row;
    this.col = col;
    this.value = value;
}
function clearstatus(list){
    while((tilestatus = list.pop()) != null){
        tilestatus.elem.className = "tile tile-" + tilestatus.row + "-" + tilestatus.col + " tile-" + tilestatus.value;
    }
}

var newTileList = new Array();
var mergedTileList = new Array();
function updateGrid() {
    var req = xmlreq;
    if(!req || req.readyState != 4) {
        return;
    }
    var object = JSON.parse(req.response);
        clearstatus(newTileList);
        clearstatus(mergedTileList);
        document.getElementById("result").innerHTML = object.Win;
    for (var row = 0;row < object.Size; row++){
        for(var col = 0; col < object.Size;col++){
            var data = object.Tile[row][col];
            var pre = data.PrePosition.Value;
            if(data.Value != 0){
                if(data.MergedFrom != null){
                    var ob1 = data.MergedFrom[0];
                    var ob2 = data.MergedFrom[1];
                        var elem1 = document.getElementsByClassName("tile tile-" + ob1.X + "-" + ob1.Y + " tile-" + ob1.Value)[0];
                        var elem2 = document.getElementsByClassName("tile tile-" + ob2.X + "-" + ob2.Y + " tile-" + ob2.Value)[0];
                        if(elem1 != null){
                            elem1.className = "tile tile-" + row + "-" + col + " tile-" + ob1.Value;
                        }
                        if(elem2 != null){
                            elem2.className = "tile tile-" + row + "-" + col + " tile-" + ob2.Value;
                        }
                        var node = document.getElementsByClassName("tile tile-" + row + "-" + col + " tile-" + pre);
                        var len = node.length;
                        for (var i = len-1; i >= 0;i--){
                            node[i].parentNode.removeChild(node[i]);
                        }
                        var newElem = document.createElement("div");
                        newElem.className =  "tile tile-" + row +"-" + col + " tile-" + data.Value + " mergedTile";
                        mergedTileList.push(new TileStatus(newElem,row,col,data.Value));
                        document.getElementsByClassName("tile-container")[0].appendChild(newElem); 
                        newElem.innerHTML = '<div class="tile-inner">'+data.Value+'</div>';
                }else{
                        if (data.IsNewTile){// New generated random tile.
                            var newElem = document.createElement("div");
                            newElem.className = "tile tile-"+ row +"-"+ col+" tile-" + data.Value + " newTile"; 
                            newTileList.push(new TileStatus(newElem,row,col,data.Value));
                            document.getElementsByClassName("tile-container")[0].appendChild(newElem); 
                            newElem.innerHTML = '<div class="tile-inner">'+data.Value+'</div>';
                        }else{
                            var elem = document.getElementsByClassName("tile tile-" + data.PrePosition.X + "-" + data.PrePosition.Y + " tile-" + pre)[0];
                            if(elem != null){
                                elem.className = "tile tile-" + row + "-" + col + " tile-" + data.Value;
                            }
                        }
                }
            }
        }
    }
    showMessage();
}
