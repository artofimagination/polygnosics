package webrtc

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pion/webrtc/v3"
)

// OfferAddr represents the address of the platform, that initiates
// webRTC with an offer.
var OfferAddr *string = nil

// AnswerAddr represents the product project address that is inactive
// until an offer arrives from the platform side.
var AnswerAddr *string = nil

// SendingFrequency determines how often a message is sent via webRTC.
var SendingFrequency time.Duration = 500 * time.Millisecond

// TurnServerAddress holds the default turn server address. Must change to a valid address when deployed.
var TurnServerAddress string = "turns:172.18.0.3:5349?transport=udp"

// StunServerAddress holds the default stun server address. Must change to a valid address when deployed.
var StunServerAddress string = "stun:172.18.0.3:5349?transport=udp"

// TurnAuthCredential holds the default long term secret. it is used to access the turn server.
// Must change to a valid address when deployed.
var TurnAuthCredential string = "default123secret"

// TurnAuthUserName holds the default long term credential user name. It is used to access the turn server.
// Must change to a valid address when deployed.
var TurnAuthUserName string = "defaultUser"

// Certificates are used to create secure webrtc clients. Certificates are generated outside the app by generateCerts.sh script
var Certificates tls.Certificate = tls.Certificate{}

var config webrtc.Configuration = webrtc.Configuration{

	ICEServers: []webrtc.ICEServer{
		{
			URLs:           []string{TurnServerAddress, StunServerAddress},
			Username:       TurnAuthUserName,
			Credential:     TurnAuthCredential,
			CredentialType: webrtc.ICECredentialTypePassword,
		},
	},

	ICETransportPolicy: webrtc.ICETransportPolicyAll,
}

func signalCandidate(client *http.Client, addr string, c *webrtc.ICECandidate) error {
	payload := []byte(c.ToJSON().Candidate)
	resp, err := client.Post(fmt.Sprintf("https://%s/candidate", addr), "application/json; charset=utf-8", bytes.NewReader(payload)) //nolint:noctx
	if err != nil {
		return err
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		return closeErr
	}

	return nil
}

