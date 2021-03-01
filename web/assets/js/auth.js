
!function($) {
  "use strict";

  var SweetAlert = function() {};
  SweetAlert.prototype.init = function() {
    $('#signup').submit(function(){
      $.ajax({
        type: 'POST',
        url: $(this).attr('action'),
        data: $(this).serialize(),
        dataType: 'json',
        error: function(data) {
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