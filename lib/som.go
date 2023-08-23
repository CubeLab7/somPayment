package lib

import (
	"context"
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

func PrepareRequest(ctx context.Context, method string, body io.Reader) (req *http.Request, err error) {
	//url := fmt.Sprintf("%v/%v", config.URI, method)

	req, err = http.NewRequest(http.MethodPost, "", body)
	if err != nil {
		err = fmt.Errorf("can't create request for Som payment system: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	//req.Header.Add("Authorization", basicAuth(config.URI, config.Pass))
	return
}

func SendRequest(req *http.Request, response interface{}) (err error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			//IdleConnTimeout: time.Second * time.Duration(config.IdleConnTimeoutSec),
		},
		//Timeout: time.Second * time.Duration(config.RequestTimeoutSec),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("can't do request: %s", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("can't read response body: %s", err)
		return
	}

	if err = jsoniter.Unmarshal(respBody, &response); err != nil {
		err = fmt.Errorf("can't unmarshall SomPayments resp: '%v'. Err: %w", string(respBody), err)
		return
	}

	return
}

func basicAuth(login, pass string) (basic string) {
	basic = "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", login, pass)))
	return
}
