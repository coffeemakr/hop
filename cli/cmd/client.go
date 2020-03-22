package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coffeemakr/wedo"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	BaseUrl string
	Client *http.Client
}

func (c *Client) getUrl(relativeUrl string) string {
	return c.BaseUrl + relativeUrl
}

func (c *Client) sendAndReceiveJson(method string, relativeUrl string, body interface{}, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.getUrl(relativeUrl), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode >= 300 || response.StatusCode < 100 {
		body, err = ioutil.ReadAll(response.Body)
 		return fmt.Errorf("request failed with status code %d \n\n%s", response.StatusCode, body)
	}
	// Decode response
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Login(credentials *wedo.Credentials) error {
	var user wedo.User
	err := c.sendAndReceiveJson("POST", "/login", credentials, &user)
	if err == nil {
		log.Println(user)
	}
	return err
}

func (c *Client) Register(request *wedo.RegistrationRequest) error {
	var user wedo.User
	err := c.sendAndReceiveJson("POST", "/register", request, &user)
	if err == nil {
		log.Println(user)
	}
	return err
}
