function queryEscape(str){
	//TODO: escape str for use in a query
	return encodeURIComponent(str);
}

function addTag(inpid, add, elem){
	//if(elem!=undefined){
	//	document.removeElement(elem);
	//}
	//TODO: add the tag from inpid
	
}

function removeTag(inpid, rem, elem){
	if(window.confirm("Are you sure you want to remove the tag '"+rem+"'?")){
		if(elem!=undefined){
			document.removeElement(elem);
		}
		//TODO: remove the tag from inpid
		
	}
}
