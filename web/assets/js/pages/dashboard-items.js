
$(function () {

  'use strict';
	
	initWebRTC("item-data")	

	var options = {
	  chart: {
		  height: 300,
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
				name: 'Developers',
				data: [70, 45, 51, 58, 59, 58, 61, 65, 60, 69, 50, 34]
	  	}, {
				name: 'Clients',
				data: [55, 71, 80, 100, 89, 98, 110, 95, 116, 90, 45, 89]
	  },],
	  xaxis: {
		  categories: ['Jan','Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
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
			  text: 'Users'
		  }
	  },
	  fill: {
		  opacity: 1
	  },
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

	var itemsPerUserChart = new ApexCharts(
	  document.querySelector("#items-per-user"),
	  options
	);
	itemsPerUserChart.render();

	var options = {
	  chart: {
		  height: 300,
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
				name: 'Developers',
				data: [70, 45, 51, 58, 59, 58, 61, 65, 60, 69, 50, 34]
	  	}, {
				name: 'Clients',
				data: [55, 71, 80, 100, 89, 98, 110, 95, 116, 90, 45, 89]
	  },],
	  xaxis: {
		  categories: ['Jan','Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
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
			  text: 'Users'
		  }
	  },
	  fill: {
		  opacity: 1
	  },
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

	var itemCountChart = new ApexCharts(
	  document.querySelector("#item-counts"),
	  options
	);
	itemCountChart.render();

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
  	
	var projectLength = {
		chart: {
			type: 'area',
			height: 220,
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
			name: "Daily avg",
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

	var projectLengthChart = new ApexCharts(document.querySelector("#project-length"), projectLength);
	projectLengthChart.render();

	var data = []
	data.push([0, 0])

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
			},
			events: {
				beforeZoom: function(chartContext, { xaxis }) {
					diff = xaxis.max - xaxis.min
					sendWebRTCMessage({"zoom": diff});
					return true
				}
			}
		},
		dataLabels: {
			enabled: false
		},
		series: [{
			name: 'Activity',
			data: data
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
			min: -100,
			max: 200,
			labels: {
				formatter: function (val) {
					return val;
				},
			},
			title: {
				text: 'Activity'
			},
		},
		xaxis: {
			type: 'datetime',
		},

		tooltip: {
			theme: 'dark',
			shared: true,
			y: {
				formatter: function (val) {
				return val
				}
			}
		}
	}

	var userActivityChart = new ApexCharts(
		document.querySelector("#user-project-activity"),
		options
	);
	userActivityChart.render();

	var loadOnce = false
	window.setInterval(function () {
		data = getWebRTCMessage()
		if (data === null) {
			return
		}

		if (loadOnce == false) {
			loadOnce = true
		} 

		// Update user count chart
		userActivityChart.updateSeries([{
			data: data["users_project_activity"]
		}])	
	}, 2000)
	
}); // End of use strict
