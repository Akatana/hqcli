package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewClient(token string) *Client {
	return &Client{token}
}

func (c *Client) request(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, "https://api.hellohq.io/v2"+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s - %s", res.Status, string(responseData))
	}

	return responseData, nil
}

func (c *Client) GetProjects() ([]Project, error) {
	data, err := c.request("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(data, &projects)
	return projects, err
}

func (c *Client) GetTasksForProject(projectID int) ([]Task, error) {
	data, err := c.request("GET", fmt.Sprintf("/projects/%d/tasks", projectID), nil)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func (c *Client) CreateUserReporting(report UserReporting) error {
	_, err := c.request("POST", "/userreportings", report) // Replace with correct endpoint
	return err
}
