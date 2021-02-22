var sendChannel
var dataJson = null

function initWebRTC(channelName) {
  /* eslint-env browser */
  let pc = new RTCPeerConnection({
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  })

  sendChannel = pc.createDataChannel(channelName)
  sendChannel.onmessage = e => {
    dataJson = JSON.parse(e.data)
  }

  pc.onicecandidate = event => {
    if (event.candidate === null) {
      value = btoa(JSON.stringify(pc.localDescription))
      var params = "offer=" + value + "&type=" + channelName;
      var http = new XMLHttpRequest();
      var url = "/user-main/stats/webrtc";
      http.open('POST', url, true);

      //Send the proper header information along with the request
      http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');

      http.onreadystatechange = function() {
        if(http.readyState == 4 && http.status == 200) {
          pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(this["response"]))))
        }
      }
      http.send(params);
    }
  }

  pc.onnegotiationneeded = e =>
    pc.createOffer().then(d => pc.setLocalDescription(d))

  waitForConnection(pc)
}

function waitForConnection(pc) {
  charts = document.getElementById("charts");
  loaderAnim = document.getElementById("webrtc-loader");

  timeout = 0
  var intervalId = setInterval(function() {
    charts.style.display = "none"
    loaderAnim.style.display = "block"
    if (pc.iceConnectionState === "connected" && dataJson != null) {
      clearInterval(intervalId);
      loaderAnim.style.display = "none"
      charts.style.display = "block"
    }

    if (timeout >= 30) {
      var params = "status=timeout";
      var http = new XMLHttpRequest();
      var url = "/user-main/stats/webrtc";
      http.open('POST', url, true);
      //Send the proper header information along with the request
      http.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  
      http.onreadystatechange = function() {
        if(http.readyState == 4 && http.status == 200) {
          pc.close()
        }
      }
      http.send(params);
      clearInterval(intervalId);
      loaderAnim.style.display = "none"
      charts.style.display = "block"
    }
    timeout++
  }, 1000);
}

function getWebRTCMessage() {
  return dataJson
}

function sendWebRTCMessage(message) {
  sendChannel.send(message)
}
