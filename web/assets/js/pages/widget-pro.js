//[Dashboard Javascript]

//Project:	Venusx Admin - Responsive Admin Template
//Primary use:   Used only for the main dashboard (index.html)


$(function () {

  'use strict';
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118, 137, 103, 168, 223]
        }],
        chart: {
          height: 325,
          type: 'bar',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'top', // top, center, bottom
            },
			  columnWidth: '40%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: true,
          formatter: function (val) {
            return val ;
          },
          offsetY: -20,
          style: {
            fontSize: '12px',
            colors: ["#2444e8"]
          }
        },
        colors:['#2444e8'],
        xaxis: {
          categories: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug"],
          position: 'top',
			
          labels: {
			show: false,
            offsetY: -18,        
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
        fill: {
          gradient: {
            shade: 'light',
            type: "horizontal",
            shadeIntensity: 0.25,
            gradientToColors: undefined,
            inverseColors: true,
            opacityFrom: 1,
            opacityTo: 1,
            stops: [50, 0, 100, 100]
          },
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

      var chart = new ApexCharts(document.querySelector("#orderchart"), options);
      chart.render();
	
	
  
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118, 137, 103, 168, 223]
        }],
        chart: {
          height: 230,
          type: 'bar',
			foreColor: '#ccc',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'bottom', // top, center, bottom
            },
			  columnWidth: '40%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: false,
        },
        colors:['#ffffff'],
        xaxis: {
          categories: ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug"],
          position: 'bottom',
			
          labels: {
			show: true,
            offsetY: 0, 
			  colors: ["#ffffff"]
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

      var chart = new ApexCharts(document.querySelector("#ratingchart"), options);
      chart.render();
	
	
	
	
	
	var options = {
        series: [{
            name: "Desktops",
            data: [189, 156, 155, 118, 167, 159, 178, 223, 195, 201, 143]
        }],
        chart: {
          height: 325,
          type: 'line',
          zoom: {
            enabled: false
          }
        },
		colors:['#cbccce'],
        dataLabels: {
          enabled: true,
			style: {
				fontSize: '16px',
				colors: ['#2444e8'],
			},
        },
        stroke: {
          	show: true,
			curve: 'smooth',
			lineCap: 'butt',
			colors: undefined,
			width: 4,
			dashArray: 0, 
        },		
		markers: {
			size: 6,
			colors: '#2444e8',
			strokeColors: '#ffffff',
			strokeWidth: 2,
			strokeOpacity: 0.9,
			strokeDashArray: 0,
			fillOpacity: 1,
			discrete: [],
			shape: "circle",
			radius: 5,
			offsetX: 0,
			offsetY: 0,
			onClick: undefined,
			onDblClick: undefined,
			hover: {
			  size: undefined,
			  sizeOffset: 3
			}
		},
        grid: {
          row: {
            colors: ['transparent'], // takes an array which will be repeated on columns
            opacity: 0
          },			
		  yaxis: {
			lines: {
			  show: false,
			},
		  },
        },
        xaxis: {
          categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep'],
		  labels: {
			show: false,
            offsetY: -18,        
          },
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false
          },
          tooltip: {
            enabled: false,        
          },
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
      var chart = new ApexCharts(document.querySelector("#visits-chart"), options);
      chart.render();
	
	
	
	
	
	
	
	var options = {
        series: [{
          name: 'series1',
          data: [189, 156, 155, 118, 167, 159, 178, 223, 195, 201, 143]
        }],
        chart: {
          height: 400,
			width: 438,
          type: 'line',
			offsetY: 0,
			offsetX: -38,
			dropShadow: {
			  enabled: true,
			  top: 20,
			  left: 0,
			  blur: 10,
			  opacity: 0.5
			}
        },
		colors:['#ffffff'],
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'smooth'
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
          }
        
        },
		
        xaxis: {
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
		grid: {
			show: true,
			borderColor: '#5578ed',
			strokeDashArray: 0,
			position: 'back',
			xaxis: {
				lines: {
					show: false
				}
			},   
			yaxis: {
				lines: {
					show: false
				}
			},  
			row: {
				colors: undefined,
				opacity: 0.5
			},  
			column: {
				colors: undefined,
				opacity: 0.1
			},  
		}
      };

      var chart = new ApexCharts(document.querySelector("#statisticschart"), options);
      chart.render();
	
	
	var options = {
        series: [{
          name: 'series1',
          data: [178, 223, 195, 201, 143, 189, 156, 155, 118, 167, 159]
        }],
        chart: {
          height: 400,
			width: 438,
          type: 'line',
			offsetY: 0,
			offsetX: -38,
			dropShadow: {
			  enabled: true,
			  top: 20,
			  left: 0,
			  blur: 10,
			  opacity: 0.5
			}
        },
		colors:['#ffffff'],
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'smooth'
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
          }
        
        },
        xaxis: {
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false,
          },
          labels: {
            show: false,
            formatter: function (val) {
              return val ;
            }
          }
        
        },
		grid: {
			show: true,
			borderColor: '#5578ed',
			strokeDashArray: 0,
			position: 'back',
			xaxis: {
				lines: {
					show: false,
				}
			},   
			yaxis: {
				lines: {
					show: false
				}
			},  
			row: {
				colors: undefined,
				opacity: 0.5
			},  
			column: {
				colors: undefined,
				opacity: 0.1
			},  
		}
      };

      var chart = new ApexCharts(document.querySelector("#statisticschart2"), options);
      chart.render();
	
	
	
	
	
	
	var options = {
        series: [{
          name: 'series1',
          data: [178, 223, 195, 201, 143, 189, 156, 155, 118, 167, 159]
        }],
        chart: {
          height: 200,
			width: 600,
          type: 'line',
			offsetY: 0,
			offsetX: -50,
        },
		colors:['#ffffff'],
        dataLabels: {
          enabled: false
        },
        stroke: {
          curve: 'smooth',
        },
			
		markers: {
			size: 0,
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
          }
        
        },
        xaxis: {
          axisBorder: {
            show: false
          },
          axisTicks: {
            show: false,
          },
          labels: {
            show: false,
            formatter: function (val) {
              return val ;
            }
          }
        
        },
		grid: {
			show: true,
			borderColor: '#5578ed',
			strokeDashArray: 0,
			position: 'back',
			xaxis: {
				lines: {
					show: false,
				}
			},   
			yaxis: {
				lines: {
					show: false
				}
			},  
			row: {
				colors: undefined,
				opacity: 0.5,
			},  
			column: {
				colors: undefined,
				opacity: 0.1
			},  
		}
      };

      var chart = new ApexCharts(document.querySelector("#statisticschart3"), options);
      chart.render();
	
	
	
	
	
	
	var myConfig = {
        "type": "line",
		"utc": true,
        "plot": {
          "animation": {
            "delay": 500,
            "effect": "ANIMATION_SLIDE_LEFT"
          }
        },
        "plotarea": {
          "margin": "10px 25px 70px 46px"
        },
        "scale-y": {
          "values": "0:100:25",
          "line-color": "none",
          "guide": {
            "line-style": "solid",
            "line-color": "#d2dae2",
            "line-width": "1px",
            "alpha": 0.5
          },
          "tick": {
            "visible": false
          },
          "item": {
            "font-color": "#8391a5",
            "font-family": "Arial",
            "font-size": "10px",
            "padding-right": "5px"
          }
        },
        "scale-x": {
          "line-color": "#d2dae2",
          "line-width": "2px",
          "values": ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"],
          "tick": {
            "line-color": "#d2dae2",
            "line-width": "1px"
          },
          "guide": {
            "visible": false
          },
          "item": {
            "font-color": "#8391a5",
            "font-family": "Arial",
            "font-size": "10px",
            "padding-top": "5px"
          }
        },
        "legend": {
          "layout": "x4",
          "background-color": "none",
          "shadow": 0,
          "margin": "auto auto 15 auto",
          "border-width": 0,
          "item": {
            "font-color": "#707d94",
            "font-family": "Arial",
            "padding": "0px",
            "margin": "0px",
            "font-size": "9px"
          },
          "marker": {
            "show-line": "true",
            "type": "match",
            "font-family": "Arial",
            "font-size": "10px",
            "size": 4,
            "line-width": 2,
            "padding": "3px"
          }
        },
        "crosshair-x": {
          "lineWidth": 1,
          "line-color": "#707d94",
          "plotLabel": {
            "shadow": false,
            "font-color": "#000",
            "font-family": "Arial",
            "font-size": "10px",
            "padding": "5px 10px",
            "border-radius": "5px",
            "alpha": 1
          },
          "scale-label": {
            "font-color": "#ffffff",
            "background-color": "#707d94",
            "font-family": "Arial",
            "font-size": "10px",
            "padding": "5px 10px",
            "border-radius": "5px"
          }
        },
        "tooltip": {
          "visible": false
        },
        "series": [{
          "values": [69, 68, 54, 48, 70, 74, 98, 70, 72, 68, 49, 69],
          "text": "Kenmore",
          "line-color": "#389f99",
          "line-width": "2px",
          "shadow": 0,
          "marker": {
            "background-color": "#fff",
            "size": 3,
            "border-width": 1,
            "border-color": "#389f99",
            "shadow": 0
          },
          "palette": 0
        }, {
          "values": [51, 53, 47, 60, 48, 52, 75, 52, 55, 47, 60, 48],
          "text": "Craftsman",
          "line-width": "2px",
          "line-color": "#38649f",
          "shadow": 0,
          "marker": {
            "background-color": "#fff",
            "size": 3,
            "border-width": 1,
            "border-color": "#38649f",
            "shadow": 0
          },
          "palette": 1,
          "visible": 1
        }, {
          "values": [42, 43, 30, 50, 31, 48, 55, 46, 48, 32, 50, 38],
          "text": "DieHard",
          "line-color": "#ee1044",
          "line-width": "2px",
          "shadow": 0,
          "marker": {
            "background-color": "#fff",
            "size": 3,
            "border-width": 1,
            "border-color": "#ee1044",
            "shadow": 0
          },
          "palette": 2,
          "visible": 1
        }, {
          "values": [25, 15, 26, 21, 24, 26, 33, 25, 15, 25, 22, 24],
          "text": "Land's End",
          "line-color": "#ff8f00",
          "line-width": "2px",
          "shadow": 0,
          "marker": {
            "background-color": "#fff",
            "size": 3,
            "border-width": 1,
            "border-color": "#ff8f00",
            "shadow": 0
          },
          "palette": 3
        }]
      };

    zingchart.render({
      id: 'myChart',
      data: myConfig,
      height: 235,
      width: '100%'
    });
	

	
	var bar = new ProgressBar.Circle(progressbar1, {
  color: '#7f4cc1',
  // This has to be the same size as the maximum width to
  // prevent clipping
  strokeWidth: 20,
  trailWidth: 1,
  easing: 'easeInOut',
  duration: 1400,
  text: {
    autoStyleContainer: false
  },
  from: { color: '#7f4cc1', width: 1 },
  to: { color: '#9c52d5', width: 4 },
  // Set default step function for all animate calls
  step: function(state, circle) {
    circle.path.setAttribute('stroke', state.color);
    circle.path.setAttribute('stroke-width', state.width);

    var value = Math.round(circle.value() * 150);
    if (value === 0) {
      circle.setText('');
    } else {
      circle.setText(value);
    }

  }
});
bar.text.style.fontFamily = '"Raleway", Helvetica, sans-serif';
bar.text.style.fontSize = '2rem';

