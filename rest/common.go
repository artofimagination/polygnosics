package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	httpModels "github.com/artofimagination/polygnosics/models/http"
)

type Server struct {
	IP   string
	Port int
	Name string
}

func (c *Server) GetAddress() string {
	return fmt.Sprintf("http://%s:%d", c.IP, c.Port)
}

type ResponseWriter struct {
	http.ResponseWriter
}

type RequestInterface interface {
	FormFile(key string) (MultiPartFileImpl, *multipart.FileHeader, error)
	ParseForm() error
	FormValue(key string) string
	ParseMultipartForm(maxMemory int64) error
}

// File is the interface class to redefine multipart.File with custom, i.e. mock implementations
type MultiPartFile interface {
	Close() error
	Read(p []byte) (n int, err error)
}

type MultiPartFileImpl struct {
	MultiPartFile
}

type Request struct {
	request *http.Request
}

func (r Request) FormFile(key string) (MultiPartFileImpl, *multipart.FileHeader, error) {
	file, handler, err := r.request.FormFile(key)
	return MultiPartFileImpl{file}, handler, err
}

func (r *Request) ParseForm() error {
	return r.request.ParseForm()
}

func (r *Request) FormValue(key string) string {
	return r.request.FormValue(key)
}

func (r *Request) ParseMultipartForm(maxMemory int64) error {
	return r.request.ParseMultipartForm(maxMemory)
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
	if err := json.NewDecoder(r.request.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (w *ResponseWriter) WriteError(message string, statusCode int) {
	response := httpModels.ResponseData{Error: message}
	b, err := json.Marshal(response)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.WriteResponse(string(b), statusCode)
}

func (w *ResponseWriter) EncodeResponse(data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteError(fmt.Sprintf("Backend -> %s", err.Error()), http.StatusInternalServerError)
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
	body, err := ioutil.ReadAll(r.request.Body)
	if err != nil {
		return nil, err
	}

	r.request.Body = ioutil.NopCloser(bytes.NewReader(body))
	proxyReq, err := http.NewRequest(r.request.Method, fmt.Sprintf("%s%s", address, r.request.RequestURI), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for header, values := range r.request.Header {
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

	responseContent := httpModels.ResponseData{}
	if err := json.Unmarshal(body, &responseContent); err != nil {
		return nil, err
	}

	if responseContent.Error != "" {
		return nil, errors.New(responseContent.Error)
	}

	return responseContent.Data, nil
}

func Post(address string, path string, parameters interface{}) (interface{}, error) {
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

	responseContent := httpModels.ResponseData{}
	if err := json.Unmarshal(body, &responseContent); err != nil {
		return nil, err
	}

	if responseContent.Error != "" {
		return nil, errors.New(responseContent.Error)
	}

	return responseContent.Data, nil
}
