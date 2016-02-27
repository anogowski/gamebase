function removeFriend(inpid, rem, elem){
	if(window.confirm("Are you sure you want to remove the friend '"+rem+"'?")){
		if(elem!=undefined){
			document.removeElement(elem);
		}
		//TODO: remove the friend from inpid
		
	}
}

function removeGame(inpid, rem, elem){
	if(window.confirm("Are you sure you want to remove the game '"+rem+"'?")){
		if(elem!=undefined){
			document.removeElement(elem);
		}
		//TODO: remove the game from inpid
		
	}
}
