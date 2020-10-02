var modal = document.getElementById('id01');
var password = document.getElementById("psw");
var passwordRepeat = document.getElementById("psw-repeat");
var email = document.getElementById("email");
var username = document.getElementById("username");

var letter = document.getElementById("letter");
var capital = document.getElementById("capital");
var number = document.getElementById("number");
var length = document.getElementById("length");
var confirmation = document.getElementById("confirmation");
var emailValid = document.getElementById("email-valid");
var usernameValid = document.getElementById("username-valid");

// Get the modal
var modalLogin = document.getElementById("id02");

// Get the <span> element that closes the modal
var span = document.getElementsByClassName("close")[0];

// When the user clicks on <span> (x), close the modal
span.onclick = function() {
  modalLogin.style.display = "none";
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modalLogin) {
    modalLogin.style.display = "none";
  }
}

// When the user clicks anywhere outside of the modal, close it
window.onclick = function(event) {
  if (event.target == modal) {
    modal.style.display = "none";
  }
}

// When the user clicks on the password field, show the message box
password.onfocus = function() {
  document.getElementById("message-psw").style.display = "block";
}

// When the user clicks outside of the password field, hide the message box
password.onblur = function() {
  document.getElementById("message-psw").style.display = "none";
}

// When the user clicks on the password confirmation field, show the message box
passwordRepeat.onfocus = function() {
  document.getElementById("message-psw-repeat").style.display = "block";
}

// When the user clicks outside of the password confirmation field, hide the message box
passwordRepeat.onblur = function() {
  document.getElementById("message-psw-repeat").style.display = "none";
}

// When the user clicks on the email field, show the message box
email.onfocus = function() {
  document.getElementById("message-email").style.display = "block";
}

// When the user clicks outside of the email field, hide the message box
email.onblur = function() {
  document.getElementById("message-email").style.display = "none";
}

// When the user clicks on the email field, show the message box
username.onfocus = function() {
  document.getElementById("message-username").style.display = "block";
}

// When the user clicks outside of the email field, hide the message box
username.onblur = function() {
  document.getElementById("message-username").style.display = "none";
}

// When the user starts to type something inside the email field
username.onkeyup = function() {
  // Validate length
  if(username.value.length >= 6) {
    usernameValid.classList.remove("invalid");
    usernameValid.classList.add("valid");
  } else {
    usernameValid.classList.remove("valid");
    usernameValid.classList.add("invalid");
  }
}

// When the user starts to type something inside the email field
email.onkeyup = function() {
  // Validate lowercase letters
  const re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

  if(re.test(String(email.value).toLowerCase())) {  
    emailValid.classList.remove("invalid");
    emailValid.classList.add("valid");
  } else {
    emailValid.classList.remove("valid");
    emailValid.classList.add("invalid");
  }
}

// When the user starts to type something inside the password field
passwordRepeat.onkeyup = function() {
  // Validate lowercase letters
  if(passwordRepeat.value === password.value) {  
    confirmation.classList.remove("invalid");
    confirmation.classList.add("valid");
    document.getElementById("message-psw-repeat").style.display = "none";
  } else {
    confirmation.classList.remove("valid");
    confirmation.classList.add("invalid");
  }
}

// When the user starts to type something inside the password field
password.onkeyup = function() {
  // Validate lowercase letters
  var lowerCaseLetters = /[a-z]/g;
  if(password.value.match(lowerCaseLetters)) {  
    letter.classList.remove("invalid");
    letter.classList.add("valid");
  } else {
    letter.classList.remove("valid");
    letter.classList.add("invalid");
  }
  
  // Validate capital letters
  var upperCaseLetters = /[A-Z]/g;
  if(password.value.match(upperCaseLetters)) {  
    capital.classList.remove("invalid");
    capital.classList.add("valid");
  } else {
    capital.classList.remove("valid");
    capital.classList.add("invalid");
  }

  // Validate numbers
  var numbers = /[0-9]/g;
  if(password.value.match(numbers)) {  
    number.classList.remove("invalid");
    number.classList.add("valid");
  } else {
    number.classList.remove("valid");
    number.classList.add("invalid");
  }
  
  // Validate length
  if(password.value.length >= 8) {
    length.classList.remove("invalid");
    length.classList.add("valid");
  } else {
    length.classList.remove("valid");
    length.classList.add("invalid");
  }
}