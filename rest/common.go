package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var FrontendAddress string = "http://172.18.0.5:8085"
var UserDBAddress string = "http://172.18.0.3:8083"

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

// PrettyPrint logs maps and structs in formatted way in the console.
func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
		return
	}
	fmt.Println("Failed to pretty print data")
}

func (r Request) DecodeRequest() (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (w *ResponseWriter) WriteError(message string, statusCode int) {
	w.WriteResponse(fmt.Sprintf("{\"error\":\"%s\"}", message), statusCode)
}

func (w *ResponseWriter) EncodeResponse(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteResponse(string(b), http.StatusOK)
}

func (w *ResponseWriter) WriteData(data string, statusCode int) {
	w.WriteResponse(fmt.Sprintf("{\"data\":\"%s\"}", data), statusCode)
}

func (w *ResponseWriter) WriteResponse(data string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, data)
}

func MakeHandler(fn func(ResponseWriter, *Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := &Request{request}
		w := ResponseWriter{writer}
		fn(w, r)
	}
}

func (r Request) ForwardRequest(address string) (interface{}, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	proxyReq, err := http.NewRequest(r.Method, fmt.Sprintf("%s%s", address, r.RequestURI), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &dataMap); err != nil {
		return nil, err
	}

	if val, ok := dataMap["error"]; ok {
		return nil, errors.New(val.(string))
	}

	if val, ok := dataMap["data"]; ok {
		return val, nil
	}

	return nil, errors.New("Invalid response")
}

func Get(address string, path string, parameters string) (interface{}, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", address, path, parameters))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &dataMap); err != nil {
		return nil, err
	}

	if val, ok := dataMap["error"]; ok {
		return nil, errors.New(val.(string))
	}

	if val, ok := dataMap["data"]; ok {
		return val, nil
	}

	return nil, errors.New("Invalid response")
}

func Post(address string, path string, parameters map[string]interface{}) (interface{}, error) {
	reqBody, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", address, path), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &dataMap); err != nil {
		return nil, err
	}

	if _, ok := dataMap["error"]; ok {
		return nil, errors.New(dataMap["error"].(string))
	}

	if val, ok := dataMap["data"]; ok {
		return val, nil
	}

	return nil, errors.New("Invalid response")
}
