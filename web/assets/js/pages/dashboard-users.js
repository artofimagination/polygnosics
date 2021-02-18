
$(function () {

  'use strict';
	
	initWebRTC("user-data")	
	
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

	var devClientHistory = new ApexCharts(
	  document.querySelector("#dev-client-history"),
	  options
	);
	devClientHistory.render();

	var options = {
		series: [17,23,40,20],
		labels: ['Active Devs', 'Inactive Devs', 'Active Clients', 'Inactive Clients'],
		chart: {
			type: 'donut',
			width: '100%',
			height: 205
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

	var devClientRatio = new ApexCharts(document.querySelector("#dev-client-ratio"), options);
	devClientRatio.render();
	
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
  	
	var avgOnline = {
		chart: {
			type: 'area',
			height: 170,
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

	var avgOnlineChart = new ApexCharts(document.querySelector("#avg-online"), avgOnline);
	avgOnlineChart.render();

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
			name: 'Total',
			data: data
		},
		{
			name: 'Online',
			data: data
		},
		{
			name: 'Deleted',
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
				text: 'Count'
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

	var userCountChart = new ApexCharts(
		document.querySelector("#user-count"),
		options
	);
	userCountChart.render();
	
	function loadMap() {
		$('#world-map-markers').vectorMap({
			map : 'world_mill_en',
			scaleColors : ['#eff0f1', '#eff0f1'],
			normalizeFunction : 'polynomial',
			hoverOpacity : 0.7,
			hoverColor : false,
			regionStyle : {
				initial : {
					fill : '#e0e7fd'
				}
			},
	
			markerStyle: {
				initial: {
					stroke: "transparent"
				},
				hover: {
					stroke: "rgba(112, 112, 112, 0.30)"
				}
			},
			backgroundColor : 'transparent',
	
			markers: [
			{
				latLng: [37.090240, -95.712891],
				name: "USA",
				style: {
				fill: "#4d79f6"
				}
			},
			{
				latLng: [71.706940, -42.604301],
				name: "Greenland",
				style: {
				fill: "#bfd0ff"
				}
			},
			{
				latLng: [-21.943369, 123.102198],
				name: "Australia",
				style: {
				fill: "#3066ff"
				}
			}
			],
			series: {
				regions: [{
					values: {
						"AU": '#bfd0ff',
						"US": '#a2bafd',
						"GL": '#688df7',
					},
					attribute: 'fill'
				}]
			},
		});
	}

	var onlinePeriod = new Chart(document.getElementById("online-period"), {
	  type: 'line',
	  data: {
		labels: ["00h","01h","02h","03h","04h","05h","06h","07h","08h","09h","10h","11h","12h","13h","14h","15h","16h","17h","18h","19h","20h","21h","22h","23h"],
		datasets: [{ 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Jan",
				borderColor: "#be2525",
				fill: false
		  }, { 
				data: [13,44,67,89,97,78,67,45,230,345,1000,1300,1234,789,345,200,100,90,78,67,46,23,23,15],
				label: "Feb",
				borderColor: "#be5825",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Mar",
				borderColor: "#be7e25",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Apr",
				borderColor: "#bea325",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "May",
				borderColor: "#b3be25",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Jun",
				borderColor: "#8ebe25",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Jul",
				borderColor: "#68be25",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Aug",
				borderColor: "#43be25",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Sept",
				borderColor: "#25be2d",
				fill: false
		  }, { 
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Oct",
				borderColor: "#25be53",
				fill: false
		  }, {
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Nov",
				borderColor: "#25be78",
				fill: false
			}, {
				data: [23,34,200,450,350,100,345,111,23,34,567,345,876,234,11,45,67,230,450,345,123,100,23,45],
				label: "Dec",
				borderColor: "#25be9e",
				fill: false
			}
		]},
	  options: {
			title: {
				display: true,
				text: 'Online periods'
			}
	  }
	});

	// Update online period chart
	function addData(chart, data) {
		var month = 0;
		chart.data.datasets.forEach((dataset) => {
			dataset.data = data[month];
			month++;
		});
		chart.update();
	}

	var loadOnce = false
	window.setInterval(function () {
		data = getWebRTCMessage()
		if ( data === null) {
			return
		}

		if (loadOnce == false) {
			loadMap()
			loadOnce = true
		} 

		// Update user count chart
		userCountChart.updateSeries([{
			data: data["users_total"]
		},
		{
			data: data["users_online"]
		},
		{
			data: data["users_deleted"]
		}])

		// Update dev client ratio
		devClientRatio.updateSeries(data["users_dev_client"])

		// Update online peak bar chart
		var values = []
		values.push(data["users_online_peaks"]["max_percent"])
		values.push(data["users_online_peaks"]["min_percent"])
		values.push(data["users_online_peaks"]["avg_percent"])
		var onlinePeakChart = $(".bar").peity("bar")
		onlinePeakChart
		.text(values.join(","))
		.change()
		var maxColor = "text-danger"
		if (data["users_online_peaks"]["max_trend"] === "up") {
			maxColor = "text-success"
		}
		var minColor = "text-danger"
		if (data["users_online_peaks"]["min_trend"] === "up") {
			minColor = "text-success"
		}
		var avgColor = "text-danger"
		if (data["users_online_peaks"]["avg_trend"] === "up") {
			avgColor = "text-success"
		}
		document.getElementById("max-text").innerHTML = "<i class=\"ti-arrow-" + data["users_online_peaks"]["max_trend"] + " " + maxColor + "\"></i> " + data["users_online_peaks"]["max"] + " <br><small class=\"text-fade\">  Max</small>"
		document.getElementById("min-text").innerHTML = "<i class=\"ti-arrow-" + data["users_online_peaks"]["min_trend"] + " " + minColor + "\"></i> " + data["users_online_peaks"]["min"] + " <br><small class=\"text-fade\">  Min</small>"
		document.getElementById("avg-text").innerHTML = "<i class=\"ti-arrow-" + data["users_online_peaks"]["avg_trend"] + " " + avgColor + "\"></i> " + data["users_online_peaks"]["avg"] + " <br><small class=\"text-fade\">  Avg</small>"
		//$("#online-users-chart").load(" #online-users-chart > *");

		// Update online period chart
		addData(onlinePeriod, data["users_online_period"])
		
	}, 2000)
	
}); // End of use strict
