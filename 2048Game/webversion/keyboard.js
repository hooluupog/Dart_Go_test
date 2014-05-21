document.onkeydown=function(e){
	var isIE = (document.all) ? true:false;
	var key;
	if(isIE){
		key = window.event.keyCode;
	}else{
		key = e.which;
	} 
	if(key==37){
		moveLeft();
	}
	if(key==38){
		moveUp();
	}
	if(key==39){
		moveRight();
	}
	if(key==40){
		moveDown();
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

function moveLeft(){
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = updateLeft;
	req.open("POST","/left",true);
	req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
	req.send(null);
}

function moveUp(){
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = updateLeft;
	req.open("POST","/up",true);
	req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
	req.send(null);

}

function moveRight(){
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = updateLeft;
	req.open("POST","/right",true);
	req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
	req.send(null);

}

function moveDown(){
	var req = new XMLHttpRequest();
	xmlreq = req;
	req.onreadystatechange = updateLeft;
	req.open("POST","/down",true);
	req.setRequestHeader("Content-Type","text/plain; charset=utf-8");
	req.send(null);

}
function updateLeft() {
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	document.getElementById("output").innerHTML = req.responseText;
	showMessage();
}
function updateUp() {
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	document.getElementById("output").innerHTML = req.responseText;
	showMessage();
}
function updateRight(){
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	document.getElementById("output").innerHTML = req.responseText;
	showMessage();
}
function updateDown() {
	var req = xmlreq;
	if(!req || req.readyState != 4) {
		return;
	}
	document.getElementById("output").innerHTML = req.responseText;
	showMessage();
}
