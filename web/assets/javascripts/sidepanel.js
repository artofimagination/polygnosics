/* Set the width of the sidebar to 250px (show it) */
function openNav() {
  document.getElementById("sidePanel").style.width = "250px";
  document.getElementById("main").style.marginLeft = "250px";
}

/* Set the width of the sidebar to 0 (hide it) */
function closeNav() {
  document.getElementById("sidePanel").style.width = "0";
  document.getElementById("main").style.marginLeft = "0";
}