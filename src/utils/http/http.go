package http

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendJsonRequest(
	method string, 
	url string, 
	data map[string]any, 
	headers map[string]string,
) (*http.Response, error)  {
	reqData, err := json.Marshal(data)
	if (err != nil) {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(reqData))
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
    }
	res, err := new(http.Client).Do(req)
	if (err != nil) {
		return nil, err
	}

	return res, nil
}