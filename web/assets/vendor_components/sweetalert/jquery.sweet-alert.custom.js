
!function($) {
    "use strict";

    var SweetAlert = function() {};

    SweetAlert.prototype.init = function() {
      //Parameter
      var i;
      for (i = 0; i < deleteLinks.length; i++)
      {	
        var stringVal = '#delete-item-' + i;
        var item = deleteLinks[i]
        $(stringVal).click(function(){
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
                  var url = deleteUrl;
                  var params = 'item-id=' + item;
                  http.open('POST', url, true);

                  //Send the proper header information along with the request
                  http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

                  http.onreadystatechange = function() { //Call a function when the state changes.
                      if(http.readyState == 4 && http.status == 200) {
                        swal("Deleted!", deleteSuccessText, "success"); 
                        location.reload()
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