
!function($) {
  "use strict";
  $("#show-loader").get(0).style.display = "block"
  $("#show-frame").get(0).style.display = "none"
  
	var timer = window.setInterval(function () {
		$.ajax({
      type: 'GET',
      url: "/check-state?item-id=" + project,
      statusCode: {
        204: function() {
        },
        200: function() {
          $("#show-loader").get(0).style.display = "none"
          $("#show-frame").get(0).style.display = "block"
          clearInterval(timer);
        }
      }
    })
	}, 2000)
  
}(window.jQuery),

//initializing 
function($) {
  "use strict";
}(window.jQuery);