function removeFriend(inpid, rem, name, elem){
	if(window.confirm("Are you sure you want to remove the friend '"+decodeURIComponent(name)+"'?")){
		if(elem!=undefined && elem.parentElement!=undefined){
			elem.parentElement.removeChild(elem);
		}
		var inp = document.getElementById(inpid);
		var arr = eval(inp.value);
		if(arr==undefined){
			arr = [];
		}
		var ind = arr.indexOf(rem);
		if(ind>=0){
			arr.splice(ind,1);
		}
		inp.value = "[";
		for(var i=0; i<arr.length; ++i){
			inp.value += (i>0?', "':'"')+arr[i]+'"';
		}
		inp.value += "]";
	}
}

function removeGame(inpid, rem, name, elem){
	if(window.confirm("Are you sure you want to remove the game '"+decodeURIComponent(name)+"'?")){
		if(elem!=undefined && elem.parentElement!=undefined){
			elem.parentElement.removeChild(elem);
		}
		var inp = document.getElementById(inpid);
		var arr = eval(inp.value);
		if(arr==undefined){
			arr = [];
		}
		var ind = arr.indexOf(rem);
		if(ind>=0){
			arr.splice(ind,1);
		}
		inp.value = "[";
		for(var i=0; i<arr.length; ++i){
			inp.value += (i>0?', "':'"')+arr[i]+'"';
		}
		inp.value += "]";
	}
}
