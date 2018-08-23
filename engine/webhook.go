package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// this is only for demo purposes, real event would go here
const (
	WebhookEventUserDetail = "user_detail"
)

// Post makes an HTTP POST request
func Post(url string, data interface{}, result interface{}, headers map[string]string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "applicaiton/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("failed post, error: ", err)
		} else {
			log.Println("failed post, recv: ", string(b))
		}
		return fmt.Errorf("error requesting %s returned %s", url, resp.Status)
	}

	if result != nil {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(result); err != nil {
			return err
		}
	}
	return nil
}
