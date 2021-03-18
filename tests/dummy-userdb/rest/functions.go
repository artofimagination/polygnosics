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

func prettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
		return
	}
	fmt.Println("Failed to pretty print data")
}

func writeError(message string, w http.ResponseWriter, statusCode int) {
	writeResponse(fmt.Sprintf("{\"error\":\"%s\"}", message), w, statusCode)
}

func writeData(data string, w http.ResponseWriter, statusCode int) {
	writeResponse(fmt.Sprintf("{\"data\":\"%s\"}", data), w, statusCode)
}

func writeResponse(data string, w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, data)
}

func (c *Controller) getRequestData(w http.ResponseWriter, r *http.Request) {
	encodeResponse(c.RequestData, w)
}

func (c Controller) ParseForm(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	c.setRequestData(r)
	return nil
}

func (c *Controller) setRequestData(r *http.Request) {
	for k := range c.RequestData {
		delete(c.RequestData, k)
	}

	c.RequestData["uri"] = r.RequestURI
}

func (c *Controller) decodeRequest(r *http.Request) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	c.RequestData = data
	return data, nil
}

func encodeResponse(data map[string]interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(data)
	if err != nil {
		writeError(fmt.Sprintf("Backend: %s", err.Error()), w, http.StatusInternalServerError)
		return
	}
	writeResponse(string(b), w, http.StatusOK)
}

func (c *Controller) updateTestData(w http.ResponseWriter, r *http.Request) {
	requestData, err := c.decodeRequest(r)
	if err != nil {
		writeError(fmt.Sprintf("Backend: %s", err.Error()), w, http.StatusInternalServerError)
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
