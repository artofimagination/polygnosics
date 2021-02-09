
!function($) {
    "use strict";

    var SweetAlert = function() {};

    //examples 
    SweetAlert.prototype.init = function() {
      //Parameter
      var i;
      for (i = 0; i < deleteLinks.length; i++)
      {
        var stringVal = '#delete-product-' + i;
        var product = deleteLinks[i]
        $(stringVal).click(function(){
            swal({   
                title: "Are you sure?",   
                text: "It will delete all projects started from this product as well",   
                type: "warning",   
                showCancelButton: true,   
                confirmButtonColor: "#DD6B55",   
                confirmButtonText: "Yes, delete it!",   
                cancelButtonText: "No, cancel!",   
                closeOnConfirm: false
            }, function(){   
                  var http = new XMLHttpRequest();
                  var url = '/user-main/my-products/delete';
                  var params = 'product=' + product;
                  http.open('POST', url, true);

                  //Send the proper header information along with the request
                  http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

                  http.onreadystatechange = function() { //Call a function when the state changes.
                      if(http.readyState == 4 && http.status == 200) {
                        swal("Deleted!", "Your product has been deleted", "success"); 
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