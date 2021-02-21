
$(function () {

  'use strict';
	
	initWebRTC("product-data")	

	function initScoreChart() {
		var plot1 = $.plot('#scoreChart', [{
				data: flotSampleData5,
				color: '#6610f2'
			},{
				data: flotSampleData3,
				color: '#00cccc'
			}], {
			series: {
				shadowSize: 0,
				lines: {
					show: true,
					lineWidth: 2,
					fill: true,
					fillColor: { colors: [ { opacity: 0 }, { opacity: 0.0 } ] }
				}
			},
			grid: {
				borderWidth: 0,
				borderColor: '#464d5c',
				labelMargin: 5,
				markings: [{
					xaxis: { from: 10, to: 20 },
					color: '#1a233a'
				}]
			},
			yaxis: {
				show: false,
				color: '#ced4da',
				tickLength: 10,
				min: 0,
				max: 110,
				font: {
					size: 11,
					color: '#969dab'
				},
				tickFormatter: function formatter(val, axis) {
					return val + 'k';
				}
			},
			xaxis: {
				show: false,
				position: 'top',
				color: 'rgba(255,255,255,0.1)'
			}
		});

		var mqSM = window.matchMedia('(min-width: 576px)');
		var mqSMMD = window.matchMedia('(min-width: 576px) and (max-width: 991px)');
		var mqLG = window.matchMedia('(min-width: 992px)');

		function screenCheck() {
			if (mqSM.matches) {
				plot1.getAxes().yaxis.options.show = true;
				plot1.getAxes().xaxis.options.show = true;
			} else {
				plot1.getAxes().yaxis.options.show = false;
				plot1.getAxes().xaxis.options.show = false;
			}

			if (mqSMMD.matches) {
				var tick = [
					[0, '<span>Oct</span><span>10</span>'],
					[20, '<span>Oct</span><span>12</span>'],
					[40, '<span>Oct</span><span>14</span>'],
					[60, '<span>Oct</span><span>16</span>'],
					[80, '<span>Oct</span><span>18</span>'],
					[100, '<span>Oct</span><span>19</span>'],
					[120, '<span>Oct</span><span>20</span>'],
					[140, '<span>Oct</span><span>23</span>']
				];

				plot1.getAxes().xaxis.options.ticks = tick;
			}

			if (mqLG.matches) {
				var tick = [
					[10, '<span>Oct</span><span>10</span>'],
					[20, '<span>Oct</span><span>11</span>'],
					[30, '<span>Oct</span><span>12</span>'],
					[40, '<span>Oct</span><span>13</span>'],
					[50, '<span>Oct</span><span>14</span>'],
					[60, '<span>Oct</span><span>15</span>'],
					[70, '<span>Oct</span><span>16</span>'],
					[80, '<span>Oct</span><span>17</span>'],
					[90, '<span>Oct</span><span>18</span>'],
					[100, '<span>Oct</span><span>19</span>'],
					[110, '<span>Oct</span><span>20</span>'],
					[120, '<span>Oct</span><span>21</span>'],
					[130, '<span>Oct</span><span>22</span>'],
					[140, '<span>Oct</span><span>23</span>']
				];

				plot1.getAxes().xaxis.options.ticks = tick;
			}
		}

		screenCheck();
		mqSM.addEventListener("change", screenCheck);
		mqSMMD.addEventListener("change" ,screenCheck);
		mqLG.addEventListener("change", screenCheck);

		plot1.setupGrid();
		plot1.draw();
	}

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
	var sparklineData = [47, 10, 54, 38, 56, 24, 65, 31, 37, 39, 62, 51, 35, 41, 35, 27, 93, 53, 61, 27, 54, 43, 19, 46];
  	
	var avgProjectGen = {
		chart: {
			type: 'area',
			height: 290,
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

	var avgProjectGenChart = new ApexCharts(document.querySelector("#avg-project-gen"), avgProjectGen);
	avgProjectGenChart.render();

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
			name: 'Views',
			data: data
		},
		{
			name: 'Pinnes',
			data: data
		},
		{
			name: 'Purchases',
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

	var details = new ApexCharts(
		document.querySelector("#detail-stats"),
		options
	);
	details.render();

	var loadOnce = false
	window.setInterval(function () {
		var data = getWebRTCMessage()
		if (data === null) {
			return
		}

		if (loadOnce == false) {
			initScoreChart()
			loadOnce = true
		}

		details.updateSeries([{
			data: data["product_viewers"]
		},
		{
			data: data["product_purchase"]
		},
		{
			data: data["product_pins"]
		}])	 
	}, 2000)
	
}); // End of use strict
