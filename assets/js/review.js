function drawReviewStars(parelem, valelem, amnt){
	valelem.value = amnt;
	parelem.innerHTML = '';
	var non = document.createElement('a');
	non.className = 'white';
	non.innerText = '0';
	non.href = '#';
	non.value = 0;
	non.onclick = function(){drawReviewStars(this.parentElement, valelem, this.value);};
	parelem.appendChild(non);
	for(var i=1; i<=amnt; ++i){
		var star = document.createElement('a');
		star.className = 'gold';
		star.innerHTML = '&#x2605;'
		star.href = '#';
		star.value = i;
		star.onclick = function(){drawReviewStars(this.parentElement, valelem, this.value);};
		parelem.appendChild(star);
	}
	for(var i=amnt+1; i<=5; ++i){
		var star = document.createElement('a');
		star.className = 'white';
		star.innerHTML = '&#x2606;'
		star.href = '#';
		star.value = i;
		star.onclick = function(){drawReviewStars(this.parentElement, valelem, this.value);};
		parelem.appendChild(star);
	}
}
