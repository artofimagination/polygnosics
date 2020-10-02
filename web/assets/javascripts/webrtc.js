/* eslint-env browser */
let pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
})

let sendChannel = pc.createDataChannel('foo')
sendChannel.onmessage = e => {
  dataJson = JSON.parse(e.data)
  entityData[dataJson.id] = e.data
}

pc.onicecandidate = event => {
  if (event.candidate === null) {
    value = btoa(JSON.stringify(pc.localDescription))
    var params = "offer=" + value;
    var http = new XMLHttpRequest();
    var url = "/user-main/" + project_id + "/webrtc";
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

window.sendMessage = () => {
  var message = ""
  sendChannel.send(message)
}