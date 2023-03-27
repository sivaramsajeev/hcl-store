package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sivaramsajeev/terraform-provider-student/api/server"
)

type Client struct {
	hostname   string
	port       int
	authToken  string
	httpClient *http.Client
}

func NewClient(hostname string, port int, token string) *Client {
	return &Client{
		hostname:   hostname,
		port:       port,
		authToken:  token,
		httpClient: &http.Client{},
	}
}

func (c *Client) httpRequest(path, method string, body bytes.Buffer) (closer io.ReadCloser, err error) {
	req, err := http.NewRequest(method, c.requestPath(path), &body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.authToken)
	switch method {
	case "GET":
	case "DELETE":
	default:
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("got a non 200 status code: %v", resp.StatusCode)
		}
		return nil, fmt.Errorf("got a non 200 status code: %v - %s", resp.StatusCode, respBody.String())
	}
	return resp.Body, nil
}

func (c *Client) requestPath(path string) string {
	return fmt.Sprintf("http://%s:%v/%s", c.hostname, c.port, path)
}

func (c *Client) GetAll() (*map[string]server.Student, error) {
	body, err := c.httpRequest("student", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	students := map[string]server.Student{}
	err = json.NewDecoder(body).Decode(&students)
	if err != nil {
		return nil, err
	}
	return &students, nil
}

func (c *Client) GetStudent(name string) (*server.Student, error) {
	body, err := c.httpRequest(fmt.Sprintf("student/%v", name), "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	student := &server.Student{}
	err = json.NewDecoder(body).Decode(student)
	if err != nil {
		return nil, err
	}
	return student, nil
}

func (c *Client) NewStudent(student *server.Student) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(student)
	if err != nil {
		return err
	}
	_, err = c.httpRequest("student", "POST", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateStudent(student *server.Student) error {
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(student)
	if err != nil {
		return err
	}
	_, err = c.httpRequest(fmt.Sprintf("student/%s", student.Name), "PUT", buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteStudent(studentName string) error {
	_, err := c.httpRequest(fmt.Sprintf("student/%s", studentName), "DELETE", bytes.Buffer{})
	if err != nil {
		return err
	}
	return nil
}