// SetupProductSide sets up the webrtc connection on the product project side
// Each product project will have its own candidate id that is known for the platform side
// for identification purposes.
func SetupProductSide(client *http.Client, businessLogic func() string, candiateID int) {
	flag.Parse()

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	// Create a new RTCPeerConnection
	cert, err := webrtc.NewCertificate(Certificates.PrivateKey, *Certificates.Leaf)
	if err != nil {
		panic(err)
	}
	config.Certificates = []webrtc.Certificate{*cert}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(client, *OfferAddr, c); onICECandidateErr != nil {
			panic(onICECandidateErr)
		}
	})

	// A HTTP handler that allows the other Pion instance to send us ICE candidates
	// This allows us to add ICE candidates faster, we don't have to wait for STUN or TURN
	// candidates which may be slower
	http.HandleFunc("/candidate", func(w http.ResponseWriter, r *http.Request) {
		candidate, candidateErr := ioutil.ReadAll(r.Body)
		if candidateErr != nil {
			panic(candidateErr)
		}
		if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(candidate)}); candidateErr != nil {
			panic(candidateErr)
		}
	})

	// A HTTP handler that processes a SessionDescription given to us from the other Pion process
	http.HandleFunc("/sdp", func(w http.ResponseWriter, r *http.Request) {
		sdp := webrtc.SessionDescription{}
		if err := json.NewDecoder(r.Body).Decode(&sdp); err != nil {
			panic(err)
		}

		if err := peerConnection.SetRemoteDescription(sdp); err != nil {
			panic(err)
		}

		// Create an answer to send to the other process
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			panic(err)
		}

		// Send our answer to the HTTP server listening in the other process
		payload, err := json.Marshal(answer)
		if err != nil {
			panic(err)
		}
		resp, err := client.Post(fmt.Sprintf("https://%s/sdp%d", *OfferAddr, candiateID), "application/json; charset=utf-8", bytes.NewReader(payload)) // nolint:noctx
		if err != nil {
			panic(err)
		} else if closeErr := resp.Body.Close(); closeErr != nil {
			panic(closeErr)
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			panic(err)
		}

		candidatesMux.Lock()
		for _, c := range pendingCandidates {
			onICECandidateErr := signalCandidate(client, *OfferAddr, c)
			if onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
		candidatesMux.Unlock()
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", d.Label(), d.ID())

			for range time.NewTicker(SendingFrequency).C {
				message := businessLogic()
				fmt.Printf("Sending '%s'\n", message)

				// Send the message as text
				sendTextErr := d.SendText(message)
				if sendTextErr != nil {
					panic(sendTextErr)
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})
}

// SetupPlatformSide initializes the webrtc connection and sends the initial offer
// to the product project candidates identified by candidate id.
func SetupPlatformSide(client *http.Client, businessLogic func() string, candiateID int) {
	flag.Parse()

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	cert, err := webrtc.NewCertificate(Certificates.PrivateKey, *Certificates.Leaf)
	if err != nil {
		panic(err)
	}
	config.Certificates = []webrtc.Certificate{*cert}
	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else if onICECandidateErr := signalCandidate(client, *AnswerAddr, c); err != nil {
			panic(onICECandidateErr)
		}
	})

	// A HTTP handler that allows the other Pion instance to send us ICE candidates
	// This allows us to add ICE candidates faster, we don't have to wait for STUN or TURN
	// candidates which may be slower
	candidatePath := fmt.Sprintf("/candidate%d", candiateID)
	http.HandleFunc(candidatePath, func(w http.ResponseWriter, r *http.Request) {
		candidate, candidateErr := ioutil.ReadAll(r.Body)
		if candidateErr != nil {
			panic(candidateErr)
		}
		if candidateErr := peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(candidate)}); candidateErr != nil {
			panic(candidateErr)
		}
	})

	// A HTTP handler that processes a SessionDescription given to us from the other Pion process
	sdpPath := fmt.Sprintf("/sdp%d", candiateID)
	http.HandleFunc(sdpPath, func(w http.ResponseWriter, r *http.Request) {
		sdp := webrtc.SessionDescription{}
		if sdpErr := json.NewDecoder(r.Body).Decode(&sdp); sdpErr != nil {
			panic(sdpErr)
		}

		if sdpErr := peerConnection.SetRemoteDescription(sdp); sdpErr != nil {
			panic(sdpErr)
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		for _, c := range pendingCandidates {
			if onICECandidateErr := signalCandidate(client, *AnswerAddr, c); onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
	})

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open. Random messages will now be sent to any connected DataChannels every 5 seconds\n", dataChannel.Label(), dataChannel.ID())

		for range time.NewTicker(SendingFrequency).C {
			message := businessLogic()
			fmt.Printf("Sending '%s'\n", message)

			// Send the message as text
			sendTextErr := dataChannel.SendText(message)
			if sendTextErr != nil {
				panic(sendTextErr)
			}
		}
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Message from DataChannel '%s': '%s'\n", dataChannel.Label(), string(msg.Data))
	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Send our offer to the HTTP server listening in the other process
	payload, err := json.Marshal(offer)
	if err != nil {
		panic(err)
	}
	resp, err := client.Post(fmt.Sprintf("https://%s/sdp", *AnswerAddr), "application/json; charset=utf-8", bytes.NewReader(payload)) // nolint:noctx
	if err != nil {
		panic(err)
	} else if err := resp.Body.Close(); err != nil {
		panic(err)
	}

}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b)
}

func SetupFrontend(w http.ResponseWriter, r *http.Request, offerStr string, dataProvider func() ([]byte, error)) error {
	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		log.Printf("New DataChannel %s %d\n", d.Label(), d.ID())

		// Register channel opening handling
		d.OnOpen(func() {
			log.Printf("Data channel '%s'-'%d' open.\n", d.Label(), d.ID())

			for range time.NewTicker(2 * time.Second).C {
				jsonData, err := dataProvider()
				if err != nil {
					log.Println(err)
					peerConnection.Close()
					break
				}
				if err := d.SendText(string(jsonData)); err != nil {
					log.Println(err)
					peerConnection.Close()
					break
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			log.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	decode(offerStr, &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		peerConnection.Close()
		return err
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		peerConnection.Close()
		return err
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		peerConnection.Close()
		return err
	}

	<-gatherComplete

	_, err = w.Write([]byte(encode(*peerConnection.LocalDescription())))
	if err != nil {
		peerConnection.Close()
		return err
	}
	return nil
}
