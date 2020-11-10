//[Dashboard Javascript]

//Project:	Venusx Admin - Responsive Admin Template
//Primary use:   Used only for the main dashboard (index.html)


$(function () {

  'use strict';
	
	var plot1 = $.plot('#flotChart', [{
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
        mqSM.addListener(screenCheck);
        mqSMMD.addListener(screenCheck);
        mqLG.addListener(screenCheck);

        plot1.setupGrid();
        plot1.draw();
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118]
        }],
        chart: {
          height: 150,
          type: 'bar',
		  foreColor: '#8a99b5',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'top', // top, center, bottom
            },
			  columnWidth: '15%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: false,
        },
        colors:['#45b6c6'],
        xaxis: {
          categories: ["Jan", "Feb", "Mar", "Apr"],
          position: 'bottom',
			
          labels: {
			show: true, 
          },
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false
          },
          tooltip: {
            enabled: false,        
          }
        },
		grid: {
		  yaxis: {
			lines: {
			  show: false
			}
		  }
		},
        yaxis: {
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false,
          },
          labels: {
            show: false,
            formatter: function (val) {
              return val + "%";
            }
          }        
        },	
      };

      var chart = new ApexCharts(document.querySelector("#meetingschart3"), options);
      chart.render();
	
	
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118]
        }],
        chart: {
          height: 120,
          type: 'bar',
		  foreColor: '#8a99b5',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'top', // top, center, bottom
            },
			  columnWidth: '15%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: false,
        },
        colors:['#2444e8'],
        xaxis: {
          categories: ["Jan", "Feb", "Mar", "Apr"],
          position: 'bottom',
			
          labels: {
			show: true, 
          },
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false
          },
          tooltip: {
            enabled: false,        
          }
        },
		grid: {
		  yaxis: {
			lines: {
			  show: false
			}
		  }
		},
        yaxis: {
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false,
          },
          labels: {
            show: false,
            formatter: function (val) {
              return val + "%";
            }
          }
        
        },
      };

      var chart = new ApexCharts(document.querySelector("#meetingschart4"), options);
      chart.render();
	
	
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
	
		
}); // End of use strict
