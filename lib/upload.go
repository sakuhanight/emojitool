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
	filename := filepath.Base(path)
	endpoint := fmt.Sprintf("https://%s/api/drive/files/create", host)

	file, err := os.Open(path)
	if err != nil {
		zap.S().Panicf("Could not open '%s': %v", path, err)
	}

	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)

	{
		tokenPart, _ := mw.CreateFormField("i")
		_, err = tokenPart.Write([]byte(token))
		if err != nil {
			zap.S().Panicf("CreatePart Failed")
		}
	}
	{
		filePart, _ := mw.CreateFormFile("file", filename)
		_, err = io.Copy(filePart, file)
		if err != nil {
			zap.S().Panicf("CreatePart Failed")
		}
	}

	mw.Close()

	return RequestRaw(endpoint, mw.FormDataContentType(), body)
}
