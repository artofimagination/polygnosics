
!function($) {
  "use strict";
  $("#signup-loader").get(0).style.display = "none"
  $("#signup-dialog").get(0).style.display = "block"

  var SweetAlert = function() {};
  SweetAlert.prototype.init = function() {
    $('#signup').submit(function(){
      $("#signup-loader").get(0).style.display = "block"
      $("#signup-dialog").get(0).style.display = "none"
      $.ajax({
        type: 'POST',
        url: $(this).attr('action'),
        data: $(this).serialize(),
        dataType: 'json',
        error: function(data) {
          $("#signup-loader").get(0).style.display = "none"
          $("#signup-dialog").get(0).style.display = "block"
          var alertType  = "error"
          if (data.responseText == "Registration successful")
          {
            alertType = "success"
          }

          swal({
            title: "Registration", 
            text: data.responseText, 
            type: alertType
            },
            function(){
              window.location.href = "/index";
            }) 
        }
      })
      return false;        
    });
  },
  //init
  $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
}(window.jQuery),

//initializing 
function($) {
  "use strict";
  $.SweetAlert.init()
}(window.jQuery);