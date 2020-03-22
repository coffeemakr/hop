package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coffeemakr/wedo"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
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
	token string
}

func (c *Client) getUrl(relativeUrl string) string {
	return c.BaseUrl + relativeUrl
}

func (c *Client) newRequest(method string, relativeUrl string, authenticationToken string, body io.Reader) (*http.Request, error){
	req, err := http.NewRequest(method, c.getUrl(relativeUrl), body)
	if err != nil {
		return nil, err
	}
	if authenticationToken != "" {
		req.Header.Set("Authorization", "Bearer " + authenticationToken)
	}
	return req, err
}

func (c *Client) sendJson(method string, relativeUrl string, authenticationToken string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("JSON creation failed: %s", err)
	}
	req, err := c.newRequest(method, relativeUrl, authenticationToken, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request sending failed: %s", err)
	}
	if response.StatusCode >= 300 || response.StatusCode < 100 {
		contentType := response.Header.Get("Content-Type")
		var description string
		if contentType == "application/json" {
			var jsonError map[string]interface{}
			decoder := json.NewDecoder(response.Body)
			err := decoder.Decode(&jsonError)
			if err == nil {
				description = jsonError["description"].(string)
			}
		} else {
			descriptionBytes, err := ioutil.ReadAll(response.Body)
			if err == nil {
				description = string(descriptionBytes[:])
			}
		}
		return nil, fmt.Errorf("request failed with status code %d \n\n%s", response.StatusCode, description)
	}
	return response, nil
}

func (c *Client) receiveJson(method string, relativeUrl string, authenticationToken string, result interface{}) error {
	req, err := c.newRequest(method, relativeUrl, authenticationToken, nil)
	if err != nil {
		return err
	}
	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	// Decode response
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) sendAndReceiveJson(method string, relativeUrl string, authenticationToken string, body interface{}, result interface{}) error {
	response, err := c.sendJson(method, relativeUrl, authenticationToken, body)
	if err != nil {
		return fmt.Errorf("failed to send JSON: %s", err)
	}
	// Decode response
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("failed to read response JSON: %s", err)
	}
	return nil
}

func (c *Client) Login(credentials *wedo.Credentials) error {
	var authenticationResult wedo.AuthenticationResult
	err := c.sendAndReceiveJson("POST", "/login", "", credentials, &authenticationResult)
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
	err := c.sendAndReceiveJson("POST", "/register", "", request, &user)
	return &user, err
}

func (c *Client) LoadToken() error {
	token, err := c.TokenStore.GetToken()
	if err != nil {
		return err
	}
	c.token = token
	return nil
}

func (c *Client) CreateGroup(name string) error {
	group := wedo.Group{
		Name: name,
	}
	err := c.sendAndReceiveJson("POST", "/groups", c.Token(), &group, &group)
	if err != nil {
		return fmt.Errorf("failed to create group: %s", err)
	}
	return nil
}

func (c *Client) ListGroup() ([]*wedo.Group, error) {
	var results []*wedo.Group
	err := c.receiveJson("GET", "/groups", c.Token(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) Token() string {
	return c.token
}

func joinUrl(parts ...string) string {
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return "/" + path.Join(parts...)
}



func (c *Client) send(method string, relativeUrl string, authenticationToken string) error {
	request, err := c.newRequest(method, relativeUrl, authenticationToken, nil)
	if err != nil {
		return err
	}
	resp, err := c.Client.Do(request)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 || resp.StatusCode < 100 {
		return fmt.Errorf("request failed with status '%s'", resp.Status)
	}
	return nil
}

func (c *Client) DeleteGroupByID(id string) error {
	err := c.send("DELETE", joinUrl("groups", id), c.Token())
	if err != nil {
		return fmt.Errorf("group deletion failed: %s", err)
	}
	return nil
}

func (c *Client) JoinGroup(groupId string) error {
	err := c.send("POST", joinUrl("groups", groupId, "join"), c.Token())
	if err != nil {
		return fmt.Errorf("failed to join group: %s", err)
	}
	return nil
}

func (c *Client) CreateTask(task *wedo.Task) error {
	_, err := c.sendJson("POST", "/tasks", c.Token(), task)
	if err != nil {
		return fmt.Errorf("creation of task failed: %s", err)
	}
}
