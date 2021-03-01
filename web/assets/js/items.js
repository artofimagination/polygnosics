
!function($) {
  "use strict";

  var SweetAlert = function() {};
  SweetAlert.prototype.init = function() {
    var i;
    for (i = 0; i < deleteLinks.length; i++)
    {	
      var stringVal = '#delete-item-' + i;
      $(stringVal).click(function(e){
          swal({   
              title: "Are you sure?",   
              text: deleteText,   
              type: "warning",   
              showCancelButton: true,   
              confirmButtonColor: "#DD6B55",   
              confirmButtonText: "Yes, delete it!",   
              cancelButtonText: "No, cancel!",   
              closeOnConfirm: false
          }, function(){   
                var http = new XMLHttpRequest();
                var params = 'item-id=' + $(e.currentTarget).attr("product");
                
                http.open('POST', deleteUrl, true);

                //Send the proper header information along with the request
                http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

                http.onreadystatechange = function() { //Call a function when the state changes.
                    if(http.readyState == 4 && http.status == 200) {
                      console.log(params, i)
                      swal({
                        title: "Deleted!", 
                        text: deleteSuccessText, 
                        type: "success" 
                        }, 
                        function(){
                          location.reload()
                        }); 
                    }else if(http.readyState == 4 && http.status != 200){
                      swal("Failed to delete!", http.response, "error");
                    }
                }
                http.send(params);         
          });
      });
    }
  },
  //init
  $.SweetAlert = new SweetAlert, $.SweetAlert.Constructor = SweetAlert
}(window.jQuery),

//initializing 
function($) {
  "use strict";
  $.SweetAlert.init()
}(window.jQuery);