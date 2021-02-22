
$(function () {

  'use strict';
	
	initWebRTC("sys-mon-data")	

	var data = []
	data.push([0, 0])

	var options = {
		chart: {
			type: 'area',
			stacked: false,
				foreColor: '#8a99b5',
			height: 250,
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
			name: 'CPU Load',
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
				text: 'Usage'
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

	var cpuStats = new ApexCharts(
		document.querySelector("#cpu-stats"),
		options
	);
	cpuStats.render();

	var options = {
		chart: {
			type: 'area',
			stacked: false,
				foreColor: '#8a99b5',
			height: 250,
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
			name: 'Memory',
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
				text: 'Usage'
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

	var memStats = new ApexCharts(
		document.querySelector("#mem-stats"),
		options
	);
	memStats.render();

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
			name: 'Storage',
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
				text: 'Usage'
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

	var storageStats = new ApexCharts(
		document.querySelector("#storage-stats"),
		options
	);
	storageStats.render();

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
			name: 'Requests',
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
				text: 'Load'
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

	var requestsStats = new ApexCharts(
		document.querySelector("#requests-stats"),
		options
	);
	requestsStats.render();

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
			name: 'Network',
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
				text: 'Load'
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

	var networkStats = new ApexCharts(
		document.querySelector("#network-stats"),
		options
	);
	networkStats.render();

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
			name: 'DB Requests',
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
				text: 'Load'
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

	var dbStats = new ApexCharts(
		document.querySelector("#db-stats"),
		options
	);
	dbStats.render();

	var loadOnce = false
	window.setInterval(function () {
		var data = getWebRTCMessage()
		if (data === null) {
			return
		}

		if (loadOnce == false) {
			loadOnce = true
		}

		cpuStats.updateSeries([{
			data: data["cpu"]
		}])

		memStats.updateSeries([{
			data: data["memory"]
		}])

		storageStats.updateSeries([{
			data: data["storage"]
		}])

		requestsStats.updateSeries([{
			data: data["requests"]
		}])

		networkStats.updateSeries([{
			data: data["network"]
		}])

		dbStats.updateSeries([{
			data: data["db_requests"]
		}])
	}, 2000)
	
}); // End of use strict
