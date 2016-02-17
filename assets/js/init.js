function initShortened(){
  var shorts = document.getElementsByClassName('shortened');
  for(var i=0; i<shorts.length; ++i){
    if(shorts[i].clientHeight<shorts[i].scrollHeight || shorts[i].clientWidth<shorts[i].scrollWidth){
      shorts[i].style.position = 'relative';
      var el = document.createElement('button');
      el.style.position = "absolute";
      el.style.top = "0px";
      el.style.right = "0px";
      el.innerText = 'Show more';
      el.onclick = function(){this.parentNode.className=this.parentNode.className.replace(/(?:^|\s)shortened(?!\S)/,''); this.parentNode.removeChild(this);};
      shorts[i].appendChild(el);
    }
  }
}

