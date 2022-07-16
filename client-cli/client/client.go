package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func postWrapper(url string, body string) (string, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func New() Client {
	var client Client
	client.URL = "http://localhost:8080/api"

	return client
}

type Client struct {
	URL string
}

func (c *Client) TransactionGet(name string) (string, error) {
	result, err := postWrapper(c.URL+"/transaction/get", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	if err != nil {
		return "", err
	}

	return result, nil
}

func (c *Client) TransactionAdd(name string) error {
	_, err := postWrapper(c.URL+"/transaction/add", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	return err
}

func (c *Client) TransactionRemove(name string) error {
	_, err := postWrapper(c.URL+"/transaction/remove", fmt.Sprintf("{ \"name\": \"%s\" }", name))
	return err
}

func (c *Client) TransactionRename(oldName string, newName string) error {
	_, err := postWrapper(c.URL+"/transaction/rename", fmt.Sprintf("{ \"oldName\": \"%s\", \"newName\": \"%s\" }", oldName, newName))
	return err
}

func (c *Client) TransactionList() (string, error) {
	result, err := postWrapper(c.URL+"/transaction/list", "{}")
	return result, err
}

func (c *Client) RepoAdd(transaction string, repo string) (string, error) {
	result, err := postWrapper(c.URL+"/repo/add", fmt.Sprintf("{\"transaction\": \"%s\", \"repo\": \"%s\"}", transaction, repo))
	return result, err
}

func (c *Client) RepoRemove(transaction string, repo string) (string, error) {
	result, err := postWrapper(c.URL+"/repos/remove", fmt.Sprintf("{\"transaction\": \"%s\", \"repo\": \"%s\"}", transaction, repo))
	return result, err
}

func (c *Client) RepoList(transaction string) (string, error) {
	result, err := postWrapper(c.URL+"/repos/list", fmt.Sprintf("{\"transaction\": \"%s\"}", transaction))
	return result, err
}
