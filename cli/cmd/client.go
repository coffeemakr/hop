package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/coffeemakr/ruck"
)

var (
	ErrNoTokenSaved = errors.New("no saved token")
	ErrNotFound     = errors.New("item not found")
)

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
	token      string
}

func (c *Client) getUrl(relativeUrl string) string {
	return c.BaseUrl + relativeUrl
}

func (c *Client) newRequest(method string, relativeUrl string, authenticationToken string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.getUrl(relativeUrl), body)
	if err != nil {
		return nil, err
	}
	if authenticationToken != "" {
		req.Header.Set("Authorization", "Bearer "+authenticationToken)
	}
	return req, err
}

func checkResponse(response *http.Response) error {
	if response.StatusCode >= 300 || response.StatusCode < 100 {
		if response.StatusCode == 404 {
			return ErrNotFound
		}
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
		return fmt.Errorf("request failed with status code %d \n\n%s", response.StatusCode, description)
	}
	return nil
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
	if err := checkResponse(response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Client) receiveJsonAuthenticated(method string, relativeUrl string, result interface{}) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	req, err := c.newRequest(method, relativeUrl, token, nil)
	if err != nil {
		return err
	}
	response, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if err := checkResponse(response); err != nil {
		return err
	}
	// Decode response
	err = json.NewDecoder(response.Body).Decode(result)
	if err != nil {
		return fmt.Errorf("parsing response: %s", err)
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

func (c *Client) Login(credentials *ruck.Credentials) error {
	var authenticationResult ruck.AuthenticationResult
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

func (c *Client) Register(request *ruck.RegistrationRequest) (*ruck.User, error) {
	var user ruck.User
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
	group := ruck.Group{
		Name: name,
	}
	token, err := c.Token()
	if err != nil {
		return err
	}
	err = c.sendAndReceiveJson("POST", "/groups", token, &group, &group)
	if err != nil {
		return fmt.Errorf("failed to create group: %s", err)
	}
	return nil
}

func (c *Client) ListGroup() (results []*ruck.Group, err error) {
	err = c.receiveJsonAuthenticated("GET", "/groups", &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) Token() (string, error) {
	if c.token == "" {
		err := c.LoadToken()
		if err != nil {
			return "", err
		}
	}
	return c.token, nil
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
	token, err := c.Token()
	if err != nil {
		return err
	}
	err = c.send("DELETE", joinUrl("groups", id), token)
	if err != nil {
		return fmt.Errorf("group deletion failed: %s", err)
	}
	return nil
}

func (c *Client) JoinGroup(groupId string) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	err = c.send("POST", joinUrl("groups", groupId, "join"), token)
	if err != nil {
		return fmt.Errorf("failed to join group: %s", err)
	}
	return nil
}

func (c *Client) CreateTask(task *ruck.Task) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	groupId := task.GroupID
	if groupId == "" {
		return errors.New("group ID not set")
	}
	err = c.sendAndReceiveJson("POST", joinUrl("groups", groupId, "tasks"), token, task, task)
	if err != nil {
		return fmt.Errorf("creation of task failed: %s", err)
	}
	return nil
}

func (c *Client) GetTaskList() (tasks []*ruck.Task, err error) {
	err = c.receiveJsonAuthenticated("GET", "/tasks", &tasks)
	if err != nil {
		err = fmt.Errorf("failed to get list of tasks: %s", err)
	}
	return
}

func (c *Client) GetTaskDetails(taskId string) (*ruck.Task, error) {
	var task ruck.Task
	var err error
	err = c.receiveJsonAuthenticated("GET", joinUrl("tasks", taskId), &task)
	if err != nil {
		err = fmt.Errorf("failed to get task: %s", err)
		return nil, err
	}
	return &task, nil
}

func (c *Client) CompleteTask(taskID string) (execution *ruck.TaskExecution, err error) {
	execution = new(ruck.TaskExecution)
	err = c.receiveJsonAuthenticated("POST", joinUrl("tasks", taskID, "complete"), execution)
	return
}
