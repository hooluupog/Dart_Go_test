document.onkeydown=function(e){
	var isIE = (document.all) ? true:false;
	var key;
	if(isIE){
		key = window.event.keyCode;
	}else{
		key = e.which;
	} 
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
};

var continuePlay = false;
function showMessage(){
	var elem = document.getElementById("result");
	if (navigator.appName.indexOf("Explorer")>-1){
		vartext = elem.innerText;
	}else{
		vartext = elem.textContent;
	}
	if(vartext === "lose"){
		elem.style.display='block';
	}
	if(vartext === "win"){
		if (!continuePlay){
			elem.style.display='block';
		}else{
			elem.style.display='none';
		}
		continuePlay = true;
	}
}

var xmlreq;
function move(direction){
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = updateGrid;
	req.open("POST","/"+ direction,true);
	req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
	req.send(direction);
}

function updateGrid() {
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	document.getElementById("output").innerHTML = req.responseText;
	showMessage();
}
