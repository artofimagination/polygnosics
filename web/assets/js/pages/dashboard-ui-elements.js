
$(function () {

  'use strict';
	
	initWebRTC("ui-data")	

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
			name: 'Popularity',
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
				text: 'Rate'
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

	var popularity = new ApexCharts(
		document.querySelector("#popularity-stats"),
		options
	);
	popularity.render();

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
			name: 'Clicks',
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

	var clicks = new ApexCharts(
		document.querySelector("#click-stats"),
		options
	);
	clicks.render();

	var loadOnce = false
	window.setInterval(function () {
		var data = getWebRTCMessage()
		if (data === null) {
			return
		}

		if (loadOnce == false) {
			loadOnce = true
		}

		clicks.updateSeries([{
			data: data["clicks"]
		}])

		popularity.updateSeries([{
			data: data["popularity"]
		}]) 
	}, 2000)
	
}); // End of use strict
