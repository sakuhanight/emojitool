package lib

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func UploadExample(host string, token string, path string) (*http.Response, error) {
	endpoint := fmt.Sprintf("https://%s/api/drive/files/create", host)
	zap.S().Debugln("UploadExample called")
	var options []MultipartRequestOption
	options = append(options, SetMultipartField("i", []byte(token)))
	options = append(options, SetMultipartField("force", []byte("true")))
	options = append(options, SetMultipartFile("file", path))

	return MultipartRequest(endpoint, options...)
}

type (
	MultipartRequestOption func(writer *multipart.Writer)
)

func MultipartRequest(endpoint string, options ...MultipartRequestOption) (*http.Response, error) {
	zap.S().Debugln("MultipartRequest called")
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for _, opt := range options {
		opt(w)
	}
	err := w.Close()
	if err != nil {
		zap.S().Fatalf("MultipartRequest writer close failed: %+v", err)
	}
	zap.S().Debugln("MultipartRequest for end")

	return RequestRaw(endpoint, w.FormDataContentType(), body)
}

func SetMultipartField(field string, data []byte) MultipartRequestOption {
	return func(writer *multipart.Writer) {
		part, _ := writer.CreateFormField(field)
		_, err := part.Write(data)
		if err != nil {
			zap.S().Panicf("CreatePart Failed")
		}
	}
}

func SetMultipartFile(field string, path string) MultipartRequestOption {
	return func(writer *multipart.Writer) {
		filename := filepath.Base(path)
		file, err := os.Open(path)
		if err != nil {
			zap.S().Panicf("Could not open '%s': %v", path, err)
		}
		part, _ := writer.CreateFormFile(field, filename)
		_, err = io.Copy(part, file)
		if err != nil {
			zap.S().Panicf("CreatePart Failed")
		}
	}
}
