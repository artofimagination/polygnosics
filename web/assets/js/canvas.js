var c = document.getElementById("world");
var ctx = c.getContext("2d");
ctx.canvas.width  = window.innerWidth;
ctx.canvas.height = window.innerHeight;

function setColor(entityType){
  if (entityType == 0) {
    return 'red'
  } else if (entityType == 1) {
    return 'blue'
  } else if (entityType == 2) {
    return 'green'
  } else if (entityType == 3) {
    return 'yellow'
  } else if (entityType == 4) {
    return 'cyan'
  } else if (entityType == 5) {
    return 'magenta'
  } else if (entityType == 1) {
    return 'black'
  } else {
    return 'white'
  } 
}

function animate(){
  ctx.clearRect(0,0,window.innerWidth,window.innerHeight)
  for (var dataKey in entityData) {
    ctx.beginPath();
    dataJson = JSON.parse(entityData[dataKey])
    ctx.arc(dataJson.posx,dataJson.posy,dataJson.size,0,2*Math.PI);
    ctx.fillStyle = setColor(dataJson.type)
    ctx.fill()
    ctx.stroke();
  }
  requestAnimationFrame(animate)
}

animate()
