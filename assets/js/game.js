function queryEscape(str){
	return encodeURIComponent(str);
}
function queryUnescape(str){
	return decodeURIComponent(str);
}

function addTag(inpid, add, elem, outpid, outpidnext){
	if(outpid!=undefined){
		var outelem = document.createElement('a');
		outelem.href = '#';
		if(outpidnext==undefined && elem!=undefined && elem.parentElement!=undefined){
			outpidnext = elem.parentElement.id;
		}
		var newout = outpidnext;
		outelem.onclick = function(){removeTag(inpid, add, this, newout);}
		outelem.innerHTML = queryUnescape(add);
		document.getElementById(outpid).appendChild(outelem);
	}
	if(elem!=undefined && elem.parentElement!=undefined){
		elem.parentElement.removeChild(elem);
	}
	var inp = document.getElementById(inpid);
	var arr = eval(inp.value)
	if(arr==undefined){
		arr = [];
	}
	arr.push(add);
	inp.value = "[";
	for(var i=0; i<arr.length; ++i){
		inp.value += (i>0?", '":"'")+arr[i]+"'";
	}
	inp.value += "]";
}

function removeTag(inpid, rem, elem, outpid, outpidnext){
	if(window.confirm("Are you sure you want to remove the tag '"+rem+"'?")){
		if(outpid!=undefined){
			var outelem = document.createElement('a');
			outelem.href = '#';
			if(outpidnext==undefined && elem!=undefined && elem.parentElement!=undefined){
				outpidnext = elem.parentElement.id;
			}
			var newout = outpidnext;
			outelem.onclick = function(){addTag(inpid, rem, this, newout);}
			outelem.innerHTML = queryUnescape(rem);
			document.getElementById(outpid).appendChild(outelem);
		}
		var newout = outpidnext;
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
			inp.value += (i>0?", '":"'")+arr[i]+"'";
		}
		inp.value += "]";
	}
}