bar.animate(0.78);
	

	
	
	
	var bar = new ProgressBar.Circle(progressbar2, {
  color: '#ffffff',
  // This has to be the same size as the maximum width to
  // prevent clipping
  strokeWidth: 20,
  trailWidth: 10,
  trailColor: '#3c55b9',
  easing: 'easeInOut',
  duration: 1400,
  text: {
    autoStyleContainer: false
  },
  from: { color: '#ffffff', width: 10 },
  to: { color: '#ffffff', width: 10 },
  // Set default step function for all animate calls
  step: function(state, circle) {
    circle.path.setAttribute('stroke', state.color);
    circle.path.setAttribute('stroke-width', state.width);

    var value = Math.round(circle.value() * 150);
    if (value === 0) {
      circle.setText('');
    } else {
      circle.setText(value);
    }

  }
});
bar.text.style.fontFamily = '"Raleway", Helvetica, sans-serif';
bar.text.style.fontSize = '2rem';

bar.animate(0.78);
	
	
	
	
	
	
	
	
	
	// Apex  start
  if($('#apexChart2').length) {
    var options2 = {
      chart: {
        type: "bar",
        height: 100,
        sparkline: {
          enabled: !0
        }
      },
      plotOptions: {
        bar: {
          columnWidth: "25%"
        }
      },
      colors: ["#ffffff"],
      series: [{
        data: [36, 77, 52, 90, 74, 35, 55, 23, 47, 10, 63, 36, 77, 52, 90, 74, 35, 55, 23, 47]
      }],
      labels: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
      xaxis: {
        crosshairs: {
          width: 2
        }
      },
      tooltip: {
        fixed: {
          enabled: !1
        },
        x: {
          show: !1
        },
        y: {
          title: {
            formatter: function(e) {
              return ""
            }
          }
        },
        marker: {
          show: !1
        }
      }
    };
    new ApexCharts(document.querySelector("#apexChart2"),options2).render();
  }
  // Apex  end
	
	
	
	// Apex  start
  if($('#apexChart3').length) {
    var options2 = {
      chart: {
        type: "bar",
        height: 100,
        sparkline: {
          enabled: !0
        }
      },
      plotOptions: {
        bar: {
          columnWidth: "25%"
        }
      },
      colors: ["#ffffff"],
      series: [{
        data: [36, 77, 52, 90, 74, 35, 55, 23, 47, 10, 63, 36, 77, 52, 90, 74, 35, 55, 23, 47]
      }],
      labels: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20],
      xaxis: {
        crosshairs: {
          width: 2
        }
      },
      tooltip: {
        fixed: {
          enabled: !1
        },
        x: {
          show: !1
        },
        y: {
          title: {
            formatter: function(e) {
              return ""
            }
          }
        },
        marker: {
          show: !1
        }
      }
    };
    new ApexCharts(document.querySelector("#apexChart3"),options2).render();
  }
  // Apex chart2 end
	
	
	var options = {
        series: [{
          name: 'Net Profit',
          data: [44, 55, 57, 56, 61, 58, 63, 60, 66, 55, 57, 56]
        }],
        chart: {
          type: 'bar',
          height: 198
        },
		colors:['#2444e8'],
        plotOptions: {
          bar: {
            horizontal: false,
            columnWidth: '40%',
            endingShape: 'rounded'
          },
        },
        dataLabels: {
          enabled: false
        },
		grid: {
			show: false,  
		},
        stroke: {
          show: false,
          width: 0,
          colors: ['transparent']
        },
        xaxis: {
          categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan'],
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
          }
        
        },
        fill: {
          opacity: 1
        },
        tooltip: {
          y: {
            formatter: function (val) {
              return "$ " + val + " thousands"
            }
          }
        }
      };

      var chart = new ApexCharts(document.querySelector("#profit"), options);
      chart.render();
	
	
	
	
	
	
	
		var options = {
        series: [{
          name: 'This Year',
          data: [44, 55, 41, 67, 22, 43, 41, 12, 56, 51, 42, 44]
        }, {
          name: 'Past Year',
          data: [13, 23, 20, 8, 13, 27, 22, 17, 28, 14, 9, 12]
        }],
        chart: {
          type: 'bar',
          height: 313,
          stacked: true,
          toolbar: {
            show: false
          },
          zoom: {
            enabled: false
          }
        },
		dataLabels: {
          enabled: false
        },
		colors:['#ef0753', '#2444e8'],
        responsive: [{
          breakpoint: 480,
        }],
        plotOptions: {
          bar: {
            horizontal: false,
			  columnWidth: '20%',
			  endingShape: 'rounded',
          },
        },
        xaxis: {
          categories: ['Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan'],
        },
        legend: {
          position: 'top',
           horizontalAlign: 'right',
        },
        fill: {
          opacity: 1
        }
      };

      var chart = new ApexCharts(document.querySelector("#yearly-comparison"), options);
      chart.render();
	
	
	
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118]
        }],
        chart: {
          height: 195,
          type: 'bar',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'top', // top, center, bottom
            },
			  columnWidth: '30%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: false,
        },
        colors:['#ffffff'],
        xaxis: {
          categories: ["Jan", "Feb", "Mar", "Apr"],
          position: 'bottom',
			
          labels: {
			show: true, 
			  style: {
				colors: ['#ffffff', '#ffffff', '#ffffff', '#ffffff']
			  },
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

      var chart = new ApexCharts(document.querySelector("#meetingschart"), options);
      chart.render();
	
	
	
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118]
        }],
        chart: {
          height: 195,
          type: 'bar',
        },
        plotOptions: {
          bar: {
            dataLabels: {
              position: 'top', // top, center, bottom
            },
			  columnWidth: '30%',
			  endingShape: 'rounded',
          }
        },
        dataLabels: {
          enabled: false,
        },
        colors:['#6b6b78'],
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

      var chart = new ApexCharts(document.querySelector("#meetingschart2"), options);
      chart.render();
	
	
	
	var options = {
        series: [{
          name: 'Inflation',
          data: [189, 156, 123, 118]
        }],
        chart: {
          height: 150,
          type: 'bar',
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
	
	
	
	
	
	
	var options = {
        series: [17, 22, 19, 47],
        chart: {
          type: 'donut',
			width: '100%',
      		height: 250
        },
		colors:['#fda44c', '#4cdaa7', '#5193ff', '#eaeaea'],
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
	
	
	
	
	  var options1 = {
        series: [{
          data: [25, 66, 41, 89, 63, 25, 44, 12, 36, 9, 54]
        }],
        chart: {
          type: 'line',
          width: 150,
          height: 50,
          sparkline: {
            enabled: true
          }
        },
		 stroke: {
          curve: 'smooth',
			 width: 3,
        },
		 
		markers: {
			size: 0,
		},
        tooltip: {
          fixed: {
            enabled: false
          },
          x: {
            show: false
          },
          y: {
            title: {
              formatter: function (seriesName) {
                return ''
              }
            }
          },
          marker: {
            show: false
          }
        }
      };

      var chart1 = new ApexCharts(document.querySelector("#visitors-char"), options1);
      chart1.render();
	
	
	
	 // ------------------------------
    // Basic line chart
    // ------------------------------
    // based on prepared DOM, initialize echarts instance
        var myChart = echarts.init(document.getElementById('basic-line'));

        // specify chart configuration item and data
        var option = {
                // Setup grid
                grid: {
                     left: '1%',
                    right: '2%',
                    bottom: '3%',
                    containLabel: true
                },

                // Add Tooltip
                tooltip : {
                    trigger: 'axis'
                },

                // Add Legend
                legend: {
                    data:['Online']
                },

                // Add custom colors
                color: ['#4974e0'],

                // Enable drag recalculate
                calculable : true,

                // Horizontal Axiz
                xAxis : [
                    {
                        type : 'category',
                        boundaryGap : false,
                        data : ['Mon','Tue','Wed','Thu','Fri','Sat','Sun']
                    }
                ],

                // Vertical Axis
                yAxis : [
                    {
                        type : 'value',
                    }
                ],

                // Add Series
                series : [
                    {
                        name:'Online',
                        type:'line',
                        data:[10, 8, 14, 19, 17, 12, 14],
                        markPoint : {
                            data : [
                                {type : 'max', name: 'Max'},
                                {type : 'min', name: 'Min'}
                            ]
                        },
                        markLine : {
                            data : [
                                {type : 'average', name: 'Average'}
                            ]
                        },
                        lineStyle: {
                            normal: {
                                width: 3,
                                shadowColor: 'rgba(0,0,0,0.1)',
                                shadowBlur: 10,
                                shadowOffsetY: 10
                            }
                        },
                    },
                ]
            };
        // use configuration item and data specified to show chart
        myChart.setOption(option);
	
	
	
	
	
		// area chart
		 Morris.Area({
				element: 'area-chart3',
				data: [{
							period: '2013',
							data1: 0,
							data2: 0
						}, {
							period: '2014',
							data1: 55,
							data2: 20
						}, {
							period: '2015',
							data1: 25,
							data2: 55
						}, {
							period: '2016',
							data1: 65,
							data2: 17
						}, {
							period: '2017',
							data1: 35,
							data2: 25
						}, {
							period: '2018',
							data1: 30,
							data2: 85
						}, {
							period: '2019',
							data1: 15,
							data2: 15
						}


						],
						lineColors: ['#e00051', '#f9b423'],
						xkey: 'period',
						ykeys: ['data1', 'data2'],
						labels: ['Data 1', 'Data 2'],
						pointSize: 0,
			 			padding: 1,
						lineWidth: 0,
						resize:true,
						fillOpacity: 1,
						behaveLikeLine: true,
						gridLineColor: '#ffffff0',
						hideHover: 'auto',
			 			axes: false,

			});
	
	
	
	
	
	// bar chart
		$(".bar").peity("bar");	
	
}); // End of use strict


$(function () {
  'use strict';
	
	
	
	
	window.Apex = {
		  stroke: {
			width: 3
		  },
		  markers: {
			size: 0
		  },
		  tooltip: {
			fixed: {
			  enabled: true,
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
		var sparklineData = [178, 223, 195, 201, 143, 189, 156, 155, 118, 167, 159, 178, 223, 195, 201, 143, 189, 156, 155, 118, 167, 159, 178, 223 ];

		var spark2 = {
		  chart: {
			type: 'area',
			height: 415,
			sparkline: {
			  enabled: true
			},
			  
			dropShadow: {
			  enabled: true,
			  top: 20,
			  left: 0,
			  blur: 10,
			  opacity: 0.5
			}
		  },
		  stroke: {		  
			show: true,
			width: 5,
			curve: 'smooth'
		  },
		  fill: {
			opacity: 0.1,
			gradient: {
			  enabled: false
			}
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
		  colors: ['#ffffff'],
		  subtitle: {
			offsetX: 0,
			style: {
			  fontSize: '14px',
			  cssClass: 'apexcharts-yaxis-title'
			}
		  }
		}


		var spark2 = new ApexCharts(document.querySelector("#spark2"), spark2);
		spark2.render();
	

}); // End of use strict
