package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func New() Client {
	var client Client
	client.URL = "http://localhost:8080/api"

	return client
}

type Client struct {
	URL string
}

func (c *Client) RepoList() (string, error) {
	resp, err := http.Post(c.URL+"/repos/list", "application/json", strings.NewReader(""))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (c *Client) RepoAdd(repos string) (string, error) {
	resp, err := http.Post(c.URL+"/repos/add", "application/json", strings.NewReader(fmt.Sprintf("{\"name\": \"%s\"}", repos)))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (c *Client) RepoRemove(repos string) (string, error) {
	resp, err := http.Post(c.URL+"/repos/remove", "application/json", strings.NewReader(fmt.Sprintf("{\"name\": \"%s\"}", repos)))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
