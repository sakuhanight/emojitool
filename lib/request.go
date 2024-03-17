package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

func Request(endpoint string, body interface{}) (*http.Response, error) {
	// bodyをjsonに変換
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("failed to marshal body: " + err.Error())
	}
	return RequestRaw(endpoint, "application/json", bytes.NewBuffer(jsonBody))
}

func RequestRaw(endpoint string, contentType string, body *bytes.Buffer) (*http.Response, error) {
	// リクエストを作成
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	req.Header.Set("Content-Type", contentType)

	// リクエストを送信
	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}
