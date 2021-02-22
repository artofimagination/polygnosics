
$(function () {

  'use strict';
	
	initWebRTC("accounting-data")	

	function initIncomeFiatChart() {
		var fiatPlot = $.plot('#income-fiat-chart', [{
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
				fiatPlot.getAxes().yaxis.options.show = true;
				fiatPlot.getAxes().xaxis.options.show = true;
			} else {
				fiatPlot.getAxes().yaxis.options.show = false;
				fiatPlot.getAxes().xaxis.options.show = false;
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

				fiatPlot.getAxes().xaxis.options.ticks = tick;
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

				fiatPlot.getAxes().xaxis.options.ticks = tick;
			}
		}

		screenCheck();
		mqSM.addEventListener("change", screenCheck);
		mqSMMD.addEventListener("change" ,screenCheck);
		mqLG.addEventListener("change", screenCheck);

		fiatPlot.setupGrid();
		fiatPlot.draw();
	}

	function initIncomeCryptoChart() {
		var cryptoPlot = $.plot('#income-crypto-chart', [{
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
				cryptoPlot.getAxes().yaxis.options.show = true;
				cryptoPlot.getAxes().xaxis.options.show = true;
			} else {
				cryptoPlot.getAxes().yaxis.options.show = false;
				cryptoPlot.getAxes().xaxis.options.show = false;
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

				cryptoPlot.getAxes().xaxis.options.ticks = tick;
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

				cryptoPlot.getAxes().xaxis.options.ticks = tick;
			}
		}

		screenCheck();
		mqSM.addEventListener("change", screenCheck);
		mqSMMD.addEventListener("change" ,screenCheck);
		mqLG.addEventListener("change", screenCheck);

		cryptoPlot.setupGrid();
		cryptoPlot.draw();
	}

	var options = {
		series: [40,40,20],
		labels: ['Subscriptions', 'Purchases', "Donations"],
		chart: {
			type: 'donut',
			width: '100%',
			height: 300
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

	var currencyRatio = new ApexCharts(document.querySelector("#income-ratio"), options);
	currencyRatio.render();

	var options = {
		series: [17,83],
		labels: ['Crypto', 'Fiat'],
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

	var incomeRatio = new ApexCharts(document.querySelector("#currency-ratio"), options);
	incomeRatio.render();

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
			name: 'Expenses',
			data: data
		},
		{
			name: 'Revenue',
			data: data
		},
		{
			name: 'Profit',
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
				text: 'Amount'
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
			name: 'Expense',
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
				text: 'Amount'
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

	var expenses = new ApexCharts(
		document.querySelector("#expenses"),
		options
	);
	expenses.render();

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
			name: 'Income',
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
				text: 'Amount'
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

	var forecast = new ApexCharts(
		document.querySelector("#forecast"),
		options
	);
	forecast.render();

	var loadOnce = false
	window.setInterval(function () {
		var data = getWebRTCMessage()
		if (data === null) {
			return
		}

		if (loadOnce == false) {
			initIncomeFiatChart()
			initIncomeCryptoChart()
			loadOnce = true
		}

		details.updateSeries([{
			data: data["expenses"]
		},
		{
			data: data["revenue"]
		},
		{
			data: data["profit"]
		}])

		expenses.updateSeries([{
			data: data["expense_item"]
		}])
		
		forecast.updateSeries([{
			data: data["forecast"]
		}])	
	}, 2000)
	
}); // End of use strict
