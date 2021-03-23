package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type Controller struct {
	TestData    map[string]interface{}
	RequestData map[string]interface{}
}

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
		return
	}
	fmt.Println("Failed to pretty print data")
}

func (w *ResponseWriter) writeError(message string, statusCode int) {
	w.writeResponse(fmt.Sprintf("{\"error\":\"%s\"}", message), statusCode)
}

func (w *ResponseWriter) writeData(data string, statusCode int) {
	w.writeResponse(fmt.Sprintf("{\"data\":\"%s\"}", data), statusCode)
}

func (w *ResponseWriter) writeResponse(data string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, data)
}

func makeHandler(fn func(ResponseWriter, *Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := &Request{request}
		w := ResponseWriter{writer}
		fn(w, r)
	}
}

func (c *Controller) clearRequestData(w ResponseWriter, r *Request) {
	for k := range c.RequestData {
		delete(c.RequestData, k)
	}
}

func (c *Controller) getRequestData(w ResponseWriter, r *Request) {
	w.encodeResponse(c.RequestData, http.StatusOK)
}

func (c Controller) ParseForm(r *Request, requestPath string) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if requestPath != "" {
		c.RequestData[requestPath] = r.RequestURI
	}
	return nil
}

func (c *Controller) decodeRequest(r *Request, requestPath string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if requestPath != "" {
		c.RequestData[requestPath] = data
	}
	return data, nil
}

func (w *ResponseWriter) encodeResponse(data interface{}, statusCode int) {
	b, err := json.Marshal(data)
	if err != nil {
		w.writeError(fmt.Sprintf("Backend: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.writeResponse(string(b), statusCode)
}

func (c *Controller) updateTestData(w ResponseWriter, r *Request) {
	requestData, err := c.decodeRequest(r, "")
	if err != nil {
		w.writeError(fmt.Sprintf("Backend: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	c.TestData = requestData
}

func uploadFile(key string, fileName string, r *http.Request) error {
	file, handler, err := r.FormFile(key)
	if err == http.ErrMissingFile {
		return nil
	}

	if err != nil {
		return err
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	defer file.Close()

	// Create file
	dst, err := os.Create(fmt.Sprintf("/user-assets/uploads/%s", fileName))
	if err != nil {
		return err
	}

	// Copy the uploaded file to the created file on the file system.
	if _, err := io.Copy(dst, file); err != nil {
		if err2 := dst.Close(); err2 != nil {
			err = errors.Wrap(errors.WithStack(err), err2.Error())
		}
		return err
	}
	dst.Close()

	return nil
}
