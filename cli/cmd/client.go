package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coffeemakr/wedo"
	"io/ioutil"
	"net/http"
	"os"
)

var ErrNoTokenSaved = errors.New("no saved token")

type TokenStore interface {
	SaveToken(token string) error
	GetToken() (string, error)
}

type fileTokenStore struct {
	Path string
}

func NewFileTokenStore(path string) TokenStore {
	return &fileTokenStore{Path: path}
}

func (s *fileTokenStore) SaveToken(token string) error {
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	_, err = file.WriteString(token)
	if err != nil {
		return err
	}
	return nil
}

func (s *fileTokenStore) GetToken() (string, error) {
	file, err := os.Open(s.Path)
	if err != nil {
		if err == os.ErrNotExist {
			err = ErrNoTokenSaved
		}
		return "", err
	}
	token, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(token[:]), nil
}

type Client struct {
	BaseUrl    string
	Client     *http.Client
	TokenStore TokenStore
	Token string
}

func (c *Client) getUrl(relativeUrl string) string {
	return c.BaseUrl + relativeUrl
}

func (c *Client) sendAndReceiveJson(method string, relativeUrl string, body interface{}, result interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, c.getUrl(relativeUrl), bytes.NewReader(bodyBytes))
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
	var authenticationResult wedo.AuthenticationResult
	err := c.sendAndReceiveJson("POST", "/login", credentials, &authenticationResult)
	if err != nil {
		return err
	}
	err = c.TokenStore.SaveToken(authenticationResult.Token)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Register(request *wedo.RegistrationRequest) (*wedo.User, error) {
	var user wedo.User
	err := c.sendAndReceiveJson("POST", "/register", request, &user)
	return &user, err
}

func (c *Client) LoadToken() error {
	token, err := c.TokenStore.GetToken()
	if err != nil {
		return err
	}
	c.Token = token
	return nil
}
