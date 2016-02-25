function removeFriend(inpid, rem, elem){
	if(window.confirm("Are you sure you want to remove the friend '"+rem+"'?")){
		document.removeElement(elem);
		//TODO: remove the friend from inpid
		
	}
}

function removeGame(inpid, rem, elem){
	if(window.confirm("Are you sure you want to remove the game '"+rem+"'?")){
		document.removeElement(elem);
		//TODO: remove the game from inpid
		
	}
}
