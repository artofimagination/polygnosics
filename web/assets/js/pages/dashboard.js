//[Dashboard Javascript]

//Project:	Venusx Admin - Responsive Admin Template
//Primary use:   Used only for the main dashboard (index.html)


$(function () {

  'use strict';
	
		
	var options = {
	  chart: {
		  height: 383,
		  type: 'bar',
		  foreColor: '#8a99b5',
		  toolbar: {
			  show: false
		  }
	  },
	  plotOptions: {
		  bar: {
			  horizontal: false,
			  endingShape: 'rounded',
			  columnWidth: '35%',
		  },
	  },
	  dataLabels: {
		  enabled: false
	  },
	  stroke: {
		  show: true,
		  width: 2,
		  colors: ['transparent']
	  },
	  colors: ["#2444e8", "#c6cffb"],
	  series: [{
		  name: 'New Visitors',
		  data: [70, 45, 51, 58, 59, 58, 61, 65, 60, 69]
	  }, {
		  name: 'Unique Visitors',
		  data: [55, 71, 80, 100, 89, 98, 110, 95, 116, 90]
	  },],
	  xaxis: {
		  categories: ['Jan','Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
		  axisBorder: {
			show: true,
			color: '#bec7e0',
		  },  
		  axisTicks: {
			show: true,
			color: '#bec7e0',
		  },    
	  },
	  legend: {
          position: 'top',
           horizontalAlign: 'right',
        },
	  yaxis: {
		  title: {
			  text: 'Visitors'
		  }
	  },
	  fill: {
		  opacity: 1

	  },
	  // legend: {
	  //     floating: true
	  // },
	  grid: {
		  row: {
			  colors: ['transparent', 'transparent'], // takes an array which will be repeated on columns
			  opacity: 0.2
		  },
		  borderColor: '#f1f3fa'
	  },
	  tooltip: {
			theme: 'dark',
		  y: {
			  formatter: function (val) {
				  return "" + val + "k"
			  }
		  }
	  }
	}

	var chart = new ApexCharts(
	  document.querySelector("#yearly-comparison"),
	  options
	);

	chart.render();
	
	
	var options = {
        series: [17, 22, 19, 47],
        chart: {
          type: 'donut',
			width: '100%',
      		height: 250
        },
		colors:['#fda44c', '#4cdaa7', '#5193ff', '#1a233a'],
		legend: {
		  show: false,
		},
		dataLabels: {
			enabled: false,
		  },
        responsive: [{
          breakpoint: 480,
          options: {
            chart: {
              width: 200
            },
          }
        }]
      };

      var chart = new ApexCharts(document.querySelector("#earning-chart"), options);
      chart.render();
	
	
	 window.Apex = {
      stroke: {
        width: 3
      },
      markers: {
        size: 0
      },
      tooltip: {
        fixed: {
          enabled: false,
        }
      }
    };
    
    var randomizeArray = function (arg) {
      var array = arg.slice();
      var currentIndex = array.length,
        temporaryValue, randomIndex;

      while (0 !== currentIndex) {

        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;

        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
      }

      return array;
    }

    // data for the sparklines that appear below header area
    var sparklineData = [47, 45, 54, 38, 56, 24, 65, 31, 37, 39, 62, 51, 35, 41, 35, 27, 93, 53, 61, 27, 54, 43, 19, 46];

    	
	var spark4 = {
      chart: {
        type: 'area',
        height: 195,
        sparkline: {
          enabled: true
        },
      },
      stroke: {
        curve: 'smooth'
      },
      fill: {
        opacity: 1,
        type: 'gradient',
		gradient: {
		  gradientToColors: ['#38649f', '#38649f']
		},
      },
      series: [{
        data: randomizeArray(sparklineData)
      }],
	  labels: [...Array(24).keys()].map(n => `2018-09-0${n+1}`),
      yaxis: {
        min: 0
      },
	  xaxis: {
		type: 'datetime',
	  },
      colors: ['#38649f'],
		tooltip: {
			theme: 'dark',
  		},
    };
	
	
    var spark4 = new ApexCharts(document.querySelector("#spark4"), spark4);
    spark4.render();
	
	 var ts2 = 1484418600000;
		var dates = [];
		var spikes = [5, -5, 3, -3, 8, -8]
		for (var i = 0; i < 120; i++) {
		  ts2 = ts2 + 86400000;
		  var innerArr = [ts2, dataSeries[1][i].value];
		  dates.push(innerArr)
		}

		var options = {
		  chart: {
			type: 'area',
			stacked: false,
			  foreColor: '#8a99b5',
			height: 330,
			zoom: {
			  type: 'x',
			  enabled: true
			},
			toolbar: {
			  autoSelected: 'zoom'
			}
		  },
		  dataLabels: {
			enabled: false
		  },
		  series: [{
			name: 'Balance ($)',
			data: dates
		  }],
		  markers: {
			size: 0,
		  },
		  fill: {
			type: 'gradient',
			gradient: {
			  shadeIntensity: 1,
			  inverseColors: false,
			  opacityFrom: 0.5,
			  opacityTo: 0,
			  stops: [0, 90, 100]
			},
		  },
		  yaxis: {
			min: 20000000,
			max: 250000000,
			labels: {
			  formatter: function (val) {
				return (val / 1000000).toFixed(0);
			  },
			},
			title: {
			  text: 'Price'
			},
		  },
		  xaxis: {
			type: 'datetime',
		  },

		  tooltip: {
			theme: 'dark',
			shared: false,
			y: {
			  formatter: function (val) {
				return (val / 1000000).toFixed(0)
			  }
			}
		  }
		}

		var chart = new ApexCharts(
		  document.querySelector("#balancehistory"),
		  options
		);

		chart.render();
	
	
	
	
	// table
	$('#invoice-list').DataTable({
	  'paging'      : true,
	  'lengthChange': false,
	  'searching'   : false,
	  'ordering'    : true,
	  'info'        : true,
	  'autoWidth'   : true,
	});
	
	// bar chart
		$(".bar").peity("bar");	
	
}); // End of use strict
