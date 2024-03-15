package misskey

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

	// リクエストを作成
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, errors.New("failed to create request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	// リクエストを送信
	client := &http.Client{}
	res, err := client.Do(req)
	return res, err
}
